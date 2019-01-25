package main

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/prataprc/goparsec"
)

func main() {
	text := []byte("λ(foo : bar) → -- foo \n baz\n")
	root, _ := parser.Expression(parsec.NewScanner(text))
	t := root.(*ast.LambdaExpr)
	fmt.Printf("%+v\n", t)
}
