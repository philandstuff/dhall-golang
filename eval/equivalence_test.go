package eval

import (
	. "github.com/philandstuff/dhall-golang/core"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("judgmentallyEqual",
	func(in, out Term, expected bool) {
		Expect(judgmentallyEqual(in, out)).To(Equal(expected))
	},
	Entry("Unequal things", Bool, Natural, false),
	Entry("Lambda terms with the same label",
		Mkλ("a", Natural, Bound("a")),
		Mkλ("a", Natural, Bound("a")),
		true),
	Entry("Lambda terms with different labels",
		Mkλ("a", Natural, Bound("a")),
		Mkλ("b", Natural, Bound("b")),
		true),
	Entry("Lambda terms with different bodies",
		Mkλ("a", Natural, Bound("a")),
		Mkλ("b", Natural, NaturalLit(3)),
		false),
	Entry("Pi types with the same label",
		MkΠ("a", Type, Apply(List, Bound("a"))),
		MkΠ("a", Type, Apply(List, Bound("a"))),
		true),
	Entry("Pi types with different labels",
		MkΠ("a", Type, Apply(List, Bound("a"))),
		MkΠ("b", Type, Apply(List, Bound("b"))),
		true),
)
