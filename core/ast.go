package core

import (
	"fmt"
	"math"
	"net/url"
	"strings"
)

// Terms in the dhall internal language
type Term interface {
	isTerm()
}

// values
type Value interface {
	isValue()
}

type Universe int

const (
	Type Universe = iota
	Kind
	Sort
)

type Builtin string

const (
	Double   = Builtin("Double")
	Text     = Builtin("Text")
	Bool     = Builtin("Bool")
	Natural  = Builtin("Natural")
	Integer  = Builtin("Integer")
	List     = Builtin("List")
	Optional = Builtin("Optional")
	None     = Builtin("None")

	NaturalBuild     = Builtin("Natural/build")
	NaturalFold      = Builtin("Natural/fold")
	NaturalIsZero    = Builtin("Natural/isZero")
	NaturalEven      = Builtin("Natural/even")
	NaturalOdd       = Builtin("Natural/odd")
	NaturalToInteger = Builtin("Natural/toInteger")
	NaturalShow      = Builtin("Natural/show")
	NaturalSubtract  = Builtin("Natural/subtract")

	IntegerToDouble = Builtin("Integer/toDouble")
	IntegerShow     = Builtin("Integer/show")

	DoubleShow = Builtin("Double/show")

	TextShow = Builtin("Text/show")

	ListBuild   = Builtin("List/build")
	ListFold    = Builtin("List/fold")
	ListLength  = Builtin("List/length")
	ListHead    = Builtin("List/head")
	ListLast    = Builtin("List/last")
	ListIndexed = Builtin("List/indexed")
	ListReverse = Builtin("List/reverse")

	OptionalBuild = Builtin("Optional/build")
	OptionalFold  = Builtin("Optional/fold")

	True  = BoolLit(true)
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
	Var struct {
		Name  string
		Index int
	}

	// a LocalVar is an internal sentinel value used by TypeOf() in the process
	// of typechecking the body of lambdas and pis
	// FIXME: should this be unexported?
	LocalVar struct {
		Name  string
		Index int
	}

	// a QuoteVar is an internal sentinel value used by Quote() in the process
	// of converting Values back to Terms
	// FIXME: should this be unexported?
	QuoteVar struct {
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

func (LocalVar) isTerm()  {}
func (LocalVar) isValue() {}

func (QuoteVar) isValue() {}

func Bound(name string) Var {
	return Var{Name: name}
}

type (
	LambdaTerm struct {
		Label string
		Type  Term
		Body  Term
	}

	PiTerm struct {
		Label string
		Type  Term
		Body  Term
	}

	AppTerm struct {
		Fn  Term
		Arg Term
	}

	OpTerm struct {
		OpCode int
		L      Term
		R      Term
	}
)

func (LambdaTerm) isTerm() {}
func (PiTerm) isTerm()     {}
func (AppTerm) isTerm()    {}
func (OpTerm) isTerm()     {}

func Lambda(label string, t Term, body Term) LambdaTerm {
	return LambdaTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

func MkLambdaTerm(label string, t Term, body Term) LambdaTerm {
	return LambdaTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

func MkPiTerm(label string, t Term, body Term) PiTerm {
	return PiTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

// FnType returns a non-dependent function type
func FnType(domain Term, codomain Term) PiTerm {
	return PiTerm{
		Label: "_",
		Type:  domain,
		Body:  codomain,
	}
}

func Apply(fn Term, args ...Term) Term {
	out := fn
	for _, arg := range args {
		out = AppTerm{Fn: out, Arg: arg}
	}
	return out
}

// Opcodes for use in the OpTerm type
// These numbers match the binary encoding label numbers
const (
	OrOp = iota
	AndOp
	EqOp
	NeOp
	PlusOp
	TimesOp
	TextAppendOp
	ListAppendOp
	RecordMergeOp
	RightBiasedRecordMergeOp
	RecordTypeMergeOp
	ImportAltOp
	EquivOp
	CompleteOp
)

func NaturalPlus(l, r Term) Term {
	return OpTerm{OpCode: PlusOp, L: l, R: r}
}

func NaturalTimes(l, r Term) Term {
	return OpTerm{OpCode: TimesOp, L: l, R: r}
}

func BoolOr(l, r Term) Term {
	return OpTerm{OpCode: OrOp, L: l, R: r}
}

func BoolAnd(l, r Term) Term {
	return OpTerm{OpCode: AndOp, L: l, R: r}
}

func ListAppend(l, r Term) Term {
	return OpTerm{OpCode: ListAppendOp, L: l, R: r}
}

func TextAppend(l, r Term) Term {
	return OpTerm{OpCode: TextAppendOp, L: l, R: r}
}

// Callable is a function Value that can be called with one Value
// argument.  Call() may return nil if normalization isn't possible
// (for example, `Natural/even x` does not normalize)
type Callable interface {
	Value
	Call(Value) Value
}

// Call implements Callable
func (l LambdaValue) Call(a Value) Value {
	return l.Fn(a)
}

var (
	_ Callable = LambdaValue{}
)

type (
	// A LambdaValue is a go function
	LambdaValue struct {
		Label  string
		Domain Value
		Fn     func(Value) Value
	}

	// A PiValue is: the type of the domain, and a go function capturing the
	// range
	PiValue struct {
		Label  string
		Domain Value
		Range  func(Value) Value
	}

	AppValue struct {
		Fn  Value
		Arg Value
	}

	OpValue struct {
		OpCode int
		L      Value
		R      Value
	}
)

func (LambdaValue) isValue() {}

func (PiValue) isValue() {}

func (AppValue) isValue() {}

func (OpValue) isValue() {}

func MkPiVal(label string, d Value, r func(Value) Value) PiValue {
	return PiValue{
		Label:  label,
		Domain: d,
		Range:  r,
	}
}

// FnTypeVal returns a non-dependent function type Value
func FnTypeVal(l string, d Value, r Value) PiValue {
	return PiValue{
		Label:  l,
		Domain: d,
		Range:  func(Value) Value { return r },
	}
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

func MakeLet(body Term, bindings ...Binding) Let {
	return Let{Bindings: bindings, Body: body}
}

type (
	// 0, 1, 2, 3...
	NaturalLit uint

	// [] : List a
	EmptyList    struct{ Type Term }
	EmptyListVal struct{ Type Value }

	NonEmptyList    []Term
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

	BoolLit bool
	IfTerm  struct {
		Cond Term
		T    Term
		F    Term
	}
	IfVal struct {
		Cond Value
		T    Value
		F    Value
	}
	DoubleLit  float64
	IntegerLit int

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
	ToMapVal struct {
		Record Value
		Type   Value // optional
	}

	Field struct {
		Record    Term
		FieldName string
	}
	FieldVal struct {
		Record    Value
		FieldName string
	}

	Project struct {
		Record     Term
		FieldNames []string
	}
	ProjectVal struct {
		Record     Value
		FieldNames []string
	}

	ProjectType struct {
		Record   Term
		Selector Term
	}
	// no ProjectTypeVal because it cannot be in a normal form

	UnionType    map[string]Term
	UnionTypeVal map[string]Value

	Merge struct {
		Handler    Term
		Union      Term
		Annotation Term // optional
	}
	MergeVal struct {
		Handler    Value
		Union      Value
		Annotation Value // optional
	}

	Assert    struct{ Annotation Term }
	AssertVal struct{ Annotation Value }
)

func (NaturalLit) isTerm()  {}
func (NaturalLit) isValue() {}

func (EmptyList) isTerm()     {}
func (EmptyListVal) isValue() {}

func (NonEmptyList) isTerm()     {}
func (NonEmptyListVal) isValue() {}

func MakeList(first Term, rest ...Term) NonEmptyList {
	return append(NonEmptyList{first}, rest...)
}

func (TextLitTerm) isTerm() {}
func (TextLitVal) isValue() {}

func (BoolLit) isTerm()  {}
func (BoolLit) isValue() {}
func (IfTerm) isTerm()   {}
func (IfVal) isValue()   {}

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
func (ToMapVal) isValue()      {}
func (Field) isTerm()          {}
func (FieldVal) isValue()      {}
func (Project) isTerm()        {}
func (ProjectVal) isValue()    {}
func (ProjectType) isTerm()    {}
func (UnionType) isTerm()      {}
func (UnionTypeVal) isValue()  {}
func (Merge) isTerm()          {}
func (MergeVal) isValue()      {}
func (Assert) isTerm()         {}
func (AssertVal) isValue()     {}

type (
	Import struct {
		ImportHashed
		ImportMode
	}
	ImportHashed struct {
		Fetchable
		Hash []byte // stored in multihash form - ie first two bytes are 0x12 0x20
	}
	ImportMode byte
)

const (
	Code ImportMode = iota
	RawText
	Location
)

func (Import) isTerm() {}

func MakeImport(fetchable Fetchable, mode ImportMode) Import {
	return Import{
		ImportHashed: ImportHashed{
			Fetchable: fetchable,
		},
		ImportMode: mode,
	}
}
func MakeEnvVarImport(envvar string, mode ImportMode) Import {
	return MakeImport(EnvVar(envvar), mode)
}

func MakeLocalImport(path string, mode ImportMode) Import {
	return MakeImport(Local(path), mode)
}

// only for generating test data - discards errors
func MakeRemoteImport(uri string, mode ImportMode) Import {
	parsedURI, _ := url.ParseRequestURI(uri)
	remote := MakeRemote(parsedURI)
	return MakeImport(remote, mode)
}

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

func (v LocalVar) String() string {
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
