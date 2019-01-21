package pprint_test

import (
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/philandstuff/dhall-golang/pprint"

	// . "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/prataprc/goparsec"
)

var _ = DescribeTable("PrettyPrint",
	func(input string, expected_regex string) {
		buf := gbytes.NewBuffer()
		node, _ := parser.CompleteExpression(parsec.NewScanner([]byte(input)))
		pprint.PrettyPrint(node, buf)
		Expect(buf).To(gbytes.Say(expected_regex))
	},
	Entry("converts ascii to unicode", "\\(foo : bar) -> baz", `λ\(foo : bar\) → baz`),
	Entry("preserves unicode", "λ(foo : bar) → baz", `λ\(foo : bar\) → baz`),
	Entry("preserves leading comments", "-- comment\nλ(foo : bar) → baz", `-- comment\nλ\(foo : bar\) → baz`),
	Entry("preserves interstitial comments", "λ(foo : bar) --comment\n → baz", `λ\(foo : bar\) --comment\n→ baz`),
)
