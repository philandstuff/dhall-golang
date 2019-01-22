package format

import (
	"fmt"
	"io"

	"github.com/philandstuff/dhall-golang/parser"
	"github.com/prataprc/goparsec"
)

func formatComment(writer io.Writer, comment parser.Comment) {
	if comment.IsBlock {
		// unimplemented
	} else {
		fmt.Fprintf(writer, "--%s\n", comment.Value)
	}
}

func Format(writer io.Writer, node parsec.ParsecNode) {
	switch typedNode := node.(type) {
	case []parsec.ParsecNode:
		for _, n := range typedNode {
			Format(writer, n)
		}
	case *parser.LambdaExpr:
		fmt.Fprint(writer, "λ(")
		Format(writer, typedNode.Label)
		fmt.Fprint(writer, " : ")
		Format(writer, typedNode.Type)
		fmt.Fprint(writer, ") → ")
		Format(writer, typedNode.Body)
	case *parser.LabelNode:
		fmt.Fprint(writer, typedNode.Value)
	case []parser.Comment:
		for _, comment := range typedNode {
			formatComment(writer, comment)
		}

	}
}
