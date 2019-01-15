package main

import "fmt"

import "github.com/philandstuff/dhall-golang/parser"
import "github.com/prataprc/goparsec"

func main() {
	text := []byte(`\(foo : bar) -> baz`)
	root, _ := parser.Expression(parsec.NewScanner(text))
	// nodes := root.([]parsec.ParsecNode)
	// t := nodes[0].(*parsec.Terminal)
	t := root.(*parser.LambdaExpr)
	fmt.Printf("%+v\n", t)
}
