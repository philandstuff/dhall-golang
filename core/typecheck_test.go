package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
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

func typecheckTest(t Term, expectedType Term) {
	actualType, err := TypeOf(t)
	Ω(err).ShouldNot(HaveOccurred())
	Ω(actualType).Should(
		Equal(expectedType))
}

var _ = Describe("TypeOf", func() {
	DescribeTable("Universe",
		typecheckTest,
		Entry("Type : Kind", Type, Kind),
		Entry("Kind : Sort", Kind, Sort),
	)
	DescribeTable("Builtin",
		typecheckTest,
		Entry(`Natural : Type`, Natural, Type),
		Entry(`List : Type -> Type`, List, FnType(Type, Type)),
	)
	DescribeTable("Lambda",
		typecheckTest,
		Entry("λ(x : Natural) → x : ∀(x : Natural) → Natural",
			Mkλ("x", Natural, Bound("x")),
			MkΠ("x", Natural, Natural)),
		Entry("λ(a : Type) → ([] : List a) : ∀(a : Type) → List a -- check presence of variables in resulting type",
			Mkλ("a", Type,
				EmptyList{Apply(List, Bound("a"))}),
			MkΠ("a", Type, Apply(List, Bound("a")))),
	)
	DescribeTable("Pi",
		typecheckTest,
		Entry(`Natural → Natural : Type`, FnType(Natural, Natural), Type),
	)
	DescribeTable("Application",
		typecheckTest,
		Entry(`List Natural : Type`, Apply(List, Natural), Type),
	)
	DescribeTable("Others",
		typecheckTest,
		Entry(`3 : Natural`, NaturalLit(3), Natural),
		Entry(`[] : List Natural : List Natural`,
			EmptyList{Apply(List, Natural)}, Apply(List, Natural)),
	)
	DescribeTable("Expected failures",
		func(t Term) {
			_, err := TypeOf(t)
			Ω(err).Should(HaveOccurred())
		},
		// Universe
		Entry(`Sort -- Sort has no type`,
			Sort),
		// EmptyList
		Entry(`[] : List 3 -- not a valid list type`,
			EmptyList{Apply(List, NaturalLit(3))}),
		Entry(`[] : Natural -- not in form "List a"`,
			EmptyList{Natural}),

		// AppTerm
		Entry(`Sort Type -- Fn of AppTerm doesn't typecheck`,
			Apply(Sort, Type)),
		Entry(`List Sort -- Arg of AppTerm doesn't typecheck`,
			Apply(List, Sort)),
		Entry(`List 3 -- Arg of AppTerm doesn't match function input type`,
			Apply(List, NaturalLit(3))),
		Entry(`Natural Natural -- Fn of AppTerm isn't of function type`,
			Apply(Natural, Natural)),
	)
})
