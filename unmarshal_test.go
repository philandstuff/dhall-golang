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
		Entry("unmarshals DoubleLit into *float32",
			ast.DoubleLit(3.5), new(float32), float32(3.5)),
		Entry("unmarshals DoubleLit into *float64",
			ast.DoubleLit(3.5), new(float64), float64(3.5)),
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
		Entry("unmarshals List Integer into int slice",
			ast.MakeList(ast.IntegerLit(5)),
			new([]int),
			[]int{5}),
		Entry("unmarshals List Integer into int64 slice",
			ast.MakeList(ast.IntegerLit(5)),
			new([]int64),
			[]int64{5}),
		Entry("unmarshals List Bool into slice",
			ast.MakeList(ast.True, ast.False),
			new([]bool),
			[]bool{true, false}),
		Entry("unmarshals List (List Bool) into slice",
			ast.MakeList(ast.MakeList(ast.True, ast.False)),
			new([][]bool),
			[][]bool{{true, false}}),
	)
	Describe("Function types", func() {
		It("Unmarshals the identity int function", func() {
			var fn func(int) int
			dhallFn := &ast.LambdaExpr{
				Label: "x",
				Type:  ast.Natural,
				Body:  ast.MkVar("x"),
			}
			Unmarshal(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(3)).To(Equal(3))
		})
		It("Unmarshals the identity int64 function", func() {
			var fn func(int64) int64
			dhallFn := &ast.LambdaExpr{
				Label: "x",
				Type:  ast.Natural,
				Body:  ast.MkVar("x"),
			}
			Unmarshal(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(int64(3))).To(Equal(int64(3)))
		})
		It("Unmarshals the int successor function", func() {
			var fn func(int) int
			dhallFn := &ast.LambdaExpr{
				Label: "x",
				Type:  ast.Natural,
				Body: ast.NaturalPlus(
					ast.MkVar("x"),
					ast.NaturalLit(1),
				),
			}
			Unmarshal(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(3)).To(Equal(4))
		})
		It("Unmarshals the natural sum function", func() {
			var fn func(int, int) int
			dhallFn := &ast.LambdaExpr{
				Label: "x",
				Type:  ast.Natural,
				Body: &ast.LambdaExpr{
					Label: "y",
					Type:  ast.Natural,
					Body: ast.NaturalPlus(
						ast.MkVar("x"),
						ast.MkVar("y"),
					),
				},
			}
			Unmarshal(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(3, 4)).To(Equal(7))
		})
	})
})
