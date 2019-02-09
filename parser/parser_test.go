package parser_test

import (
	"math"

	. "github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func ParseAndCompare(input []byte, expected interface{}) {
	root, err := parser.Parse("test", input)
	Expect(err).ToNot(HaveOccurred())
	Expect(root).To(Equal(expected))
}

func ParseAndFail(input []byte) {
	_, err := parser.Parse("test", input)
	Expect(err).To(HaveOccurred())
}

var _ = Describe("Expression", func() {
	DescribeTable("simple expressions", ParseAndCompare,
		Entry("Type", []byte(`Type`), Type),
		Entry("Kind", []byte(`Kind`), Kind),
		Entry("Sort", []byte(`Sort`), Sort),
		Entry("Double", []byte(`Double`), Double),
		Entry("DoubleLit", []byte(`3.0`), DoubleLit(3.0)),
		Entry("DoubleLit with exponent", []byte(`3E5`), DoubleLit(3e5)),
		Entry("DoubleLit with sign", []byte(`+3.0`), DoubleLit(3.0)),
		Entry("DoubleLit with everything", []byte(`-5.0e1`), DoubleLit(-50.0)),
		Entry("Infinity", []byte(`Infinity`), DoubleLit(math.Inf(1))),
		Entry("-Infinity", []byte(`-Infinity`), DoubleLit(math.Inf(-1))),
		Entry("Integer", []byte(`Integer`), Integer),
		Entry("IntegerLit", []byte(`+1234`), IntegerLit(1234)),
		Entry("IntegerLit", []byte(`-3`), IntegerLit(-3)),
		Entry("Identifier", []byte(`x`), Var{"x", 0}),
		Entry("Identifier with index", []byte(`x@1`), Var{"x", 1}),
	)
	DescribeTable("naturals", ParseAndCompare,
		Entry("Natural", []byte(`Natural`), Natural),
		Entry("NaturalLit", []byte(`1234`), NaturalLit(1234)),
		Entry("NaturalLit", []byte(`3`), NaturalLit(3)),
		Entry("NaturalPlus", []byte(`3 + 5`), NaturalPlus{NaturalLit(3), NaturalLit(5)}),
		// Check that if we skip whitespace, it parses
		// correctly as function application, not natural
		// addition
		Entry("Plus without whitespace", []byte(`3 +5`), &App{NaturalLit(3), IntegerLit(5)}),
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
			[]byte(`λ(foo : bar) → baz`),
			&LambdaExpr{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry(`simple \`,
			[]byte(`\(foo : bar) → baz`),
			&LambdaExpr{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("with line comment",
			[]byte("λ(foo : bar) --asdf\n → baz"),
			&LambdaExpr{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("with block comment",
			[]byte("λ(foo : bar) {-asdf\n-} → baz"),
			&LambdaExpr{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("simple ∀",
			[]byte(`∀(foo : bar) → baz`),
			&Pi{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry(`simple forall`,
			[]byte(`forall(foo : bar) → baz`),
			&Pi{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
		Entry("with line comment",
			[]byte("∀(foo : bar) --asdf\n → baz"),
			&Pi{
				"foo", Var{"bar", 0}, Var{"baz", 0}}),
	)
	DescribeTable("applications", ParseAndCompare,
		Entry("identifier application",
			[]byte(`foo bar`),
			&App{
				Var{"foo", 0},
				Var{"bar", 0},
			}),
		Entry("lambda application",
			[]byte(`(λ(foo : bar) → baz) quux`),
			&App{
				&LambdaExpr{
					"foo", Var{"bar", 0}, Var{"baz", 0}},
				Var{"quux", 0}}),
	)
	Describe("Expected failures", func() {
		// these keywords should fail to parse unless they're part of
		// a larger expression
		DescribeTable("keywords", ParseAndFail,
			Entry("if", []byte(`if`)),
			Entry("then", []byte(`then`)),
			Entry("else", []byte(`else`)),
			Entry("let", []byte(`let`)),
			Entry("in", []byte(`in`)),
			Entry("as", []byte(`as`)),
			Entry("using", []byte(`using`)),
			Entry("merge", []byte(`merge`)),
			Entry("constructors", []byte(`constructors`)),
			Entry("Some", []byte(`Some`)),
		)
	})
})
