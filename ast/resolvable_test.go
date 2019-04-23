package ast_test

import (
	. "github.com/philandstuff/dhall-golang/ast"

	// . "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("ChainOnto", func(resolvable, base, expected Resolvable) {
	actual, err := resolvable.ChainOnto(base)
	if expected == nil {
		Expect(err).To(HaveOccurred())
	} else {
		Expect(actual).To(Equal(expected))
		Expect(err).ToNot(HaveOccurred())
	}
},
	Entry("Missing onto EnvVar", Missing{}, EnvVar(""), Missing{}),
	Entry("Missing onto Local", Missing{}, Local(""), Missing{}),
	Entry("Missing onto Remote", Missing{}, Remote(""), Missing{}),
	Entry("Missing onto Missing", Missing{}, Missing{}, Missing{}),
	Entry("EnvVar onto EnvVar", EnvVar("foo"), EnvVar("bar"), EnvVar("foo")),
	Entry("EnvVar onto Local", EnvVar("foo"), Local(""), EnvVar("foo")),
	Entry("EnvVar onto Remote", EnvVar("foo"), Remote(""), nil),
	Entry("EnvVar onto Missing", EnvVar("foo"), Missing{}, EnvVar("foo")),
	Entry("Relative local onto EnvVar", Local("foo"), EnvVar("bar"), Local("foo")),
	Entry("Relative local onto Local", Local("foo"), Local("/bar/baz"), Local("/bar/foo")),
	Entry("Relative local onto Remote", Local("foo"), Remote("https://example.com/bar/baz"), Remote("https://example.com/bar/foo")),
	Entry("Parent-relative local onto Remote", Local("../foo"), Remote("https://example.com/bar/baz/quux"), Remote("https://example.com/bar/foo")),
	Entry("Relative local onto Missing", Local("foo"), Missing{}, Local("foo")),
	Entry("Absolute local onto EnvVar", Local("/foo"), EnvVar("bar"), Local("/foo")),
	Entry("Absolute local onto Local", Local("/foo"), Local("/bar/baz"), Local("/foo")),
	Entry("Absolute local onto Remote", Local("/foo"), Remote("https://example.com/bar/baz"), nil),
	Entry("Absolute local onto Missing", Local("/foo"), Missing{}, Local("/foo")),
)
