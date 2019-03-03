package ast_test

import (
	. "github.com/philandstuff/dhall-golang/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func expectType(in, expectedType Expr) {
	annot := Annot{Expr: in, Annotation: expectedType}
	_, err := annot.TypeWith(EmptyContext())
	Expect(err).ToNot(HaveOccurred())
}

var _ = Describe("TypeCheck in empty context", func() {
	It("Kind : Sort", func() {
		// We can't typecheck this using expectType because it
		// will fail by not finding a type for Sort
		// so we have a custom test for this case

		typ, err := Kind.TypeWith(EmptyContext())
		Expect(err).ToNot(HaveOccurred())
		Expect(typ).To(Equal(Sort))
	})
	DescribeTable("Simple types",
		expectType,
		Entry("Type : Kind", Type, Kind),
		Entry("True : Bool", BoolLit(true), Bool),
		Entry("3 : Natural", NaturalLit(3), Natural),
		Entry("(3 : Natural) : Natural", Annot{NaturalLit(3), Natural}, Natural),
		Entry("(3 : (λ(x : Type) → x) Natural) : Natural", Annot{NaturalLit(3), &App{&LambdaExpr{"x", Type, x(0)}, Natural}}, Natural),
		Entry("3 + 5 : Natural", NaturalPlus{NaturalLit(3), NaturalLit(5)}, Natural),
	)
	DescribeTable("Function types",
		expectType,
		Entry("λ(x : Natural) → x : ∀(x : Natural) → Natural",
			&LambdaExpr{"x", Natural, x(0)},
			&Pi{"x", Natural, Natural}),
		Entry("λ(x : Natural) → x : Natural → Natural",
			&LambdaExpr{"x", Natural, x(0)},
			&Pi{"_", Natural, Natural}),
		Entry("λ(x : Natural) → x : Natural → Natural",
			&LambdaExpr{"x", Natural, x(0)},
			&Pi{"_", Natural, Natural}),
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
	)
	DescribeTable("List types",
		expectType,
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
		Entry("let x = 3 in x : Natural",
			MakeLet(Var{"x", 0}, Binding{Variable: "x", Value: NaturalLit(3)}),
			Natural),
		Entry("let x = True in 3 : Natural",
			MakeLet(NaturalLit(3), Binding{Variable: "x", Value: True}),
			Natural),
		// Below case from https://github.com/dhall-lang/dhall-lang/pull/69
		Entry("(let x = Natural in 1 : x) : Natural",
			MakeLet(Annot{NaturalLit(1), Var{"x", 0}}, Binding{Variable: "x", Value: Natural}),
			Natural),
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
		Entry("if True then Type else (Type -> Type)", BoolIf{True, Type, &Pi{"_", Type, Type}}),
		Entry("let x : Bool = 3 in 5", MakeLet(NaturalLit(5),
			Binding{Variable: "x", Annotation: Bool, Value: NaturalLit(3)})),
	)
})
