package ast

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type EnvVar string
type Local string
type Remote struct{ url *url.URL }
type Missing struct{}

const NullOrigin = "null"

type Fetchable interface {
	Name() string
	Origin() string
	//Fetch
	// fetches the import
	// the `origin` parameter should be `scheme://authority` or NullOrigin
	Fetch(origin string) (string, error)
	ChainOnto(base Fetchable) (Fetchable, error)
	String() string
}

var _ Fetchable = EnvVar("")
var _ Fetchable = Local("")
var _ Fetchable = Remote{}
var _ Fetchable = Missing{}

func (e EnvVar) Name() string { return string(e) }
func (EnvVar) Origin() string { return NullOrigin }
func (e EnvVar) String() string {
	return "env:" + string(e)
}
func (e EnvVar) Fetch(origin string) (string, error) {
	if origin != NullOrigin {
		return "", errors.New("Can't access environment variable from remote import")
	}
	val, ok := os.LookupEnv(string(e))
	if !ok {
		return "", fmt.Errorf("Unset environment variable %s", string(e))
	}
	return val, nil
}
func (e EnvVar) ChainOnto(base Fetchable) (Fetchable, error) {
	return e, nil
}

func (l Local) Name() string { return string(l) }
func (Local) Origin() string { return NullOrigin }
func (l Local) String() string {
	if l.IsAbs() || l.IsRelativeToHome() || l.IsRelativeToParent() {
		return string(l)
	} else {
		return "./" + string(l)
	}
}
func (l Local) Fetch(origin string) (string, error) {
	if origin != NullOrigin {
		return "", fmt.Errorf("Can't get %s from remote import at %s", l, origin)
	}
	bytes, err := ioutil.ReadFile(string(l))
	return string(bytes), err
}
func (l Local) ChainOnto(base Fetchable) (Fetchable, error) {
	switch r := base.(type) {
	case Local:
		if l.IsAbs() || l.IsRelativeToHome() {
			return l, nil
		}
		return Local(path.Join(path.Dir(string(r)), string(l))), nil
	case Remote:
		if l.IsAbs() {
			return nil, errors.New("Can't get absolute path from remote import")
		}
		if l.IsRelativeToHome() {
			return nil, errors.New("Can't get home-relative path from remote import")
		}
		relativeURL, err := url.Parse(string(l))
		if err != nil {
			return nil, err
		}
		newURL := r.url.ResolveReference(relativeURL)
		return Remote{url: newURL}, nil
	default:
		return l, nil
	}
}

func (l Local) IsAbs() bool              { return path.IsAbs(string(l)) }
func (l Local) IsRelativeToParent() bool { return strings.HasPrefix(string(l), "..") }
func (l Local) IsRelativeToHome() bool   { return string(l)[0] == '~' }

func (l Local) PathComponents() []string {
	if l.IsAbs() || l.IsRelativeToHome() || l.IsRelativeToParent() {
		return strings.Split(string(l), "/")[1:]
	} else {
		return strings.Split(string(l), "/")
	}
}

func MakeRemote(u *url.URL) (Remote, error) {
	if u.EscapedPath() == "/" || u.EscapedPath() == "" {
		return Remote{}, errors.New("URLs must have a nonempty path")
	}
	return Remote{url: u}, nil
}

var client http.Client

func (r Remote) Name() string   { return r.url.String() }
func (r Remote) Origin() string { return fmt.Sprintf("%s://%s", r.url.Scheme, r.Authority()) }
func (r Remote) String() string { return fmt.Sprintf("%v", r.url) }
func (r Remote) Fetch(origin string) (string, error) {
	req, err := http.NewRequest("GET", r.url.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "dhall-golang")
	corsFlag := origin != NullOrigin && origin != r.Origin()
	if corsFlag {
		req.Header.Set("Origin", origin)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Got status %d from URL %s", resp.StatusCode, r.url)
	}
	if corsFlag &&
		resp.Header.Get("Access-Control-Allow-Origin") != "*" &&
		resp.Header.Get("Access-Control-Allow-Origin") != origin {
		return "", fmt.Errorf("URL %s does not permit CORS requests from %s", r.url, origin)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	return string(bodyBytes), err
}
func (r Remote) ChainOnto(base Fetchable) (Fetchable, error) {
	return r, nil
}
func (r Remote) IsPlainHttp() bool { return r.url.Scheme == "http" }
func (r Remote) Authority() string {
	if r.url.User != nil {
		return fmt.Sprintf("%s@%s", r.url.User.String(), r.url.Host)
	}
	return r.url.Host
}
func (r Remote) PathComponents() []string {
	escapedComps := strings.Split(r.url.EscapedPath()[1:], "/")
	unescapedComps := make([]string, len(escapedComps))
	for i, comp := range escapedComps {
		var err error
		unescapedComps[i], err = url.PathUnescape(comp)
		if err != nil {
			// can't happen, surely
			panic(fmt.Sprintf("Got error %v", err))
		}
	}
	return unescapedComps
}
func (r Remote) Query() *string {
	if r.url.RawQuery == "" && !r.url.ForceQuery {
		return nil
	}
	return &r.url.RawQuery
}

func (Missing) Name() string   { return "" }
func (Missing) Origin() string { return NullOrigin }
func (Missing) String() string { return "missing" }
func (Missing) Fetch(origin string) (string, error) {
	return "", errors.New("Cannot resolve missing import")
}
func (Missing) ChainOnto(base Fetchable) (Fetchable, error) {
	return Missing{}, nil
}
