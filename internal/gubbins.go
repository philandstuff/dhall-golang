/*
Package internal contains internal names, not for use by library
consumers.

The exported names within this package are subject to change without
warning.
*/
package internal

import (
	"net/url"

	"github.com/philandstuff/dhall-golang/core"
)

func NewImport(fetchable core.Fetchable, mode core.ImportMode) core.Import {
	return core.Import{
		ImportHashed: core.ImportHashed{
			Fetchable: fetchable,
		},
		ImportMode: mode,
	}
}
func NewEnvVarImport(envvar string, mode core.ImportMode) core.Import {
	return NewImport(core.EnvVar(envvar), mode)
}

func NewLocalImport(path string, mode core.ImportMode) core.Import {
	return NewImport(core.Local(path), mode)
}

// only for generating test data - discards errors
func NewRemoteImport(uri string, mode core.ImportMode) core.Import {
	parsedURI, _ := url.ParseRequestURI(uri)
	remote := core.NewRemote(parsedURI)
	return NewImport(remote, mode)
}
