package ast_test

import (
	. "github.com/philandstuff/dhall-golang/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var Error Const = Const(-1)

var _ = DescribeTable("Rule",
	func(a Const, b Const, expected Const) {
		actual, err := Rule(a, b)
		if expected == Error {
			Expect(err).To(HaveOccurred())
		} else {
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		}
	},
	Entry("Type → Type : Type", Type, Type, Type),
	Entry("Kind → Type : Type", Kind, Type, Type),
	Entry("Sort → Type : Type", Sort, Type, Type),
	Entry("Type → Kind : !!!!", Type, Kind, Error),
	Entry("Kind → Kind : Kind", Kind, Kind, Kind),
	Entry("Sort → Kind : Sort", Sort, Kind, Sort),
	Entry("Type → Sort : !!!!", Type, Sort, Error),
	Entry("Kind → Sort : !!!!", Kind, Sort, Error),
	Entry("Sort → Sort : Type", Sort, Sort, Sort),
)

func x(i int) Var { return Var{Name: "x", Index: i} }
func y(i int) Var { return Var{Name: "y", Index: i} }

var _ = DescribeTable("Shift",
	func(actual Expr, expected Expr) {
		Expect(actual).To(Equal(expected))
	},
	Entry("Shift(_, _, Type) = Type", Shift(10, x(0), Type), Type),
	Entry("Shift(_, x, y) = y", Shift(10, x(0), y(0)), y(0)),
	Entry("Shift(5, x, x) = x@5", Shift(5, x(0), x(0)), x(5)),
	Entry("Shift(5, x@1, x) = x", Shift(5, x(1), x(0)), x(0)),
	Entry("Shift(5, x@1, x@1) = x@6", Shift(5, x(1), x(1)), x(6)),
	Entry("Shift(-1, x@1, x@1) = x", Shift(-1, x(1), x(1)), x(0)),
	Entry("Shift(5, x, λ(x : Natural) -> x)) = no change",
		Shift(5, x(0), &LambdaExpr{Label: "x", Type: Natural, Body: x(0)}),
		&LambdaExpr{Label: "x", Type: Natural, Body: x(0)}),
	Entry("Shift(5, x, λ(x : Natural) -> x@1)) = λ(x : Natural) -> x@6)",
		Shift(5, x(0), &LambdaExpr{Label: "x", Type: Natural, Body: x(1)}),
		&LambdaExpr{Label: "x", Type: Natural, Body: x(6)}),
	Entry("Shift(5, x, ∀(x : Natural) -> x)) = no change",
		Shift(5, x(0), &Pi{Label: "x", Type: Natural, Body: x(0)}),
		&Pi{Label: "x", Type: Natural, Body: x(0)}),
	Entry("Shift(5, x, ∀(x : Natural) -> x@1)) = λ(x : Natural) -> x@6)",
		Shift(5, x(0), &Pi{Label: "x", Type: Natural, Body: x(1)}),
		&Pi{Label: "x", Type: Natural, Body: x(6)}),
	Entry("Shift(1, x, Natural) = Natural", Shift(1, x(0), Natural), Natural),
	Entry("Shift(1, x, 3) = 3", Shift(1, x(0), NaturalLit(3)), NaturalLit(3)),
	Entry("Shift(1, x, x + 3) = x@1 + 3", Shift(1, x(0), NaturalPlus{L: x(0), R: NaturalLit(3)}), NaturalPlus{L: x(1), R: NaturalLit(3)}),
)

var _ = DescribeTable("Subst",
	func(actual Expr, expected Expr) {
		Expect(actual).To(Equal(expected))
	},
	Entry("Subst(_, _, Type) = Type", Subst(x(0), Sort, Type), Type),
	Entry("Subst(_, _, Natural) = Natural", Subst(x(0), Sort, Natural), Natural),
	Entry("Subst(_, _, 3) = Natural", Subst(x(0), Sort, NaturalLit(3)), NaturalLit(3)),
	Entry("Subst(x, 10, x + 3) = 10 + 3", Subst(x(0), NaturalLit(10), NaturalPlus{L: x(0), R: NaturalLit(3)}), NaturalPlus{L: NaturalLit(10), R: NaturalLit(3)}),
	Entry("Subst(x, Natural, x) = Natural", Subst(x(0), Natural, x(0)), Natural),
	Entry("Subst(x@1, Natural, x) = x", Subst(x(1), Natural, x(0)), x(0)),
	Entry("Subst(x, Natural, x@1) = x@1", Subst(x(0), Natural, x(1)), x(1)),
	Entry("Subst(x, Natural, λ(x: x) -> x) = λ(x: Natural) -> x",
		Subst(x(0), Natural, &LambdaExpr{Label: "x", Type: x(0), Body: x(0)}),
		&LambdaExpr{Label: "x", Type: Natural, Body: x(0)}),
	Entry("Subst(x, Natural, (λ(x: x) -> x) 3) = (λ(x: Natural) -> x) 3",
		Subst(x(0), Natural, &App{Fn: &LambdaExpr{Label: "x", Type: x(0), Body: x(0)}, Arg: NaturalLit(3)}),
		&App{Fn: &LambdaExpr{Label: "x", Type: Natural, Body: x(0)}, Arg: NaturalLit(3)}),
)

var _ = Describe("TypeCheck in empty context", func() {
	DescribeTable("Successful typechecks",
		func(in Expr, expectedType Expr) {
			actualType, err := in.TypeWith(EmptyContext())
			Expect(err).ToNot(HaveOccurred())
			Expect(actualType).To(Equal(expectedType))
		},
		Entry("Type : Kind", Type, Kind),
		Entry("Kind : Sort", Kind, Sort),
		Entry("3 : Natural", NaturalLit(3), Natural),
		Entry("3 + 5 : Natural", NaturalPlus{L: NaturalLit(3), R: NaturalLit(5)}, Natural),
		Entry("λ(x : Natural) → x : ∀(x : Natural) → Natural",
			&LambdaExpr{Label: "x", Type: Natural, Body: Var{Name: "x"}},
			&Pi{Label: "x", Type: Natural, Body: Natural}),
		Entry("(λ(x : Natural) → x) 3 : Natural",
			&App{Fn: &LambdaExpr{Label: "x", Type: Natural, Body: Var{Name: "x"}},
				Arg: NaturalLit(3)},
			Natural),
		Entry("(λ(x : Type) → λ(x : x) → x) : ∀(x : Type) → ∀(x : x) → x@1",
			&LambdaExpr{
				Label: "x",
				Type:  Type,
				Body: &LambdaExpr{
					Label: "x",
					Type:  x(0),
					Body:  x(0)}},
			&Pi{Label: "x", Type: Type,
				Body: &Pi{Label: "x", Type: x(0),
					Body: x(1)}}),
		Entry("(λ(x : Type) → λ(x : x) → x) Natural : ∀(x : Natural) → Natural",
			&App{
				Fn: &LambdaExpr{
					Label: "x",
					Type:  Type,
					Body: &LambdaExpr{
						Label: "x",
						Type:  x(0),
						Body:  x(0)}},
				Arg: Natural,
			},
			&Pi{Label: "x", Type: Natural, Body: Natural}),
		Entry("(λ(x : Type) → λ(x : x) → x) Natural 3 : Natural",
			&App{
				Fn: &App{
					Fn: &LambdaExpr{
						Label: "x",
						Type:  Type,
						Body: &LambdaExpr{
							Label: "x",
							Type:  Var{Name: "x"},
							Body:  Var{Name: "x"}}},
					Arg: Natural,
				},
				Arg: NaturalLit(3),
			},
			Natural),
	)
	DescribeTable("Type errors",
		func(in Expr) {
			_, err := in.TypeWith(EmptyContext())
			Expect(err).To(HaveOccurred())
		},
		Entry("3 +5", &App{Fn: NaturalLit(3), Arg: IntegerLit(5)}),
	)
})

// TODO: Normalize() when it does anything interesting
