package pprint

import (
	"fmt"
	"io"

	"github.com/philandstuff/dhall-golang/parser"
	"github.com/prataprc/goparsec"
)

func prettyPrintComment(writer io.Writer, comment parser.Comment) {
	if comment.IsBlock {
		// unimplemented
	} else {
		fmt.Fprintf(writer, "--%s\n", comment.Value)
	}
}

func PrettyPrint(writer io.Writer, node parsec.ParsecNode) {
	switch typedNode := node.(type) {
	case []parsec.ParsecNode:
		for _, n := range typedNode {
			PrettyPrint(writer, n)
		}
	case *parser.LambdaExpr:
		fmt.Fprint(writer, "λ(")
		PrettyPrint(writer, typedNode.Label)
		fmt.Fprint(writer, " : ")
		PrettyPrint(writer, typedNode.Type)
		fmt.Fprint(writer, ") → ")
		PrettyPrint(writer, typedNode.Body)
	case *parser.LabelNode:
		fmt.Fprint(writer, typedNode.Value)
	case []parser.Comment:
		for _, comment := range typedNode {
			prettyPrintComment(writer, comment)
		}

	}
}
