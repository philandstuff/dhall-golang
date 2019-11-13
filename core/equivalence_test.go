package core

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("judgmentallyEqual",
	func(in, out Term, expected bool) {
		Expect(judgmentallyEqual(in, out)).To(Equal(expected))
	},
	Entry("Unequal things", Bool, Natural, false),
	Entry("Lambda terms with the same label",
		MkLambdaTerm("a", Natural, Bound("a")),
		MkLambdaTerm("a", Natural, Bound("a")),
		true),
	Entry("Lambda terms with different labels",
		MkLambdaTerm("a", Natural, Bound("a")),
		MkLambdaTerm("b", Natural, Bound("b")),
		true),
	Entry("Lambda terms with different bodies",
		MkLambdaTerm("a", Natural, Bound("a")),
		MkLambdaTerm("b", Natural, NaturalLit(3)),
		false),
	Entry("Pi types with the same label",
		MkPiTerm("a", Type, Apply(List, Bound("a"))),
		MkPiTerm("a", Type, Apply(List, Bound("a"))),
		true),
	Entry("Pi types with different labels",
		MkPiTerm("a", Type, Apply(List, Bound("a"))),
		MkPiTerm("b", Type, Apply(List, Bound("b"))),
		true),
)
