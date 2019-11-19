package examples_test

import (
	. "github.com/philandstuff/dhall-golang/examples"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DecodeMessage", func() {
	It("Returns a sensible message", func() {
		Expect(DecodeMessage()).To(Equal(&Message{
			Name: "Alice",
			Body: "Hello",
			Time: 1294706395881547000,
		}))
	})
})
