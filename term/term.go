/*
Package term contains data types for Dhall terms.

A Term is a data type which represents a Dhall expression.  Terms only
capture abstract syntax (ie how Terms are built up from each other),
not semantics (ie how to evaluate Terms).
*/
package term

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// A Term is an arbitrary Dhall expression.  When you parse text into
// Dhall, you get a value of type Term.
//
// The Term interface is just a marker; there are no interesting
// methods.
type Term interface {
	isTerm()
}

// A Universe is a type of types.
type Universe int

func (Universe) isTerm() {}

// These are the valid Universes.
const (
	Type Universe = iota
	Kind
	Sort
)

// Builtin is the type of Dhall's builtin constants.
type Builtin string

func (Builtin) isTerm() {}

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
)

// A BoolLit is a Dhall boolean literal.
type BoolLit bool

func (BoolLit) isTerm() {}

// Naturally, it is True or False.
const (
	True  BoolLit = true
	False BoolLit = false
)

// A Var is a variable, either bound or free.  A free Var is a
// valid Value; a bound Var is not.
//
// The Index is a de Bruijn index.  In an expression such as:
//
//  λ(x : Natural) → λ(x : Natural) → x@1
//
// x@1 refers to the outer bound variable x.  x@1 is represented
// by Var{"x", 1}.
type Var struct {
	Name  string
	Index int
}

func (Var) isTerm() {}

// NewVar returns a new Var.
func NewVar(name string) Var {
	return Var{Name: name}
}

// A LocalVar is an internal sentinel value used by TypeOf() in
// the process of typechecking the body of lambdas and pis.
// Essentially it lets us convert de Bruijn indices (which count
// how many binders we need to pass from the variable to the
// correct binder) to de Bruijn levels (which count how many
// binders we passed before binding this variable)
type LocalVar struct {
	Name  string
	Index int
}

func (LocalVar) isTerm() {}

type (
	// A Lambda is a lambda abstraction.
	Lambda struct {
		Label string
		Type  Term
		Body  Term
	}

	// A Pi is a function type.
	Pi struct {
		Label string
		Type  Term
		Body  Term
	}

	// An App is a Term formed by applying a function to an
	// argument.
	App struct {
		Fn  Term
		Arg Term
	}

	// An Op is two Terms combined by an operator.  The OpCode
	// determines the type of operator.
	Op struct {
		OpCode OpCode
		L      Term
		R      Term
	}
)

func (Lambda) isTerm() {}
func (Pi) isTerm()     {}
func (App) isTerm()    {}
func (Op) isTerm()     {}

