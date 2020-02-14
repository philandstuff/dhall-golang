package term

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
type LocalFile string
type RemoteFile struct{ url *url.URL }
type Missing struct{}

// NullOrigin is used in Fetchable.Fetch() to indicate no origin.
const NullOrigin = "null"

var LocationType = UnionType{
	"Local":       Text,
	"Remote":      Text,
	"Environment": Text,
	"Missing":     nil,
}

type Fetchable interface {
	Name() string
	Origin() string
	//Fetch
	// fetches the import
	// the `origin` parameter should be `scheme://authority` or NullOrigin
	Fetch(origin string) (string, error)
	ChainOnto(base Fetchable) (Fetchable, error)
	String() string
	AsLocation() Term
}

var _ Fetchable = EnvVar("")
var _ Fetchable = LocalFile("")
var _ Fetchable = RemoteFile{}
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
func (e EnvVar) AsLocation() Term {
	return Apply(Field{LocationType, "Environment"}, TextLit{Suffix: e.String()})
}

func (l LocalFile) Name() string { return string(l) }
func (LocalFile) Origin() string { return NullOrigin }
func (l LocalFile) String() string {
	if l.IsAbs() || l.IsRelativeToHome() || l.IsRelativeToParent() {
		return string(l)
	} else {
		return "./" + string(l)
	}
}
func (l LocalFile) Fetch(origin string) (string, error) {
	if origin != NullOrigin {
		return "", fmt.Errorf("Can't get %s from remote import at %s", l, origin)
	}
	bytes, err := ioutil.ReadFile(string(l))
	return string(bytes), err
}
func (l LocalFile) ChainOnto(base Fetchable) (Fetchable, error) {
	switch r := base.(type) {
	case LocalFile:
		if l.IsAbs() || l.IsRelativeToHome() {
			return l, nil
		}
		return LocalFile(path.Join(path.Dir(string(r)), string(l))), nil
	case RemoteFile:
		if l.IsAbs() {
			return nil, errors.New("Can't get absolute path from remote import")
		}
		if l.IsRelativeToHome() {
			return nil, errors.New("Can't get home-relative path from remote import")
		}
		newURL := r.url.ResolveReference(l.asRelativeRef())
		return RemoteFile{url: newURL}, nil
	default:
		return l, nil
	}
}

func (l LocalFile) IsAbs() bool              { return path.IsAbs(string(l)) }
func (l LocalFile) IsRelativeToParent() bool { return strings.HasPrefix(string(l), "..") }
func (l LocalFile) IsRelativeToHome() bool   { return string(l)[0] == '~' }

//asRelativeRef converts a local path to a relative reference
func (l LocalFile) asRelativeRef() *url.URL {
	if l.IsAbs() || l.IsRelativeToHome() {
		panic("Can't convert absolute or home-relative path to relative reference")
	}
	var s strings.Builder
	if l.IsRelativeToParent() {
		s.WriteString("..")
	} else {
		s.WriteString(".")
	}
	for _, segment := range l.PathComponents() {
		s.WriteString("/")
		s.WriteString(url.PathEscape(segment))
	}
	u, err := url.Parse(s.String())
	if err != nil {
		panic(err)
	}
	return u
}

func (l LocalFile) PathComponents() []string {
	if l.IsAbs() || l.IsRelativeToHome() || l.IsRelativeToParent() {
		return strings.Split(string(l), "/")[1:]
	} else {
		return strings.Split(string(l), "/")
	}
}
func (l LocalFile) AsLocation() Term {
	return Apply(Field{LocationType, "Local"}, TextLit{Suffix: l.String()})
}
func NewRemoteFile(u *url.URL) RemoteFile {
	return RemoteFile{url: u}
}

var client http.Client

func (r RemoteFile) Name() string   { return r.url.String() }
func (r RemoteFile) Origin() string { return fmt.Sprintf("%s://%s", r.url.Scheme, r.Authority()) }
func (r RemoteFile) String() string { return fmt.Sprintf("%v", r.url) }
func (r RemoteFile) Fetch(origin string) (string, error) {
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
func (r RemoteFile) ChainOnto(base Fetchable) (Fetchable, error) {
	return r, nil
}
func (r RemoteFile) IsPlainHTTP() bool { return r.url.Scheme == "http" }
func (r RemoteFile) Authority() string {
	if r.url.User != nil {
		return fmt.Sprintf("%s@%s", r.url.User.String(), r.url.Host)
	}
	return r.url.Host
}
func (r RemoteFile) PathComponents() []string {
	if r.url.Path == "" || r.url.Path == "/" {
		return []string{""}
	}
	return strings.Split(r.url.EscapedPath()[1:], "/")
}
func (r RemoteFile) Query() *string {
	if r.url.RawQuery == "" && !r.url.ForceQuery {
		return nil
	}
	return &r.url.RawQuery
}
func (r RemoteFile) AsLocation() Term {
	return Apply(Field{LocationType, "Remote"}, TextLit{Suffix: r.String()})
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
func (Missing) AsLocation() Term {
	return Field{LocationType, "Missing"}
}
