package dhall_test

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

type testStruct struct {
	Foo int
	Bar string
}

var _ = Describe("Unmarshal", func() {
	DescribeTable("Simple types", UnmarshalAndCompare,
		Entry("unmarshals DoubleLit into float32",
			ast.DoubleLit(3.5), new(float32), float32(3.5)),
		Entry("unmarshals DoubleLit into float64",
			ast.DoubleLit(3.5), new(float64), float64(3.5)),
		Entry("unmarshals True into bool",
			ast.True, new(bool), true),
		Entry("unmarshals NaturalLit into int",
			ast.NaturalLit(5), new(int), 5),
		Entry("unmarshals NaturalLit into int64",
			ast.NaturalLit(5), new(int64), int64(5)),
		Entry("unmarshals IntegerLit into int",
			ast.IntegerLit(5), new(int), 5),
		Entry("unmarshals IntegerLit into int",
			ast.IntegerLit(5), new(int64), int64(5)),
		Entry("unmarshals TextLit into string",
			ast.TextLit{Suffix: "lalala"}, new(string), "lalala"),
	)
	DescribeTable("Simple types into interface{}", UnmarshalAndCompare,
		Entry("unmarshals DoubleLit into interface{}",
			ast.DoubleLit(3.5), new(interface{}), float64(3.5)),
		Entry("unmarshals True into interface{}",
			ast.True, new(interface{}), true),
		Entry("unmarshals NaturalLit into interface{}",
			ast.NaturalLit(5), new(interface{}), uint(5)),
		Entry("unmarshals IntegerLit into interface{}",
			ast.IntegerLit(5), new(interface{}), int(5)),
		Entry("unmarshals TextLit into interface{}",
			ast.TextLit{Suffix: "lalala"}, new(interface{}), "lalala"),
	)
	DescribeTable("Compound types", UnmarshalAndCompare,
		Entry("unmarshals Some 5 into int",
			ast.Some{ast.NaturalLit(5)},
			new(int),
			5),
		Entry("unmarshals None Natural into int",
			ast.Apply(ast.None, ast.Natural),
			new(int),
			0),
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
		Entry("unmarshals empty List Bool into slice",
			ast.EmptyList{ast.Bool},
			new([]bool),
			[]bool{}),
		Entry("unmarshals None (List Bool) into slice",
			ast.Apply(ast.None, ast.EmptyList{ast.Bool}),
			new([]bool),
			[]bool(nil)),
		Entry("unmarshals List (List Bool) into slice",
			ast.MakeList(ast.MakeList(ast.True, ast.False)),
			new([][]bool),
			[][]bool{{true, false}}),
		Entry("unmarshals {Foo : Natural, Bar : Text} into struct",
			ast.RecordLit{"Foo": ast.NaturalLit(3), "Bar": ast.TextLit{Suffix: "xyzzy"}},
			new(testStruct),
			testStruct{Foo: 3, Bar: "xyzzy"}),
		Entry("unmarshals None {Foo : Natural, Bar : Text} into struct",
			ast.Apply(ast.None, ast.Record{"Foo": ast.Natural, "Bar": ast.Text}),
			new(testStruct),
			testStruct{}),
		Entry("unmarshals List {mapKey : Natural, mapValue : Text} into map",
			ast.MakeList(ast.RecordLit{"mapKey": ast.NaturalLit(3), "mapValue": ast.TextLit{Suffix: "fizz"}},
				ast.RecordLit{"mapKey": ast.NaturalLit(5), "mapValue": ast.TextLit{Suffix: "buzz"}}),
			new(map[int]string),
			map[int]string{3: "fizz", 5: "buzz"}),
		Entry("unmarshals None List {mapKey : Natural, mapValue : Text} into map",
			ast.Apply(ast.None, ast.Apply(ast.List, ast.Record{"mapKey": ast.Natural, "mapValue": ast.Text})),
			new(map[int]string),
			map[int]string(nil)),
		Entry("unmarshals empty List {mapKey : Natural, mapValue : Text} into map",
			ast.EmptyList{ast.Record{"mapKey": ast.Natural, "mapValue": ast.Text}},
			new(map[int]string),
			map[int]string{}),
	)
	// EmptyList into interface{} deserves its own tests because we have to
	// convert the list type to a reflect.Type
	DescribeTable("Empty list into interface{}", UnmarshalAndCompare,
		Entry("unmarshals empty List Bool into interface{}",
			ast.EmptyList{ast.Bool},
			new(interface{}),
			[]bool{}),
		Entry("unmarshals empty List Natural into interface{}",
			ast.EmptyList{ast.Natural},
			new(interface{}),
			[]uint{}),
		Entry("unmarshals empty List Integer into interface{}",
			ast.EmptyList{ast.Integer},
			new(interface{}),
			[]int{}),
		Entry("unmarshals empty List Double into interface{}",
			ast.EmptyList{ast.Double},
			new(interface{}),
			[]float64{}),
		Entry("unmarshals empty List Text into interface{}",
			ast.EmptyList{ast.Text},
			new(interface{}),
			[]string{}),
		Entry("unmarshals empty List {Foo : Integer, Bar : Text} into interface{}",
			ast.EmptyList{ast.Record{"Foo": ast.Integer, "Bar": ast.Text}},
			new(interface{}),
			[]struct {
				Bar string
				Foo int
			}{}),
		Entry("unmarshals empty List {foo : Integer, bar : Text} into interface{}",
			ast.EmptyList{ast.Record{"foo": ast.Integer, "bar": ast.Text}},
			new(interface{}),
			[]struct {
				Bar string
				Foo int
			}{}),
		Entry("unmarshals empty List {mapKey : Integer, mapValue : Text} into interface{}",
			ast.EmptyList{ast.Record{"mapKey": ast.Integer, "mapValue": ast.Text}},
			new(interface{}),
			map[int]string{}))
	DescribeTable("Compound types into interface{}", UnmarshalAndCompare,
		XEntry("unmarshals List Integer into interface{}",
			ast.MakeList(ast.IntegerLit(5)),
			new(interface{}),
			[]int{5}),
		XEntry("unmarshals {Foo : Integer, Bar : Text} into interface{}",
			ast.RecordLit{"Foo": ast.IntegerLit(3), "Bar": ast.TextLit{Suffix: "xyzzy"}},
			new(interface{}),
			testStruct{Foo: 3, Bar: "xyzzy"}),
		XEntry("unmarshals None {Foo : Integer, Bar : Text} into struct",
			ast.Apply(ast.None, ast.Record{"Foo": ast.Integer, "Bar": ast.Text}),
			new(interface{}),
			testStruct{}),
		XEntry("unmarshals List {mapKey : Natural, mapValue : Text}",
			ast.MakeList(ast.RecordLit{"mapKey": ast.NaturalLit(3), "mapValue": ast.TextLit{Suffix: "fizz"}},
				ast.RecordLit{"mapKey": ast.NaturalLit(5), "mapValue": ast.TextLit{Suffix: "buzz"}}),
			new(interface{}),
			map[uint]string{3: "fizz", 5: "buzz"}),
		XEntry("unmarshals None List {mapKey : Natural, mapValue : Text}",
			ast.Apply(ast.None, ast.Apply(ast.List, ast.Record{"mapKey": ast.Natural, "mapValue": ast.Text})),
			new(interface{}),
			map[uint]string(nil)),
		Entry("unmarshals empty List {mapKey : Natural, mapValue : Text}",
			ast.EmptyList{ast.Record{"mapKey": ast.Natural, "mapValue": ast.Text}},
			new(interface{}),
			map[uint]string{}),
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
		It("Unmarshals the identity string function", func() {
			var fn func(string) string
			dhallFn := &ast.LambdaExpr{
				Label: "x",
				Type:  ast.Text,
				Body:  ast.MkVar("x"),
			}
			Unmarshal(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn("foo")).To(Equal("foo"))
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