// NewLambda constructs a new lambda Term.
func NewLambda(label string, t Term, body Term) Lambda {
	return Lambda{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

// NewPi constructs a new pi Term.
func NewPi(label string, t Term, body Term) Pi {
	return Pi{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

// NewAnonPi returns a pi Term with label "_", typically used for
// non-dependent function types.
func NewAnonPi(domain Term, codomain Term) Pi {
	return NewPi("_", domain, codomain)
}

// Apply takes fn and applies it successively to each arg in args.
func Apply(fn Term, args ...Term) Term {
	out := fn
	for _, arg := range args {
		out = App{Fn: out, Arg: arg}
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

// NaturalPlus takes Terms l and r and returns (l + r).
func NaturalPlus(l, r Term) Op {
	return Op{OpCode: PlusOp, L: l, R: r}
}

// NaturalTimes takes Terms l and r and returns (l * r).
func NaturalTimes(l, r Term) Op {
	return Op{OpCode: TimesOp, L: l, R: r}
}

// BoolOr takes Terms l and r and returns (l || r).
func BoolOr(l, r Term) Op {
	return Op{OpCode: OrOp, L: l, R: r}
}

// BoolAnd takes Terms l and r and returns (l && r).
func BoolAnd(l, r Term) Op {
	return Op{OpCode: AndOp, L: l, R: r}
}

// ListAppend takes Terms l and r and returns (l # r).
func ListAppend(l, r Term) Op {
	return Op{OpCode: ListAppendOp, L: l, R: r}
}

// TextAppend takes Terms l and r and returns (l ++ r).
func TextAppend(l, r Term) Op {
	return Op{OpCode: TextAppendOp, L: l, R: r}
}

// Equivalent takes Terms l and r and returns (l ≡ r).
func Equivalent(l, r Term) Op {
	return Op{OpCode: EquivOp, L: l, R: r}
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

	Annot struct {
		Expr       Term
		Annotation Term
	}
)

func (Let) isTerm()   {}
func (Annot) isTerm() {}

// NewLet returns a let Term.
func NewLet(body Term, bindings ...Binding) Let {
	return Let{Bindings: bindings, Body: body}
}

type (
	// A NaturalLit is a literal of type Natural.
	NaturalLit uint

	// An EmptyList is an empty list literal Term of the given type.
	EmptyList struct{ Type Term }

	// A NonEmptyList is a non-empty list literal Term with the given contents.
	NonEmptyList []Term

	Chunk struct {
		Prefix string
		Expr   Term
	}
	Chunks  []Chunk
	TextLit struct {
		Chunks Chunks
		Suffix string
	}

	If struct {
		Cond Term
		T    Term
		F    Term
	}

	// A DoubleLit is a literal of type Double.
	DoubleLit float64

	// A IntegerLit is a literal of type Integer.
	IntegerLit int

	// Some represents a Term which is present in an Optional type.
	Some struct{ Val Term }

	RecordType map[string]Term

	RecordLit map[string]Term

	ToMap struct {
		Record Term
		Type   Term // optional
	}

	Field struct {
		Record    Term
		FieldName string
	}

	Project struct {
		Record     Term
		FieldNames []string
	}

	ProjectType struct {
		Record   Term
		Selector Term
	}

	UnionType map[string]Term

	Merge struct {
		Handler    Term
		Union      Term
		Annotation Term // optional
	}

	Assert struct{ Annotation Term }
)

func (NaturalLit) isTerm() {}

func (EmptyList) isTerm() {}

func (NonEmptyList) isTerm() {}

// NewList returns a non-empty list Term from the given Terms.
func NewList(first Term, rest ...Term) NonEmptyList {
	return append(NonEmptyList{first}, rest...)
}

// PlainText returns an uninterpolated text literal containing the
// given string as text.
func PlainText(content string) TextLit {
	return TextLit{Suffix: content}
}

func (TextLit) isTerm() {}

func (t TextLit) String() string {
	var out strings.Builder
	out.WriteString(`"`)
	for _, chunk := range t.Chunks {
		for _, r := range chunk.Prefix {
			switch r {
			case '"':
				out.WriteString(`\"`)
			case '$':
				out.WriteString(`\u0024`)
			case '\\':
				out.WriteString(`\\`)
			case '\b':
				out.WriteString(`\b`)
			case '\f':
				out.WriteString(`\f`)
			case '\n':
				out.WriteString(`\n`)
			case '\r':
				out.WriteString(`\r`)
			case '\t':
				out.WriteString(`\t`)
			default:
				if r < 0x1f {
					out.WriteString(fmt.Sprintf(`\u%04x`, r))
				} else {
					out.WriteRune(r)
				}
			}
		}
		out.WriteString("${")
		out.WriteString(fmt.Sprint(chunk.Expr))
		out.WriteString("}")
	}
	// TODO: properly deserialise string here
	out.WriteString(t.Suffix)
	out.WriteString(`"`)
	return out.String()
}

func (If) isTerm() {}

func (DoubleLit) isTerm()  {}
func (IntegerLit) isTerm() {}

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

func (Some) isTerm() {}

func (RecordType) isTerm()  {}
func (RecordLit) isTerm()   {}
func (ToMap) isTerm()       {}
func (Field) isTerm()       {}
func (Project) isTerm()     {}
func (ProjectType) isTerm() {}
func (UnionType) isTerm()   {}
func (Merge) isTerm()       {}
func (Assert) isTerm()      {}

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
)

// An ImportMode encodes how an import should be processed once
// fetched.
type ImportMode byte

// These are the valid ImportModes.
const (
	Code     ImportMode = iota // Import as Dhall code.
	RawText                    // Import as a Text value - `as Text`.
	Location                   // Import as a Location - `as Location`.
)

func (Import) isTerm() {}

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

func (lam Lambda) String() string {
	return fmt.Sprintf("(λ(%s : %v) → %v)", lam.Label, lam.Type, lam.Body)
}

func (pi Pi) String() string {
	if pi.Label == "_" {
		return fmt.Sprintf("%v → %v", pi.Type, pi.Body)
	}
	return fmt.Sprintf("∀(%s : %v) → %v", pi.Label, pi.Type, pi.Body)
}

func (app App) String() string {
	if subApp, ok := app.Fn.(App); ok {
		return fmt.Sprintf("(%v %v)", subApp.stringNoParens(), app.Arg)
	}
	return fmt.Sprintf("(%v %v)", app.Fn, app.Arg)
}

func (app App) stringNoParens() string {
	if subApp, ok := app.Fn.(App); ok {
		return fmt.Sprintf("%v %v", subApp.stringNoParens(), app.Arg)
	}
	return fmt.Sprintf("%v %v", app.Fn, app.Arg)
}

// higher precedence binds tighter
func (op Op) precedence() int {
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

func (op Op) operatorStr() string {
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

func (op Op) String() string {
	prec := op.precedence()
	l := fmt.Sprint(op.L)
	r := fmt.Sprint(op.R)

	var buf strings.Builder

	if lop, ok := op.L.(Op); ok {
		if prec > lop.precedence() {
			buf.WriteRune('(')
		}
		buf.WriteString(l)
		if prec > lop.precedence() {
			buf.WriteRune(')')
		}
	} else if _, ok := op.L.(App); ok {
		buf.WriteString(l)

	} else {
		buf.WriteRune('(')
		buf.WriteString(l)
		buf.WriteRune(')')
	}
	buf.WriteString(op.operatorStr())
	if rop, ok := op.R.(Op); ok {
		if prec > rop.precedence() {
			buf.WriteRune('(')
		}
		buf.WriteString(r)
		if prec > rop.precedence() {
			buf.WriteRune(')')
		}
	} else if _, ok := op.R.(App); ok {
		buf.WriteString(r)

	} else {
		buf.WriteRune('(')
		buf.WriteString(r)
		buf.WriteRune(')')
	}
	return buf.String()
}

func (e EmptyList) String() string {
	return fmt.Sprintf("[] : (%v)", e.Type)
}

func (r RecordType) String() string {
	if len(r) == 0 {
		return "{}"
	}
	// get keys in sorted order
	var sortedKeys []string
	for k := range r {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var buf strings.Builder
	first := true
	buf.WriteString("{ ")
	for _, name := range sortedKeys {
		if !first {
			buf.WriteString(", ")
		}
		buf.WriteString(name)
		buf.WriteString(" : ")
		buf.WriteString(fmt.Sprintf("%v", r[name]))
		first = false
	}
	buf.WriteString("}")
	return buf.String()
}

func (r RecordLit) String() string {
	if len(r) == 0 {
		return "{=}"
	}
	// get keys in sorted order
	var sortedKeys []string
	for k := range r {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var buf strings.Builder
	first := true
	buf.WriteString("{ ")
	for _, name := range sortedKeys {
		if !first {
			buf.WriteString(", ")
		}
		buf.WriteString(name)
		buf.WriteString(" = ")
		buf.WriteString(fmt.Sprintf("%v", r[name]))
		first = false
	}
	buf.WriteString("}")
	return buf.String()
}

func (f Field) String() string {
	return fmt.Sprintf("(%v).%s", f.Record, f.FieldName)
}

func (u UnionType) String() string {
	if len(u) == 0 {
		return "<>"
	}
	var buf strings.Builder
	first := true
	buf.WriteString("< ")

	// get keys in sorted order
	var sortedKeys []string
	for k := range u {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, name := range sortedKeys {
		if !first {
			buf.WriteString(" | ")
		}
		buf.WriteString(name)
		buf.WriteString(" : ")
		buf.WriteString(fmt.Sprintf("%v", u[name]))
		first = false
	}
	buf.WriteString(" >")
	return buf.String()
}

func (m Merge) String() string {
	if m.Annotation != nil {
		return fmt.Sprintf("merge (%s) (%s) : (%s)", m.Handler, m.Union, m.Annotation)
	}
	return fmt.Sprintf("merge (%s) (%s)", m.Handler, m.Union)
}
