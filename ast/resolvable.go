package ast

import (
	"errors"
	"fmt"
	"os"
)

type EnvVar string
type Missing struct{}

type Resolvable interface {
	Name() string
	Resolve() (string, error)
}

var _ Resolvable = EnvVar("")
var _ Resolvable = Missing(struct{}{})

func (e EnvVar) Name() string { return string(e) }
func (e EnvVar) Resolve() (string, error) {
	val, ok := os.LookupEnv(string(e))
	if !ok {
		return "", fmt.Errorf("Unset environment variable %s", string(e))
	}
	return val, nil
}

func (Missing) Name() string { return "" }
func (Missing) Resolve() (string, error) {
	return "", errors.New("Cannot resolve missing import")
}
