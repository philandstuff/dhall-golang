package parser

import (
	"strconv"

	"github.com/prataprc/goparsec"
)

type LambdaExpr struct {
	Label string
	Type  interface{}
	Body  interface{}
}

var ExprA parsec.Parser

func Pure(name string, value string) parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		cursor := s.GetCursor()
		return parsec.NewTerminal(name, value, cursor), s
	}
}

var WhitespaceChunk = parsec.OrdChoice(
	nil,
	parsec.TokenExact(`[ \t\n]|\r\n`, "WS"),
	// LineComment
	// BlockComment
)

var Whitespace = parsec.Kleene(nil, WhitespaceChunk)

func SkipWSAfter(p parsec.Parser) parsec.Parser {
	parseFirst := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		return ns[0]
	}
	return parsec.And(parseFirst, p, Whitespace)
}

var Lambda = SkipWSAfter(parsec.TokenExact(`[λ\\]`, "LAMBDA"))
var OpenParens = SkipWSAfter(parsec.AtomExact(`(`, "OPAREN"))
var CloseParens = SkipWSAfter(parsec.AtomExact(`)`, "CPAREN"))
var Colon = SkipWSAfter(parsec.AtomExact(`:`, "COLON"))
var Arrow = SkipWSAfter(parsec.OrdChoice(nil,
	parsec.AtomExact(`->`, "ARROW"),
	parsec.AtomExact(`→`, "ARROW")))

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

func Skip(s parsec.Scanner, p parsec.Parser) parsec.Scanner {
	news := s.Clone()
	_, news = p(news)
	return news
}

func parseLambda(ns []parsec.ParsecNode) parsec.ParsecNode {
	label := ns[2]
	t := ns[4]
	body := ns[7]
	return &LambdaExpr{
		Label: label.(string),
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
