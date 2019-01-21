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
	case []parser.Comment:
		for _, comment := range typedNode {
			prettyPrintComment(writer, comment)
		}

	}
}
