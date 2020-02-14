package term_test

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/philandstuff/dhall-golang/internal"
	. "github.com/philandstuff/dhall-golang/term"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func makeRemoteFile(u string) RemoteFile {
	parsed, _ := url.ParseRequestURI(u)
	remote := NewRemoteFile(parsed)
	return remote
}

var _ = DescribeTable("ChainOnto", func(fetchable, base, expected Fetchable) {
	actual, err := fetchable.ChainOnto(base)
	if expected == nil {
		Expect(actual).To(BeNil())
		Expect(err).To(HaveOccurred())
	} else {
		Expect(actual).To(Equal(expected))
		Expect(err).ToNot(HaveOccurred())
	}
},
	Entry("Missing onto EnvVar", Missing{}, EnvVar(""), Missing{}),
	Entry("Missing onto LocalFile", Missing{}, LocalFile(""), Missing{}),
	Entry("Missing onto RemoteFile", Missing{}, RemoteFile{}, Missing{}),
	Entry("Missing onto Missing", Missing{}, Missing{}, Missing{}),
	Entry("EnvVar onto EnvVar", EnvVar("foo"), EnvVar("bar"), EnvVar("foo")),
	Entry("EnvVar onto LocalFile", EnvVar("foo"), LocalFile(""), EnvVar("foo")),
	Entry("EnvVar onto RemoteFile", EnvVar("foo"), RemoteFile{}, EnvVar("foo")),
	Entry("EnvVar onto Missing", EnvVar("foo"), Missing{}, EnvVar("foo")),
	Entry("Relative local onto EnvVar", LocalFile("foo"), EnvVar("bar"), LocalFile("foo")),
	Entry("Relative local onto LocalFile", LocalFile("foo"), LocalFile("/bar/baz"), LocalFile("/bar/foo")),
	Entry("Relative local onto RemoteFile", LocalFile("foo"), makeRemoteFile("https://example.com/bar/baz"), makeRemoteFile("https://example.com/bar/foo")),
	Entry("Relative local with tricky chars onto RemoteFile", LocalFile("foo:bar#[☃"), makeRemoteFile("https://example.com/bar/baz"), makeRemoteFile("https://example.com/bar/foo:bar%23%5B%E2%98%83")),
	Entry("Relative local onto Missing", LocalFile("foo"), Missing{}, LocalFile("foo")),
	Entry("Parent-relative local onto EnvVar", LocalFile("../foo"), EnvVar("bar"), LocalFile("../foo")),
	Entry("Parent-relative local onto LocalFile", LocalFile("../foo"), LocalFile("/bar/baz/quux"), LocalFile("/bar/foo")),
	Entry("Parent-relative local onto RemoteFile", LocalFile("../foo"), makeRemoteFile("https://example.com/bar/baz/quux"), makeRemoteFile("https://example.com/bar/foo")),
	Entry("Parent-relative local with tricky chars onto RemoteFile", LocalFile("../foo#[☃"), makeRemoteFile("https://example.com/bar/baz/quux"), makeRemoteFile("https://example.com/bar/foo%23%5B%E2%98%83")),
	Entry("Parent-relative local onto Missing", LocalFile("../foo"), Missing{}, LocalFile("../foo")),
	Entry("Home-relative local onto EnvVar", LocalFile("~/foo"), EnvVar("bar"), LocalFile("~/foo")),
	Entry("Home-relative local onto LocalFile", LocalFile("~/foo"), LocalFile("/bar/baz"), LocalFile("~/foo")),
	Entry("Home-relative local onto RemoteFile", LocalFile("~/foo"), makeRemoteFile("https://example.com/bar/baz"), nil),
	Entry("Home-relative local onto Missing", LocalFile("~/foo"), Missing{}, LocalFile("~/foo")),
	Entry("Absolute local onto EnvVar", LocalFile("/foo"), EnvVar("bar"), LocalFile("/foo")),
	Entry("Absolute local onto LocalFile", LocalFile("/foo"), LocalFile("/bar/baz"), LocalFile("/foo")),
	Entry("Absolute local onto RemoteFile", LocalFile("/foo"), makeRemoteFile("https://example.com/bar/baz"), nil),
	Entry("Absolute local onto Missing", LocalFile("/foo"), Missing{}, LocalFile("/foo")),
	Entry("RemoteFile onto EnvVar", makeRemoteFile("https://example.com/foo"), EnvVar("bar"), makeRemoteFile("https://example.com/foo")),
	Entry("RemoteFile onto LocalFile", makeRemoteFile("https://example.com/foo"), LocalFile(""), makeRemoteFile("https://example.com/foo")),
	Entry("RemoteFile onto RemoteFile", makeRemoteFile("https://example.com/foo"), RemoteFile{}, makeRemoteFile("https://example.com/foo")),
	Entry("RemoteFile onto Missing", makeRemoteFile("https://example.com/foo"), Missing{}, makeRemoteFile("https://example.com/foo")),
)

