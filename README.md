# dhall-golang

[![GoDoc](https://godoc.org/github.com/philandstuff/dhall-golang?status.svg)][dhall-golang godoc]

Go bindings for the [dhall configuration language][dhall].

[dhall]: https://dhall-lang.org/

## Quick start

Here's a minimal example of how you might use dhall-golang to load a
Dhall file into your own struct:

```golang
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/philandstuff/dhall-golang/v5"
)

// Config can be a fairly arbitrary Go datatype.  You would put your
// application configuration in this struct.
type Config struct {
	Port int
	Name string
}

func main() {
	var config Config
	bytes, err := ioutil.ReadFile("/path/to/config.dhall")
	if err != nil {
		panic(err)
	}
	err = dhall.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded config: %#v\n", config)
}
```

## Documentation

You can find more documentation in the [dhall-golang godoc][].

[dhall-golang godoc]: https://godoc.org/github.com/philandstuff/dhall-golang

## Development

This is a fairly standard Go project.  It uses go modules, so no
vendoring of dependencies is required.

### Running the tests

    git submodule update --init --recursive

    go test ./...

    go test -short ./... # skips long-running tests

### Making changes to the PEG grammar

Dhall-golang uses [pigeon][] to generate the parser source file
`parser/internal/dhall.go` from the PEG grammar at
`parser/internal/dhall.peg`.  If you change the PEG grammar, you need
to first install the pigeon binary if you don't already have it:

    # either outside a module directory, or with GO111MODULE=off
    go get github.com/mna/pigeon

Then, to regenerate the parser:

    go generate ./parser

[pigeon]: https://godoc.org/github.com/mna/pigeon

## Support

Issues and pull requests are welcome on this repository.  If you have
a question, you can ask it on the [Dhall discourse][].

[Dhall discourse]: https://discourse.dhall-lang.org/
