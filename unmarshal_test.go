package dhall_test

import (
	"reflect"

	. "github.com/philandstuff/dhall-golang/v5"
	"github.com/philandstuff/dhall-golang/v5/core"
	"github.com/philandstuff/dhall-golang/v5/term"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func DecodeAndCompare(input core.Value, ptr interface{}, expected interface{}) {
	err := Decode(input, ptr)
	Expect(err).ToNot(HaveOccurred())
	// use reflect to dereference a pointer of unknown type
	Expect(reflect.ValueOf(ptr).Elem().Interface()).
		To(Equal(expected))
}

type testStruct struct {
	Foo uint
	Bar string
}

type testTaggedStruct struct {
	Foo uint `dhall:"baz"`
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
			core.PlainTextLit("lalala"), new(string), "lalala"),
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
			core.RecordLit{"Foo": core.NaturalLit(3), "Bar": core.PlainTextLit("xyzzy")},
			new(testStruct),
			testStruct{Foo: 3, Bar: "xyzzy"}),
		Entry("unmarshals {Foo : Natural, Bar : Text} into *struct",
			core.RecordLit{"Foo": core.NaturalLit(3), "Bar": core.PlainTextLit("xyzzy")},
			new(*testStruct),
			&testStruct{Foo: 3, Bar: "xyzzy"}),
		Entry("unmarshals {baz : Natural, Bar : Text} into tagged struct",
			core.RecordLit{"baz": core.NaturalLit(3), "Bar": core.PlainTextLit("xyzzy")},
			new(testTaggedStruct),
			testTaggedStruct{Foo: 3, Bar: "xyzzy"}),
		Entry("unmarshals None {Foo : Natural, Bar : Text} into struct",
			core.NoneOf{core.RecordType{"Foo": core.Natural, "Bar": core.Text}},
			new(testStruct),
			testStruct{}),
		Entry("unmarshals List {mapKey : Natural, mapValue : Text} into map",
			core.NonEmptyList{core.RecordLit{"mapKey": core.NaturalLit(3), "mapValue": core.PlainTextLit("fizz")},
				core.RecordLit{"mapKey": core.NaturalLit(5), "mapValue": core.PlainTextLit("buzz")}},
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
			err := Decode(id, fnPtr.Interface())
			Expect(err).ToNot(HaveOccurred())
			fnVal := fnPtr.Elem()
			result := fnVal.Call([]reflect.Value{arg})
			Expect(result[0].Interface()).To(Equal(arg.Interface()))
		},
		Entry("Bool into bool", term.Bool, true),
		Entry("Natural into uint", term.Natural, uint(1)),
		Entry("Natural into uint64", term.Natural, uint64(1)),
		Entry("Integer into int", term.Integer, int(1)),
		Entry("Integer into int64", term.Integer, int64(1)),
		Entry("Text into string", term.Text, "foo"),
		Entry("List Natural into []uint",
			term.Apply(term.List, term.Natural), []uint{1, 2, 3}),
		Entry("Optional Natural into uint",
			term.Apply(term.Optional, term.Natural), uint(1)),
		Entry("Optional Natural into *uint with nil",
			term.Apply(term.Optional, term.Natural), (*uint)(nil)),
		Entry("Optional Natural into *uint with value",
			term.Apply(term.Optional, term.Natural), new(uint)),
		Entry("Optional (List Natural) into []uint with nil",
			term.Apply(term.Optional, term.Apply(term.List, term.Natural)), []uint(nil)),
		Entry("Optional (List Natural) into []uint with empty slice",
			term.Apply(term.Optional, term.Apply(term.List, term.Natural)), []uint{}),
		Entry("Optional (List Natural) into []uint with nonempty slice",
			term.Apply(term.Optional, term.Apply(term.List, term.Natural)), []uint{1, 2}),
		Entry("record into struct",
			term.RecordType{"Foo": term.Natural, "Bar": term.Text},
			testStruct{Foo: 1, Bar: "howdy"},
		),
		Entry("record into tagged struct",
			term.RecordType{"baz": term.Natural, "Bar": term.Text},
			testTaggedStruct{Foo: 1, Bar: "howdy"},
		),
		Entry("record into struct ptr",
			term.RecordType{"Foo": term.Natural, "Bar": term.Text},
			&testStruct{Foo: 1, Bar: "howdy"},
		),
		Entry("map into map",
			term.Apply(term.List,
				term.RecordType{"mapKey": term.Text, "mapValue": term.Natural}),
			map[string]uint{"foo": 1, "bar": 2},
		),
	)
	Describe("Function types", func() {
		It("Decodes the Natural successor function", func() {
			var fn func(uint) uint
			dhallFn := core.Eval(term.Lambda{
				Label: "x",
				Type:  term.Natural,
				Body: term.NaturalPlus(
					term.NewVar("x"),
					term.NaturalLit(1),
				),
			})
			err := Decode(dhallFn, &fn)
			Expect(err).ToNot(HaveOccurred())
			Expect(fn).ToNot(BeNil())
			Expect(fn(uint(3))).To(Equal(uint(4)))
		})
		It("Decodes the natural sum function", func() {
			var fn func(uint, uint) uint
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
			err := Decode(dhallFn, &fn)
			Expect(err).ToNot(HaveOccurred())
			Expect(fn).ToNot(BeNil())
			Expect(fn(3, 4)).To(Equal(uint(7)))
		})
		It("Decodes the Natural/subtract builtin as a function", func() {
			var fn func(uint, uint) uint
			err := Decode(core.NaturalSubtract, &fn)
			Expect(err).ToNot(HaveOccurred())

			Expect(fn).ToNot(BeNil())
			Expect(fn(uint(1), uint(3))).To(Equal(uint(2)))
		})
		It("Decodes the Natural/subtract builtin as a curried function", func() {
			var fn func(uint) func(uint) uint
			err := Decode(core.NaturalSubtract, &fn)
			Expect(err).ToNot(HaveOccurred())

			Expect(fn).ToNot(BeNil())
			Expect(fn(uint(1))(uint(3))).To(Equal(uint(2)))
		})
	})
	// TODO expected errors
})