const ExampleRemoteOrigin = "http://example.com"

var _ = Describe("Fetch", func() {
	DescribeTable("LocalFile fetching", func(fetchable Fetchable, origin string, expected string) {
		os.Setenv("foo", "Value of envvar foo")
		actual, err := fetchable.Fetch(origin)
		if expected == "" {
			Expect(err).To(HaveOccurred())
		} else {
			Expect(actual).To(Equal(expected))
			Expect(err).ToNot(HaveOccurred())
		}
	},
		Entry("Missing from local returns error", Missing{}, NullOrigin, ""),
		Entry("Missing from remote returns error", Missing{}, ExampleRemoteOrigin, ""),
		Entry("EnvVar from local is allowed", EnvVar("foo"), NullOrigin, "Value of envvar foo"),
		Entry("EnvVar from remote returns error", EnvVar("foo"), ExampleRemoteOrigin, ""),
		Entry("Local from local is allowed", LocalFile("./testdata/foo"), NullOrigin, "Content of file 'foo'\n"),
		Entry("Local from remote returns error", LocalFile("./testdata/foo"), ExampleRemoteOrigin, ""),
	)
	Describe("RemoteFile fetching", func() {
		var server *ghttp.Server
		AfterEach(func() {
			server.Close()
		})
		BeforeEach(func() {
			server = ghttp.NewServer()
			server.RouteToHandler("GET", "/no-cors.dhall",
				func(w http.ResponseWriter, r *http.Request) {
					// no Access-Control-Allow-Origin header
					io.WriteString(w, "this content only allows the same origin")
				},
			)
			server.RouteToHandler("GET", "/cors-ok-with-star.dhall",
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					io.WriteString(w, "this content allows origin *")
				},
			)
			server.RouteToHandler("GET", "/cors-ok-with-origin.dhall",
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
					io.WriteString(w, "this content allows origin "+r.Header.Get("Origin"))
				},
			)
		})
		It("is allowed from local", func() {
			actual, err := internal.NewRemoteImport(server.URL()+"/no-cors.dhall", Code).Fetch(NullOrigin)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal("this content only allows the same origin"))
		})
		It("is allowed from same origin, even if CORS fails", func() {
			actual, err := internal.NewRemoteImport(server.URL()+"/no-cors.dhall", Code).Fetch(server.URL())

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal("this content only allows the same origin"))
		})
		Context("when fetching from different origin", func() {
			It("returns error if CORS fails", func() {
				_, err := internal.NewRemoteImport(server.URL()+"/no-cors.dhall", Code).Fetch("http://example.com")

				Expect(err).To(HaveOccurred())
			})
			It("is allowed if Access-Control-Allow-Origin is '*'", func() {
				actual, err := internal.NewRemoteImport(server.URL()+"/cors-ok-with-star.dhall", Code).Fetch("http://example.com")

				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal("this content allows origin *"))
			})
			It("is allowed if Access-Control-Allow-Origin matches the Origin header", func() {
				actual, err := internal.NewRemoteImport(server.URL()+"/cors-ok-with-origin.dhall", Code).Fetch("http://example.com")

				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal("this content allows origin http://example.com"))
			})
		})
	})
})
