package ast_test

import (
	. "github.com/philandstuff/dhall-golang/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeCheck in empty context", func() {
	DescribeTable("Successful typechecks",
		func(in Expr, expectedType Expr) {
			actualType, err := in.TypeWith(EmptyContext())
			Expect(err).ToNot(HaveOccurred())
			Expect(actualType.Normalize()).To(Equal(expectedType))
		},
		Entry("Type : Kind", Type, Kind),
		Entry("Kind : Sort", Kind, Sort),
		Entry("True : Bool", BoolLit(true), Bool),
		Entry("3 : Natural", NaturalLit(3), Natural),
		Entry("(3 : Natural) : Natural", Annot{NaturalLit(3), Natural}, Natural),
		Entry("(3 : (λ(x : Type) → x) Natural) : Natural", Annot{NaturalLit(3), &App{&LambdaExpr{"x", Type, x(0)}, Natural}}, Natural),
		Entry("3 + 5 : Natural", NaturalPlus{NaturalLit(3), NaturalLit(5)}, Natural),
		Entry("λ(x : Natural) → x : ∀(x : Natural) → Natural",
			&LambdaExpr{"x", Natural, x(0)},
			&Pi{"x", Natural, Natural}),
		Entry("(λ(x : Natural) → x) 3 : Natural",
			&App{&LambdaExpr{"x", Natural, x(0)},
				NaturalLit(3)},
			Natural),
		Entry("(λ(x : Type) → λ(x : x) → x) : ∀(x : Type) → ∀(x : x) → x@1",
			&LambdaExpr{
				"x",
				Type,
				&LambdaExpr{
					"x",
					x(0),
					x(0)}},
			&Pi{"x", Type,
				&Pi{"x", x(0), x(1)}}),
		Entry("(λ(x : Type) → λ(x : x) → x) Natural : ∀(x : Natural) → Natural",
			&App{
				&LambdaExpr{"x", Type,
					&LambdaExpr{"x", x(0), x(0)}},
				Natural,
			},
			&Pi{"x", Natural, Natural}),
		Entry("(λ(x : Type) → λ(x : x) → x) Natural 3 : Natural",
			&App{
				&App{
					&LambdaExpr{"x", Type,
						&LambdaExpr{"x", x(0), x(0)}},
					Natural,
				},
				NaturalLit(3),
			},
			Natural),
		Entry("([] : List Natural) : List Natural",
			EmptyList{Natural},
			&App{List, Natural}),
		Entry("[3] : List Natural",
			MakeList(NaturalLit(3)),
			&App{List, Natural}),
		Entry("([3] : List Natural) : List Natural",
			Annot{MakeList(NaturalLit(3)), &App{List, Natural}},
			&App{List, Natural}),
		Entry("if True then 3 else 4 : Natural",
			BoolIf{True, NaturalLit(3), NaturalLit(4)}, Natural),
	)
	DescribeTable("Type errors",
		func(in Expr) {
			_, err := in.TypeWith(EmptyContext())
			Expect(err).To(HaveOccurred())
		},
		Entry("3 +5", &App{NaturalLit(3), IntegerLit(5)}),
		Entry("3 : Integer", Annot{NaturalLit(3), Integer}),
		Entry("if True then 3 else +4", BoolIf{True, NaturalLit(3), IntegerLit(4)}),
		Entry("if 2 then 3 else 4", BoolIf{NaturalLit(3), NaturalLit(3), NaturalLit(4)}),
	)
})
