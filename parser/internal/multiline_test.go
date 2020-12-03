package internal

import (
	. "github.com/philandstuff/dhall-golang/v6/term"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("removeLeadingCommonIndent removes leading common indent", func() {
	DescribeTable("when the TextLit has no interpolations", func(input, expected string) {
		actual := removeLeadingCommonIndent(PlainText(input))
		Expect(actual).To(Equal(PlainText(expected)))
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
		actual := removeLeadingCommonIndent(input)
		Expect(actual).To(Equal(expected))
	},
		Entry("when every line has a 3-space prefix",
			TextLit{Chunks: Chunks{{Prefix: "   ", Expr: True}}, Suffix: "\n   foo"},
			TextLit{Chunks: Chunks{{Prefix: "", Expr: True}}, Suffix: "\nfoo"}),
		Entry("when there is trailing space after an interpolation",
			TextLit{Chunks: Chunks{{Prefix: "   ", Expr: True}}, Suffix: "   foo"},
			TextLit{Chunks: Chunks{{Prefix: "", Expr: True}}, Suffix: "   foo"}),
		Entry("when there is trailing space after every interpolation",
			// this is ''
			//    ${True}   ${True}   foo''
			TextLit{Chunks: Chunks{{Prefix: "   ", Expr: True}, {Prefix: "   ", Expr: True}}, Suffix: "   foo"},
			TextLit{Chunks: Chunks{{Prefix: "", Expr: True}, {Prefix: "   ", Expr: True}}, Suffix: "   foo"}),
		Entry("when there are multiple interpolations on the same line",
			// this is ''
			//    ${True} ${True}
			//    ''
			TextLit{Chunks: Chunks{{Prefix: "   ", Expr: True}, {Prefix: " ", Expr: True}}, Suffix: "\n   "},
			TextLit{Chunks: Chunks{{Prefix: "", Expr: True}, {Prefix: " ", Expr: True}}, Suffix: "\n"}),
		Entry("when there is no trailing newline",
			// this is ''
			//    ${True}''
			TextLit{Chunks: Chunks{{Prefix: "   ", Expr: True}}, Suffix: ""},
			TextLit{Chunks: Chunks{{Prefix: "", Expr: True}}, Suffix: ""}),
	)
})
