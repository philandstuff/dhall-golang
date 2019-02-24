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

func x(i int) Var { return Var{"x", i} }
func y(i int) Var { return Var{"y", i} }

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
		Shift(5, x(0), &LambdaExpr{"x", Natural, x(0)}),
		&LambdaExpr{"x", Natural, x(0)}),
	Entry("Shift(5, x, λ(x : Natural) -> x@1)) = λ(x : Natural) -> x@6)",
		Shift(5, x(0), &LambdaExpr{"x", Natural, x(1)}),
		&LambdaExpr{"x", Natural, x(6)}),
	Entry("Shift(5, x, ∀(x : Natural) -> x)) = no change",
		Shift(5, x(0), &Pi{"x", Natural, x(0)}),
		&Pi{"x", Natural, x(0)}),
	Entry("Shift(5, x, ∀(x : Natural) -> x@1)) = λ(x : Natural) -> x@6)",
		Shift(5, x(0), &Pi{"x", Natural, x(1)}),
		&Pi{"x", Natural, x(6)}),
	Entry("Shift(1, x, Natural) = Natural", Shift(1, x(0), Natural), Natural),
	Entry("Shift(1, x, 3) = 3", Shift(1, x(0), NaturalLit(3)), NaturalLit(3)),
	Entry("Shift(1, x, x : Natural) = x@1 : Natural", Shift(1, x(0), Annot{x(0), Natural}), Annot{x(1), Natural}),
	Entry("Shift(1, x, x + 3) = x@1 + 3", Shift(1, x(0), NaturalPlus{x(0), NaturalLit(3)}), NaturalPlus{x(1), NaturalLit(3)}),
	Entry("Shift(1, x, [] : List Natural) = [] : List Natural", Shift(1, x(0), EmptyList{Natural}), EmptyList{Natural}),
	Entry("Shift(1, x, [x]) = [x@1]", Shift(1, x(0), MakeList(Var{"x", 0})), MakeList(Var{"x", 1})),
	Entry("Shift(1, x, if x then x else x) = if x@1 then x@1 else x@1", Shift(1, x(0), BoolIf{x(0), x(0), x(0)}), BoolIf{x(1), x(1), x(1)}),
	Entry("Shift(1, x, let x = 3 in x) = let x = 3 in x",
		Shift(1, x(0), MakeLet(x(0), Binding{Variable: "x", Value: NaturalLit(3)})),
		MakeLet(x(0), Binding{Variable: "x", Value: NaturalLit(3)})),
	Entry("Shift(1, x, let x = 3 in x@1) = let x = 3 in x@2",
		Shift(1, x(0), MakeLet(x(1), Binding{Variable: "x", Value: NaturalLit(3)})),
		MakeLet(x(2), Binding{Variable: "x", Value: NaturalLit(3)})),
	Entry("Shift(1, x, let y : x = 3 in y) = let y : x@1 = 3 in y",
		Shift(1, x(0), MakeLet(y(0), Binding{Variable: "y", Annotation: x(0), Value: NaturalLit(3)})),
		MakeLet(y(0), Binding{Variable: "y", Annotation: x(1), Value: NaturalLit(3)})),
)

