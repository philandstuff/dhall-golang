package main_test

import (
	. "github.com/philandstuff/dhall-golang"
	"github.com/philandstuff/dhall-golang/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func UnmarshalAndCompare(input ast.Expr, expected interface{}) {
	var actual interface{}
	Unmarshal(input, actual)
	Expect(actual).To(Equal(expected))
}

var _ = Describe("Unmarshal", func() {
	Describe("Simple types", func() {
		It("unmarshals NaturalLit into *int", func() {
			var i int
			Unmarshal(ast.NaturalLit(5), &i)
			Expect(i).To(Equal(5))
		})
		It("unmarshals IntegerLit into *int", func() {
			var i int
			Unmarshal(ast.IntegerLit(5), &i)
			Expect(i).To(Equal(5))
		})
	})
	Describe("Compound types", func() {
		It("unmarshals List Integer into slice", func() {
			var i []int
			Unmarshal(ast.MakeList(ast.IntegerLit(5)), &i)
			Expect(i).To(Equal([]int{5}))
		})
		It("unmarshals List Bool into slice", func() {
			var i []bool
			Unmarshal(ast.MakeList(ast.True, ast.False), &i)
			Expect(i).To(Equal([]bool{true, false}))
		})
	})
})
