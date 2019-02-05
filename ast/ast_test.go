package ast_test

import (
	"github.com/philandstuff/dhall-golang/ast"

	// . "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var Error ast.Const = ast.Const(-1)

var _ = DescribeTable("Rule",
	func(a ast.Const, b ast.Const, expected ast.Const) {
		actual, err := ast.Rule(a, b)
		if expected == Error {
			Expect(err).To(HaveOccurred())
		} else {
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		}
	},
	Entry("Type → Type : Type", ast.Type, ast.Type, ast.Type),
	Entry("Kind → Type : Type", ast.Kind, ast.Type, ast.Type),
	Entry("Sort → Type : Type", ast.Sort, ast.Type, ast.Type),
	Entry("Type → Kind : !!!!", ast.Type, ast.Kind, Error),
	Entry("Kind → Kind : Kind", ast.Kind, ast.Kind, ast.Kind),
	Entry("Sort → Kind : Sort", ast.Sort, ast.Kind, ast.Sort),
	Entry("Type → Sort : !!!!", ast.Type, ast.Sort, Error),
	Entry("Kind → Sort : !!!!", ast.Kind, ast.Sort, Error),
	Entry("Sort → Sort : Type", ast.Sort, ast.Sort, ast.Sort),
)

func x(i int) ast.Var { return ast.Var{Name: "x", Index: i} }
func y(i int) ast.Var { return ast.Var{Name: "y", Index: i} }

var _ = DescribeTable("Shift",
	func(actual ast.Expr, expected ast.Expr) {
		Expect(actual).To(Equal(expected))
	},
	Entry("Shift(_, _, Type) = Type", ast.Shift(10, x(0), ast.Type), ast.Type),
	Entry("Shift(_, x, y) = y", ast.Shift(10, x(0), y(0)), y(0)),
	Entry("Shift(5, x, x) = x@5", ast.Shift(5, x(0), x(0)), x(5)),
	Entry("Shift(5, x@1, x) = x", ast.Shift(5, x(1), x(0)), x(0)),
	Entry("Shift(5, x@1, x@1) = x@6", ast.Shift(5, x(1), x(1)), x(6)),
	Entry("Shift(-1, x@1, x@1) = x", ast.Shift(-1, x(1), x(1)), x(0)),
	Entry("Shift(5, x, λ(x : Natural) -> x)) = no change",
		ast.Shift(5, x(0), &ast.LambdaExpr{Label: "x", Type: ast.Natural, Body: x(0)}),
		&ast.LambdaExpr{Label: "x", Type: ast.Natural, Body: x(0)}),
	Entry("Shift(5, x, λ(x : Natural) -> x@1)) = λ(x : Natural) -> x@6)",
		ast.Shift(5, x(0), &ast.LambdaExpr{Label: "x", Type: ast.Natural, Body: x(1)}),
		&ast.LambdaExpr{Label: "x", Type: ast.Natural, Body: x(6)}),
	Entry("Shift(5, x, ∀(x : Natural) -> x)) = no change",
		ast.Shift(5, x(0), &ast.Pi{Label: "x", Type: ast.Natural, Body: x(0)}),
		&ast.Pi{Label: "x", Type: ast.Natural, Body: x(0)}),
	Entry("Shift(5, x, ∀(x : Natural) -> x@1)) = λ(x : Natural) -> x@6)",
		ast.Shift(5, x(0), &ast.Pi{Label: "x", Type: ast.Natural, Body: x(1)}),
		&ast.Pi{Label: "x", Type: ast.Natural, Body: x(6)}),
)

var _ = DescribeTable("Subst",
	func(actual ast.Expr, expected ast.Expr) {
		Expect(actual).To(Equal(expected))
	},
	Entry("Subst(_, _, Type) = Type", ast.Subst(x(0), ast.Natural, ast.Type), ast.Type),
	Entry("Subst(x, Natural, x) = Natural", ast.Subst(x(0), ast.Natural, x(0)), ast.Natural),
	Entry("Subst(x@1, Natural, x) = x", ast.Subst(x(1), ast.Natural, x(0)), x(0)),
	Entry("Subst(x, Natural, x@1) = x@1", ast.Subst(x(0), ast.Natural, x(1)), x(1)),
	Entry("Subst(x, Natural, λ(x: x) -> x) = λ(x: Natural) -> x",
		ast.Subst(x(0), ast.Natural, &ast.LambdaExpr{Label: "x", Type: x(0), Body: x(0)}),
		&ast.LambdaExpr{Label: "x", Type: ast.Natural, Body: x(0)}),
	Entry("Subst(x, Natural, (λ(x: x) -> x) 3) = (λ(x: Natural) -> x) 3",
		ast.Subst(x(0), ast.Natural, &ast.App{Fn: &ast.LambdaExpr{Label: "x", Type: x(0), Body: x(0)}, Arg: ast.NaturalLit(3)}),
		&ast.App{Fn: &ast.LambdaExpr{Label: "x", Type: ast.Natural, Body: x(0)}, Arg: ast.NaturalLit(3)}),
)

var _ = DescribeTable("TypeCheck in empty context",
	func(in ast.Expr, expectedType ast.Expr) {
		actualType, err := in.TypeWith(ast.EmptyContext())
		Expect(err).ToNot(HaveOccurred())
		Expect(actualType).To(Equal(expectedType))
	},
	Entry("Type : Kind", ast.Type, ast.Kind),
	Entry("Kind : Sort", ast.Kind, ast.Sort),
	Entry("3 : Natural", ast.NaturalLit(3), ast.Natural),
	Entry("λ(x : Natural) → x : ∀(x : Natural) → Natural",
		&ast.LambdaExpr{Label: "x", Type: ast.Natural, Body: ast.Var{Name: "x"}},
		&ast.Pi{Label: "x", Type: ast.Natural, Body: ast.Natural}),
	Entry("(λ(x : Natural) → x) 3 : Natural",
		&ast.App{Fn: &ast.LambdaExpr{Label: "x", Type: ast.Natural, Body: ast.Var{Name: "x"}},
			Arg: ast.NaturalLit(3)},
		ast.Natural),
	Entry("(λ(x : Type) → λ(x : x) → x) : ∀(x : Type) → ∀(x : x) → x@1",
		&ast.LambdaExpr{
			Label: "x",
			Type:  ast.Type,
			Body: &ast.LambdaExpr{
				Label: "x",
				Type:  x(0),
				Body:  x(0)}},
		&ast.Pi{Label: "x", Type: ast.Type,
			Body: &ast.Pi{Label: "x", Type: x(0),
				Body: x(1)}}),
	Entry("(λ(x : Type) → λ(x : x) → x) Natural : ∀(x : Natural) → Natural",
		&ast.App{
			Fn: &ast.LambdaExpr{
				Label: "x",
				Type:  ast.Type,
				Body: &ast.LambdaExpr{
					Label: "x",
					Type:  x(0),
					Body:  x(0)}},
			Arg: ast.Natural,
		},
		&ast.Pi{Label: "x", Type: ast.Natural, Body: ast.Natural}),
	Entry("(λ(x : Type) → λ(x : x) → x) Natural 3 : Natural",
		&ast.App{
			Fn: &ast.App{
				Fn: &ast.LambdaExpr{
					Label: "x",
					Type:  ast.Type,
					Body: &ast.LambdaExpr{
						Label: "x",
						Type:  ast.Var{Name: "x"},
						Body:  ast.Var{Name: "x"}}},
				Arg: ast.Natural,
			},
			Arg: ast.NaturalLit(3),
		},
		ast.Natural),
)

// TODO: Normalize() when it does anything interesting
