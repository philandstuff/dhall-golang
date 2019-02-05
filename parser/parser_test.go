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
		Entry("Double", []byte(`Double`), ast.Double),
		Entry("DoubleLit", []byte(`3.0`), ast.DoubleLit(3.0)),
		Entry("DoubleLit with exponent", []byte(`3E5`), ast.DoubleLit(3e5)),
		Entry("DoubleLit with sign", []byte(`+3.0`), ast.DoubleLit(3.0)),
		Entry("DoubleLit with everything", []byte(`-5.0e1`), ast.DoubleLit(-50.0)),
		Entry("Natural", []byte(`Natural`), ast.Natural),
		Entry("NaturalLit", []byte(`1234`), ast.NaturalLit(1234)),
		Entry("NaturalLit", []byte(`3`), ast.NaturalLit(3)),
		Entry("Integer", []byte(`Integer`), ast.Integer),
		Entry("IntegerLit", []byte(`+1234`), ast.IntegerLit(1234)),
		Entry("IntegerLit", []byte(`-3`), ast.IntegerLit(-3)),
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
})
