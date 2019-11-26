package dhall_test

import (
	"reflect"

	. "github.com/philandstuff/dhall-golang"
	"github.com/philandstuff/dhall-golang/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func DecodeAndCompare(input core.Value, ptr interface{}, expected interface{}) {
	Decode(input, ptr)
	// use reflect to dereference a pointer of unknown type
	Expect(reflect.ValueOf(ptr).Elem().Interface()).
		To(Equal(expected))
}

type testStruct struct {
	Foo int
	Bar string
}

var _ = Describe("Decode", func() {
	DescribeTable("Simple types", DecodeAndCompare,
		Entry("unmarshals DoubleLit into float32",
			core.DoubleLit(3.5), new(float32), float32(3.5)),
		Entry("unmarshals DoubleLit into float64",
			core.DoubleLit(3.5), new(float64), float64(3.5)),
		Entry("unmarshals True into bool",
			core.True, new(bool), true),
		Entry("unmarshals NaturalLit into int",
			core.NaturalLit(5), new(int), 5),
		Entry("unmarshals NaturalLit into int64",
			core.NaturalLit(5), new(int64), int64(5)),
		Entry("unmarshals IntegerLit into int",
			core.IntegerLit(5), new(int), 5),
		Entry("unmarshals IntegerLit into int",
			core.IntegerLit(5), new(int64), int64(5)),
		Entry("unmarshals TextLit into string",
			core.TextLitVal{Suffix: "lalala"}, new(string), "lalala"),
	)
	DescribeTable("Simple types into interface{}", DecodeAndCompare,
		Entry("unmarshals DoubleLit into interface{}",
			core.DoubleLit(3.5), new(interface{}), float64(3.5)),
		Entry("unmarshals True into interface{}",
			core.True, new(interface{}), true),
		Entry("unmarshals NaturalLit into interface{}",
			core.NaturalLit(5), new(interface{}), uint(5)),
		Entry("unmarshals IntegerLit into interface{}",
			core.IntegerLit(5), new(interface{}), int(5)),
		Entry("unmarshals TextLit into interface{}",
			core.TextLitVal{Suffix: "lalala"}, new(interface{}), "lalala"),
	)
	DescribeTable("Compound types", DecodeAndCompare,
		Entry("unmarshals Some 5 into int",
			core.SomeVal{core.NaturalLit(5)},
			new(int),
			5),
		Entry("unmarshals None Natural into int",
			core.AppValue{core.None, core.Natural},
			new(int),
			0),
		Entry("unmarshals List Integer into int slice",
			core.NonEmptyListVal{core.IntegerLit(5)},
			new([]int),
			[]int{5}),
		Entry("unmarshals List Integer into int64 slice",
			core.NonEmptyListVal{core.IntegerLit(5)},
			new([]int64),
			[]int64{5}),
		Entry("unmarshals List Bool into slice",
			core.NonEmptyListVal{core.True, core.False},
			new([]bool),
			[]bool{true, false}),
		Entry("unmarshals empty List Bool into slice",
			core.EmptyListVal{core.Bool},
			new([]bool),
			[]bool{}),
		Entry("unmarshals None (List Bool) into slice",
			core.AppValue{core.None, core.EmptyListVal{core.Bool}},
			new([]bool),
			[]bool(nil)),
		Entry("unmarshals List (List Bool) into slice",
			core.NonEmptyListVal{
				core.NonEmptyListVal{core.True, core.False}},
			new([][]bool),
			[][]bool{{true, false}}),
		Entry("unmarshals {Foo : Natural, Bar : Text} into struct",
			core.RecordLitVal{"Foo": core.NaturalLit(3), "Bar": core.TextLitVal{Suffix: "xyzzy"}},
			new(testStruct),
			testStruct{Foo: 3, Bar: "xyzzy"}),
		Entry("unmarshals None {Foo : Natural, Bar : Text} into struct",
			core.AppValue{core.None, core.RecordTypeVal{"Foo": core.Natural, "Bar": core.Text}},
			new(testStruct),
			testStruct{}),
		Entry("unmarshals List {mapKey : Natural, mapValue : Text} into map",
			core.NonEmptyListVal{core.RecordLitVal{"mapKey": core.NaturalLit(3), "mapValue": core.TextLitVal{Suffix: "fizz"}},
				core.RecordLitVal{"mapKey": core.NaturalLit(5), "mapValue": core.TextLitVal{Suffix: "buzz"}}},
			new(map[int]string),
			map[int]string{3: "fizz", 5: "buzz"}),
		Entry("unmarshals None List {mapKey : Natural, mapValue : Text} into map",
			core.AppValue{core.None, core.AppValue{core.List, core.RecordTypeVal{"mapKey": core.Natural, "mapValue": core.Text}}},
			new(map[int]string),
			map[int]string(nil)),
		Entry("unmarshals empty List {mapKey : Natural, mapValue : Text} into map",
			core.EmptyListVal{core.RecordTypeVal{"mapKey": core.Natural, "mapValue": core.Text}},
			new(map[int]string),
			map[int]string{}),
	)
	// EmptyListVal into interface{} deserves its own tests
	// because we have to convert the list type to a reflect.Type
	DescribeTable("Empty list into interface{}", DecodeAndCompare,
		Entry("unmarshals empty List Bool into interface{}",
			core.EmptyListVal{core.Bool},
			new(interface{}),
			[]bool{}),
		Entry("unmarshals empty List Natural into interface{}",
			core.EmptyListVal{core.Natural},
			new(interface{}),
			[]uint{}),
		Entry("unmarshals empty List Integer into interface{}",
			core.EmptyListVal{core.Integer},
			new(interface{}),
			[]int{}),
		Entry("unmarshals empty List Double into interface{}",
			core.EmptyListVal{core.Double},
			new(interface{}),
			[]float64{}),
		Entry("unmarshals empty List Text into interface{}",
			core.EmptyListVal{core.Text},
			new(interface{}),
			[]string{}),
		Entry("unmarshals empty List {Foo : Integer, Bar : Text} into interface{}",
			core.EmptyListVal{core.RecordTypeVal{"Foo": core.Integer, "Bar": core.Text}},
			new(interface{}),
			[]struct {
				Bar string
				Foo int
			}{}),
		Entry("unmarshals empty List {foo : Integer, bar : Text} into interface{}",
			core.EmptyListVal{core.RecordTypeVal{"foo": core.Integer, "bar": core.Text}},
			new(interface{}),
			[]struct {
				Bar string
				Foo int
			}{}),
		Entry("unmarshals empty List {mapKey : Integer, mapValue : Text} into interface{}",
			core.EmptyListVal{core.RecordTypeVal{"mapKey": core.Integer, "mapValue": core.Text}},
			new(interface{}),
			map[int]string{}))
	DescribeTable("Compound types into interface{}", DecodeAndCompare,
		XEntry("unmarshals List Integer into interface{}",
			core.NonEmptyListVal{core.IntegerLit(5)},
			new(interface{}),
			[]int{5}),
		XEntry("unmarshals {Foo : Integer, Bar : Text} into interface{}",
			core.RecordLitVal{"Foo": core.IntegerLit(3), "Bar": core.TextLitVal{Suffix: "xyzzy"}},
			new(interface{}),
			testStruct{Foo: 3, Bar: "xyzzy"}),
		XEntry("unmarshals None {Foo : Integer, Bar : Text} into struct",
			core.AppValue{core.None, core.RecordTypeVal{"Foo": core.Integer, "Bar": core.Text}},
			new(interface{}),
			testStruct{}),
		XEntry("unmarshals List {mapKey : Natural, mapValue : Text}",
			core.NonEmptyListVal{core.RecordLitVal{"mapKey": core.NaturalLit(3), "mapValue": core.TextLitVal{Suffix: "fizz"}},
				core.RecordLitVal{"mapKey": core.NaturalLit(5), "mapValue": core.TextLitVal{Suffix: "buzz"}}},
			new(interface{}),
			map[uint]string{3: "fizz", 5: "buzz"}),
		XEntry("unmarshals None List {mapKey : Natural, mapValue : Text}",
			core.AppValue{core.None, core.AppValue{core.List, core.RecordTypeVal{"mapKey": core.Natural, "mapValue": core.Text}}},
			new(interface{}),
			map[uint]string(nil)),
		Entry("unmarshals empty List {mapKey : Natural, mapValue : Text}",
			core.EmptyListVal{core.RecordTypeVal{"mapKey": core.Natural, "mapValue": core.Text}},
			new(interface{}),
			map[uint]string{}),
	)
	Describe("Function types", func() {
		It("Decodes the identity int function", func() {
			var fn func(int) int
			dhallFn := core.Eval(core.LambdaTerm{
				Label: "x",
				Type:  core.Natural,
				Body:  core.NewVar("x"),
			})
			Decode(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(3)).To(Equal(3))
		})
		It("Decodes the identity int64 function", func() {
			var fn func(int64) int64
			dhallFn := core.Eval(core.LambdaTerm{
				Label: "x",
				Type:  core.Natural,
				Body:  core.NewVar("x"),
			})
			Decode(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(int64(3))).To(Equal(int64(3)))
		})
		It("Decodes the identity string function", func() {
			var fn func(string) string
			dhallFn := core.Eval(core.LambdaTerm{
				Label: "x",
				Type:  core.Text,
				Body:  core.NewVar("x"),
			})
			Decode(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn("foo")).To(Equal("foo"))
		})
		It("Decodes the int successor function", func() {
			var fn func(int) int
			dhallFn := core.LambdaValue{
				Label:  "x",
				Domain: core.Natural,
				Fn: func(x core.Value) core.Value {
					return core.Eval(core.NaturalPlus(
						core.Quote(x),
						core.NaturalLit(1),
					))
				},
			}
			Decode(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(3)).To(Equal(4))
		})
		It("Decodes the natural sum function", func() {
			var fn func(int, int) int
			dhallFn := core.LambdaValue{
				Label:  "x",
				Domain: core.Natural,
				Fn: func(x core.Value) core.Value {
					return core.LambdaValue{
						Label:  "y",
						Domain: core.Natural,
						Fn: func(y core.Value) core.Value {
							return core.Eval(core.NaturalPlus(core.Quote(x), core.Quote(y)))
						},
					}
				},
			}
			Decode(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(3, 4)).To(Equal(7))
		})
	})
})
