package eval_test

import (
	. "github.com/philandstuff/dhall-golang/core"
	. "github.com/philandstuff/dhall-golang/eval"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeOf", func() {
	Describe("Const", func() {
		It("⊦ Type : Kind", func() {
			Expect(TypeOf(Context{}, Type)).
				To(Equal(Kind))
		})
		It("⊦ Kind : Sort", func() {
			Expect(TypeOf(Context{}, Kind)).
				To(Equal(Sort))
		})
	})
	DescribeTable("Builtin",
		func(t Term, expectedType Term) {
			actualType, err := TypeOf(Context{}, t)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(Quote(actualType)).Should(
				Equal(expectedType))
		},
		Entry(`Natural : Type`, Natural, Type),
		Entry(`List : Type -> Type`, List, FnType(Type, Type)),
	)
	Describe("Lambda", func() {
		It("⊦ λ(x : Natural) → x : ∀(x : Natural) → Natural", func() {
			t, err := TypeOf(Context{}, Mkλ("x", Natural, Bound("x")))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(Quote(t)).Should(
				Equal(MkΠ("x", Natural, Natural)))
		})
		XIt("⊦ λ(a : Type) → [] : List a : ∀(a : Type) → List a", func() {
			t, err := TypeOf(Context{}, Mkλ("x", Natural, Bound("x")))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(Quote(t)).Should(
				Equal(MkΠ("x", Type, Bound("x"))))
		})
	})
	DescribeTable("Application",
		func(t Term, expectedType Term) {
			actualType, err := TypeOf(Context{}, t)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(Quote(actualType)).Should(
				Equal(expectedType))
		},
		Entry(`List Natural : Type`, Apply(List, Natural), Type),
	)
	DescribeTable("Others",
		func(t Term, expectedType Term) {
			actualType, err := TypeOf(Context{}, t)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(Quote(actualType)).Should(
				Equal(expectedType))
		},
		Entry(`3 : Natural`, NaturalLit(3), Natural),
		Entry(`[] : List Natural : List Natural`,
			EmptyList{Apply(List, Natural)}, Apply(List, Natural)),
	)
	DescribeTable("Expected failures",
		func(t Term) {
			_, err := TypeOf(Context{}, t)
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
