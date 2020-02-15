package dhall_test

import (
	"reflect"

	. "github.com/philandstuff/dhall-golang"
	"github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/term"

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

type testTaggedStruct struct {
	Foo int `json:"baz,string"`
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
			core.TextLit{Suffix: "lalala"}, new(string), "lalala"),
	)
	DescribeTable("Compound types", DecodeAndCompare,
		Entry("unmarshals Some 5 into int",
			core.Some{core.NaturalLit(5)},
			new(int),
			5),
		Entry("unmarshals None Natural into int",
			core.NoneOf{core.Natural},
			new(int),
			0),
		Entry("unmarshals List Integer into int slice",
			core.NonEmptyList{core.IntegerLit(5)},
			new([]int),
			[]int{5}),
		Entry("unmarshals List Integer into int64 slice",
			core.NonEmptyList{core.IntegerLit(5)},
			new([]int64),
			[]int64{5}),
		Entry("unmarshals List Bool into slice",
			core.NonEmptyList{core.True, core.False},
			new([]bool),
			[]bool{true, false}),
		Entry("unmarshals empty List Bool into slice",
			core.EmptyList{core.Bool},
			new([]bool),
			[]bool{}),
		Entry("unmarshals None (List Bool) into slice",
			core.NoneOf{core.ListOf{core.Bool}},
			new([]bool),
			[]bool(nil)),
		Entry("unmarshals List (List Bool) into slice",
			core.NonEmptyList{
				core.NonEmptyList{core.True, core.False}},
			new([][]bool),
			[][]bool{{true, false}}),
		Entry("unmarshals {Foo : Natural, Bar : Text} into struct",
			core.RecordLit{"Foo": core.NaturalLit(3), "Bar": core.TextLit{Suffix: "xyzzy"}},
			new(testStruct),
			testStruct{Foo: 3, Bar: "xyzzy"}),
		Entry("unmarshals {baz : Natural, Bar : Text} into tagged struct",
			core.RecordLit{"baz": core.NaturalLit(3), "Bar": core.TextLit{Suffix: "xyzzy"}},
			new(testTaggedStruct),
			testTaggedStruct{Foo: 3, Bar: "xyzzy"}),
		Entry("unmarshals None {Foo : Natural, Bar : Text} into struct",
			core.NoneOf{core.RecordType{"Foo": core.Natural, "Bar": core.Text}},
			new(testStruct),
			testStruct{}),
		Entry("unmarshals List {mapKey : Natural, mapValue : Text} into map",
			core.NonEmptyList{core.RecordLit{"mapKey": core.NaturalLit(3), "mapValue": core.TextLit{Suffix: "fizz"}},
				core.RecordLit{"mapKey": core.NaturalLit(5), "mapValue": core.TextLit{Suffix: "buzz"}}},
			new(map[int]string),
			map[int]string{3: "fizz", 5: "buzz"}),
		Entry("unmarshals None List {mapKey : Natural, mapValue : Text} into map",
			core.NoneOf{core.ListOf{core.RecordType{"mapKey": core.Natural, "mapValue": core.Text}}},
			new(map[int]string),
			map[int]string(nil)),
		Entry("unmarshals empty List {mapKey : Natural, mapValue : Text} into map",
			core.EmptyList{core.RecordType{"mapKey": core.Natural, "mapValue": core.Text}},
			new(map[int]string),
			map[int]string{}),
	)
	// Testing various identity functions ensures we support both
	// encoding and decoding each particular type
	DescribeTable("Identity functions of various Dhall and Go types",
		func(inputType term.Term, testValue interface{}) {
			id := core.Eval(term.Lambda{
				Label: "x",
				Type:  inputType,
				Body:  term.NewVar("x"),
			})
			arg := reflect.ValueOf(testValue)
			fnPtr := reflect.New(
				reflect.FuncOf(
					[]reflect.Type{arg.Type()},
					[]reflect.Type{arg.Type()},
					false))
			Decode(id, fnPtr.Interface())
			fnVal := fnPtr.Elem()
			result := fnVal.Call([]reflect.Value{arg})
			Expect(result[0].Interface()).To(Equal(arg.Interface()))
		},
		Entry("Bool into bool", term.Bool, true),
		Entry("Natural into int", term.Natural, 1),
		Entry("Natural into uint", term.Natural, uint(1)),
		Entry("Natural into int64", term.Natural, int64(1)),
		Entry("Integer into int", term.Integer, 1),
		Entry("Integer into int64", term.Integer, int64(1)),
		Entry("Text into string", term.Text, "foo"),
		Entry("List Natural into []int",
			term.Apply(term.List, term.Natural), []int{1, 2, 3}),
		XEntry("Optional Natural into int",
			term.Apply(term.Optional, term.Natural), 1),
		Entry("record into struct",
			term.RecordType{"Foo": term.Natural, "Bar": term.Text},
			testStruct{Foo: 1, Bar: "howdy"},
		),
		Entry("record into tagged struct",
			term.RecordType{"baz": term.Natural, "Bar": term.Text},
			testTaggedStruct{Foo: 1, Bar: "howdy"},
		),
		Entry("map into map",
			term.Apply(term.List,
				term.RecordType{"mapKey": term.Text, "mapValue": term.Natural}),
			map[string]int{"foo": 1, "bar": 2},
		),
	)
	Describe("Function types", func() {
		It("Decodes the int successor function", func() {
			var fn func(int) int
			dhallFn := core.Eval(term.Lambda{
				Label: "x",
				Type:  term.Natural,
				Body: term.NaturalPlus(
					term.NewVar("x"),
					term.NaturalLit(1),
				),
			})
			Decode(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(3)).To(Equal(4))
		})
		It("Decodes the natural sum function", func() {
			var fn func(int, int) int
			dhallFn := core.Eval(term.Lambda{
				Label: "x",
				Type:  term.Natural,
				Body: term.Lambda{
					Label: "y",
					Type:  term.Natural,
					Body: term.NaturalPlus(
						term.NewVar("x"), term.NewVar("y")),
				},
			})
			Decode(dhallFn, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(3, 4)).To(Equal(7))
		})
		It("Decodes the Natural/subtract builtin as a function", func() {
			var fn func(int, int) int
			Decode(core.NaturalSubtract, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(1, 3)).To(Equal(2))
		})
		It("Decodes the Natural/subtract builtin as a curried function", func() {
			var fn func(int) func(int) int
			Decode(core.NaturalSubtract, &fn)
			Expect(fn).ToNot(BeNil())
			Expect(fn(1)(3)).To(Equal(2))
		})
	})
})

var _ = Describe("Unmarshal", func() {
	It("Parses 1 + 2", func() {
		var actual uint
		err := Unmarshal([]byte("1 + 2"), &actual)
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(uint(3)))
	})
	It("Throws a type error for `1 + -2`", func() {
		var actual uint
		err := Unmarshal([]byte("1 + -2"), &actual)
		Expect(err).To(HaveOccurred())
	})
	It("Fetches imports", func() {
		type Config struct {
			Port int
			Name string
		}
		var actual Config
		err := Unmarshal([]byte("./testdata/unmarshal-test.dhall"), &actual)
		Expect(err).ToNot(HaveOccurred())
		Expect(actual).To(Equal(Config{Port: 5050, Name: "inetd"}))
	})
})