func ExpectUnmarshalError(source string, targetVar interface{}) {
	err := Unmarshal([]byte(source), targetVar)
	Expect(err).To(HaveOccurred())
}

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
	Context("Unmarshalling functions", func() {
		DescribeTable("Expected successes",
			func(source string, targetVar interface{}, testInput interface{}, expectedOutput interface{}) {
				err := Unmarshal([]byte(source), targetVar)
				Expect(err).ToNot(HaveOccurred())
				Expect(reflect.ValueOf(targetVar).Elem().Call([]reflect.Value{reflect.ValueOf(testInput)})[0].Interface()).
					To(Equal(expectedOutput))
			},
			Entry("Natural identity", `
				λ(x : Natural) → x
			`, new(func(uint) uint), uint(37), uint(37)),
			Entry("Can return Natural as int", `
				λ(x : Natural) → x
			`, new(func(uint) int), uint(37), int(37)),
			Entry("Natural successor", `
				λ(x : Natural) → x + 1
			`, new(func(uint) uint), uint(37), uint(38)),
			Entry("Natural/even", `
				Natural/even
			`, new(func(uint) bool), uint(2), true),
			Entry("Natural/show", `
				Natural/show
			`, new(func(uint) string), uint(37), "37"),
			Entry("Text greet", `
				λ(x : Text) → "Hello, " ++ x ++ "!"
			`, new(func(string) string), "Brian", "Hello, Brian!"),
		)
		DescribeTable("Expected type errors", ExpectUnmarshalError,
			Entry("Incompatible output parameter type", `
				λ(x : Natural) → x
			`, new(func(uint) string)),
			Entry("No output parameters", `
				λ(x : Natural) → x
			`, new(func(uint))),
			Entry("Multiple output parameters", `
				λ(x : Natural) → x
			`, new(func(uint) (uint, error))),
			Entry("No input parameters", `
				λ(x : Natural) → x
			`, new(func() uint)),
			Entry("Incompatible input parameter type", `
				λ(x : Natural) → x
			`, new(func(string) uint)),
		)
	})
})
