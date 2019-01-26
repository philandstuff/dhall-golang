package main

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/parser"
)

func main() {
	text := []byte("λ(foo : bar) → -- foo \n baz\n")
	root := parser.ParseExpression(text)
	fmt.Printf("%+v\n", root.Normalize())
}
