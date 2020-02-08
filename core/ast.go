package core

import (
	"fmt"
	"math"
	"strings"
)

// A Term is an arbitrary Dhall expression.  When you parse text into
// Dhall, you get a value of type Term.
//
// To evaluate a Term, you call the Eval() function, which returns a
// Value.
type Term interface {
	isTerm()
}

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

const (
	// Type is the type of all types
	Type Universe = iota
	// Kind is the type of all kinds (including Type)
	Kind
	// Sort is the type of all sorts (including Kind)
	Sort
)

// Builtin is the type of Dhall's builtin constants
type Builtin string

const (
	// Double is the type of doubles
	Double = Builtin("Double")
	// Text is the type of Text
	Text = Builtin("Text")
	// Bool is the type of booleans
	Bool = Builtin("Bool")
	// Natural is the type of natural numbers
	Natural = Builtin("Natural")
	// Integer is the type of integers
	Integer = Builtin("Integer")
	// List is a function that takes a type to the type of lists of
	// that type
	List = Builtin("List")
	// Optional is a function that takes a type to the type of
	// optional values of that type
	Optional = Builtin("Optional")
	// None takes a type to the empty optional value of that type
	None = Builtin("None")

	// NaturalBuild is Natural/build
	NaturalBuild = Builtin("Natural/build")
	// NaturalFold is Natural/fold
	NaturalFold = Builtin("Natural/fold")
	// NaturalZero is Natural/isZero
	NaturalIsZero = Builtin("Natural/isZero")
	// NaturalEven is Natural/even
	NaturalEven = Builtin("Natural/even")
	// NaturalOdd is Natural/odd
	NaturalOdd = Builtin("Natural/odd")
	// NaturalToInteger is Natural/toInteger
	NaturalToInteger = Builtin("Natural/toInteger")
	// NaturalShow is Natural/show
	NaturalShow = Builtin("Natural/show")
	// NaturalSubtract is Natural/subtract
	NaturalSubtract = Builtin("Natural/subtract")

	// IntegerClamp is Integer/clamp
	IntegerClamp = Builtin("Integer/clamp")
	// IntegerNegate is Integer/negate
	IntegerNegate = Builtin("Integer/negate")
	// IntegerToDouble is Integer/toDouble
	IntegerToDouble = Builtin("Integer/toDouble")
	// IntegerShow is Integer/show
	IntegerShow = Builtin("Integer/show")

	// DoubleShow is Double/show
	DoubleShow = Builtin("Double/show")

	// TextShow is Text/show
	TextShow = Builtin("Text/show")

	// ListBuild is List/build
	ListBuild = Builtin("List/build")
	// ListFold is List/fold
	ListFold = Builtin("List/fold")
	// ListLength is List/length
	ListLength = Builtin("List/length")
	// ListHead is List/head
	ListHead = Builtin("List/head")
	// ListLast is List/last
	ListLast = Builtin("List/last")
	// ListIndexed is List/indexed
	ListIndexed = Builtin("List/indexed")
	// ListReverse is List/reverse
	ListReverse = Builtin("List/reverse")

	// OptionalBuild is Optional/build
	OptionalBuild = Builtin("Optional/build")
	// OptionalFold is Optional/fold
	OptionalFold = Builtin("Optional/fold")

	// True is the true Bool value
	True = BoolLit(true)
	// False is the false Bool value
	False = BoolLit(false)
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

	optionalBuildVal struct{ typ Value }
	optionalFoldVal  struct {
		typ1 Value
		opt  Value
		typ2 Value
		some Value
		// none Value
	}

	textShowVal struct{}

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

func (optionalBuildVal) isValue() {}
func (optionalFoldVal) isValue()  {}

func (textShowVal) isValue() {}

func (listBuildVal) isValue()   {}
func (listFoldVal) isValue()    {}
func (listLengthVal) isValue()  {}
func (listHeadVal) isValue()    {}
func (listLastVal) isValue()    {}
func (listIndexedVal) isValue() {}
func (listReverseVal) isValue() {}

type (
	// A Var is a variable, either bound or free.  A free Var is a
	// valid Value; a bound Var is not.
	//
	// The Index is a de Bruijn index.  In an expression such as:
	//
	//  λ(x : Natural) → λ(x : Natural) → x@1
	//
	// x@1 refers to the outer bound variable x.  x@1 is represented
	// by Var{"x", 1}.
	Var struct {
		Name  string
		Index int
	}

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

func (Universe) isTerm()  {}
func (Universe) isValue() {}

func (Builtin) isTerm()  {}
func (Builtin) isValue() {}

func (Var) isTerm()  {}
func (Var) isValue() {}

func (localVar) isTerm()  {}
func (localVar) isValue() {}

func (quoteVar) isValue() {}

// NewVar returns a new Var Term
func NewVar(name string) Term {
	return Var{Name: name}
}

type (
	// A LambdaTerm is a lambda abstraction.
	LambdaTerm struct {
		Label string
		Type  Term
		Body  Term
	}

	// A PiTerm is a function type.
	PiTerm struct {
		Label string
		Type  Term
		Body  Term
	}

	// An AppTerm is a Term formed by applying a function to an
	// argument.
	AppTerm struct {
		Fn  Term
		Arg Term
	}

	// An OpTerm is two Terms combined by an operator.  The OpCode
	// determines the type of operator.
	OpTerm struct {
		OpCode OpCode
		L      Term
		R      Term
	}
)

func (LambdaTerm) isTerm() {}
func (PiTerm) isTerm()     {}
func (AppTerm) isTerm()    {}
func (OpTerm) isTerm()     {}

// NewLambda constructs a new lambda Term.
func NewLambda(label string, t Term, body Term) Term {
	return LambdaTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

// NewPi constructs a new pi Term.
func NewPi(label string, t Term, body Term) Term {
	return PiTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

// NewAnonPi returns a pi Term with label "_", typically used for
// non-dependent function types.
func NewAnonPi(domain Term, codomain Term) Term {
	return NewPi("_", domain, codomain)
}

// Apply takes fn and applies it successively to each arg in args.
func Apply(fn Term, args ...Term) Term {
	out := fn
	for _, arg := range args {
		out = AppTerm{Fn: out, Arg: arg}
	}
	return out
}

// An OpCode encodes the type of an operator.
type OpCode int

// These are the valid OpCodes.  Note that the actual integer values
// have been chosen to match the binary encoding label numbers.
const (
	OrOp  OpCode = iota // a || b
	AndOp               // a && b
	EqOp                // a == b
	NeOp                // a != b

	PlusOp  // a + b
	TimesOp // a * b

	TextAppendOp // a ++ b
	ListAppendOp // a # b

	RecordMergeOp            // a /\ b
	RightBiasedRecordMergeOp // a // b
	RecordTypeMergeOp        // a //\\ b

	ImportAltOp // a ? b
	EquivOp     // a === b
	CompleteOp  // A::b
)

// NaturalPlus takes Terms l and r and returns (l + r)
func NaturalPlus(l, r Term) Term {
	return OpTerm{OpCode: PlusOp, L: l, R: r}
}

// NaturalTimes takes Terms l and r and returns (l * r)
func NaturalTimes(l, r Term) Term {
	return OpTerm{OpCode: TimesOp, L: l, R: r}
}

// BoolOr takes Terms l and r and returns (l || r)
func BoolOr(l, r Term) Term {
	return OpTerm{OpCode: OrOp, L: l, R: r}
}

// BoolAnd takes Terms l and r and returns (l && r)
func BoolAnd(l, r Term) Term {
	return OpTerm{OpCode: AndOp, L: l, R: r}
}

// ListAppend takes Terms l and r and returns (l # r)
func ListAppend(l, r Term) Term {
	return OpTerm{OpCode: ListAppendOp, L: l, R: r}
}

// TextAppend takes Terms l and r and returns (l ++ r)
func TextAppend(l, r Term) Term {
	return OpTerm{OpCode: TextAppendOp, L: l, R: r}
}

// Equivalent takes Terms l and r and returns (l ≡ r)
func Equivalent(l, r Term) Term {
	return OpTerm{OpCode: EquivOp, L: l, R: r}
}

// Callable is a function Value that can be called with one Value
// argument.  Call() may return nil if normalization isn't possible
// (for example, `Natural/even x` does not normalize)
// ArgType() returns the declared type of Call()'s parameter
type Callable interface {
	Value
	Call(Value) Value
	ArgType() Value
}

// Call implements Callable
func (l lambdaValue) Call(a Value) Value {
	return l.Fn(a)
}

// ArgType implements Callable
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
	AppValue struct {
		Fn  Value
		Arg Value
	}

	opValue struct {
		OpCode OpCode
		L      Value
		R      Value
	}
)

func (lambdaValue) isValue() {}

func (PiValue) isValue() {}

func (AppValue) isValue() {}

func (opValue) isValue() {}

// NewPiVal returns a new pi Value.
func NewPiVal(label string, d Value, r func(Value) Value) Value {
	return PiValue{
		Label:  label,
		Domain: d,
		Range:  r,
	}
}

// NewFnTypeVal returns a non-dependent function type Value.
func NewFnTypeVal(l string, d Value, r Value) Value {
	return NewPiVal(l, d, func(Value) Value { return r })
}

type (
	Binding struct {
		Variable   string
		Annotation Term // may be nil
		Value      Term
	}
	Let struct {
		Bindings []Binding
		Body     Term
	}

	// no LetValue, since normalization removes all lets

	Annot struct {
		Expr       Term
		Annotation Term
	}

	// no AnnotValue either
)

func (Let) isTerm()   {}
func (Annot) isTerm() {}

// NewLet returns a let Term
func NewLet(body Term, bindings ...Binding) Term {
	return Let{Bindings: bindings, Body: body}
}

type (
	// A NaturalLit is a literal of type Natural
	NaturalLit uint

	// An EmptyList is an empty list literal Term of the given type
	EmptyList struct{ Type Term }
	// An EmptyListVal is an empty list literal Value of the given type
	EmptyListVal struct{ Type Value }

	// A NonEmptyList is a non-empty list literal Term with the given contents
	NonEmptyList []Term
	// A NonEmptyListVal is a non-empty list literal Value with the given contents
	NonEmptyListVal []Value

	Chunk struct {
		Prefix string
		Expr   Term
	}
	Chunks      []Chunk
	TextLitTerm struct {
		Chunks Chunks
		Suffix string
	}

	ChunkVal struct {
		Prefix string
		Expr   Value
	}
	ChunkVals  []ChunkVal
	TextLitVal struct {
		Chunks ChunkVals
		Suffix string
	}

	// A BoolLit is a literal of type Bool.
	BoolLit bool

	// An IfTerm is an if-then-else Term.
	IfTerm struct {
		Cond Term
		T    Term
		F    Term
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

	// Some represents a Term which is present in an Optional type.
	Some    struct{ Val Term }
	SomeVal struct{ Val Value }

	RecordType    map[string]Term
	RecordTypeVal map[string]Value

	RecordLit    map[string]Term
	RecordLitVal map[string]Value

	ToMap struct {
		Record Term
		Type   Term // optional
	}
	toMapVal struct {
		Record Value
		Type   Value // optional
	}

	Field struct {
		Record    Term
		FieldName string
	}
	fieldVal struct {
		Record    Value
		FieldName string
	}

	Project struct {
		Record     Term
		FieldNames []string
	}
	projectVal struct {
		Record     Value
		FieldNames []string
	}

	ProjectType struct {
		Record   Term
		Selector Term
	}
	// no ProjectTypeVal because it cannot be in a normal form

	UnionType    map[string]Term
	unionTypeVal map[string]Value

	Merge struct {
		Handler    Term
		Union      Term
		Annotation Term // optional
	}
	mergeVal struct {
		Handler    Value
		Union      Value
		Annotation Value // optional
	}

	Assert    struct{ Annotation Term }
	assertVal struct{ Annotation Value }
)

func (NaturalLit) isTerm()  {}
func (NaturalLit) isValue() {}

func (EmptyList) isTerm()     {}
func (EmptyListVal) isValue() {}

func (NonEmptyList) isTerm()     {}
func (NonEmptyListVal) isValue() {}

// NewList returns a non-empty list Term from the given Terms
func NewList(first Term, rest ...Term) Term {
	return append(NonEmptyList{first}, rest...)
}

// PlainText returns an uninterpolated text literal containing the
// given string as text.
func PlainText(content string) Term {
	return TextLitTerm{Suffix: content}
}

func (TextLitTerm) isTerm() {}
func (TextLitVal) isValue() {}

func (BoolLit) isTerm()  {}
func (BoolLit) isValue() {}
func (IfTerm) isTerm()   {}
func (ifVal) isValue()   {}

func (DoubleLit) isTerm()   {}
func (DoubleLit) isValue()  {}
func (IntegerLit) isTerm()  {}
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

func (Some) isTerm()     {}
func (SomeVal) isValue() {}

func (RecordType) isTerm()     {}
func (RecordTypeVal) isValue() {}
func (RecordLit) isTerm()      {}
func (RecordLitVal) isValue()  {}
func (ToMap) isTerm()          {}
func (toMapVal) isValue()      {}
func (Field) isTerm()          {}
func (fieldVal) isValue()      {}
func (Project) isTerm()        {}
func (projectVal) isValue()    {}
func (ProjectType) isTerm()    {}
func (UnionType) isTerm()      {}
func (unionTypeVal) isValue()  {}
func (Merge) isTerm()          {}
func (mergeVal) isValue()      {}
func (Assert) isTerm()         {}
func (assertVal) isValue()     {}

type (
	// An Import is an import Term.
	Import struct {
		ImportHashed
		ImportMode
	}

	// ImportHashed is a Fetchable with an optional hash for integrity
	// protection.
	ImportHashed struct {
		Fetchable
		Hash []byte // stored in multihash form - ie first two bytes are 0x12 0x20
	}

	// ImportMode can be normal (ie code import), "as Text" or "as
	// Location".
	ImportMode byte
)

const (
	// Code says to import as Dhall code.
	Code ImportMode = iota
	// RawText says to import as a Text value.
	RawText
	// Location says to import as a Location.
	Location
)

func (Import) isTerm() {}

// Decent output
func (c Universe) String() string {
	if c == Type {
		return "Type"
	} else if c == Kind {
		return "Kind"
	} else {
		return "Sort"
	}
}

func (v Var) String() string {
	if v.Index == 0 {
		return v.Name
	}
	return fmt.Sprint(v.Name, "@", v.Index)
}

func (v localVar) String() string {
	return fmt.Sprint("local:", v.Name, "/", v.Index)
}

func (lam LambdaTerm) String() string {
	return fmt.Sprintf("(λ(%s : %v) → %v)", lam.Label, lam.Type, lam.Body)
}

func (pi PiTerm) String() string {
	if pi.Label == "_" {
		return fmt.Sprintf("%v → %v", pi.Type, pi.Body)
	}
	return fmt.Sprintf("∀(%s : %v) → %v", pi.Label, pi.Type, pi.Body)
}

func (app AppTerm) String() string {
	if subApp, ok := app.Fn.(AppTerm); ok {
		return fmt.Sprintf("(%v %v)", subApp.stringNoParens(), app.Arg)
	}
	return fmt.Sprintf("(%v %v)", app.Fn, app.Arg)
}

func (app AppTerm) stringNoParens() string {
	if subApp, ok := app.Fn.(AppTerm); ok {
		return fmt.Sprintf("%v %v", subApp.stringNoParens(), app.Arg)
	}
	return fmt.Sprintf("%v %v", app.Fn, app.Arg)
}

// higher precedence binds tighter
func (op OpTerm) precedence() int {
	switch op.OpCode {
	case ImportAltOp:
		return 1
	case OrOp:
		return 2
	case PlusOp:
		return 3
	case TextAppendOp:
		return 4
	case ListAppendOp:
		return 5
	case AndOp:
		return 6
	case RecordMergeOp:
		return 7
	case RightBiasedRecordMergeOp:
		return 8
	case RecordTypeMergeOp:
		return 9
	case TimesOp:
		return 10
	case EqOp:
		return 11
	case NeOp:
		return 12
	case EquivOp:
		return 13
	case CompleteOp:
		return 14
	default:
		panic("unknown opcode")
	}
}

func (op OpTerm) operatorStr() string {
	switch op.OpCode {
	case ImportAltOp:
		return " ? "
	case OrOp:
		return " || "
	case PlusOp:
		return " + "
	case TextAppendOp:
		return " ++ "
	case ListAppendOp:
		return " # "
	case AndOp:
		return " && "
	case RecordMergeOp:
		return " ∧ "
	case RightBiasedRecordMergeOp:
		return " ⫽ "
	case RecordTypeMergeOp:
		return " ⩓ "
	case TimesOp:
		return " * "
	case EqOp:
		return " == "
	case NeOp:
		return " != "
	case EquivOp:
		return " ≡ "
	case CompleteOp:
		return "::"
	default:
		panic("unknown opcode")
	}
}

func (op OpTerm) String() string {
	prec := op.precedence()
	l := fmt.Sprint(op.L)
	r := fmt.Sprint(op.R)

	var buf strings.Builder

	if lop, ok := op.L.(OpTerm); ok {
		if prec > lop.precedence() {
			buf.WriteRune('(')
		}
		buf.WriteString(l)
		if prec > lop.precedence() {
			buf.WriteRune(')')
		}
	} else if _, ok := op.L.(AppTerm); ok {
		buf.WriteString(l)

	} else {
		buf.WriteRune('(')
		buf.WriteString(l)
		buf.WriteRune(')')
	}
	buf.WriteString(op.operatorStr())
	if rop, ok := op.R.(OpTerm); ok {
		if prec > rop.precedence() {
			buf.WriteRune('(')
		}
		buf.WriteString(r)
		if prec > rop.precedence() {
			buf.WriteRune(')')
		}
	} else if _, ok := op.R.(AppTerm); ok {
		buf.WriteString(r)

	} else {
		buf.WriteRune('(')
		buf.WriteString(r)
		buf.WriteRune(')')
	}
	return buf.String()
}
