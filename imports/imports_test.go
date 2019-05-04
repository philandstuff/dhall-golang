package imports_test

import (
	"io"
	"net/http"
	"os"

	. "github.com/philandstuff/dhall-golang/ast"
	. "github.com/philandstuff/dhall-golang/imports"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func expectResolves(input, expected Expr) {
	os.Setenv("FOO", "abcd")
	actual, err := Load(input)

	Expect(err).ToNot(HaveOccurred())
	Expect(actual).To(Equal(expected))
}

var importFooAsText = Embed(MakeEnvVarImport("FOO", RawText))
var resolvedFooAsText = TextLit{Suffix: "abcd"}

var _ = Describe("Import resolution", func() {
	Describe("Environment varibles", func() {
		It("Resolves as Text", func() {
			os.Setenv("FOO", "abcd")
			actual, err := Load(importFooAsText)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(resolvedFooAsText))
		})
		It("Resolves as code", func() {
			os.Setenv("FOO", "3 : Natural")
			actual, err := Load(Embed(MakeEnvVarImport("FOO", Code)))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(Annot{Expr: NaturalLit(3), Annotation: Natural}))
		})
		It("Fails to resolve code with free variables", func() {
			os.Setenv("FOO", "x")
			_, err := Load(Embed(MakeEnvVarImport("FOO", Code)))

			Expect(err).To(HaveOccurred())
		})
		It("Performs import chaining", func() {
			os.Setenv("CHAIN1", "env:CHAIN2")
			os.Setenv("CHAIN2", "2 + 2")
			actual, err := Load(Embed(MakeEnvVarImport("CHAIN1", Code)))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(NaturalPlus(NaturalLit(2), NaturalLit(2))))
		})
		It("Rejects import cycles", func() {
			os.Setenv("CYCLE", "env:CYCLE")
			_, err := Load(Embed(MakeEnvVarImport("CYCLE", Code)))

			Expect(err).To(HaveOccurred())
		})
	})
	Describe("http imports", func() {
		var server *ghttp.Server
		BeforeEach(func() {
			server = ghttp.NewServer()
		})
		AfterEach(func() {
			server.Close()
		})
		It("Resolves as Text", func() {
			server.RouteToHandler("GET", "/foo.dhall",
				ghttp.RespondWith(http.StatusOK, "abcd"),
			)
			actual, err := Load(Embed(MakeRemoteImport(server.URL()+"/foo.dhall", RawText)))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(TextLit{Suffix: "abcd"}))
		})
		It("Resolves as code", func() {
			server.RouteToHandler("GET", "/foo.dhall",
				ghttp.RespondWith(http.StatusOK, "3 : Natural"),
			)
			actual, err := Load(Embed(MakeRemoteImport(server.URL()+"/foo.dhall", Code)))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(Annot{Expr: NaturalLit(3), Annotation: Natural}))
		})
		It("Fails to resolve code with free variables", func() {
			server.RouteToHandler("GET", "/foo.dhall",
				ghttp.RespondWith(http.StatusOK, "x"),
			)
			_, err := Load(Embed(MakeRemoteImport(server.URL()+"/foo.dhall", Code)))

			Expect(err).To(HaveOccurred())
		})
		Describe("CORS checks", func() {
			BeforeEach(func() {
				server.RouteToHandler("GET", "/no-cors.dhall",
					func(w http.ResponseWriter, r *http.Request) {
						// no Access-Control-Allow-Origin header
						io.WriteString(w, "3 : Natural")
					},
				)
				server.RouteToHandler("GET", "/cors-ok-with-star.dhall",
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Access-Control-Allow-Origin", "*")
						io.WriteString(w, "3 : Natural")
					},
				)
				server.RouteToHandler("GET", "/cors-ok-with-origin.dhall",
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
						io.WriteString(w, "3 : Natural")
					},
				)
			})
			Context("when remote import fetches same origin", func() {
				It("allows the request", func() {
					server.RouteToHandler("GET", "/same-origin.dhall",
						ghttp.RespondWith(http.StatusOK, "./no-cors.dhall"),
					)

					actual, err := Load(Embed(MakeRemoteImport(server.URL()+"/same-origin.dhall", Code)))

					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(Annot{Expr: NaturalLit(3), Annotation: Natural}))
				})
			})
			Context("when remote import fetches different origin", func() {
				It("refuses if CORS fails", func() {
					otherOrigin := ghttp.NewServer()
					otherOrigin.RouteToHandler("GET", "/other-origin.dhall",
						ghttp.RespondWith(http.StatusOK, server.URL()+"/no-cors.dhall"),
					)

					_, err := Load(Embed(MakeRemoteImport(otherOrigin.URL()+"/other-origin.dhall", Code)))
					Expect(err).To(HaveOccurred())
				})
				It("allows if Access-Control-Allow-Origin is '*'", func() {
					otherOrigin := ghttp.NewServer()
					otherOrigin.RouteToHandler("GET", "/other-origin.dhall",
						ghttp.RespondWith(http.StatusOK, server.URL()+"/cors-ok-with-star.dhall"),
					)

					actual, err := Load(Embed(MakeRemoteImport(otherOrigin.URL()+"/other-origin.dhall", Code)))

					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(Annot{Expr: NaturalLit(3), Annotation: Natural}))
				})
				It("allows if Access-Control-Allow-Origin matches the Origin header", func() {
					otherOrigin := ghttp.NewServer()
					otherOrigin.RouteToHandler("GET", "/other-origin.dhall",
						ghttp.RespondWith(http.StatusOK, server.URL()+"/cors-ok-with-origin.dhall"),
					)

					actual, err := Load(Embed(MakeRemoteImport(otherOrigin.URL()+"/other-origin.dhall", Code)))

					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(Annot{Expr: NaturalLit(3), Annotation: Natural}))
				})
			})
			Context("when local import fetches remote", func() {
				It("allows the request", func() {
					actual, err := Load(Embed(MakeRemoteImport(server.URL()+"/no-cors.dhall", Code)))

					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(Annot{Expr: NaturalLit(3), Annotation: Natural}))
				})
			})
		})
	})
	Describe("local imports", func() {
		It("Resolves as Text", func() {
			actual, err := Load(Embed(MakeLocalImport("./testdata/just_text.txt", RawText)))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(TextLit{Suffix: "here is some text\n"}))
		})
		It("Resolves as code", func() {
			actual, err := Load(Embed(MakeLocalImport("./testdata/natural.dhall", Code)))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(Annot{Expr: NaturalLit(3), Annotation: Natural}))
		})
		It("Fails to resolve code with free variables", func() {
			_, err := Load(Embed(MakeLocalImport("./testdata/free_variable.dhall", Code)))

			Expect(err).To(HaveOccurred())
		})
		It("Performs import chaining", func() {
			actual, err := Load(Embed(MakeLocalImport("./testdata/chain1.dhall", Code)))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(NaturalPlus(NaturalLit(2), NaturalLit(2))))
		})
		It("Rejects import cycles", func() {
			_, err := Load(Embed(MakeLocalImport("./testdata/cycle1.dhall", Code)))

			Expect(err).To(HaveOccurred())
		})
	})
	DescribeTable("Other subexpressions", expectResolves,
		Entry("Literal expression", NaturalLit(3), NaturalLit(3)),
		Entry("Simple import", importFooAsText, resolvedFooAsText),
		Entry("Import within lambda type",
			&LambdaExpr{Type: importFooAsText},
			&LambdaExpr{Type: resolvedFooAsText},
		),
		Entry("Import within lambda body",
			&LambdaExpr{Body: importFooAsText},
			&LambdaExpr{Body: resolvedFooAsText},
		),
		Entry("Import within pi type",
			&Pi{Type: importFooAsText},
			&Pi{Type: resolvedFooAsText},
		),
		Entry("Import within pi body",
			&Pi{Body: importFooAsText},
			&Pi{Body: resolvedFooAsText},
		),
		Entry("Import within app fn",
			&App{Fn: importFooAsText},
			&App{Fn: resolvedFooAsText},
		),
		Entry("Import within app arg",
			&App{Arg: importFooAsText},
			&App{Arg: resolvedFooAsText},
		),
		Entry("Import within let binding value",
			MakeLet(Natural, Binding{Value: importFooAsText}),
			MakeLet(Natural, Binding{Value: resolvedFooAsText}),
		),
		Entry("Import within let body",
			MakeLet(importFooAsText, Binding{}),
			MakeLet(resolvedFooAsText, Binding{}),
		),
		Entry("Import within annotated expression",
			Annot{importFooAsText, Text},
			Annot{resolvedFooAsText, Text},
		),
		Entry("Import within annotation",
			// these don't typecheck but we're just
			// checking the imports here
			Annot{Natural, importFooAsText},
			Annot{Natural, resolvedFooAsText},
		),
		Entry("Import within TextLit",
			TextLit{
				Chunks: []Chunk{
					Chunk{
						Prefix: "foo",
						Expr:   importFooAsText,
					}},
				Suffix: "baz",
			},
			TextLit{
				Chunks: []Chunk{
					Chunk{
						Prefix: "foo",
						Expr:   resolvedFooAsText,
					},
				},
				Suffix: "baz",
			},
		),
		Entry("Import within if condition",
			BoolIf{Cond: importFooAsText},
			BoolIf{Cond: resolvedFooAsText},
		),
		Entry("Import within if true branch",
			BoolIf{T: importFooAsText},
			BoolIf{T: resolvedFooAsText},
		),
		Entry("Import within if false branch",
			BoolIf{F: importFooAsText},
			BoolIf{F: resolvedFooAsText},
		),
		Entry("Import within Operator (left side)",
			Operator{L: importFooAsText},
			Operator{L: resolvedFooAsText},
		),
		Entry("Import within natural plus (right side)",
			Operator{R: importFooAsText},
			Operator{R: resolvedFooAsText},
		),
		Entry("Import within empty list type",
			EmptyList{Type: importFooAsText},
			EmptyList{Type: resolvedFooAsText},
		),
		Entry("Import within list",
			MakeList(importFooAsText),
			MakeList(resolvedFooAsText),
		),
		Entry("Import within record type",
			Record(map[string]Expr{"foo": importFooAsText}),
			Record(map[string]Expr{"foo": resolvedFooAsText}),
		),
		Entry("Import within record literal",
			RecordLit(map[string]Expr{"foo": importFooAsText}),
			RecordLit(map[string]Expr{"foo": resolvedFooAsText}),
		),
		Entry("Import within field extract",
			Field{Record: importFooAsText},
			Field{Record: resolvedFooAsText},
		),
	)
})
