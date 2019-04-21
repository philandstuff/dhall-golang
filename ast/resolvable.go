package ast

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type EnvVar string
type Remote string
type Missing struct{}

type Resolvable interface {
	Name() string
	Resolve() (string, error)
}

var _ Resolvable = EnvVar("")
var _ Resolvable = Missing(struct{}{})
var _ Resolvable = Remote("")

func (e EnvVar) Name() string { return string(e) }
func (e EnvVar) Resolve() (string, error) {
	val, ok := os.LookupEnv(string(e))
	if !ok {
		return "", fmt.Errorf("Unset environment variable %s", string(e))
	}
	return val, nil
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

func (Missing) Name() string { return "" }
func (Missing) Resolve() (string, error) {
	return "", errors.New("Cannot resolve missing import")
}
