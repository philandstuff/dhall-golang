package parser_test

import (
	"github.com/philandstuff/dhall-golang/parser"

	"github.com/prataprc/goparsec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func label(value string) *parser.LabelNode {
	return &parser.LabelNode{Value: value, Comments: []parser.Comment{}}
}

func labelWithComment(value string, comment string) *parser.LabelNode {
	return &parser.LabelNode{
		Value: value,
		Comments: []parser.Comment{
			parser.Comment{
				Value: comment,
			},
		}}
}

var _ = Describe("Expression", func() {
	DescribeTable("Simple tokens",
		func(p parsec.Parser, text []byte, name string) {
			root, news := p(parsec.NewScanner(text))
			Expect(news.GetCursor()).To(Equal(len(text)), "Should parse all input")
			t := root.(*parser.AtomNode)
			Expect(t.Value).To(Equal(name))
		},
		Entry("ASCII lambda", parser.Lambda, []byte(`\`), "LAMBDA"),
		Entry("Unicode lambda", parser.Lambda, []byte(`λ`), "LAMBDA"),
		Entry("Open parens", parser.OpenParens, []byte(`(`), "OPAREN"),
		Entry("Close parens", parser.CloseParens, []byte(`)`), "CPAREN"),
		Entry("Colon", parser.Colon, []byte(`:`), "COLON"),
		Entry("ASCII arrow", parser.Arrow, []byte(`->`), "ARROW"),
		Entry("Unicode arrow", parser.Arrow, []byte(`→`), "ARROW"),
	)
	DescribeTable("lambda expressions",
		func(text []byte, expected parser.LambdaExpr) {
			root, news := parser.Expression(parsec.NewScanner(text))
			Expect(news.GetCursor()).To(Equal(len(text)), "Should parse all input")
			var t *parser.LambdaExpr
			Expect(root).To(BeAssignableToTypeOf(t))
			t = root.(*parser.LambdaExpr)
			Expect(*t).To(Equal(expected))
		},
		Entry("simple",
			[]byte(`λ(foo : bar) → baz`),
			parser.LambdaExpr{
				Label: label("foo"),
				Type:  label("bar"),
				Body:  label("baz"),
			}),
		Entry("with line comment",
			[]byte("λ(foo : bar) --asdf\n → baz"),
			parser.LambdaExpr{
				Label: label("foo"),
				Type:  label("bar"),
				Body:  label("baz"),
			}),
	)
	DescribeTable("line comments",
		func(text []byte, expected interface{}) {
			root, news := parser.LineComment(parsec.NewScanner(text))
			Expect(news.GetCursor()).To(Equal(len(text)), "Should parse all input")
			Expect(root).To(Equal(expected))
		},
		Entry("minimal line comment",
			[]byte("--\n"),
			&parser.Comment{
				Value: "",
			}),
		Entry("line comment with string",
			[]byte("-- foobar\n"),
			&parser.Comment{
				Value: " foobar",
			}),
	)
	DescribeTable("whitespace chunk",
		func(text []byte, expected interface{}) {
			root, news := parser.WhitespaceChunk(parsec.NewScanner(text))
			Expect(news.GetCursor()).To(Equal(len(text)), "Should parse all input")
			Expect(root.([]parsec.ParsecNode)[0]).To(Equal(expected))
		},
		Entry("simple space",
			[]byte(" "),
			parsec.NewTerminal("WS", " ", 0)),
		Entry("simple tab",
			[]byte("\t"),
			parsec.NewTerminal("WS", "\t", 0)),
		Entry("simple newline",
			[]byte("\n"),
			parsec.NewTerminal("WS", "\n", 0)),
		Entry("windows newline",
			[]byte("\r\n"),
			parsec.NewTerminal("WS", "\r\n", 0)),
		Entry("minimal line comment",
			[]byte("--\n"),
			&parser.Comment{
				Value: "",
			}),
		Entry("line comment with string",
			[]byte("-- foobar\n"),
			&parser.Comment{
				Value: " foobar",
			}),
	)
	DescribeTable("whitespace",
		func(text []byte, expected interface{}) {
			root, news := parser.Whitespace(parsec.NewScanner(text))
			Expect(news.GetCursor()).To(Equal(len(text)), "Should parse all input")
			Expect(root).To(Equal(expected))
		},
		Entry("simple space",
			[]byte(`    `),
			[]parser.Comment{}),
		Entry("newlines and space",
			[]byte("\n\n\n\n    "),
			[]parser.Comment{}),
		Entry("minimal line comment",
			[]byte("--\n"),
			[]parser.Comment{
				parser.Comment{
					Value: "",
				},
			}),
		Entry("line comment with string",
			[]byte("-- foobar\n"),
			[]parser.Comment{
				parser.Comment{
					Value: " foobar",
				},
			}),
	)
	DescribeTable("leading & trailing whitespace",
		func(text []byte, expected interface{}) {
			root, news := parser.Expression(parsec.NewScanner(text))
			Expect(news.GetCursor()).To(Equal(len(text)), "Should parse all input")
			Expect(root).To(Equal(expected))
		},
		Entry("simple lambda",
			[]byte(`λ(foo : bar) → baz`),
			&parser.LambdaExpr{
				Label: label("foo"),
				Type:  label("bar"),
				Body:  label("baz"),
			}),
		Entry("lambda with trailing comment",
			[]byte("λ(foo : bar) → baz -- bar\n"),
			&parser.LambdaExpr{
				Label: label("foo"),
				Type:  label("bar"),
				Body:  labelWithComment("baz", " bar"),
			}),
	)
})
