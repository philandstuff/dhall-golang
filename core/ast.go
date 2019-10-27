package core

import (
	"fmt"
	"math"
	"net/url"
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
	naturalBuildVal     struct{ hasCall1 }
	naturalEvenVal      struct{ hasCall1 }
	naturalFoldVal      struct{ hasCall4 }
	naturalIsZeroVal    struct{ hasCall1 }
	naturalOddVal       struct{ hasCall1 }
	naturalShowVal      struct{ hasCall1 }
	naturalSubtractVal  struct{ hasCall2 }
	naturalToIntegerVal struct{ hasCall1 }

	integerShowVal     struct{ hasCall1 }
	integerToDoubleVal struct{ hasCall1 }

	doubleShowVal struct{ hasCall1 }

	optionalBuildVal struct{ hasCall2 }
	optionalFoldVal  struct{ hasCall5 }

	textShowVal struct{ hasCall1 }

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
	listLengthVal  struct{ hasCall2 }
	listHeadVal    struct{ hasCall2 }
	listLastVal    struct{ hasCall2 }
	listIndexedVal struct{ hasCall2 }
	listReverseVal struct{ hasCall2 }
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
	BoundVar struct {
		Name  string
		Index int
	}

	FreeVar struct {
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

func (BoundVar) isTerm() {}

func (FreeVar) isTerm()  {}
func (FreeVar) isValue() {}

func (LocalVar) isTerm()  {}
func (LocalVar) isValue() {}

func (QuoteVar) isValue() {}

func Bound(name string) BoundVar {
	return BoundVar{Name: name}
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

func Mkλ(label string, t Term, body Term) LambdaTerm {
	return LambdaTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

func MkΠ(label string, t Term, body Term) PiTerm {
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

// Callable1 is a function Value that can be called with one Value
// argument.  Call1() may return nil if normalization isn't possible
// (for example, `Natural/even x` does not normalize)
type Callable1 interface {
	Callable2
	Call1(Value) Value
}

type hasCall1 func(Value) Value

// Call1 implements Callable1
func (s hasCall1) Call1(a Value) Value {
	return s(a)
}

// Call2 implements Callable2
func (s hasCall1) Call2(a, b Value) Value {
	return s.Call1(a).(Callable1).Call1(b)
}

// Call3 implements Callable3
func (s hasCall1) Call3(a, b, c Value) Value {
	return s.Call1(a).(Callable1).Call1(b).(Callable1).Call1(c)
}

// Call4 implements Callable4
func (s hasCall1) Call4(a, b, c, d Value) Value {
	return s.Call1(a).(Callable1).Call1(b).(Callable1).Call1(c).(Callable1).Call1(d)
}

// Call5 implements Callable5
func (s hasCall1) Call5(a, b, c, d, e Value) Value {
	return s.Call1(a).(Callable1).Call1(b).(Callable1).Call1(c).(Callable1).Call1(d).(Callable1).Call1(e)
}

// Callable2 is a function Value that can be called with one Value
// argument.  Call2() may return nil if normalization isn't possible
// (for example, `Natural/subtract x y` does not normalize)
type Callable2 interface {
	Callable3
	Call2(Value, Value) Value
}

type hasCall2 func(Value, Value) Value

// Call2 implements Callable2
func (s hasCall2) Call2(a, b Value) Value {
	return s(a, b)
}

// Call3 implements Callable3
func (s hasCall2) Call3(a, b, c Value) Value {
	return s.Call2(a, b).(Callable1).Call1(c)
}

// Call4 implements Callable4
func (s hasCall2) Call4(a, b, c, d Value) Value {
	return s.Call2(a, b).(Callable1).Call1(c).(Callable1).Call1(d)
}

// Call5 implements Callable5
func (s hasCall2) Call5(a, b, c, d, e Value) Value {
	return s.Call2(a, b).(Callable1).Call1(c).(Callable1).Call1(d).(Callable1).Call1(e)
}

// Callable3 is a function Value that can be called with one Value
// argument.  Call3() may return nil if normalization isn't possible.
type Callable3 interface {
	Callable4
	Call3(Value, Value, Value) Value
}

type hasCall3 func(Value, Value, Value) Value

// Call3 implements Callable3
func (s hasCall3) Call3(a, b, c Value) Value {
	return s(a, b, c)
}

// Call4 implements Callable4
func (s hasCall3) Call4(a, b, c, d Value) Value {
	return s.Call3(a, b, c).(Callable1).Call1(d)
}

// Call5 implements Callable5
func (s hasCall3) Call5(a, b, c, d, e Value) Value {
	return s.Call3(a, b, c).(Callable1).Call1(d).(Callable1).Call1(e)
}

// Callable4 is a function Value that can be called with one Value
// argument.  Call4() may return nil if normalization isn't possible.
type Callable4 interface {
	Callable5
	Call4(Value, Value, Value, Value) Value
}

type hasCall4 func(Value, Value, Value, Value) Value

// Call4 implements Callable4
func (s hasCall4) Call4(a, b, c, d Value) Value {
	return s(a, b, c, d)
}

// Call5 implements Callable5
func (s hasCall4) Call5(a, b, c, d, e Value) Value {
	return s(a, b, c, d).(Callable1).Call1(e)
}

// Callable5 is a function Value that can be called with one Value
// argument.  Call5() may return nil if normalization isn't possible.
type Callable5 interface {
	Value
	Call5(Value, Value, Value, Value, Value) Value
}

type hasCall5 func(Value, Value, Value, Value, Value) Value

// Call5 implements Callable5
func (s hasCall5) Call5(a, b, c, d, e Value) Value {
	return s(a, b, c, d, e)
}

var (
	_ Callable1 = LambdaValue{}
	_ Callable2 = LambdaValue{}
	_ Callable5 = LambdaValue{}
)

type (
	// A LambdaValue is a go function
	LambdaValue struct {
		hasCall1
		Label  string
		Domain Value
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

func MkΠval(label string, d Value, r func(Value) Value) PiValue {
	return PiValue{
		Label:  label,
		Domain: d,
		Range:  r,
	}
}

// FnTypeVal returns a non-dependent function type Value
func FnTypeVal(d Value, r Value) PiValue {
	return PiValue{
		Label:  "_",
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
		Annotaiton Value // optional
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

func (v BoundVar) String() string {
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
