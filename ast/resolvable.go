package ast

import (
	"fmt"
	"os"
)

type EnvVar string

type Resolvable interface {
	Name() string
	Resolve() (string, error)
}

var _ Resolvable = EnvVar("")

func (e EnvVar) Name() string { return string(e) }
func (e EnvVar) Resolve() (string, error) {
	val, ok := os.LookupEnv(string(e))
	if !ok {
		return "", fmt.Errorf("Unset environment variable %s", string(e))
	}
	return val, nil
}
