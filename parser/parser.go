package parser

import (
	"strconv"

	"github.com/prataprc/goparsec"
)

type LambdaExpr struct {
	Label parsec.ParsecNode
	Type  interface{}
	Body  interface{}
}

var ExprA parsec.Parser

func parseLineComment(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 3 {
		return nil
	}
	return &Comment{Value: ns[1].(*parsec.Terminal).Value}
}

var LineComment = parsec.And(parseLineComment,
	parsec.AtomExact(`--`, "LINECOMMENTMARK"),
	parsec.TokenExact(`[\x20-\x{10ffff}\t]*`, "LINECOMMENT"),
	parsec.TokenExact(`\n|\r\n`, "EOL"),
)

type Comment struct {
	IsBlock bool
	Value   string
}

type WhitespaceNode struct {
	Comments []Comment
}

var WhitespaceChunk = parsec.OrdChoice(
	nil,
	parsec.TokenExact(`[ \t\n]|\r\n`, "WS"),
	LineComment,
	// BlockComment,
)

func parseWhitespace(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return &WhitespaceNode{}
	}
	var comments []Comment
	for _, n := range ns {
		chunk := n.([]parsec.ParsecNode)[0]
		switch t := chunk.(type) {
		case *parsec.Terminal:
			continue
		case *Comment:
			comments = append(comments, *t)
		}
	}
	return &WhitespaceNode{Comments: comments}
}

var Whitespace = parsec.Kleene(parseWhitespace, WhitespaceChunk)

func SkipWSAfter(p parsec.Parser) parsec.Parser {
	return parsec.And(nil, p, Whitespace)
}

var Lambda = SkipWSAfter(parsec.TokenExact(`[λ\\]`, "LAMBDA"))
var OpenParens = SkipWSAfter(parsec.AtomExact(`(`, "OPAREN"))
var CloseParens = SkipWSAfter(parsec.AtomExact(`)`, "CPAREN"))
var Colon = SkipWSAfter(parsec.AtomExact(`:`, "COLON"))
var Arrow = SkipWSAfter(parsec.TokenExact(`(->|→)`, "ARROW"))

var SimpleLabel = parsec.TokenExact(`[A-Za-z_][0-9a-zA-Z_/-]*`, "SIMPLE")

var Label = SkipWSAfter(parsec.OrdChoice(parseLabel,
	SimpleLabel,
))

func parseLabel(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	switch n := ns[0].(type) {
	case *parsec.Terminal:
		switch n.Name {
		case "SIMPLE":
			return n.Value
		}
	}
	return nil
}

func parseNatural(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	switch n := ns[0].(type) {
	case *parsec.Terminal:
		switch n.Name {
		case "DEC":
			val, _ := strconv.ParseUint(n.Value, 10, 64)
			return val
		case "OCT":
			val, _ := strconv.ParseUint(n.Value[1:], 8, 64)
			return val
		case "HEX":
			val, _ := strconv.ParseUint(n.Value[2:], 16, 64)
			return val
		}
	}
	return nil
}

var Natural = parsec.OrdChoice(parseNatural, parsec.Hex(), parsec.Oct(), parsec.Token(`[1-9][0-9]+`, "DEC"))
var Identifier = parsec.OrdChoice(nil, SimpleLabel)

func parseLambda(ns []parsec.ParsecNode) parsec.ParsecNode {
	label := ns[2]
	t := ns[4]
	body := ns[7]
	return &LambdaExpr{
		Label: label.(parsec.ParsecNode),
		Type:  t,
		Body:  body,
	}
}

func parseExpr(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	return ns[0]
}

func Expression(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	var Expr parsec.Parser = Expression
	lambdaAbstraction := parsec.And(parseLambda,
		Lambda,
		OpenParens,
		Label,
		Colon,
		Expr,
		CloseParens,
		Arrow,
		Expr,
	)

	expr := parsec.OrdChoice(parseExpr,
		lambdaAbstraction,
		Label,
	)
	return expr(s)
}
