package main_test

import (
	"reflect"

	. "github.com/philandstuff/dhall-golang"
	"github.com/philandstuff/dhall-golang/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func UnmarshalAndCompare(input ast.Expr, ptr interface{}, expected interface{}) {
	Unmarshal(input, ptr)
	// use reflect to dereference a pointer of unknown type
	Expect(reflect.ValueOf(ptr).Elem().Interface()).
		To(Equal(expected))
}

var _ = Describe("Unmarshal", func() {
	DescribeTable("Simple types", UnmarshalAndCompare,
		Entry("unmarshals True into *bool",
			ast.True, new(bool), true),
		Entry("unmarshals NaturalLit into *int",
			ast.NaturalLit(5), new(int), 5),
		Entry("unmarshals NaturalLit into *int64",
			ast.NaturalLit(5), new(int64), int64(5)),
		Entry("unmarshals IntegerLit into *int",
			ast.IntegerLit(5), new(int), 5),
		Entry("unmarshals IntegerLit into *int",
			ast.IntegerLit(5), new(int64), int64(5)),
	)
	DescribeTable("Compound types", UnmarshalAndCompare,
		Entry("unmarshals List Integer into slice",
			ast.MakeList(ast.IntegerLit(5)),
			new([]int),
			[]int{5}),
		Entry("unmarshals List Bool into slice",
			ast.MakeList(ast.True, ast.False),
			new([]bool),
			[]bool{true, false}),
	)
})
