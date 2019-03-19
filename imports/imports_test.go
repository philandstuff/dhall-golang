package imports_test

import (
	"os"

	. "github.com/philandstuff/dhall-golang/ast"
	. "github.com/philandstuff/dhall-golang/imports"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func expectResolves(input, expected Expr) {
	os.Setenv("FOO", "abcd")
	actual, err := Load(input)

	Expect(err).ToNot(HaveOccurred())
	Expect(actual).To(Equal(expected))
}

var importFoo = Embed(Import{EnvVar: "FOO"})
var resolvedFoo = TextLit{Suffix: "abcd"}

var _ = Describe("Import resolution", func() {
	It("Leaves a literal expression untouched", func() {
		actual, err := Load(NaturalLit(3))

		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(NaturalLit(3)))
	})
	DescribeTable("Subexpressions to resolve", expectResolves,
		Entry("Simple import", importFoo, resolvedFoo),
		Entry("Import within lambda type",
			&LambdaExpr{Type: importFoo},
			&LambdaExpr{Type: resolvedFoo},
		),
		Entry("Import within lambda body",
			&LambdaExpr{Body: importFoo},
			&LambdaExpr{Body: resolvedFoo},
		),
		Entry("Import within pi type",
			&Pi{Type: importFoo},
			&Pi{Type: resolvedFoo},
		),
		Entry("Import within pi body",
			&Pi{Body: importFoo},
			&Pi{Body: resolvedFoo},
		),
		Entry("Import within app fn",
			&App{Fn: importFoo},
			&App{Fn: resolvedFoo},
		),
		Entry("Import within app arg",
			&App{Arg: importFoo},
			&App{Arg: resolvedFoo},
		),
		Entry("Import within let binding value",
			MakeLet(Natural, Binding{Value: importFoo}),
			MakeLet(Natural, Binding{Value: resolvedFoo}),
		),
		Entry("Import within let body",
			MakeLet(importFoo, Binding{}),
			MakeLet(resolvedFoo, Binding{}),
		),
		Entry("Import within annotated expression",
			Annot{importFoo, Text},
			Annot{resolvedFoo, Text},
		),
		Entry("Import within annotation",
			// these don't typecheck but we're just
			// checking the imports here
			Annot{Natural, importFoo},
			Annot{Natural, resolvedFoo},
		),
		Entry("Import within TextLit",
			TextLit{
				Chunks: []Chunk{
					Chunk{
						Prefix: "foo",
						Expr:   importFoo,
					}},
				Suffix: "baz",
			},
			TextLit{
				Chunks: []Chunk{
					Chunk{
						Prefix: "foo",
						Expr:   resolvedFoo,
					},
				},
				Suffix: "baz",
			},
		),
		Entry("Import within if condition",
			BoolIf{Cond: importFoo},
			BoolIf{Cond: resolvedFoo},
		),
		Entry("Import within if true branch",
			BoolIf{T: importFoo},
			BoolIf{T: resolvedFoo},
		),
		Entry("Import within if false branch",
			BoolIf{F: importFoo},
			BoolIf{F: resolvedFoo},
		),
		Entry("Import within natural plus (left side)",
			NaturalPlus{L: importFoo},
			NaturalPlus{L: resolvedFoo},
		),
		Entry("Import within natural plus (right side)",
			NaturalPlus{R: importFoo},
			NaturalPlus{R: resolvedFoo},
		),
		Entry("Import within natural times (left side)",
			NaturalTimes{L: importFoo},
			NaturalTimes{L: resolvedFoo},
		),
		Entry("Import within natural times (right side)",
			NaturalTimes{R: importFoo},
			NaturalTimes{R: resolvedFoo},
		),
		Entry("Import within empty list type",
			EmptyList{Type: importFoo},
			EmptyList{Type: resolvedFoo},
		),
		Entry("Import within list",
			MakeList(importFoo),
			MakeList(resolvedFoo),
		),
		Entry("Import within record type",
			Record(map[string]Expr{"foo": importFoo}),
			Record(map[string]Expr{"foo": resolvedFoo}),
		),
		Entry("Import within record literal",
			RecordLit(map[string]Expr{"foo": importFoo}),
			RecordLit(map[string]Expr{"foo": resolvedFoo}),
		),
		Entry("Import within field extract",
			Field{Record: importFoo},
			Field{Record: resolvedFoo},
		),
	)
})
