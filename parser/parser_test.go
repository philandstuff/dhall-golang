package parser_test

import (
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

var _ = Describe("Expression", func() {
	DescribeTable("simple expressions", ParseAndCompare,
		Entry("Type", []byte(`Type`), ast.Type),
		Entry("Kind", []byte(`Kind`), ast.Kind),
		Entry("Sort", []byte(`Sort`), ast.Sort),
		Entry("Natural", []byte(`Natural`), ast.Natural),
		Entry("NaturalLit decimal", []byte(`1234`), ast.NaturalLit(1234)),
		Entry("NaturalLit decimal", []byte(`3`), ast.NaturalLit(3)),
	)
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
})
