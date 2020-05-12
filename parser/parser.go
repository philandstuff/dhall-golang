package parser

import (
	"errors"
	"io"

	"github.com/philandstuff/dhall-golang/v3/parser/internal"
	"github.com/philandstuff/dhall-golang/v3/term"
)

//go:generate pigeon -optimize-grammar -optimize-parser -o internal/dhall.go internal/dhall.peg

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte) (term.Term, error) {
	result, err := internal.Parse(filename, b)
	if err != nil {
		return nil, err
	}
	term, ok := result.(term.Term)
	if !ok {
		// can't happen if the PEG is correct
		return nil, errors.New("dhall-golang internal error: parser returned a non-Term")
	}
	return term, nil
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string) (term.Term, error) {
	result, err := internal.ParseFile(filename)
	if err != nil {
		return nil, err
	}
	term, ok := result.(term.Term)
	if !ok {
		// can't happen if the PEG is correct
		return nil, errors.New("dhall-golang internal error: parser returned a non-Term")
	}
	return term, nil
}

// ParseReader parses the data from r using filename as information in
// the error messages.
func ParseReader(filename string, r io.Reader) (term.Term, error) {
	result, err := internal.ParseReader(filename, r)
	if err != nil {
		return nil, err
	}
	term, ok := result.(term.Term)
	if !ok {
		// can't happen if the PEG is correct
		return nil, errors.New("dhall-golang internal error: parser returned a non-Term")
	}
	return term, nil
}
