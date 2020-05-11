package dhall

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/philandstuff/dhall-golang/v3/core"
	"github.com/philandstuff/dhall-golang/v3/parser"
)

func parseAndTypecheckTest(source string, expectedTypeSource string) {
	parsed, err := parser.Parse("-", []byte(source))
	Ω(err).ShouldNot(HaveOccurred())
	parsedType, err := parser.Parse("-", []byte(expectedTypeSource))
	Ω(err).ShouldNot(HaveOccurred())
	Ω(TypeOf(parsed)).Should(BeAlphaEquivalentTo(parsedType))
}

var _ = Describe("Regression tests", func() {
	DescribeTable("TypeOf",
		parseAndTypecheckTest,
		// The problem here was really hard to spot.. but it was
		// because rebindLocal() had a bug for handling NonEmptyLists.
		// The only time rebindLocal() is called is on types, and the
		// only types which can contain a NonEmptyList term at the
		// moment are equivalences.
		Entry("Rebinding problem in lists", `
  λ(x : Bool)
→   assert
  :   [True , x]
	≡ [True , x]
`, `
      ∀(x : Bool)
    → [True , x]
	≡ [True , x]
`),
	)
})
