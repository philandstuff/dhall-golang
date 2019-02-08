package parser_test

import (
	"math"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func Var(value string) ast.Var {
	return ast.Var{Name: value}
}

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
		Entry("Type", []byte(`Type`), ast.Type),
		Entry("Kind", []byte(`Kind`), ast.Kind),
		Entry("Sort", []byte(`Sort`), ast.Sort),
		Entry("Double", []byte(`Double`), ast.Double),
		Entry("DoubleLit", []byte(`3.0`), ast.DoubleLit(3.0)),
		Entry("DoubleLit with exponent", []byte(`3E5`), ast.DoubleLit(3e5)),
		Entry("DoubleLit with sign", []byte(`+3.0`), ast.DoubleLit(3.0)),
		Entry("DoubleLit with everything", []byte(`-5.0e1`), ast.DoubleLit(-50.0)),
		Entry("Infinity", []byte(`Infinity`), ast.DoubleLit(math.Inf(1))),
		Entry("-Infinity", []byte(`-Infinity`), ast.DoubleLit(math.Inf(-1))),
		Entry("Natural", []byte(`Natural`), ast.Natural),
		Entry("NaturalLit", []byte(`1234`), ast.NaturalLit(1234)),
		Entry("NaturalLit", []byte(`3`), ast.NaturalLit(3)),
		Entry("Integer", []byte(`Integer`), ast.Integer),
		Entry("IntegerLit", []byte(`+1234`), ast.IntegerLit(1234)),
		Entry("IntegerLit", []byte(`-3`), ast.IntegerLit(-3)),
		Entry("Identifier", []byte(`x`), ast.Var{Name: "x"}),
		Entry("Identifier with index", []byte(`x@1`), ast.Var{Name: "x", Index: 1}),
	)
	// can't test NaN using ParseAndCompare because NaN ≠ NaN
	It("handles NaN correctly", func() {
		root, err := parser.Parse("test", []byte(`NaN`))
		Expect(err).ToNot(HaveOccurred())
		f := float64(root.(ast.DoubleLit))
		Expect(math.IsNaN(f)).To(BeTrue())
	})
	DescribeTable("lambda expressions", ParseAndCompare,
		Entry("simple λ",
			[]byte(`λ(foo : bar) → baz`),
			&ast.LambdaExpr{
				Label: "foo", Type: Var("bar"), Body: Var("baz")}),
		Entry(`simple \`,
			[]byte(`\(foo : bar) → baz`),
			&ast.LambdaExpr{
				Label: "foo", Type: Var("bar"), Body: Var("baz")}),
		Entry("with line comment",
			[]byte("λ(foo : bar) --asdf\n → baz"),
			&ast.LambdaExpr{
				Label: "foo", Type: Var("bar"), Body: Var("baz")}),
		Entry("with block comment",
			[]byte("λ(foo : bar) {-asdf\n-} → baz"),
			&ast.LambdaExpr{
				Label: "foo", Type: Var("bar"), Body: Var("baz")}),
		Entry("simple ∀",
			[]byte(`∀(foo : bar) → baz`),
			&ast.Pi{
				Label: "foo", Type: Var("bar"), Body: Var("baz")}),
		Entry(`simple forall`,
			[]byte(`forall(foo : bar) → baz`),
			&ast.Pi{
				Label: "foo", Type: Var("bar"), Body: Var("baz")}),
		Entry("with line comment",
			[]byte("∀(foo : bar) --asdf\n → baz"),
			&ast.Pi{
				Label: "foo", Type: Var("bar"), Body: Var("baz")}),
	)
	DescribeTable("applications", ParseAndCompare,
		Entry("identifier application",
			[]byte(`foo bar`),
			&ast.App{
				Fn:  Var("foo"),
				Arg: Var("bar"),
			}),
		Entry("lambda application",
			[]byte(`(λ(foo : bar) → baz) quux`),
			&ast.App{
				Fn: &ast.LambdaExpr{
					Label: "foo", Type: Var("bar"), Body: Var("baz")},
				Arg: Var("quux")}),
	)
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
