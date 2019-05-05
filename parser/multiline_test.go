package parser_test

import (
	. "github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemoveLeadingCommonIndent removes leading common indent", func() {
	DescribeTable("when the TextLit has no interpolations", func(input, expected string) {
		actual := parser.RemoveLeadingCommonIndent(TextLit{Suffix: input})
		Expect(actual).To(Equal(TextLit{Suffix: expected}))
	},
		Entry("when every line has a 3-space prefix",
			`   foo
   bar`, "foo\nbar"),
		Entry("when every line has a 3-space prefix except on the last (empty) line",
			`   foo
   bar
`, "   foo\n   bar\n"),
		Entry("when every line has a 2-tab prefix",
			`		foo
		bar`, "foo\nbar"),
		Entry("when every line has a one-tab-one-space prefix",
			`	 foo
	 bar`, "foo\nbar"),
		Entry("when one line is tab-then-space but the other is space-then-tab (so there is no common prefix)",
			`	 foo
 	bar`, "\t foo\n \tbar"),
		Entry("when there is a blank line in the middle, it's ignored for prefix comparison",
			`   foo

   bar`, "foo\n\nbar"),
	)
	DescribeTable("when the TextLit has interpolations", func(input, expected TextLit) {
		actual := parser.RemoveLeadingCommonIndent(input)
		Expect(actual).To(Equal(expected))
	},
		Entry("when every line has a 3-space prefix",
			TextLit{Chunks: Chunks{{Prefix: "   ", Expr: True}}, Suffix: "   foo"},
			TextLit{Chunks: Chunks{{Prefix: "", Expr: True}}, Suffix: "foo"}),
	)
})
