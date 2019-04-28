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
type Remote string
type Missing struct{}

type Resolvable interface {
	Name() string
	Resolve() (string, error)
	ChainOnto(base Resolvable) (Resolvable, error)
}

var _ Resolvable = EnvVar("")
var _ Resolvable = Local("")
var _ Resolvable = Remote("")
var _ Resolvable = Missing(struct{}{})

func (e EnvVar) Name() string { return string(e) }
func (e EnvVar) Resolve() (string, error) {
	val, ok := os.LookupEnv(string(e))
	if !ok {
		return "", fmt.Errorf("Unset environment variable %s", string(e))
	}
	return val, nil
}
func (e EnvVar) ChainOnto(base Resolvable) (Resolvable, error) {
	if _, ok := base.(Remote); ok {
		return nil, errors.New("Can't access environment variable from remote import")
	}
	return e, nil
}

func (l Local) Name() string { return string(l) }
func (l Local) Resolve() (string, error) {
	bytes, err := ioutil.ReadFile(string(l))
	return string(bytes), err
}
func (l Local) ChainOnto(base Resolvable) (Resolvable, error) {
	switch r := base.(type) {
	case Local:
		if l.IsAbs() {
			return l, nil
		}
		return Local(path.Join(path.Dir(string(r)), string(l))), nil
	case Remote:
		if path.IsAbs(string(l)) {
			return nil, errors.New("Can't get absolute path from remote import")
		}
		baseURL := r.URL()
		relativeURL, err := url.Parse(string(l))
		if err != nil {
			return nil, err
		}
		newURL := baseURL.ResolveReference(relativeURL)
		return Remote(newURL.String()), nil
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

func (r Remote) Name() string { return string(r) }
func (r Remote) Resolve() (string, error) {
	resp, err := http.Get(string(r))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Got status %d from URL %s", resp.StatusCode, string(r))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	return string(bodyBytes), err
}
func (r Remote) ChainOnto(base Resolvable) (Resolvable, error) {
	return r, nil
}
func (r Remote) URL() *url.URL     { url, _ := url.Parse(string(r)); return url }
func (r Remote) IsPlainHttp() bool { return r.URL().Scheme == "http" }
func (r Remote) Authority() string { return r.URL().Host }
func (r Remote) PathComponents() []string {
	return strings.Split(r.URL().EscapedPath()[1:], "/")
}
func (r Remote) Query() string { return r.URL().RawQuery }

func (Missing) Name() string { return "" }
func (Missing) Resolve() (string, error) {
	return "", errors.New("Cannot resolve missing import")
}
func (Missing) ChainOnto(base Resolvable) (Resolvable, error) {
	return Missing{}, nil
}
