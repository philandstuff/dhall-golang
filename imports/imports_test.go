package imports_test

import (
	"io"
	"net/http"
	"os"

	. "github.com/philandstuff/dhall-golang/v4/imports"
	. "github.com/philandstuff/dhall-golang/v4/internal"
	. "github.com/philandstuff/dhall-golang/v4/term"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func expectResolves(input, expected Term) {
	os.Setenv("FOO", "abcd")
	actual, err := Load(input)

	Expect(err).ToNot(HaveOccurred())
	Expect(actual).To(Equal(expected))
}

var importFooAsText = NewEnvVarImport("FOO", RawText)
var resolvedFooAsText = PlainText("abcd")

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
			actual, err := Load(NewEnvVarImport("FOO", Code))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(NaturalLit(3)))
		})
		It("Fails to resolve code with free variables", func() {
			os.Setenv("FOO", "x")
			_, err := Load(NewEnvVarImport("FOO", Code))

			Expect(err).To(HaveOccurred())
		})
		It("Performs import chaining", func() {
			os.Setenv("CHAIN1", "env:CHAIN2")
			os.Setenv("CHAIN2", "2 + 2")
			actual, err := Load(NewEnvVarImport("CHAIN1", Code))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(NaturalLit(4)))
		})
		It("Rejects import cycles", func() {
			result := make(chan error)
			go func() {
				os.Setenv("CYCLE", "env:CYCLE")
				_, err := Load(NewEnvVarImport("CYCLE", Code))
				if err != nil {
					result <- err
				}
			}()
			Eventually(result).Should(Receive())
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
			actual, err := Load(NewRemoteImport(server.URL()+"/foo.dhall", RawText))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(PlainText("abcd")))
		})
		It("Resolves as code", func() {
			server.RouteToHandler("GET", "/foo.dhall",
				ghttp.RespondWith(http.StatusOK, "3 : Natural"),
			)
			actual, err := Load(NewRemoteImport(server.URL()+"/foo.dhall", Code))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(NaturalLit(3)))
		})
		It("Fails to resolve code with free variables", func() {
			server.RouteToHandler("GET", "/foo.dhall",
				ghttp.RespondWith(http.StatusOK, "x"),
			)
			_, err := Load(NewRemoteImport(server.URL()+"/foo.dhall", Code))

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

					actual, err := Load(NewRemoteImport(server.URL()+"/same-origin.dhall", Code))

					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(NaturalLit(3)))
				})
			})
			Context("when remote import fetches different origin", func() {
				It("refuses if CORS fails", func() {
					otherOrigin := ghttp.NewServer()
					otherOrigin.RouteToHandler("GET", "/other-origin.dhall",
						ghttp.RespondWith(http.StatusOK, server.URL()+"/no-cors.dhall"),
					)

					_, err := Load(NewRemoteImport(otherOrigin.URL()+"/other-origin.dhall", Code))
					Expect(err).To(HaveOccurred())
				})
				It("allows if Access-Control-Allow-Origin is '*'", func() {
					otherOrigin := ghttp.NewServer()
					otherOrigin.RouteToHandler("GET", "/other-origin.dhall",
						ghttp.RespondWith(http.StatusOK, server.URL()+"/cors-ok-with-star.dhall"),
					)

					actual, err := Load(NewRemoteImport(otherOrigin.URL()+"/other-origin.dhall", Code))

					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(NaturalLit(3)))
				})
				It("allows if Access-Control-Allow-Origin matches the Origin header", func() {
					otherOrigin := ghttp.NewServer()
					otherOrigin.RouteToHandler("GET", "/other-origin.dhall",
						ghttp.RespondWith(http.StatusOK, server.URL()+"/cors-ok-with-origin.dhall"),
					)

					actual, err := Load(NewRemoteImport(otherOrigin.URL()+"/other-origin.dhall", Code))

					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(NaturalLit(3)))
				})
			})
			Context("when local import fetches remote", func() {
				It("allows the request", func() {
					actual, err := Load(NewRemoteImport(server.URL()+"/no-cors.dhall", Code))

					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(NaturalLit(3)))
				})
			})
		})
	})
	Describe("local imports", func() {
		It("Resolves as Text", func() {
			actual, err := Load(NewLocalImport("./testdata/just_text.txt", RawText))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(PlainText("here is some text\n")))
		})
		It("Resolves as code", func() {
			actual, err := Load(NewLocalImport("./testdata/natural.dhall", Code))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(NaturalLit(3)))
		})
		It("Fails to resolve code with free variables", func() {
			_, err := Load(NewLocalImport("./testdata/free_variable.dhall", Code))

			Expect(err).To(HaveOccurred())
		})
		It("Performs import chaining", func() {
			actual, err := Load(NewLocalImport("./testdata/chain1.dhall", Code))

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(NaturalLit(4)))
		})
		It("Rejects import cycles", func() {
			result := make(chan error)
			go func() {
				_, err := Load(NewLocalImport("./testdata/cycle1.dhall", Code))
				if err != nil {
					result <- err
				}
			}()
			Eventually(result).Should(Receive())
		})
	})
})
