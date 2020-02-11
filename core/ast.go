package core

import (
	"fmt"
	"math"

	"github.com/philandstuff/dhall-golang/term"
)

// A Value is a Dhall value in beta-normal form.  You can think of
// Values as the subset of Dhall which cannot be beta-reduced any
// further.  Valid Values include 3, "foo" and 5 + x.
//
// You create a Value by calling Eval() on a Term.  You can convert a
// Value back to a Term with Quote().
type Value interface {
	isValue()
}

// A Universe is a type of types.
type Universe int

// These are the valid Universes.
const (
	Type Universe = iota
	Kind
	Sort
)

// Builtin is the type of Dhall's builtin constants.
type Builtin string

// These are the Builtins.
const (
	Double   Builtin = "Double"
	Text     Builtin = "Text"
	Bool     Builtin = "Bool"
	Natural  Builtin = "Natural"
	Integer  Builtin = "Integer"
	List     Builtin = "List"
	Optional Builtin = "Optional"
	None     Builtin = "None"

	NaturalBuild     Builtin = "Natural/build"
	NaturalFold      Builtin = "Natural/fold"
	NaturalIsZero    Builtin = "Natural/isZero"
	NaturalEven      Builtin = "Natural/even"
	NaturalOdd       Builtin = "Natural/odd"
	NaturalToInteger Builtin = "Natural/toInteger"
	NaturalShow      Builtin = "Natural/show"
	NaturalSubtract  Builtin = "Natural/subtract"

	IntegerClamp    Builtin = "Integer/clamp"
	IntegerNegate   Builtin = "Integer/negate"
	IntegerToDouble Builtin = "Integer/toDouble"
	IntegerShow     Builtin = "Integer/show"

	DoubleShow Builtin = "Double/show"

	TextShow Builtin = "Text/show"

	ListBuild   Builtin = "List/build"
	ListFold    Builtin = "List/fold"
	ListLength  Builtin = "List/length"
	ListHead    Builtin = "List/head"
	ListLast    Builtin = "List/last"
	ListIndexed Builtin = "List/indexed"
	ListReverse Builtin = "List/reverse"

	OptionalBuild Builtin = "Optional/build"
	OptionalFold  Builtin = "Optional/fold"
)

// A BoolLit is a Dhall boolean literal.
type BoolLit bool

func (BoolLit) isValue() {}

// Naturally, it is True or False.
const (
	True  BoolLit = true
	False BoolLit = false
)

type (
	naturalBuildVal struct{}
	naturalEvenVal  struct{}
	naturalFoldVal  struct {
		n    Value
		typ  Value
		succ Value
		// zero Value
	}
	naturalIsZeroVal   struct{}
	naturalOddVal      struct{}
	naturalShowVal     struct{}
	naturalSubtractVal struct {
		a Value
		// b Value
	}
	naturalToIntegerVal struct{}

	integerClampVal    struct{}
	integerNegateVal   struct{}
	integerShowVal     struct{}
	integerToDoubleVal struct{}

	doubleShowVal struct{}

	optionalVal      struct{}
	optionalBuildVal struct{ typ Value }
	optionalFoldVal  struct {
		typ1 Value
		opt  Value
		typ2 Value
		some Value
		// none Value
	}
	noneVal struct{}

	textShowVal struct{}

	listVal      struct{}
	listBuildVal struct {
		typ Value
		// fn  Value
	}
	listFoldVal struct {
		typ1 Value
		list Value
		typ2 Value
		cons Value
		// empty Value
	}
	listLengthVal  struct{ typ Value }
	listHeadVal    struct{ typ Value }
	listLastVal    struct{ typ Value }
	listIndexedVal struct{ typ Value }
	listReverseVal struct{ typ Value }
)

func (naturalBuildVal) isValue()     {}
func (naturalEvenVal) isValue()      {}
func (naturalFoldVal) isValue()      {}
func (naturalIsZeroVal) isValue()    {}
func (naturalOddVal) isValue()       {}
func (naturalShowVal) isValue()      {}
func (naturalSubtractVal) isValue()  {}
func (naturalToIntegerVal) isValue() {}

func (integerClampVal) isValue()    {}
func (integerNegateVal) isValue()   {}
func (integerShowVal) isValue()     {}
func (integerToDoubleVal) isValue() {}

func (doubleShowVal) isValue() {}

func (optionalVal) isValue()      {}
func (optionalBuildVal) isValue() {}
func (optionalFoldVal) isValue()  {}
func (noneVal) isValue()          {}

func (textShowVal) isValue() {}

func (listVal) isValue()        {}
func (listBuildVal) isValue()   {}
func (listFoldVal) isValue()    {}
func (listLengthVal) isValue()  {}
func (listHeadVal) isValue()    {}
func (listLastVal) isValue()    {}
func (listIndexedVal) isValue() {}
func (listReverseVal) isValue() {}

type (
	// OptionalOf is the Value version of `Optional a`
	OptionalOf struct{ Type Value }

	// ListOf is the Value version of `List a`
	ListOf struct{ Type Value }

	// NoneOf is the Value version of `None a`
	NoneOf struct{ Type Value }
)

func (OptionalOf) isValue() {}
func (ListOf) isValue()     {}
func (NoneOf) isValue()     {}