var _ = DescribeTable("Subst",
	func(actual Expr, expected Expr) {
		Expect(actual).To(Equal(expected))
	},
	Entry("Subst(_, _, Type) = Type", Subst(x(0), Sort, Type), Type),
	Entry("Subst(_, _, Natural) = Natural", Subst(x(0), Sort, Natural), Natural),
	Entry("Subst(_, _, 3) = Natural", Subst(x(0), Sort, NaturalLit(3)), NaturalLit(3)),
	Entry("Subst(x, 10, x + 3) = 10 + 3", Subst(x(0), NaturalLit(10), NaturalPlus{x(0), NaturalLit(3)}), NaturalPlus{NaturalLit(10), NaturalLit(3)}),
	Entry("Subst(x, Natural, x) = Natural", Subst(x(0), Natural, x(0)), Natural),
	Entry("Subst(x@1, Natural, x) = x", Subst(x(1), Natural, x(0)), x(0)),
	Entry("Subst(x, Natural, x@1) = x@1", Subst(x(0), Natural, x(1)), x(1)),
	Entry("Subst(x, Natural, λ( x) -> x) = λ( Natural) -> x",
		Subst(x(0), Natural, &LambdaExpr{"x", x(0), x(0)}),
		&LambdaExpr{"x", Natural, x(0)}),
	Entry("Subst(x, Natural, (λ(x: x) -> x) 3) = (λ(x: Natural) -> x) 3",
		Subst(x(0), Natural, &App{&LambdaExpr{"x", x(0), x(0)}, NaturalLit(3)}),
		&App{&LambdaExpr{"x", Natural, x(0)}, NaturalLit(3)}),
	Entry("Subst(x, Natural, y : x) = y : Natural", Subst(x(0), Natural, Annot{y(0), x(0)}), Annot{y(0), Natural}),
	Entry("Subst(x, Natural, x : y) = Natural : y", Subst(x(0), Natural, Annot{x(0), y(0)}), Annot{Natural, y(0)}),
	Entry("Subst(x, 3, [] : List Natural) = [] : List Type", Subst(x(0), NaturalLit(3), EmptyList{Natural}), EmptyList{Natural}),
	Entry("Subst(x, 3, [x]) = [3]", Subst(x(0), NaturalLit(3), MakeList(x(0))), MakeList(NaturalLit(3))),
	Entry("Subst(x, 3, if True then x else x) = if True then 3 else 3", Subst(x(0), NaturalLit(3), BoolIf{True, x(0), x(0)}), BoolIf{True, NaturalLit(3), NaturalLit(3)}),
	Entry("Subst(x, True, if x then 3 else 4) = if True then 3 else 4", Subst(x(0), True, BoolIf{x(0), NaturalLit(3), NaturalLit(4)}), BoolIf{True, NaturalLit(3), NaturalLit(4)}),
	Entry("Subst(x, 10, let x = 3 in x) = let x = 3 in x",
		Subst(x(0), NaturalLit(10), MakeLet(x(0), Binding{Variable: "x", Value: NaturalLit(3)})),
		MakeLet(x(0), Binding{Variable: "x", Value: NaturalLit(3)})),
	Entry("Subst(x, 10, let x = 3 in x@1) = let x = 3 in 10",
		Subst(x(0), NaturalLit(10), MakeLet(x(1), Binding{Variable: "x", Value: NaturalLit(3)})),
		MakeLet(NaturalLit(10), Binding{Variable: "x", Value: NaturalLit(3)})),
	Entry("Subst(x, Natural, let y : x = 3 in y) = let y : Natural = 3 in y",
		Subst(x(0), Natural, MakeLet(y(0), Binding{Variable: "y", Annotation: x(0), Value: NaturalLit(3)})),
		MakeLet(y(0), Binding{Variable: "y", Annotation: Natural, Value: NaturalLit(3)})),
)

var _ = Describe("Normalize", func() {
	DescribeTable("Normalize",
		func(in Expr, expected Expr) {
			actual := in.Normalize()
			Expect(actual).To(Equal(expected))
		},
		Entry("Sort", Sort, Sort),
		Entry("Kind", Kind, Kind),
		Entry("Type", Type, Type),
		Entry("Natural", Natural, Natural),
		Entry("Integer", Integer, Integer),
		Entry("List Integer", &App{List, Integer}, &App{List, Integer}),
		Entry("3", NaturalLit(3), NaturalLit(3)),
		Entry("3 + 5 » 8", NaturalPlus{NaturalLit(3), NaturalLit(5)}, NaturalLit(8)),
		Entry("3 + 5 : Natural » 8", Annot{NaturalPlus{NaturalLit(3), NaturalLit(5)}, Natural}, NaturalLit(8)),
		Entry("3 + x » 3 + x", NaturalPlus{NaturalLit(3), x(0)}, NaturalPlus{NaturalLit(3), x(0)}),
		Entry("x + 0 » x", NaturalPlus{x(0), NaturalLit(0)}, x(0)),
		Entry("0 + x » x", NaturalPlus{NaturalLit(0), x(0)}, x(0)),
		Entry("(λ(x : Natural) → x) 3 » 3", &App{&LambdaExpr{"x", Natural, x(0)}, NaturalLit(3)}, NaturalLit(3)),
		Entry("λ(x : Natural) → 2 + 3 » λ(x : Natural) → 5", &LambdaExpr{"x", Natural, NaturalPlus{NaturalLit(2), NaturalLit(3)}}, &LambdaExpr{"x", Natural, NaturalLit(5)}),
		Entry("[3 + 5] » [8]", MakeList(NaturalPlus{NaturalLit(3), NaturalLit(5)}), MakeList(NaturalLit(8))),
		Entry("[] : List Natural", EmptyList{Natural}, EmptyList{Natural}),
		Entry("[] : List ((λ(x : Type) → x) Natural) → [] : List Natural",
			EmptyList{&App{&LambdaExpr{"x", Type, x(0)}, Natural}},
			EmptyList{Natural}),
		Entry("if True then 3 else 4 » 3", BoolIf{True, NaturalLit(3), NaturalLit(4)}, NaturalLit(3)),
		Entry("if False then 3 else 4 » 4", BoolIf{False, NaturalLit(3), NaturalLit(4)}, NaturalLit(4)),
		Entry("let x = 3 in x » 3", MakeLet(x(0), Binding{Variable: "x", Value: NaturalLit(3)}), NaturalLit(3)),
		Entry("let x = 3 let y = x in y » 3",
			MakeLet(y(0), Binding{Variable: "x", Value: NaturalLit(3)}, Binding{Variable: "y", Value: x(0)}), NaturalLit(3)),
	)
})
