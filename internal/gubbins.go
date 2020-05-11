/*
Package internal contains internal names, not for use by library
consumers.

The exported names within this package are subject to change without
warning.
*/
package internal

import (
	"net/url"

	"github.com/philandstuff/dhall-golang/v3/term"
)

func NewImport(fetchable term.Fetchable, mode term.ImportMode) term.Import {
	return term.Import{
		ImportHashed: term.ImportHashed{
			Fetchable: fetchable,
		},
		ImportMode: mode,
	}
}
func NewEnvVarImport(envvar string, mode term.ImportMode) term.Import {
	return NewImport(term.EnvVar(envvar), mode)
}

func NewLocalImport(path string, mode term.ImportMode) term.Import {
	return NewImport(term.LocalFile(path), mode)
}

// only for generating test data - discards errors
func NewRemoteImport(uri string, mode term.ImportMode) term.Import {
	parsedURI, _ := url.ParseRequestURI(uri)
	remote := term.NewRemoteFile(parsedURI)
	return NewImport(remote, mode)
}
