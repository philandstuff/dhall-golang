package main

import "fmt"

import "github.com/philandstuff/dhall-golang/parser"
import "github.com/prataprc/goparsec"

func main() {
	text := []byte("λ(foo : bar) → -- foo \n baz\n")
	root, _ := parser.Expression(parsec.NewScanner(text))
	t := root.(*parser.LambdaExpr)
	fmt.Printf("%+v\n", t)
}
