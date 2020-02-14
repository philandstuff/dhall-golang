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

// An EnvVar is a Fetchable which represents fetching the value of an
// environment variable.
type EnvVar string

// A LocalFile is a Fetchable which represents fetching the content of
// a local file.  It is defined to be one of four classes:
// here-relative (ie starts with ./), parent-relative (starts with
// ../), home-relative (starts with ~/), or absolute (starts with /).
type LocalFile string

// A RemoteFile is a Fetchable which represents fetching the content
// of a remote file (over HTTP or HTTPS).
type RemoteFile struct{ url *url.URL }

// Missing is a Fetchable which cannot be Fetched.
type Missing struct{}

// NullOrigin is used in Fetchable.Fetch() to indicate no origin.
const NullOrigin = "null"

var locationType = UnionType{
	"Local":       Text,
	"Remote":      Text,
	"Environment": Text,
	"Missing":     nil,
}

// A Fetchable is the target of a Dhall import: a remote file, local
// file, environment variable, or the special value `missing`.
//
// Fetch(origin) is the key method on this interface; it fetches the
// underlying resource, with authority from the given origin.
type Fetchable interface {
	Origin() string
	Fetch(origin string) (string, error)
	ChainOnto(base Fetchable) (Fetchable, error)
	String() string
	AsLocation() Term
}

var _ Fetchable = EnvVar("")
var _ Fetchable = LocalFile("")
var _ Fetchable = RemoteFile{}
var _ Fetchable = Missing{}

// Origin returns NullOrigin, since EnvVars do not have an origin.
func (EnvVar) Origin() string { return NullOrigin }
func (e EnvVar) String() string {
	return "env:" + string(e)
}

// Fetch reads the environment variable.  If origin is not NullOrigin,
// an error is returned, to prevent remote imports from importing
// environment variables.
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

// ChainOnto returns e.
func (e EnvVar) ChainOnto(base Fetchable) (Fetchable, error) {
	return e, nil
}

// AsLocation returns the EnvVar as a Dhall Term.  This implements the
// `env:FOO as Location` Dhall feature.
func (e EnvVar) AsLocation() Term {
	return Apply(Field{locationType, "Environment"}, TextLit{Suffix: e.String()})
}

// Origin returns NullOrigin, since LocalFiles do not have an origin.
func (LocalFile) Origin() string { return NullOrigin }
func (l LocalFile) String() string {
	if l.IsAbs() || l.IsRelativeToHome() || l.IsRelativeToParent() {
		return string(l)
	}
	return "./" + string(l)
}

// Fetch reads the local file.  If origin is not NullOrigin, an error
// is returned, to prevent remote imports from importing local files.
func (l LocalFile) Fetch(origin string) (string, error) {
	if origin != NullOrigin {
		return "", fmt.Errorf("Can't get %s from remote import at %s", l, origin)
	}
	bytes, err := ioutil.ReadFile(string(l))
	return string(bytes), err
}

// ChainOnto chains l onto the base Fetchable, according to the Dhall
// definition of import chaining:
// https://github.com/dhall-lang/dhall-lang/blob/master/standard/imports.md#chaining-imports
//
// For here- or parent-relative LocalFiles, they chain onto
// RemoteFiles using the URL reference resolution algorithm; they
// chain onto LocalFiles using filesystem path joining; they chain
// onto Missing or EnvVar by just returning the LocalFile unmodified.
//
// For home-relative or absolute LocalFiles, chaining them onto a
// RemoteFile is an error; all other cases return the LocalFile
// unmodified.
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

// IsAbs returns true if the LocalFile is an absolute path.
func (l LocalFile) IsAbs() bool { return path.IsAbs(string(l)) }

// IsRelativeToParent returns true if the LocalFile starts with "../"
func (l LocalFile) IsRelativeToParent() bool { return strings.HasPrefix(string(l), "..") }

// IsRelativeToHome returns true if the LocalFile starts with "~/"
func (l LocalFile) IsRelativeToHome() bool { return string(l)[0] == '~' }

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

// PathComponents returns a slice of strings, one for each component
// of the given path.  It excludes any leading ".", ".." or "~".
func (l LocalFile) PathComponents() []string {
	if l.IsAbs() || l.IsRelativeToHome() || l.IsRelativeToParent() {
		return strings.Split(string(l), "/")[1:]
	}
	return strings.Split(string(l), "/")
}

// AsLocation returns the LocalFile as a Dhall Term.  This implements the
// `./file as Location` Dhall feature.
func (l LocalFile) AsLocation() Term {
	return Apply(Field{locationType, "Local"}, TextLit{Suffix: l.String()})
}

// NewRemoteFile constructs a RemoteFile from a *url.URL.
func NewRemoteFile(u *url.URL) RemoteFile {
	return RemoteFile{url: u}
}

var client http.Client

// Origin returns the scheme and authority of the underlying URL of a
// RemoteFile.  For example, the Origin of
// "https://example.com/foo/bar" is "https://example.com".
func (r RemoteFile) Origin() string { return fmt.Sprintf("%s://%s", r.url.Scheme, r.Authority()) }
func (r RemoteFile) String() string { return fmt.Sprintf("%v", r.url) }

// Fetch makes an HTTP request to fetch the RemoteFile.  If origin is
// neither NullOrigin nor the same origin as this RemoteFile, this is
// considered a cross-origin request and so appropriate CORS checks
// are made; if these fail, an error is returned with no content.
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

// ChainOnto returns the RemoteFile unmodified.
func (r RemoteFile) ChainOnto(base Fetchable) (Fetchable, error) {
	return r, nil
}

// IsPlainHTTP returns true if this is an "http://" URL, and false if
// it an "https://" URL.
func (r RemoteFile) IsPlainHTTP() bool { return r.url.Scheme == "http" }

// Authority returns the authority of the URL; that is, the bit
// between the first "//" and the next "/", which includes optional
// userinfo, remote host, and optional port number.
func (r RemoteFile) Authority() string {
	if r.url.User != nil {
		return fmt.Sprintf("%s@%s", r.url.User.String(), r.url.Host)
	}
	return r.url.Host
}

// PathComponents returns a slice of strings, one for each path
// component of the given URL.
func (r RemoteFile) PathComponents() []string {
	if r.url.Path == "" || r.url.Path == "/" {
		return []string{""}
	}
	return strings.Split(r.url.EscapedPath()[1:], "/")
}

// Query returns the query string, or nil if no query string is
// present.
func (r RemoteFile) Query() *string {
	if r.url.RawQuery == "" && !r.url.ForceQuery {
		return nil
	}
	return &r.url.RawQuery
}

// AsLocation returns the RemoteFile as a Dhall Term.  This implements the
// `https://example.com/foo/bar as Location` feature.
func (r RemoteFile) AsLocation() Term {
	return Apply(Field{locationType, "Remote"}, TextLit{Suffix: r.String()})
}

// Origin returns NullOrigin, since Missing does not have an origin.
func (Missing) Origin() string { return NullOrigin }
func (Missing) String() string { return "missing" }

// Fetch always returns an error, because Missing cannot be fetched.
func (Missing) Fetch(origin string) (string, error) {
	return "", errors.New("Cannot resolve missing import")
}

// ChainOnto returns a Missing.
func (Missing) ChainOnto(base Fetchable) (Fetchable, error) {
	return Missing{}, nil
}

// AsLocation returns Missing as a Dhall Term.  This implements the
// `missing as Location` feature.
func (Missing) AsLocation() Term {
	return Field{locationType, "Missing"}
}
