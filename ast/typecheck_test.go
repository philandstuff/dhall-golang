package ast_test

import (
	. "github.com/philandstuff/dhall-golang/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Rule",
	func(a Const, b Const, expected Const) {
		actual, err := Rule(a, b)
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(expected))
	},
	Entry("Type → Type : Type", Type, Type, Type),
	Entry("Kind → Type : Type", Kind, Type, Type),
	Entry("Sort → Type : Type", Sort, Type, Type),
	Entry("Type → Kind : Kind", Kind, Kind, Kind),
	Entry("Kind → Kind : Kind", Kind, Kind, Kind),
	Entry("Sort → Kind : Sort", Sort, Kind, Sort),
	Entry("Type → Sort : Sort", Type, Sort, Sort),
	Entry("Kind → Sort : Sort", Kind, Sort, Sort),
	Entry("Sort → Sort : Sort", Sort, Sort, Sort),
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
		Entry("(3 : (λ(x : Type) → x) Natural) : Natural", Annot{NaturalLit(3), Apply(&LambdaExpr{"x", Type, x(0)}, Natural)}, Natural),
		Entry("3 + 5 : Natural", NaturalPlus(NaturalLit(3), NaturalLit(5)), Natural),
		Entry("3 * 5 : Natural", NaturalTimes(NaturalLit(3), NaturalLit(5)), Natural),
	)
	DescribeTable("Builtins",
		expectType,
		Entry("Double : Type", Double, Type),
		Entry("Text : Type", Text, Type),
		Entry("Bool : Type", Bool, Type),
		Entry("Natural : Type", Natural, Type),
		Entry("Integer : Type", Integer, Type),
		Entry("List : Type → Type", List, FnType(Type, Type)),
		Entry("Optional : Type → Type", Optional, FnType(Type, Type)),
		Entry("None : ∀(A: Type) → Optional A", None, &Pi{"A", Type, Apply(Optional, MkVar("A"))}),
	)
	DescribeTable("Function types",
		expectType,
		Entry("λ(x : Natural) → x : ∀(x : Natural) → Natural",
			&LambdaExpr{"x", Natural, x(0)},
			&Pi{"x", Natural, Natural}),
		Entry("λ(x : Natural) → x : Natural → Natural",
			&LambdaExpr{"x", Natural, x(0)},
			FnType(Natural, Natural)),
		Entry("(λ(x : Natural) → x) 3 : Natural",
			Apply(&LambdaExpr{"x", Natural, x(0)},
				NaturalLit(3)),
			Natural),
		Entry("(λ(f : Natural → Natural) → f 3) λ(n : Natural) → n+1 : Natural",
			Apply(
				&LambdaExpr{"f", FnType(Natural, Natural), Apply(MkVar("f"), NaturalLit(3))},
				&LambdaExpr{"n", Natural, NaturalPlus(MkVar("n"), NaturalLit(1))}),
			Natural),
		Entry("λ(x : Bool) → λ(x : Natural) → x : ∀(x : Bool) → ∀(x : Natural) → Natural",
			&LambdaExpr{"x", Bool,
				&LambdaExpr{"x", Natural,
					x(0)}},
			&Pi{"x", Bool,
				&Pi{"x", Natural, Natural}}),
		Entry("λ(x : Bool) → λ(x : Natural) → x@1 : ∀(x : Bool) → ∀(x : Natural) → Bool",
			&LambdaExpr{"x", Bool,
				&LambdaExpr{"x", Natural,
					x(1)}},
			&Pi{"x", Bool,
				&Pi{"x", Natural, Bool}}),
	)
	DescribeTable("List types",
		expectType,
		Entry("([] : List Natural) : List Natural",
			EmptyList{Apply(List, Natural)},
			Apply(List, Natural)),
		Entry("[3] : List Natural",
			MakeList(NaturalLit(3)),
			Apply(List, Natural)),
		Entry("([3] : List Natural) : List Natural",
			Annot{MakeList(NaturalLit(3)), Apply(List, Natural)},
			Apply(List, Natural)),
		Entry("[3] # [4] : List Natural",
			ListAppend(MakeList(NaturalLit(3)), MakeList(NaturalLit(4))),
			Apply(List, Natural)),
	)
	DescribeTable("Optional types",
		expectType,
		Entry("None Natural : Optional Natural",
			Apply(None, Natural),
			Apply(Optional, Natural)),
		Entry("Some 3 : Optional Natural",
			Some{NaturalLit(3)},
			Apply(Optional, Natural)),
	)
	DescribeTable("Records",
		expectType,
		Entry("{=} : {}",
			RecordLit{},
			Record{}),
		Entry("{} : Type",
			Record{},
			Type),
		Entry("{foo = 3} : {foo : Natural}",
			RecordLit{"foo": NaturalLit(3)},
			Record{"foo": Natural}),
		Entry("{foo = 3, bar = +4} : {foo : Natural, bar : Integer}",
			RecordLit{"foo": NaturalLit(3), "bar": IntegerLit(4)},
			Record{"foo": Natural, "bar": Integer}),
		Entry("{foo : Type} : Kind",
			Record{"foo": Type},
			Kind),
		Entry(`\(x : {y : Natural}) → x.y : {y : Natural} → Natural`,
			&LambdaExpr{"x", Record{"y": Natural}, Field{x(0), "y"}},
			FnType(Record{"y": Natural}, Natural)),
		Entry("record types can mix types and kinds",
			Record{"foo": Natural, "bar": Type}, Kind),
		Entry("record types can mix kinds and sorts",
			Record{"foo": Type, "bar": Kind}, Sort),
		Entry("record literals can mix terms and types",
			RecordLit{"foo": NaturalLit(3), "bar": Natural},
			Record{"foo": Natural, "bar": Type},
		),
		Entry("record literals can mix types and kinds",
			RecordLit{"foo": Natural, "bar": Type},
			Record{"foo": Type, "bar": Kind},
		),
	)
	DescribeTable("Other",
		expectType,
		Entry("if True then 3 else 4 : Natural",
			BoolIf{True, NaturalLit(3), NaturalLit(4)}, Natural),
		Entry("let x = 3 in x : Natural",
			MakeLet(MkVar("x"), Binding{Variable: "x", Value: NaturalLit(3)}),
			Natural),
		Entry("let x = True in 3 : Natural",
			MakeLet(NaturalLit(3), Binding{Variable: "x", Value: True}),
			Natural),
		// Below case from https://github.com/dhall-lang/dhall-lang/pull/69
		Entry("(let x = Natural in 1 : x) : Natural",
			MakeLet(Annot{NaturalLit(1), MkVar("x")}, Binding{Variable: "x", Value: Natural}),
			Natural),
	)
	DescribeTable("Sorts",
		// We can't typecheck sorts using expectType because
		// it will fail by not finding a type for Sort so we
		// have a custom function (expectSort instead of
		// expectType) for this case
		expectSort,
		Entry("Kind", Kind),
		Entry("{foo : Kind}", Record{"foo": Kind}),
	)
	DescribeTable("Type errors",
		func(in Expr) {
			_, err := in.TypeWith(EmptyContext())
			Expect(err).To(HaveOccurred())
		},
		Entry("3 +5", Apply(NaturalLit(3), IntegerLit(5))),
		Entry("3 : Integer", Annot{NaturalLit(3), Integer}),
		Entry("\"foo ${3} baz\" : Text", TextLit{Chunks{Chunk{"foo ", NaturalLit(3)}}, " baz"}),
		Entry("if True then 3 else +4", BoolIf{True, NaturalLit(3), IntegerLit(4)}),
		Entry("if 2 then 3 else 4", BoolIf{NaturalLit(3), NaturalLit(3), NaturalLit(4)}),
		Entry("if True then Type else (Type -> Type)", BoolIf{True, Type, FnType(Type, Type)}),
		Entry("[3] # [True]", ListAppend(MakeList(NaturalLit(3)), MakeList(True))),
		Entry("let x : Bool = 3 in 5", MakeLet(NaturalLit(5),
			Binding{Variable: "x", Annotation: Bool, Value: NaturalLit(3)})),
		Entry("record types can't have terms as types", Record{"foo": NaturalLit(3)}),
		Entry("you can't select from non-records", Field{NaturalLit(3), "y"}),
		Entry("you can't select from record types, only record literals", Field{Record{"y": Natural}, "y"}),
		Entry("you can't select nonexistent fields", Field{RecordLit{"foo": Natural}, "y"}),
	)
})
