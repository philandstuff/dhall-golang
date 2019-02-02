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

var _ = DescribeTable("TypeCheck in empty context",
	func(in ast.Expr, expected ast.Expr) {
		actual, err := in.TypeWith(ast.EmptyContext())
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(expected))
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
)

// TODO: Normalize() when it does anything interesting