// A freeVar is a free variable.  It can appear in a Value where we
// Eval() a sub-Term within a whole, larger Term.
type freeVar struct {
	Name  string
	Index int
}

type (
	// A localVar is an internal sentinel value used by TypeOf() in
	// the process of typechecking the body of lambdas and pis.
	// Essentially it lets us convert de Bruijn indices (which count
	// how many binders we need to pass from the variable to the
	// correct binder) to de Bruijn levels (which count how many
	// binders we passed before binding this variable)
	localVar struct {
		Name  string
		Index int
	}

	// A quoteVar is an internal sentinel value used by Quote() in the
	// process of converting Values back to Terms.
	quoteVar struct {
		Name  string
		Index int
	}
)

func (Universe) isValue() {}

func (Builtin) isValue() {}

func (freeVar) isValue() {}

func (localVar) isValue() {}

func (quoteVar) isValue() {}

// Callable is a function Value that can be called with one Value
// argument.  Call() may return nil if normalization isn't possible
// (for example, `Natural/even x` does not normalize).  ArgType()
// returns the declared type of Call()'s parameter.
type Callable interface {
	Value
	Call(Value) Value
	ArgType() Value
}

func (l lambdaValue) Call(a Value) Value {
	return l.Fn(a)
}

func (l lambdaValue) ArgType() Value {
	return l.Domain
}

var (
	_ Callable = lambdaValue{}
)

type (
	// A LambdaValue is a go function representing a Dhall function
	// which has not yet been applied to its argument
	lambdaValue struct {
		Label  string
		Domain Value
		Fn     func(Value) Value
	}

	// A PiValue is the value of a Dhall Pi type.  Domain is the type
	// of the domain, and Range is a go function which returns the
	// type of the range, given the type of the domain.
	PiValue struct {
		Label  string
		Domain Value
		Range  func(Value) Value
	}

	// An AppValue is the Value of a Fn applied to an Arg.  Note that
	// this is only a valid Value if Fn is a free variable (such as f
	// 3, with f free), or if Fn is a builtin function which hasn't
	// been applied to enough arguments (such as Natural/subtract 3).
	appValue struct {
		Fn  Value
		Arg Value
	}

	opValue struct {
		OpCode term.OpCode
		L      Value
		R      Value
	}
)

func (lambdaValue) isValue() {}

func (PiValue) isValue() {}

func (appValue) isValue() {}

func (opValue) isValue() {}

// NewPiVal returns a new pi Value.
func NewPiVal(label string, d Value, r func(Value) Value) PiValue {
	return PiValue{
		Label:  label,
		Domain: d,
		Range:  r,
	}
}

// NewFnTypeVal returns a non-dependent function type Value.
func NewFnTypeVal(l string, d Value, r Value) PiValue {
	return NewPiVal(l, d, func(Value) Value { return r })
}

type (
	// A NaturalLit is a literal of type Natural.
	NaturalLit uint

	// An EmptyListVal is an empty list literal Value of the given type.
	EmptyListVal struct{ Type Value }

	// A NonEmptyListVal is a non-empty list literal Value with the given contents.
	NonEmptyListVal []Value

	ChunkVal struct {
		Prefix string
		Expr   Value
	}
	ChunkVals  []ChunkVal
	TextLitVal struct {
		Chunks ChunkVals
		Suffix string
	}

	ifVal struct {
		Cond Value
		T    Value
		F    Value
	}

	// A DoubleLit is a literal of type Double.
	DoubleLit float64

	// A IntegerLit is a literal of type Integer.
	IntegerLit int

	// SomeVal represents a Value which is present in an Optional type.
	SomeVal struct{ Val Value }

	RecordTypeVal map[string]Value

	RecordLitVal map[string]Value

	toMapVal struct {
		Record Value
		Type   Value // optional
	}

	fieldVal struct {
		Record    Value
		FieldName string
	}

	projectVal struct {
		Record     Value
		FieldNames []string
	}

	// no ProjectTypeVal because it cannot be in a normal form

	unionTypeVal map[string]Value

	mergeVal struct {
		Handler    Value
		Union      Value
		Annotation Value // optional
	}

	assertVal struct{ Annotation Value }
)

func (NaturalLit) isValue() {}

func (EmptyListVal) isValue()    {}
func (NonEmptyListVal) isValue() {}

func (TextLitVal) isValue() {}

func (ifVal) isValue() {}

func (DoubleLit) isValue()  {}
func (IntegerLit) isValue() {}

func (d DoubleLit) String() string {
	f := float64(d)
	if math.IsInf(f, 1) {
		return "Infinity"
	}
	if math.IsInf(f, -1) {
		return "-Infinity"
	}
	// if we have a whole number, we need to append .0 to it so we get a valid
	// Double literal
	if f == float64(int64(f)) {
		return fmt.Sprintf("%#v.0", float64(d))
	}
	return fmt.Sprintf("%#v", float64(d))
}

func (SomeVal) isValue() {}

func (RecordTypeVal) isValue() {}
func (RecordLitVal) isValue()  {}
func (toMapVal) isValue()      {}
func (fieldVal) isValue()      {}
func (projectVal) isValue()    {}
func (unionTypeVal) isValue()  {}
func (mergeVal) isValue()      {}
func (assertVal) isValue()     {}
