package core

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("judgmentallyEqual",
	func(in, out Term, expected bool) {
		Expect(AlphaEquivalent(in, out)).To(Equal(expected))
	},
	Entry("Unequal things", Bool, Natural, false),
	Entry("Lambda terms with the same label",
		NewLambda("a", Natural, NewVar("a")),
		NewLambda("a", Natural, NewVar("a")),
		true),
	Entry("Lambda terms with different labels",
		NewLambda("a", Natural, NewVar("a")),
		NewLambda("b", Natural, NewVar("b")),
		true),
	Entry("Lambda terms with different bodies",
		NewLambda("a", Natural, NewVar("a")),
		NewLambda("b", Natural, NaturalLit(3)),
		false),
	Entry("Pi types with the same label",
		NewPi("a", Type, Apply(List, NewVar("a"))),
		NewPi("a", Type, Apply(List, NewVar("a"))),
		true),
	Entry("Pi types with different labels",
		NewPi("a", Type, Apply(List, NewVar("a"))),
		NewPi("b", Type, Apply(List, NewVar("b"))),
		true),
)
