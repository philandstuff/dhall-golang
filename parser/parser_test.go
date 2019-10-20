package parser_test

import (
	"math"

	. "github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func ParseAndCompare(input string, expected interface{}) {
	root, err := parser.Parse("test", []byte(input))
	Expect(err).ToNot(HaveOccurred())
	Expect(root).To(Equal(expected))
}

func ParseAndFail(input string) {
	_, err := parser.Parse("test", []byte(input))
	Expect(err).To(HaveOccurred())
}

var _ = Describe("Expression", func() {
	DescribeTable("simple expressions", ParseAndCompare,
		Entry("Type", `Type`, Type),
		Entry("Kind", `Kind`, Kind),
		Entry("Sort", `Sort`, Sort),
		Entry("Double", `Double`, Double),
		Entry("Text", `Text`, Text),
		Entry("DoubleLit", `3.0`, DoubleLit(3.0)),
		Entry("DoubleLit with exponent", (`3E5`), DoubleLit(3e5)),
		Entry("DoubleLit with sign", (`+3.0`), DoubleLit(3.0)),
		Entry("DoubleLit with everything", (`-5.0e1`), DoubleLit(-50.0)),
		Entry("Infinity", `Infinity`, DoubleLit(math.Inf(1))),
		Entry("-Infinity", `-Infinity`, DoubleLit(math.Inf(-1))),
		Entry("Integer", `Integer`, Integer),
		Entry("IntegerLit", `+1234`, IntegerLit(1234)),
		Entry("IntegerLit", `-3`, IntegerLit(-3)),
		Entry("Annotated expression", `3 : Natural`, Annot{NaturalLit(3), Natural}),
	)
	DescribeTable("bools", ParseAndCompare,
		Entry("Bool", `Bool`, Bool),
		Entry("True", `True`, BoolLit(true)),
		Entry("False", `False`, BoolLit(false)),
		Entry("if True then x else y", `if True then x else y`, IfTerm{True, Bound("x"), Bound("y")}),
	)
	DescribeTable("naturals", ParseAndCompare,
		Entry("Natural", `Natural`, Natural),
		Entry("NaturalLit", `1234`, NaturalLit(1234)),
		Entry("NaturalLit", `3`, NaturalLit(3)),
		Entry("NaturalPlus", `3 + 5`, NaturalPlus(NaturalLit(3), NaturalLit(5))),
		Entry("NaturalTimes", `3 * 5`, NaturalTimes(NaturalLit(3), NaturalLit(5))),
		Entry("Natural op order #1", `3 * 4 + 5`, NaturalPlus(NaturalTimes(NaturalLit(3), NaturalLit(4)), NaturalLit(5))),
		Entry("Natural op order #2", `3 + 4 * 5`, NaturalPlus(NaturalLit(3), NaturalTimes(NaturalLit(4), NaturalLit(5)))),
		// Check that if we skip whitespace, it parses
		// correctly as function application, not natural
		// addition
		Entry("Plus without whitespace", `3 +5`, Apply(NaturalLit(3), IntegerLit(5))),
	)
	DescribeTable("double-quoted text literals", ParseAndCompare,
		Entry("Empty TextLitTerm", `""`, TextLitTerm{}),
		Entry("Simple TextLitTerm", `"foo"`, TextLitTerm{Suffix: "foo"}),
		Entry(`TextLitTerm escape "`, `"\""`, TextLitTerm{Suffix: `"`}),
		Entry(`TextLitTerm escape $`, `"\$"`, TextLitTerm{Suffix: `$`}),
		Entry(`TextLitTerm escape \`, `"\\"`, TextLitTerm{Suffix: `\`}),
		Entry(`TextLitTerm escape /`, `"\/"`, TextLitTerm{Suffix: `/`}),
		Entry(`TextLitTerm escape \b`, `"\b"`, TextLitTerm{Suffix: "\b"}),
		Entry(`TextLitTerm escape \f`, `"\f"`, TextLitTerm{Suffix: "\f"}),
		Entry(`TextLitTerm escape \n`, `"\n"`, TextLitTerm{Suffix: "\n"}),
		Entry(`TextLitTerm escape \r`, `"\r"`, TextLitTerm{Suffix: "\r"}),
		Entry(`TextLitTerm escape \t`, `"\t"`, TextLitTerm{Suffix: "\t"}),
		Entry(`TextLitTerm escape \u2200`, `"\u2200"`, TextLitTerm{Suffix: "∀"}),
		Entry(`TextLitTerm escape \u03bb`, `"\u03bb"`, TextLitTerm{Suffix: "λ"}),
		Entry(`TextLitTerm escape \u03BB`, `"\u03BB"`, TextLitTerm{Suffix: "λ"}),
		Entry("Interpolated TextLitTerm", `"foo ${"bar"} baz"`,
			TextLitTerm{Chunks{Chunk{"foo ", TextLitTerm{Suffix: "bar"}}},
				" baz"},
		),
	)
	DescribeTable("single-quoted text literals", ParseAndCompare,
		Entry("Empty TextLitTerm", `''
''`, TextLitTerm{}),
		Entry("Simple TextLitTerm with no newlines", `''
foo''`, TextLitTerm{Suffix: "foo"}),
		Entry("Simple TextLitTerm with newlines", `''
foo
''`, TextLitTerm{Suffix: "foo\n"}),
		Entry("TextLitTerm with space indent", `''
  foo
  bar
  ''`, TextLitTerm{Suffix: "foo\nbar\n"}),
		Entry("TextLitTerm with tab indent", `''
		foo
		bar
		''`, TextLitTerm{Suffix: "foo\nbar\n"}),
		Entry("TextLitTerm with mixed tab/space indent", `''
	  foo
	  bar
	  ''`, TextLitTerm{Suffix: "foo\nbar\n"}),
		Entry("TextLitTerm with weird indenting", `''
	    foo
	  	bar
	  ''`, TextLitTerm{Suffix: "  foo\n\tbar\n"}),
		Entry(`Escape ''`, `''
'''
''`, TextLitTerm{Suffix: "''\n"}),
		Entry(`Escape ${`, `''
''${
''`, TextLitTerm{Suffix: "${\n"}),
		Entry("Interpolation", `''
foo ${"bar"}
baz
''`,
			TextLitTerm{Chunks{Chunk{"foo ", TextLitTerm{Suffix: "bar"}}},
				"\nbaz\n"},
		),
	)
	DescribeTable("simple expressions", ParseAndCompare,
		Entry("Identifier", `x`, Bound("x")),
		Entry("Identifier with index", `x@1`, BoundVar{"x", 1}),
		Entry("Identifier with reserved prefix", `Listicle`, Bound("Listicle")),
		Entry("Identifier with reserved prefix and index", `Listicle@3`, BoundVar{"Listicle", 3}),
	)
	DescribeTable("lists", ParseAndCompare,
		Entry("List Natural", `List Natural`, Apply(List, Natural)),
		Entry("[3]", `[3]`, MakeList(NaturalLit(3))),
		Entry("[3,4]", `[3,4]`, MakeList(NaturalLit(3), NaturalLit(4))),
		Entry("[] : List Natural", `[] : List Natural`, EmptyList{Apply(List, Natural)}),
		Entry("[3] : List Natural", `[3] : List Natural`, Annot{MakeList(NaturalLit(3)), Apply(List, Natural)}),
		Entry("a # b", `a # b`, ListAppend(Bound("a"), Bound("b"))),
	)
	DescribeTable("optionals", ParseAndCompare,
		Entry("Optional Natural", `Optional Natural`, Apply(Optional, Natural)),
		Entry("Some 3", `Some 3`, Some{NaturalLit(3)}),
		Entry("None Natural", `None Natural`, Apply(None, Natural)),
	)
	DescribeTable("records", ParseAndCompare,
		Entry("{}", `{}`, RecordType{}),
		Entry("{=}", `{=}`, RecordLit{}),
		Entry("{foo : Natural}", `{foo : Natural}`, RecordType{"foo": Natural}),
		Entry("{foo = 3}", `{foo = 3}`, RecordLit{"foo": NaturalLit(3)}),
		Entry("{foo : Natural, bar : Integer}", `{foo : Natural, bar: Integer}`, RecordType{"foo": Natural, "bar": Integer}),
		Entry("{foo = 3 , bar = +3}", `{foo = 3 , bar = +3}`, RecordLit{"foo": NaturalLit(3), "bar": IntegerLit(3)}),
		Entry("t.x", `t.x`, Field{Bound("t"), "x"}),
		Entry("t.x.y", `t.x.y`, Field{Field{Bound("t"), "x"}, "y"}),
	)
	DescribeTable("imports", ParseAndCompare,
		Entry("bash envvar text import", `env:FOO as Text`, MakeEnvVarImport("FOO", RawText)),
		Entry("posix envvar text import", `env:"FOO" as Text`, MakeEnvVarImport("FOO", RawText)),
		Entry("posix envvar text import", `env:"foo\nbar\a!" as Text`, MakeEnvVarImport("foo\nbar\a!", RawText)),
		Entry("bash envvar code import", `env:FOO`, MakeEnvVarImport("FOO", Code)),
		Entry("posix envvar code import", `env:"FOO"`, MakeEnvVarImport("FOO", Code)),
		Entry("posix envvar code import", `env:"foo\nbar\a!"`, MakeEnvVarImport("foo\nbar\a!", Code)),
		Entry("missing", `missing`, MakeImport(Missing(struct{}{}), Code)),
		Entry("local here-path import", `./local`, MakeLocalImport("local", Code)),
		Entry("local parent-path import", `../local`, MakeLocalImport("../local", Code)),
		Entry("local home import", `~/in/home`, MakeLocalImport("~/in/home", Code)),
		Entry("local absolute import", `/local`, MakeLocalImport("/local", Code)),
		Entry("simple remote", `https://example.com/foo`, MakeRemoteImport("https://example.com/foo", Code)),
		Entry("http remote", `http://example.com/foo`, MakeRemoteImport("http://example.com/foo", Code)),
		Entry("remote with query string", `https://example.com/foo?bar=baz&fred=jim`, MakeRemoteImport("https://example.com/foo?bar=baz&fred=jim", Code)),
		Entry("remote with port", `https://example.com:8080/foo`, MakeRemoteImport("https://example.com:8080/foo", Code)),
		Entry("remote with userinfo", `https://foo:bar@example.com/foo`, MakeRemoteImport("https://foo:bar@example.com/foo", Code)),
		Entry("remote with IPv4 address", `https://127.0.0.1/foo`, MakeRemoteImport("https://127.0.0.1/foo", Code)),
		Entry("remote with IPv6 address", `https://[cafe:d00d::1234]/foo`, MakeRemoteImport("https://[cafe:d00d::1234]/foo", Code)),
		// unimplemented yet. don't care too much about these features
		PEntry("remote with headers", ``, nil),
	)
	// can't test NaN using ParseAndCompare because NaN ≠ NaN
	It("handles NaN correctly", func() {
		root, err := parser.Parse("test", []byte(`NaN`))
		Expect(err).ToNot(HaveOccurred())
		f := float64(root.(DoubleLit))
		Expect(math.IsNaN(f)).To(BeTrue())
	})
	DescribeTable("lambda expressions", ParseAndCompare,
		Entry("simple λ",
			`λ(foo : bar) → baz`,
			LambdaTerm{
				"foo", Bound("bar"), Bound("baz")}),
		Entry(`simple \`,
			`\(foo : bar) → baz`,
			LambdaTerm{
				"foo", Bound("bar"), Bound("baz")}),
		Entry("with line comment",
			"λ(foo : bar) --asdf\n → baz",
			LambdaTerm{
				"foo", Bound("bar"), Bound("baz")}),
		Entry("with block comment",
			"λ(foo : bar) {-asdf\n-} → baz",
			LambdaTerm{
				"foo", Bound("bar"), Bound("baz")}),
		Entry("simple ∀",
			`∀(foo : bar) → baz`,
			PiTerm{"foo", Bound("bar"), Bound("baz")}),
		Entry("arrow type has implicit _ var",
			`foo → bar`,
			FnType(Bound("foo"), Bound("bar"))),
		Entry(`simple forall`,
			`forall(foo : bar) → baz`,
			PiTerm{"foo", Bound("bar"), Bound("baz")}),
		Entry("with line comment",
			"∀(foo : bar) --asdf\n → baz",
			PiTerm{"foo", Bound("bar"), Bound("baz")}),
	)
	DescribeTable("applications", ParseAndCompare,
		Entry("identifier application",
			`foo bar`,
			Apply(
				Bound("foo"),
				Bound("bar"),
			)),
		Entry("lambda application",
			`(λ(foo : bar) → baz) quux`,
			Apply(
				LambdaTerm{
					"foo", Bound("bar"), Bound("baz")},
				Bound("quux"))),
	)
	DescribeTable("lets", ParseAndCompare,
		Entry("simple let",
			`let x = y in z`,
			MakeLet(Bound("z"), Binding{
				Variable: "x", Value: Bound("y"),
			})),
		Entry("lambda application",
			`(λ(foo : bar) → baz) quux`,
			Apply(
				LambdaTerm{
					"foo", Bound("bar"), Bound("baz")},
				Bound("quux"))),
	)
	Describe("Expected failures", func() {
		// these keywords should fail to parse unless they're part of
		// a larger expression
		DescribeTable("keywords", ParseAndFail,
			Entry("if", `if`),
			Entry("then", `then`),
			Entry("else", `else`),
			Entry("let", `let`),
			Entry("in", `in`),
			Entry("using", `using`),
			Entry("as", `as`),
			Entry("merge", `merge`),
			Entry("Some", `Some`),
		)
		DescribeTable("bad URLs", ParseAndFail,
			Entry("bad IPv6", `https://[11111::22222]/abc`),
		)
		DescribeTable("other expected failures", ParseAndFail,
			Entry("annotation without required space", `3 :Natural`),
			Entry("unannotated list", `[]`),
		)
	})
})
