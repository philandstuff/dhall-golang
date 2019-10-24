package core

import (
	"fmt"
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
	NaturalEvenVal      struct{}
	NaturalFoldVal      struct{}
	NaturalIsZeroVal    struct{}
	NaturalOddVal       struct{}
	NaturalShowVal      struct{}
	NaturalSubtractVal  struct{}
	NaturalToIntegerVal struct{}
)

func (NaturalEvenVal) isValue()      {}
func (NaturalFoldVal) isValue()      {}
func (NaturalIsZeroVal) isValue()    {}
func (NaturalOddVal) isValue()       {}
func (NaturalShowVal) isValue()      {}
func (NaturalSubtractVal) isValue()  {}
func (NaturalToIntegerVal) isValue() {}

// Call1 implements Callable1
func (f NaturalEvenVal) Call1(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 0)
	}
	return AppValue{Fn: f, Arg: x}
}

// Call4 implements Callable4
func (f NaturalFoldVal) Call4(n, T, s, z Value) Value {
	if n, ok := n.(NaturalLit); ok {
		result := z
		for i := 0; i < int(n); i++ {
			if succ, ok := s.(Callable1); ok {
				result = succ.Call1(result)
			} else {
				result = AppValue{s, result}
			}
		}
		return result
	}
	return applyVal(f, n, T, s, z)
}

// Call1 implements Callable1
func (f NaturalIsZeroVal) Call1(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n == 0)
	}
	return AppValue{Fn: f, Arg: x}
}

// Call1 implements Callable1
func (f NaturalOddVal) Call1(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 1)
	}
	return AppValue{Fn: f, Arg: x}
}

// Call1 implements Callable1
func (f NaturalShowVal) Call1(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return TextLitVal{Suffix: fmt.Sprintf("%d", n)}
	}
	return AppValue{Fn: f, Arg: x}
}

// Call2 implements Callable2
func (f NaturalSubtractVal) Call2(a, b Value) Value {
	m, mok := a.(NaturalLit)
	n, nok := b.(NaturalLit)
	if mok && nok {
		if n >= m {
			return NaturalLit(n - m)
		}
		return NaturalLit(0)
	}
	if a == NaturalLit(0) {
		return b
	}
	if b == NaturalLit(0) {
		return NaturalLit(0)
	}
	if judgmentallyEqualVals(a, b) {
		return NaturalLit(0)
	}
	return applyVal(f, a, b)
}

// Call1 implements Callable1
func (f NaturalToIntegerVal) Call1(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return IntegerLit(n)
	}
	return AppValue{Fn: f, Arg: x}
}

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

func applyVal(fn Value, args ...Value) Value {
	out := fn
	for _, arg := range args {
		out = AppValue{Fn: out, Arg: arg}
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
// argument.
type Callable1 interface {
	Value
	Call1(Value) Value
}

// Callable2 is a function Value that can be called with one Value
// argument.
type Callable2 interface {
	Value
	Call2(Value, Value) Value
}

// Callable4 is a function Value that can be called with three Value
// arguments.
type Callable4 interface {
	Value
	Call4(Value, Value, Value, Value) Value
}

var (
	_ Callable1 = LambdaValue{}
	_ Callable1 = NaturalEvenVal{}
	_ Callable1 = NaturalIsZeroVal{}
	_ Callable1 = NaturalOddVal{}
	_ Callable1 = NaturalShowVal{}
	_ Callable1 = NaturalToIntegerVal{}

	_ Callable2 = NaturalSubtractVal{}

	_ Callable4 = NaturalFoldVal{}
)

type (
	// a lambda value is a go function
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
)

func (LambdaValue) isValue() {}

// Call1 implements Callable1
func (l LambdaValue) Call1(x Value) Value { return l.Fn(x) }

func (PiValue) isValue() {}

func (AppValue) isValue() {}

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
