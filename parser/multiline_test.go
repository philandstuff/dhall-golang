package parser_test

import (
	. "github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("RemoveLeadingCommonIndent removes leading common indent", func() {
	DescribeTable("when the TextLitTerm has no interpolations", func(input, expected string) {
		actual := parser.RemoveLeadingCommonIndent(TextLitTerm{Suffix: input})
		Expect(actual).To(Equal(TextLitTerm{Suffix: expected}))
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
	DescribeTable("when the TextLitTerm has interpolations", func(input, expected TextLitTerm) {
		actual := parser.RemoveLeadingCommonIndent(input)
		Expect(actual).To(Equal(expected))
	},
		Entry("when every line has a 3-space prefix",
			TextLitTerm{Chunks: Chunks{{Prefix: "   ", Expr: True}}, Suffix: "\n   foo"},
			TextLitTerm{Chunks: Chunks{{Prefix: "", Expr: True}}, Suffix: "\nfoo"}),
		Entry("when there is trailing space after an interpolation",
			TextLitTerm{Chunks: Chunks{{Prefix: "   ", Expr: True}}, Suffix: "   foo"},
			TextLitTerm{Chunks: Chunks{{Prefix: "", Expr: True}}, Suffix: "   foo"}),
		Entry("when there is trailing space after every interpolation",
			// this is ''
			//    ${True}   ${True}   foo''
			TextLitTerm{Chunks: Chunks{{Prefix: "   ", Expr: True}, {Prefix: "   ", Expr: True}}, Suffix: "   foo"},
			TextLitTerm{Chunks: Chunks{{Prefix: "", Expr: True}, {Prefix: "   ", Expr: True}}, Suffix: "   foo"}),
		Entry("when there are multiple interpolations on the same line",
			// this is ''
			//    ${True} ${True}
			//    ''
			TextLitTerm{Chunks: Chunks{{Prefix: "   ", Expr: True}, {Prefix: " ", Expr: True}}, Suffix: "\n   "},
			TextLitTerm{Chunks: Chunks{{Prefix: "", Expr: True}, {Prefix: " ", Expr: True}}, Suffix: "\n"}),
		Entry("when there is no trailing newline",
			// this is ''
			//    ${True}''
			TextLitTerm{Chunks: Chunks{{Prefix: "   ", Expr: True}}, Suffix: ""},
			TextLitTerm{Chunks: Chunks{{Prefix: "", Expr: True}}, Suffix: ""}),
	)
})
