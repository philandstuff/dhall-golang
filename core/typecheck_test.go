package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/philandstuff/dhall-golang/v5/term"
)

var _ = DescribeTable("functionCheck",
	func(in, out, expected Universe) {
		Expect(functionCheck(in, out)).To(Equal(expected))
	},
	Entry(`Type ↝ Type : Type`, Type, Type, Type),
	Entry(`Kind ↝ Type : Type`, Kind, Type, Type),
	Entry(`Sort ↝ Type : Type`, Sort, Type, Type),
	Entry(`Type ↝ Kind : Kind`, Type, Kind, Kind),
	Entry(`Kind ↝ Kind : Kind`, Kind, Kind, Kind),
	Entry(`Sort ↝ Kind : Sort`, Sort, Kind, Sort),
	Entry(`Type ↝ Sort : Sort`, Type, Sort, Sort),
	Entry(`Kind ↝ Sort : Sort`, Kind, Sort, Sort),
	Entry(`Sort ↝ Sort : Sort`, Sort, Sort, Sort),
)

func typecheckTest(t term.Term, expectedType Value) {
	Ω(TypeOf(t)).Should(BeAlphaEquivalentTo(expectedType))
}

var _ = Describe("TypeOf", func() {
	DescribeTable("Universe",
		typecheckTest,
		Entry("Type : Kind", term.Type, Kind),
		Entry("Kind : Sort", term.Kind, Sort),
	)
	DescribeTable("Builtin",
		typecheckTest,
		Entry(`Natural : Type`, term.Natural, Type),
		Entry(`List : Type -> Type`, term.List, NewFnType("_", Type, Type)),
	)
	DescribeTable("Lambda",
		typecheckTest,
		Entry("λ(x : Natural) → x : ∀(x : Natural) → Natural",
			term.NewLambda("x", term.Natural, term.NewVar("x")),
			NewPi("x", Natural, func(Value) Value { return Natural })),
		Entry("λ(a : Type) → ([] : List a) : ∀(a : Type) → List a -- check presence of variables in resulting type",
			term.NewLambda("a", term.Type,
				term.EmptyList{term.App{term.List, term.NewVar("a")}}),
			NewPi("a", Type, func(a Value) Value {
				return ListOf{a}
			})),
		Entry("λ(a : Natural) → assert : a ≡ a -- check presence of variables in resulting type",
			term.NewLambda("a", term.Natural,
				term.Assert{term.Op{term.EquivOp, term.NewVar("a"), term.NewVar("a")}}),
			NewPi("a", Natural, func(a Value) Value {
				return oper{term.EquivOp, a, a}
			})),
	)
	DescribeTable("Pi",
		typecheckTest,
		Entry(`Natural → Natural : Type`, term.NewAnonPi(term.Natural, term.Natural), Type),
	)
	DescribeTable("Application",
		typecheckTest,
		Entry(`List Natural : Type`, term.App{term.List, term.Natural}, Type),
		Entry("(λ(a : Natural) → assert : a ≡ a) 3 -- check presence of variables in resulting type",
			term.Apply(
				term.NewLambda("a", term.Natural,
					term.Assert{term.Op{term.EquivOp, term.NewVar("a"), term.NewVar("a")}}),
				term.NaturalLit(3)),
			oper{term.EquivOp, NaturalLit(3), NaturalLit(3)}),
	)
	DescribeTable("Others",
		typecheckTest,
		Entry(`3 : Natural`, term.NaturalLit(3), Natural),
		Entry(`[] : List Natural : List Natural`,
			term.EmptyList{term.Apply(term.List, term.Natural)}, ListOf{Natural}),
	)
	DescribeTable("Expected failures",
		func(t term.Term) {
			_, err := TypeOf(t)
			Ω(err).Should(HaveOccurred())
		},
		// Universe
		Entry(`Sort -- Sort has no type`,
			term.Sort),
		// EmptyList
		Entry(`[] : List 3 -- not a valid list type`,
			term.EmptyList{term.Apply(term.List, term.NaturalLit(3))}),
		Entry(`[] : Natural -- not in form "List a"`,
			term.EmptyList{term.Natural}),

		// AppTerm
		Entry(`Sort Type -- Fn of AppTerm doesn't typecheck`,
			term.Apply(term.Sort, term.Type)),
		Entry(`List Sort -- Arg of AppTerm doesn't typecheck`,
			term.Apply(term.List, term.Sort)),
		Entry(`List 3 -- Arg of AppTerm doesn't match function input type`,
			term.Apply(term.List, term.NaturalLit(3))),
		Entry(`Natural Natural -- Fn of AppTerm isn't of function type`,
			term.Apply(term.Natural, term.Natural)),
	)
	DescribeTable("Detailed expected failures",
		func(t term.Term, msg string) {
			_, err := TypeOf(t)
			Ω(err).Should(MatchError(msg))
		},
		Entry(`2 === Type`,
			term.Equivalent(term.NaturalLit(2), term.Type), "Incomparable expression"),
	)
})
