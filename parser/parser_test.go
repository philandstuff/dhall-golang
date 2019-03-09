package parser_test

import (
	"math"

	. "github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func ParseAndCompare(input string, expected interface{}) {
	root, err := parser.Parse("test", []byte(input))
	Expect(err).ToNot(HaveOccurred())
	Expect(root).To(Equal(expected))
}

func ParseAndFail(input string) {
	_, err := parser.Parse("test", []byte(input))
	Expect(err).To(HaveOccurred())
}

var _ = Describe("Expression", func() {
	DescribeTable("simple expressions", ParseAndCompare,
		Entry("Type", `Type`, Type),
		Entry("Kind", `Kind`, Kind),
		Entry("Sort", `Sort`, Sort),
		Entry("Double", `Double`, Double),
		Entry("DoubleLit", `3.0`, DoubleLit(3.0)),
		Entry("DoubleLit with exponent", (`3E5`), DoubleLit(3e5)),
		Entry("DoubleLit with sign", (`+3.0`), DoubleLit(3.0)),
		Entry("DoubleLit with everything", (`-5.0e1`), DoubleLit(-50.0)),
		Entry("Infinity", `Infinity`, DoubleLit(math.Inf(1))),
		Entry("-Infinity", `-Infinity`, DoubleLit(math.Inf(-1))),
		Entry("Integer", `Integer`, Integer),
		Entry("IntegerLit", `+1234`, IntegerLit(1234)),
		Entry("IntegerLit", `-3`, IntegerLit(-3)),
		Entry("Annotated expression", `3 : Natural`, Annot{NaturalLit(3), Natural}),
	)
	DescribeTable("bools", ParseAndCompare,
		Entry("Bool", `Bool`, Bool),
		Entry("True", `True`, BoolLit(true)),
		Entry("False", `False`, BoolLit(false)),
		Entry("if True then x else y", `if True then x else y`, BoolIf{True, Var{"x", 0}, Var{"y", 0}}),
	)
	DescribeTable("naturals", ParseAndCompare,
		Entry("Natural", `Natural`, Natural),
		Entry("NaturalLit", `1234`, NaturalLit(1234)),
		Entry("NaturalLit", `3`, NaturalLit(3)),
		Entry("NaturalPlus", `3 + 5`, NaturalPlus{NaturalLit(3), NaturalLit(5)}),
		Entry("NaturalTimes", `3 * 5`, NaturalTimes{NaturalLit(3), NaturalLit(5)}),
		Entry("Natural op order #1", `3 * 4 + 5`, NaturalPlus{NaturalTimes{NaturalLit(3), NaturalLit(4)}, NaturalLit(5)}),
		Entry("Natural op order #2", `3 + 4 * 5`, NaturalPlus{NaturalLit(3), NaturalTimes{NaturalLit(4), NaturalLit(5)}}),
		// Check that if we skip whitespace, it parses
		// correctly as function application, not natural
		// addition
		Entry("Plus without whitespace", `3 +5`, &App{NaturalLit(3), IntegerLit(5)}),
	)
	DescribeTable("text", ParseAndCompare,
		Entry("Text", `Text`, Text),
		Entry("Empty TextLit", `""`, TextLit{}),
		Entry("Simple TextLit", `"foo"`, TextLit{Suffix: "foo"}),
		Entry(`TextLit escape "`, `"\""`, TextLit{Suffix: `"`}),
		Entry(`TextLit escape $`, `"\$"`, TextLit{Suffix: `$`}),
		Entry(`TextLit escape \`, `"\\"`, TextLit{Suffix: `\`}),
		Entry(`TextLit escape /`, `"\/"`, TextLit{Suffix: `/`}),
		Entry(`TextLit escape \b`, `"\b"`, TextLit{Suffix: "\b"}),
		Entry(`TextLit escape \f`, `"\f"`, TextLit{Suffix: "\f"}),
		Entry(`TextLit escape \n`, `"\n"`, TextLit{Suffix: "\n"}),
		Entry(`TextLit escape \r`, `"\r"`, TextLit{Suffix: "\r"}),
		Entry(`TextLit escape \t`, `"\t"`, TextLit{Suffix: "\t"}),
		Entry(`TextLit escape \u2200`, `"\u2200"`, TextLit{Suffix: "∀"}),
		Entry(`TextLit escape \u03bb`, `"\u03bb"`, TextLit{Suffix: "λ"}),
		Entry(`TextLit escape \u03BB`, `"\u03BB"`, TextLit{Suffix: "λ"}),
		Entry("Interpolated TextLit", `"foo ${"bar"} baz"`,
			TextLit{Chunks{Chunk{"foo ", TextLit{Suffix: "bar"}}},
				" baz"},
		),
	)
	DescribeTable("simple expressions", ParseAndCompare,
		Entry("Identifier", `x`, Var{"x", 0}),
		Entry("Identifier with index", `x@1`, Var{"x", 1}),
		Entry("Identifier with reserved prefix", `Listicle`, Var{"Listicle", 0}),
		Entry("Identifier with reserved prefix and index", `Listicle@3`, Var{"Listicle", 3}),
	)
	DescribeTable("lists", ParseAndCompare,
		Entry("List Natural", `List Natural`, &App{List, Natural}),
		Entry("[3]", `[3]`, MakeList(NaturalLit(3))),
		Entry("[3,4]", `[3,4]`, MakeList(NaturalLit(3), NaturalLit(4))),
		Entry("[] : List Natural", `[] : List Natural`, EmptyList{Natural}),
		Entry("[3] : List Natural", `[3] : List Natural`, Annot{MakeList(NaturalLit(3)), &App{List, Natural}}),
	)
	DescribeTable("records", ParseAndCompare,
		Entry("{}", `{}`, Record(map[string]Expr{})),
		Entry("{=}", `{=}`, RecordLit(map[string]Expr{})),
		Entry("{foo : Natural}", `{foo : Natural}`, Record(map[string]Expr{"foo": Natural})),
		Entry("{foo = 3}", `{foo = 3}`, RecordLit(map[string]Expr{"foo": NaturalLit(3)})),
		Entry("{foo : Natural, bar : Integer}", `{foo : Natural, bar: Integer}`, Record(map[string]Expr{"foo": Natural, "bar": Integer})),
		Entry("{foo = 3 , bar = +3}", `{foo = 3 , bar = +3}`, RecordLit(map[string]Expr{"foo": NaturalLit(3), "bar": IntegerLit(3)})),
		Entry("t.x", `t.x`, Field{Var{"t", 0}, "x"}),
		Entry("t.x.y", `t.x.y`, Field{Field{Var{"t", 0}, "x"}, "y"}),
	)
	// can't test NaN using ParseAndCompare because NaN ≠ NaN
	It("handles NaN correctly", func() {
		root, err := parser.Parse("test", []byte(`NaN`))
		Expect(err).ToNot(HaveOccurred())
		f := float64(root.(DoubleLit))
		Expect(math.IsNaN(f)).To(BeTrue())
	})
	DescribeTable("lambda expressions", ParseAndCompare,
		Entry("simple λ",
			`λ(foo : bar) → baz`,
			&LambdaExpr{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry(`simple \`,
			`\(foo : bar) → baz`,
			&LambdaExpr{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("with line comment",
			"λ(foo : bar) --asdf\n → baz",
			&LambdaExpr{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("with block comment",
			"λ(foo : bar) {-asdf\n-} → baz",
			&LambdaExpr{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("simple ∀",
			`∀(foo : bar) → baz`,
			&Pi{"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("arrow type has implicit _ var",
			`foo → bar`,
			&Pi{"_", Var{"foo", 0}, Var{"bar", 0}}),
		Entry(`simple forall`,
			`forall(foo : bar) → baz`,
			&Pi{"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("with line comment",
			"∀(foo : bar) --asdf\n → baz",
			&Pi{"foo", Var{"bar", 0}, Var{"baz", 0}}),
	)
	DescribeTable("applications", ParseAndCompare,
		Entry("identifier application",
			`foo bar`,
			&App{
				Var{"foo", 0},
				Var{"bar", 0},
			}),
		Entry("lambda application",
			`(λ(foo : bar) → baz) quux`,
			&App{
				&LambdaExpr{
					"foo", Var{"bar", 0}, Var{"baz", 0}},
				Var{"quux", 0}}),
	)
	DescribeTable("lets", ParseAndCompare,
		Entry("simple let",
			`let x = y in z`,
			MakeLet(Var{"z", 0}, Binding{
				Variable: "x", Value: Var{"y", 0},
			})),
		Entry("lambda application",
			`(λ(foo : bar) → baz) quux`,
			&App{
				&LambdaExpr{
					"foo", Var{"bar", 0}, Var{"baz", 0}},
				Var{"quux", 0}}),
	)
	Describe("Expected failures", func() {
		// these keywords should fail to parse unless they're part of
		// a larger expression
		DescribeTable("keywords", ParseAndFail,
			Entry("if", `if`),
			Entry("then", `then`),
			Entry("else", `else`),
			Entry("let", `let`),
			Entry("in", `in`),
			Entry("as", `as`),
			Entry("using", `using`),
			Entry("merge", `merge`),
			Entry("constructors", `constructors`),
			Entry("Some", `Some`),
		)
		DescribeTable("other expected failures", ParseAndFail,
			Entry("annotation without required space", `3 :Natural`),
			Entry("unannotated list", `[]`),
		)
	})
})
