package pprint

import (
	"fmt"
	"io"

	"github.com/philandstuff/dhall-golang/parser"
	"github.com/prataprc/goparsec"
)

func PrettyPrint(node parsec.ParsecNode, writer io.Writer) {
	switch typedNode := node.(type) {
	case []parsec.ParsecNode:
		for _, n := range typedNode {
			PrettyPrint(n, writer)
		}
	case *parser.LambdaExpr:
		fmt.Fprint(writer, "λ(")
		PrettyPrint(typedNode.Label, writer)
		fmt.Fprint(writer, " : ")
		PrettyPrint(typedNode.Type, writer)
		fmt.Fprint(writer, ") → ")
		PrettyPrint(typedNode.Body, writer)
	case *parser.LabelNode:
		fmt.Fprint(writer, typedNode.Value)

	}
}
