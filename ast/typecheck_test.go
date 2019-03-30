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

func expectSort(in Expr) {
	typ, err := in.TypeWith(EmptyContext())
	Expect(err).ToNot(HaveOccurred())
	Expect(typ).To(Equal(Sort))
}

var _ = Describe("TypeCheck in empty context", func() {
	DescribeTable("Simple types",
		expectType,
		Entry("Type : Kind", Type, Kind),
		// Kind : Sort is handled using expectSort elsewhere
		Entry("True : Bool", BoolLit(true), Bool),
		Entry("3.3 : Double", DoubleLit(3.3), Double),
		Entry(`"" : Text`, TextLit{}, Text),
		Entry(`"foo ${"bar"} baz" : Text`, TextLit{Chunks{Chunk{"foo ", TextLit{Suffix: "bar"}}}, " baz"}, Text),
		Entry("3 : Natural", NaturalLit(3), Natural),
		Entry("(3 : Natural) : Natural", Annot{NaturalLit(3), Natural}, Natural),
		Entry("(3 : (λ(x : Type) → x) Natural) : Natural", Annot{NaturalLit(3), &App{&LambdaExpr{"x", Type, x(0)}, Natural}}, Natural),
		Entry("3 + 5 : Natural", NaturalPlus{NaturalLit(3), NaturalLit(5)}, Natural),
		Entry("3 * 5 : Natural", NaturalTimes{NaturalLit(3), NaturalLit(5)}, Natural),
	)
	DescribeTable("Builtins",
		expectType,
		Entry("Double : Type", Double, Type),
		Entry("Text : Type", Text, Type),
		Entry("Bool : Type", Bool, Type),
		Entry("Natural : Type", Natural, Type),
		Entry("Integer : Type", Integer, Type),
		Entry("List : Type → Type", List, &Pi{"_", Type, Type}),
		Entry("Optional : Type → Type", Optional, &Pi{"_", Type, Type}),
		Entry("None : ∀(A: Type) → Optional A", None, &Pi{"A", Type, &App{Optional, Var{"A", 0}}}),
	)
	DescribeTable("Function types",
		expectType,
		Entry("λ(x : Natural) → x : ∀(x : Natural) → Natural",
			&LambdaExpr{"x", Natural, x(0)},
			&Pi{"x", Natural, Natural}),
		Entry("λ(x : Natural) → x : Natural → Natural",
			&LambdaExpr{"x", Natural, x(0)},
			&Pi{"_", Natural, Natural}),
		Entry("(λ(x : Natural) → x) 3 : Natural",
			&App{&LambdaExpr{"x", Natural, x(0)},
				NaturalLit(3)},
			Natural),
		Entry("(λ(f : Natural → Natural) → f 3) λ(n : Natural) → n+1 : Natural",
			&App{
				Fn:  &LambdaExpr{"f", &Pi{"_", Natural, Natural}, &App{Var{"f", 0}, NaturalLit(3)}},
				Arg: &LambdaExpr{"n", Natural, NaturalPlus{Var{"n", 0}, NaturalLit(1)}}},
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
	)
	DescribeTable("Optional types",
		expectType,
		Entry("None Natural : Optional Natural",
			&App{None, Natural},
			&App{Optional, Natural}),
		Entry("Some 3 : Optional Natural",
			Some{NaturalLit(3)},
			&App{Optional, Natural}),
	)
	DescribeTable("Records",
		expectType,
		Entry("{=} : {}",
			RecordLit(map[string]Expr{}),
			Record(map[string]Expr{})),
		Entry("{} : Type",
			Record(map[string]Expr{}),
			Type),
		Entry("{foo = 3} : {foo : Natural}",
			RecordLit(map[string]Expr{"foo": NaturalLit(3)}),
			Record(map[string]Expr{"foo": Natural})),
		Entry("{foo = 3, bar = +4} : {foo : Natural, bar : Integer}",
			RecordLit(map[string]Expr{"foo": NaturalLit(3), "bar": IntegerLit(4)}),
			Record(map[string]Expr{"foo": Natural, "bar": Integer})),
		Entry("{foo : Type} : Kind",
			Record(map[string]Expr{"foo": Type}),
			Kind),
		Entry(`\(x : {y : Natural}) → x.y : {y : Natural} → Natural`,
			&LambdaExpr{"x", Record(map[string]Expr{"y": Natural}), Field{x(0), "y"}},
			&Pi{"_", Record(map[string]Expr{"y": Natural}), Natural}),
	)
	DescribeTable("Other",
		expectType,
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
	DescribeTable("Sorts",
		// We can't typecheck sorts using expectType because
		// it will fail by not finding a type for Sort so we
		// have a custom function (expectSort instead of
		// expectType) for this case
		expectSort,
		Entry("Kind", Kind),
		Entry("{foo : Kind}", Record(map[string]Expr{"foo": Kind})),
	)
	DescribeTable("Type errors",
		func(in Expr) {
			_, err := in.TypeWith(EmptyContext())
			Expect(err).To(HaveOccurred())
		},
		Entry("3 +5", &App{NaturalLit(3), IntegerLit(5)}),
		Entry("3 : Integer", Annot{NaturalLit(3), Integer}),
		Entry("\"foo ${3} baz\" : Text", TextLit{Chunks{Chunk{"foo ", NaturalLit(3)}}, " baz"}),
		Entry("if True then 3 else +4", BoolIf{True, NaturalLit(3), IntegerLit(4)}),
		Entry("if 2 then 3 else 4", BoolIf{NaturalLit(3), NaturalLit(3), NaturalLit(4)}),
		Entry("if True then Type else (Type -> Type)", BoolIf{True, Type, &Pi{"_", Type, Type}}),
		Entry("let x : Bool = 3 in 5", MakeLet(NaturalLit(5),
			Binding{Variable: "x", Annotation: Bool, Value: NaturalLit(3)})),
		Entry("record types can't have Kind->Kind fields", Record(map[string]Expr{"foo": &Pi{Label: "_", Type: Kind, Body: Kind}})),
		Entry("record types can't have terms as types", Record(map[string]Expr{"foo": NaturalLit(3)})),
		Entry("record types can't mix types and kinds", Record(map[string]Expr{"foo": Natural, "bar": Type})),
		Entry("record types can't mix kinds and sorts", Record(map[string]Expr{"foo": Type, "bar": Kind})),
		Entry("record literals can't have Kind->Kind fields", RecordLit(map[string]Expr{"foo": &LambdaExpr{Label: "x", Type: Kind, Body: x(0)}})),
		Entry("record literals can't mix types and kinds", RecordLit(map[string]Expr{"foo": NaturalLit(3), "bar": Natural})),
		Entry("record literals can't mix kinds and sorts", RecordLit(map[string]Expr{"foo": Natural, "bar": Type})),
		Entry("you can't select from non-records", Field{NaturalLit(3), "y"}),
		Entry("you can't select from record types, only record literals", Field{Record(map[string]Expr{"y": Natural}), "y"}),
		Entry("you can't select nonexistent fields", Field{RecordLit(map[string]Expr{"foo": Natural}), "y"}),
	)
})
