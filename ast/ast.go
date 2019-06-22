package ast

import (
	"fmt"
	"math"
	"net/url"
	"strings"
)

type TypeContext map[string][]Expr

func (ctx *TypeContext) Insert(name string, val Expr) *TypeContext {
	newctx := make(TypeContext)
	for k, v := range *ctx {
		newctx[k] = v
	}
	newctx[name] = append(newctx[name], val)
	return &newctx
}

func (ctx *TypeContext) Lookup(name string, i int) (Expr, bool) {
	slice := (*ctx)[name]
	if i >= len(slice) {
		return nil, false
	}
	// we index from the right of the slices
	return slice[len(slice)-1-i], true
}

//TODO: make this lazy
func (ctx *TypeContext) Map(f func(Expr) Expr) *TypeContext {
	newctx := make(TypeContext)
	for k, vs := range *ctx {
		a := make([]Expr, len(vs))
		for i, v := range vs {
			a[i] = f(v)
		}
		newctx[k] = a
	}
	return &newctx
}

func EmptyContext() *TypeContext {
	return &TypeContext{}
}

type (
	Expr interface {
		AlphaNormalize() Expr
		Normalize() Expr
		TypeWith(*TypeContext) (Expr, error)
	}

	Const int

	Var struct {
		Name  string
		Index int
	}

	LambdaExpr struct {
		Label string
		Type  Expr
		Body  Expr
	}

	Pi struct {
		Label string
		Type  Expr
		Body  Expr
	}

	App struct {
		Fn  Expr
		Arg Expr
	}

	Binding struct {
		Variable   string
		Annotation Expr // may be nil
		Value      Expr
	}
	Let struct {
		Bindings []Binding
		Body     Expr
	}

	Annot struct {
		Expr       Expr
		Annotation Expr
	}

	Builtin string

	DoubleLit float64

	Chunk struct {
		Prefix string
		Expr   Expr
	}
	Chunks  []Chunk
	TextLit struct {
		Chunks Chunks
		Suffix string
	}

	BoolLit bool
	BoolIf  struct {
		Cond Expr
		T    Expr
		F    Expr
	}

	NaturalLit uint

	IntegerLit int

	// `[] : List Natural` == EmptyList{Natural}
	EmptyList struct{ Type Expr }
	// `[2,3,4]` == NonEmptyList(2,3,4)
	NonEmptyList []Expr

	// `Some 3` == Some(3)
	Some struct{ Val Expr }
	// None is a Builtin

	Record    map[string]Expr // { x : Natural }
	RecordLit map[string]Expr // { x = 3 }

	// e.x
	Field struct {
		Record    Expr
		FieldName string
	}

	UnionType map[string]Expr // < x : Natural | y >
	Merge     struct {
		Handler    Expr
		Union      Expr
		Annotation Expr // optional
	}

	Embed Import
)

const (
	Type Const = iota
	Kind
	Sort
)

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
)

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
)

type Operator struct {
	OpCode int
	L      Expr
	R      Expr
}

func NaturalPlus(l, r Expr) Expr {
	return Operator{OpCode: PlusOp, L: l, R: r}
}

func NaturalTimes(l, r Expr) Expr {
	return Operator{OpCode: TimesOp, L: l, R: r}
}

func BoolOr(l, r Expr) Expr {
	return Operator{OpCode: OrOp, L: l, R: r}
}

func BoolAnd(l, r Expr) Expr {
	return Operator{OpCode: AndOp, L: l, R: r}
}

func ListAppend(l, r Expr) Expr {
	return Operator{OpCode: ListAppendOp, L: l, R: r}
}

func TextAppend(l, r Expr) Expr {
	return Operator{OpCode: TextAppendOp, L: l, R: r}
}

const (
	True  = BoolLit(true)
	False = BoolLit(false)
)

func MakeList(first Expr, rest ...Expr) NonEmptyList {
	return NonEmptyList(append([]Expr{first}, rest...))
}

func MakeLet(body Expr, bindings ...Binding) Let {
	return Let{Bindings: bindings, Body: body}
}

func FnType(input, output Expr) Expr {
	return &Pi{"_", input, output}
}

func MkVar(name string) Var { return Var{Name: name} }

func Apply(fn Expr, args ...Expr) Expr {
	out := fn
	for _, arg := range args {
		out = &App{Fn: out, Arg: arg}
	}
	return out
}

var (
	_ Expr = Type
	_ Expr = &Var{}
	_ Expr = &LambdaExpr{}
	_ Expr = &Pi{}
	_ Expr = &App{}
	_ Expr = Let{}
	_ Expr = Annot{}
	_ Expr = Double
	_ Expr = DoubleLit(3.0)
	_ Expr = Text
	_ Expr = TextLit{}
	_ Expr = Bool
	_ Expr = BoolLit(true)
	_ Expr = BoolIf{}
	_ Expr = Natural
	_ Expr = NaturalLit(3)
	_ Expr = Operator{}
	_ Expr = Integer
	_ Expr = IntegerLit(3)
	_ Expr = List
	_ Expr = EmptyList{Natural}
	_ Expr = NonEmptyList([]Expr{NaturalLit(3)})
	_ Expr = Optional
	_ Expr = Some{NaturalLit(3)}
	_ Expr = Record(map[string]Expr{})
	_ Expr = RecordLit(map[string]Expr{})
	_ Expr = Field{}
	_ Expr = UnionType{}
	_ Expr = Merge{}
	_ Expr = Embed(Import{})
)

type ImportHashed struct {
	Fetchable
	Hash []byte // stored in multihash form - ie first two bytes are 0x12 0x20
}

type ImportMode byte

const (
	Code ImportMode = iota
	RawText
)

type Import struct {
	ImportHashed
	ImportMode
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
	remote, _ := MakeRemote(parsedURI)
	return MakeImport(remote, mode)
}

func MakeImport(fetchable Fetchable, mode ImportMode) Import {
	return Import{
		ImportHashed: ImportHashed{
			Fetchable: fetchable,
		},
		ImportMode: mode,
	}
}

func Shift(d int, v Var, e Expr) Expr {
	switch e := e.(type) {
	case Const:
		return e
	case Var:
		if v.Name == e.Name && v.Index <= e.Index {
			if e.Index+d < 0 {
				panic("tried to shift to negative")
			}
			return Var{Name: e.Name, Index: e.Index + d}
		} else {
			return e
		}
	case *LambdaExpr:
		var body Expr
		if e.Label == v.Name {
			body = Shift(d, Var{Name: v.Name, Index: v.Index + 1}, e.Body)
		} else {
			body = Shift(d, v, e.Body)
		}
		return &LambdaExpr{
			Label: e.Label,
			Type:  Shift(d, v, e.Type),
			Body:  body,
		}
	case *Pi:
		var body Expr
		if e.Label == v.Name {
			body = Shift(d, Var{Name: v.Name, Index: v.Index + 1}, e.Body)
		} else {
			body = Shift(d, v, e.Body)
		}
		return &Pi{
			Label: e.Label,
			Type:  Shift(d, v, e.Type),
			Body:  body,
		}
	case *App:
		return Apply(
			Shift(d, v, e.Fn),
			Shift(d, v, e.Arg),
		)
	case Let:
		newBindings := make([]Binding, len(e.Bindings))
		for i, binding := range e.Bindings {
			newBindings[i].Variable = binding.Variable
			if binding.Annotation != nil {
				newBindings[i].Annotation = Shift(d, v, binding.Annotation)
			}
			newBindings[i].Value = Shift(d, v, binding.Value)
			if v.Name == binding.Variable {
				v.Index++
			}
		}
		return Let{Bindings: newBindings, Body: Shift(d, v, e.Body)}
	case Annot:
		return Annot{Shift(d, v, e.Expr), Shift(d, v, e.Annotation)}
	case Builtin:
		return e
	case DoubleLit:
		return e
	case TextLit:
		newTextLit := TextLit{make(Chunks, len(e.Chunks)), e.Suffix}
		for i, chunk := range e.Chunks {
			newTextLit.Chunks[i].Prefix = chunk.Prefix
			newTextLit.Chunks[i].Expr = Shift(d, v, chunk.Expr)
		}
		return newTextLit
	case BoolLit:
		return e
	case BoolIf:
		return BoolIf{Cond: Shift(d, v, e.Cond), T: Shift(d, v, e.T), F: Shift(d, v, e.F)}
	case NaturalLit:
		return e
	case Operator:
		return Operator{OpCode: e.OpCode, L: Shift(d, v, e.L), R: Shift(d, v, e.R)}
	case IntegerLit:
		return e
	case EmptyList:
		return EmptyList{Type: Shift(d, v, e.Type)}
	case NonEmptyList:
		exprs := make([]Expr, len([]Expr(e)))
		for i, expr := range []Expr(e) {
			exprs[i] = Shift(d, v, expr)
		}
		return NonEmptyList(exprs)
	case Some:
		return Some{Shift(d, v, e.Val)}
	case Record:
		fields := make(map[string]Expr, len(e))
		for name, val := range e {
			fields[name] = Shift(d, v, val)
		}
		return Record(fields)
	case RecordLit:
		fields := make(map[string]Expr, len(e))
		for name, val := range e {
			fields[name] = Shift(d, v, val)
		}
		return RecordLit(fields)
	case Field:
		return Field{
			Record:    Shift(d, v, e.Record),
			FieldName: e.FieldName,
		}
	case UnionType:
		fields := make(map[string]Expr, len(e))
		for name, val := range e {
			if val == nil {
				fields[name] = nil
				continue
			}
			fields[name] = Shift(d, v, val)
		}
		return UnionType(fields)
	case Merge:
		output := Merge{
			Handler: Shift(d, v, e.Handler),
			Union:   Shift(d, v, e.Union),
		}
		if e.Annotation != nil {
			output.Annotation = Shift(d, v, e.Annotation)
		}
		return output
	case Embed:
		return e
	}
	panic("missing switch case in Shift()")
}

// Subst(x, C, B) == B[x := C]
func Subst(v Var, c Expr, b Expr) Expr {
	switch e := b.(type) {
	case Const:
		return e
	case Var:
		if e == v {
			return c
		} else {
			return b
		}
	case *LambdaExpr:
		substType := Subst(v, c, e.Type)
		v2 := v
		if v.Name == e.Label {
			v2.Index++
		}
		substBody := Subst(v2, Shift(1, Var{Name: e.Label}, c), e.Body)
		return &LambdaExpr{
			Label: e.Label,
			Type:  substType,
			Body:  substBody,
		}
	case *Pi:
		substType := Subst(v, c, e.Type)
		v2 := v
		if v.Name == e.Label {
			v2.Index++
		}
		substBody := Subst(v2, Shift(1, Var{Name: e.Label}, c), e.Body)
		return &Pi{
			Label: e.Label,
			Type:  substType,
			Body:  substBody,
		}
	case *App:
		return Apply(
			Subst(v, c, e.Fn),
			Subst(v, c, e.Arg),
		)
	case Let:
		newBindings := make([]Binding, len(e.Bindings))
		for i, binding := range e.Bindings {
			newBindings[i].Variable = binding.Variable
			if binding.Annotation != nil {
				newBindings[i].Annotation = Subst(v, c, binding.Annotation)
			}
			newBindings[i].Value = Subst(v, c, binding.Value)
			if v.Name == binding.Variable {
				v.Index++
			}
		}
		return Let{Bindings: newBindings, Body: Subst(v, c, e.Body)}
	case Annot:
		return Annot{Subst(v, c, e.Expr), Subst(v, c, e.Annotation)}
	case Builtin:
		return e
	case DoubleLit:
		return e
	case TextLit:
		newTextLit := TextLit{make(Chunks, len(e.Chunks)), e.Suffix}
		for i, chunk := range e.Chunks {
			newTextLit.Chunks[i].Prefix = chunk.Prefix
			newTextLit.Chunks[i].Expr = Subst(v, c, chunk.Expr)
		}
		return newTextLit
	case BoolLit:
		return e
	case BoolIf:
		return BoolIf{Cond: Subst(v, c, e.Cond), T: Subst(v, c, e.T), F: Subst(v, c, e.F)}
	case NaturalLit:
		return e
	case Operator:
		return Operator{OpCode: e.OpCode, L: Subst(v, c, e.L), R: Subst(v, c, e.R)}
	case IntegerLit:
		return e
	case EmptyList:
		return EmptyList{Type: Subst(v, c, e.Type)}
	case NonEmptyList:
		exprs := make([]Expr, len([]Expr(e)))
		for i, expr := range []Expr(e) {
			exprs[i] = Subst(v, c, expr)
		}
		return NonEmptyList(exprs)
	case Some:
		return Some{Subst(v, c, e.Val)}
	case Record:
		fields := make(map[string]Expr, len(e))
		for name, val := range e {
			fields[name] = Subst(v, c, val)
		}
		return Record(fields)
	case RecordLit:
		fields := make(map[string]Expr, len(e))
		for name, val := range e {
			fields[name] = Subst(v, c, val)
		}
		return RecordLit(fields)
	case Field:
		return Field{
			Record:    Subst(v, c, e.Record),
			FieldName: e.FieldName,
		}
	case UnionType:
		fields := make(map[string]Expr, len(e))
		for name, val := range e {
			if val == nil {
				fields[name] = nil
				continue
			}
			fields[name] = Subst(v, c, val)
		}
		return UnionType(fields)
	case Merge:
		output := Merge{
			Handler: Subst(v, c, e.Handler),
			Union:   Subst(v, c, e.Union),
		}
		if e.Annotation != nil {
			output.Annotation = Subst(v, c, e.Annotation)
		}
		return output
	case Embed:
		return e
	}
	panic("missing switch case in Subst()")
}

func IsFreeIn(e Expr, x string) bool {
	e2 := Subst(MkVar(x), Bool, e)
	return !judgmentallyEqual(e, e2)
}

func (c Const) String() string {
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
	} else {
		return fmt.Sprintf("%s@%d", v.Name, v.Index)
	}
}

func (lam *LambdaExpr) String() string {
	return fmt.Sprintf("(λ(%s : %v) → %v)", lam.Label, lam.Type, lam.Body)
}

func (pi *Pi) String() string {
	if pi.Label == "_" {
		return fmt.Sprintf("%v → %v", pi.Type, pi.Body)
	}
	return fmt.Sprintf("∀(%s : %v) → %v", pi.Label, pi.Type, pi.Body)
}

func (app *App) String() string {
	if subApp, ok := app.Fn.(*App); ok {
		return fmt.Sprintf("(%v %v)", subApp.StringNoParens(), app.Arg)
	}
	return fmt.Sprintf("(%v %v)", app.Fn, app.Arg)
}

func (app *App) StringNoParens() string {
	if subApp, ok := app.Fn.(*App); ok {
		return fmt.Sprintf("%v %v", subApp.StringNoParens(), app.Arg)
	}
	return fmt.Sprintf("%v %v", app.Fn, app.Arg)
}

func (l Let) String() string {
	var b strings.Builder
	for _, binding := range l.Bindings {
		if binding.Annotation != nil {
			b.WriteString(fmt.Sprintf("let %v : %v = %v\n", binding.Variable, binding.Annotation, binding.Value))
		} else {
			b.WriteString(fmt.Sprintf("let %v = %v\n", binding.Variable, binding.Value))
		}
	}
	b.WriteString(fmt.Sprintf("in %v", l.Body))
	return b.String()
}

func (a Annot) String() string {
	return fmt.Sprintf("%v : %v", a.Expr, a.Annotation)
}

func (t Builtin) String() string {
	return string(t)
}

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
	} else {
		return fmt.Sprintf("%#v", float64(d))
	}
}

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

func (bl BoolLit) String() string {
	if bool(bl) {
		return "True"
	} else {
		return "False"
	}
}

func (b BoolIf) String() string {
	return fmt.Sprintf("if %v then %v else %v", b.Cond, b.T, b.F)
}

func (nl NaturalLit) String() string {
	return fmt.Sprintf("%d", nl)
}

func (op Operator) String() string {
	var opStr string
	switch op.OpCode {
	case OrOp:
		opStr = "||"
	case AndOp:
		opStr = "&&"
	case EqOp:
		opStr = "=="
	case NeOp:
		opStr = "!="
	case PlusOp:
		opStr = "+"
	case TimesOp:
		opStr = "*"
	case TextAppendOp:
		opStr = "++"
	case ListAppendOp:
		opStr = "#"
	case RecordMergeOp:
		opStr = "∧"
	case RightBiasedRecordMergeOp:
		opStr = "⫽"
	case RecordTypeMergeOp:
		opStr = "⩓"
	case ImportAltOp:
		opStr = "?"
	default:
		panic(fmt.Sprintf("unknown opcode in Operator struct %#v", op))
	}
	return fmt.Sprintf("(%v %s %v)", op.L, opStr, op.R)
}

func (i IntegerLit) String() string {
	return fmt.Sprintf("%+d", i)
}

func (l EmptyList) String() string {
	return fmt.Sprintf("[] : List %v", l.Type)
}

func (l NonEmptyList) String() string {
	var out strings.Builder
	out.WriteString("[ ")
	exprs := []Expr(l)
	out.WriteString(fmt.Sprint(exprs[0]))
	for _, expr := range exprs[1:] {
		out.WriteString(", ")
		out.WriteString(fmt.Sprint(expr))
	}
	out.WriteString(" ]")
	return out.String()
}

func (s Some) String() string {
	return fmt.Sprintf("Some %v", s.Val)
}

func (r Record) String() string {
	var out strings.Builder
	out.WriteString("{ ")
	fields := map[string]Expr(r)
	first := true
	for name, expr := range fields {
		if !first {
			out.WriteString(", ")
		}
		first = false
		out.WriteString(fmt.Sprintf("%s : %v", name, expr))
	}
	out.WriteString(" }")
	return out.String()
}

func (r RecordLit) String() string {
	var out strings.Builder
	out.WriteString("{ ")
	fields := map[string]Expr(r)
	first := true
	for name, expr := range fields {
		if !first {
			out.WriteString(", ")
		}
		first = false
		out.WriteString(fmt.Sprintf("%s = %v", name, expr))
	}
	out.WriteString(" }")
	return out.String()
}

func (f Field) String() string {
	return fmt.Sprintf("%v.%s", f.Record, f.FieldName)
}

func (u UnionType) String() string {
	var out strings.Builder
	out.WriteString("< ")
	fields := map[string]Expr(u)
	first := true
	for name, expr := range fields {
		if !first {
			out.WriteString(" | ")
		}
		first = false
		if expr == nil {
			out.WriteString(name)
		} else {
			out.WriteString(fmt.Sprintf("%s : %v", name, expr))
		}
	}
	out.WriteString(" >")
	return out.String()
}

func (m Merge) String() string {
	if m.Annotation != nil {
		return fmt.Sprintf("merge %s %s : %s", m.Handler, m.Union, m.Annotation)
	} else {
		return fmt.Sprintf("merge %s %s", m.Handler, m.Union)
	}
}

func (e Embed) String() string {
	return e.Fetchable.String()
}

func (c Const) Normalize() Expr { return c }
func (v Var) Normalize() Expr   { return v }

func (lam *LambdaExpr) Normalize() Expr {
	return &LambdaExpr{
		Label: lam.Label,
		Type:  lam.Type.Normalize(),
		Body:  lam.Body.Normalize(),
	}
}
func (pi *Pi) Normalize() Expr {
	return &Pi{
		Label: pi.Label,
		Type:  pi.Type.Normalize(),
		Body:  pi.Body.Normalize(),
	}
}
func (app *App) Normalize() Expr {
	f0 := app.Fn.Normalize()
	a0 := app.Arg.Normalize()
	switch f1 := f0.(type) {
	case *LambdaExpr:
		v := MkVar(f1.Label)
		a2 := Shift(1, v, a0)
		b1 := Subst(v, a2, f1.Body)
		b2 := Shift(-1, v, b1)
		return b2.Normalize()
	case *App:
		switch f2 := f1.Fn.(type) {
		case *App:
			switch f3 := f2.Fn.(type) {
			case *App:
				switch f4 := f3.Fn.(type) {
				case *App:
					switch f5 := f4.Fn.(type) {
					case Builtin: // five args: (Foo/bar f4.Arg f3.Arg f2.Arg f1.Arg a0)
						switch f5 {
						case ListFold:
							cons := f1.Arg
							output := a0
							switch list := f3.Arg.(type) {
							case NonEmptyList:
								for i := len(list) - 1; i >= 0; i-- {
									output = Apply(cons, list[i], output).Normalize()
								}
								return output
							case EmptyList:
								return output
							}
						case OptionalFold:
							input := f3.Arg
							someFn := f1.Arg
							noneVal := a0
							if someVal, ok := input.(Some); ok {
								return Apply(someFn, someVal.Val).Normalize()
							}
							if appVal, ok := input.(*App); ok {
								if appVal.Fn == None {
									return noneVal
								}
							}
						}
					}
				case Builtin: // four args: (Foo/bar f3.Arg f2.Arg f1.Arg a0)
					switch f4 {
					case NaturalFold:
						if n, ok := f3.Arg.(NaturalLit); ok {
							output := a0
							for i := 0; i < int(n); i++ {
								output = Apply(f1.Arg, output).Normalize()
							}
							return output
						}
					}

				}
			}
		case Builtin: // two args: (Foo/bar f1.Arg a0)
			switch f2 {
			case ListBuild:
				if ap0, ok := a0.(*App); ok {
					if ap1, ok := ap0.Fn.(*App); ok {
						if ap1.Fn == ListFold {
							// List/build A₀ (List/fold A₁ b) → b
							return ap0.Arg
						}
					}
				}
				A0 := f1.Arg
				g := a0
				A1 := Shift(1, Var{"a", 0}, A0)
				return Apply(
					g,
					Apply(List, A0),
					&LambdaExpr{"a", A0,
						&LambdaExpr{"as", Apply(List, A1),
							ListAppend(
								MakeList(MkVar("a")),
								MkVar("as"),
							)}},
					EmptyList{A0},
				).Normalize()
			case ListReverse:
				switch list := a0.(type) {
				case EmptyList:
					return EmptyList{f1.Arg}
				case NonEmptyList:
					output := make([]Expr, len(list))
					for i, a := range list {
						output[len(list)-i-1] = a
					}
					return NonEmptyList(output)
				}
			case ListLength:
				switch list := a0.(type) {
				case EmptyList:
					return NaturalLit(0)
				case NonEmptyList:
					return NaturalLit(len(list))
				}
			case ListHead, ListLast:
				switch list := a0.(type) {
				case EmptyList:
					return Apply(None, f1.Arg)
				case NonEmptyList:
					if f2 == ListHead {
						return Some{list[0]}
					} else {
						return Some{list[len(list)-1]}
					}
				}
			case ListIndexed:
				switch list := a0.(type) {
				case EmptyList:
					return EmptyList{Record{"index": Natural, "value": f1.Arg}}
				case NonEmptyList:
					output := make([]Expr, len(list))
					for i, a := range list {
						output[i] = RecordLit{"index": NaturalLit(i), "value": a}
					}
					return NonEmptyList(output)
				}
			case OptionalBuild:
				if ap0, ok := a0.(*App); ok {
					if ap1, ok := ap0.Fn.(*App); ok {
						if ap1.Fn == OptionalFold {
							// Optional/build A₀ (Optional/fold A₁ b) → b
							return ap0.Arg
						}
					}
				}
				A0 := f1.Arg
				g := a0
				return Apply(
					g,
					Apply(Optional, A0),
					&LambdaExpr{"a", A0, Some{MkVar("a")}},
					Apply(None, A0),
				).Normalize()
			}
		}
	case Builtin: // one arg: (Foo/bar a0)
		switch f1 {
		case NaturalBuild:
			if ap, ok := a0.(*App); ok {
				if ap.Fn == NaturalFold {
					// Natural/build (Natural/fold b) → b
					return ap.Arg
				}
			}
			if l, ok := a0.(*LambdaExpr); ok {
				return Apply(l,
					Natural,
					&LambdaExpr{"x", Natural, NaturalPlus(MkVar("x"), NaturalLit(1))},
					NaturalLit(0),
				).Normalize()
			}
		case NaturalIsZero:
			if n, ok := a0.(NaturalLit); ok {
				if n == 0 {
					return True
				} else {
					return False
				}
			}
		case NaturalEven:
			if n, ok := a0.(NaturalLit); ok {
				if n%2 == 0 {
					return True
				} else {
					return False
				}
			}
		case NaturalOdd:
			if n, ok := a0.(NaturalLit); ok {
				if n%2 == 0 {
					return False
				} else {
					return True
				}
			}
		case NaturalToInteger:
			if n, ok := a0.(NaturalLit); ok {
				return IntegerLit(int(n))
			}
		case NaturalShow:
			if n, ok := a0.(NaturalLit); ok {
				return TextLit{Suffix: n.String()}
			}
		case IntegerToDouble:
			if i, ok := a0.(IntegerLit); ok {
				return DoubleLit(float64(i))
			}
		case IntegerShow:
			if i, ok := a0.(IntegerLit); ok {
				return TextLit{Suffix: i.String()}
			}
		case DoubleShow:
			if d, ok := a0.(DoubleLit); ok {
				return TextLit{Suffix: d.String()}
			}
		case TextShow:
			if t, ok := a0.(TextLit); ok {
				if t.Chunks == nil || len(t.Chunks) == 0 {
					var out strings.Builder
					out.WriteRune('"')
					for _, r := range t.Suffix {
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
					out.WriteRune('"')
					return TextLit{Suffix: out.String()}
				}
			}
		}
	}
	return Apply(f0, a0)
}

func (l Let) Normalize() Expr {
	binding := l.Bindings[0]
	x := binding.Variable
	a1 := binding.Value.Normalize()
	a2 := Shift(1, Var{x, 0}, a1)

	rest := l.Body
	if len(l.Bindings) > 1 {
		rest = Let{Bindings: l.Bindings[1:], Body: l.Body}
	}
	rest = rest.Normalize()

	b1 := Subst(Var{x, 0}, a2, rest)
	b2 := Shift(-1, Var{x, 0}, b1)
	return b2.Normalize()
}

func (a Annot) Normalize() Expr { return a.Expr.Normalize() }

func (t Builtin) Normalize() Expr { return t }

func (d DoubleLit) Normalize() Expr { return d }

func (t TextLit) Normalize() Expr {
	var str strings.Builder
	var newChunks Chunks
	for _, chunk := range t.Chunks {
		str.WriteString(chunk.Prefix)
		normExpr := chunk.Expr.Normalize()
		if text, ok := normExpr.(TextLit); ok {
			if len(text.Chunks) != 0 {
				// first chunk gets the rest of str
				str.WriteString(text.Chunks[0].Prefix)
				newChunks = append(newChunks,
					Chunk{Prefix: str.String(), Expr: text.Chunks[0].Expr})
				newChunks = append(newChunks,
					text.Chunks[1:]...)
				str.Reset()
			}
			str.WriteString(text.Suffix)

		} else {
			newChunks = append(newChunks, Chunk{Prefix: str.String(), Expr: normExpr})
			str.Reset()
		}
	}
	str.WriteString(t.Suffix)
	newSuffix := str.String()

	// Special case: "${<expr>}" → <expr>
	if len(newChunks) == 1 && newChunks[0].Prefix == "" && newSuffix == "" {
		return newChunks[0].Expr
	}

	return TextLit{Chunks: newChunks, Suffix: newSuffix}
}

func (n BoolLit) Normalize() Expr { return n }
func (b BoolIf) Normalize() Expr {
	cond := b.Cond.Normalize()
	t := b.T.Normalize()
	f := b.F.Normalize()
	if cond == True {
		return t
	}
	if cond == False {
		return f
	}
	if t == True && f == False {
		return cond
	}
	if judgmentallyEqual(t, f) {
		return t
	}
	return BoolIf{Cond: cond, T: t, F: f}
}

func (n NaturalLit) Normalize() Expr { return n }

func mustMergeRecordLits(l RecordLit, r RecordLit) RecordLit {
	output := make(RecordLit)
	for k, v := range l {
		output[k] = v
	}
	for k, v := range r {
		if lField, ok := output[k]; ok {
			lSubrecord, Lok := lField.(RecordLit)
			rSubrecord, Rok := v.(RecordLit)
			if !(Lok && Rok) {
				// typecheck ought to have caught this
				panic("Record mismatch")
			}
			output[k] = mustMergeRecordLits(lSubrecord, rSubrecord)
		} else {
			output[k] = v
		}
	}
	return output
}

func (op Operator) Normalize() Expr {
	L := op.L.Normalize()
	R := op.R.Normalize()

	switch op.OpCode {
	case OrOp:
		Lb, Lok := L.(BoolLit)
		Rb, Rok := R.(BoolLit)

		if Lok && Lb == False {
			return R
		} else if Rok && Rb == False {
			return L
		} else if Lok && Lb == True {
			return True
		} else if Rok && Rb == True {
			return True
		}
		if judgmentallyEqual(L, R) {
			return L
		}
	case AndOp:
		Lb, Lok := L.(BoolLit)
		Rb, Rok := R.(BoolLit)

		if Lok && Lb == True {
			return R
		} else if Rok && Rb == True {
			return L
		} else if Lok && Lb == False {
			return False
		} else if Rok && Rb == False {
			return False
		}
		if judgmentallyEqual(L, R) {
			return L
		}
	case EqOp:
		Lb, Lok := L.(BoolLit)
		Rb, Rok := R.(BoolLit)

		if Lok && Lb == True {
			return R
		} else if Rok && Rb == True {
			return L
		}
		if judgmentallyEqual(L, R) {
			return True
		}
	case NeOp:
		Lb, Lok := L.(BoolLit)
		Rb, Rok := R.(BoolLit)

		if Lok && Lb == False {
			return R
		} else if Rok && Rb == False {
			return L
		}
		if judgmentallyEqual(L, R) {
			return False
		}
	case PlusOp:
		Ln, Lok := L.(NaturalLit)
		Rn, Rok := R.(NaturalLit)

		if Lok && Rok {
			return NaturalLit(uint(Ln) + uint(Rn))
		} else if Lok && uint(Ln) == 0 {
			return R
		} else if Rok && uint(Rn) == 0 {
			return L
		}
	case TimesOp:
		Ln, Lok := L.(NaturalLit)
		Rn, Rok := R.(NaturalLit)

		if Lok && Rok {
			return NaturalLit(uint(Ln) * uint(Rn))
		} else if Lok && uint(Ln) == 0 {
			return NaturalLit(0)
		} else if Lok && uint(Ln) == 1 {
			return R
		} else if Rok && uint(Rn) == 0 {
			return NaturalLit(0)
		} else if Rok && uint(Rn) == 1 {
			return L
		}
	case TextAppendOp:
		return TextLit{Chunks: Chunks{{Expr: L}, {Expr: R}}}.Normalize()
	case ListAppendOp:
		_, lEmpty := L.(EmptyList)
		if lEmpty {
			return R
		}
		_, rEmpty := R.(EmptyList)
		if rEmpty {
			return L
		}

		lList, lOk := L.(NonEmptyList)
		rList, rOk := R.(NonEmptyList)
		if lOk && rOk {
			return NonEmptyList(append(lList, rList...))
		}
	case RecordMergeOp:
		Lr, Lok := L.(RecordLit)
		Rr, Rok := R.(RecordLit)

		if Lok && len(Lr) == 0 {
			return R
		} else if Rok && len(Rr) == 0 {
			return L
		} else if Lok && Rok {
			return mustMergeRecordLits(Lr, Rr)
		}
	case RightBiasedRecordMergeOp:
		Lr, Lok := L.(RecordLit)
		Rr, Rok := R.(RecordLit)

		if Lok && len(Lr) == 0 {
			return R
		} else if Rok && len(Rr) == 0 {
			return L
		} else if Lok && Rok {
			output := make(RecordLit)
			for k, v := range Lr {
				output[k] = v
			}
			for k, v := range Rr {
				output[k] = v
			}
			return output
		}
	case RecordTypeMergeOp:
		Lr, Lok := L.(Record)
		Rr, Rok := R.(Record)

		if Lok && len(Lr) == 0 {
			return R
		} else if Rok && len(Rr) == 0 {
			return L
		} else if Lok && Rok {
			out, err := mergeRecords(Lr, Rr)
			if err != nil {
				panic(err)
			}
			return out
		}
	}
	return Operator{OpCode: op.OpCode, L: L, R: R}
}

func (i IntegerLit) Normalize() Expr { return i }

func (l EmptyList) Normalize() Expr { return EmptyList{l.Type.Normalize()} }
func (l NonEmptyList) Normalize() Expr {
	exprs := []Expr(l)
	vals := make([]Expr, len(exprs))
	for i, expr := range exprs {
		vals[i] = expr.Normalize()
	}
	return NonEmptyList(vals)
}

func (s Some) Normalize() Expr {
	return Some{s.Val.Normalize()}
}

func (r Record) Normalize() Expr {
	fields := make(map[string]Expr, len(r))
	for name, val := range r {
		fields[name] = val.Normalize()
	}
	return Record(fields)
}

func (r RecordLit) Normalize() Expr {
	fields := make(map[string]Expr, len(r))
	for name, val := range r {
		fields[name] = val.Normalize()
	}
	return RecordLit(fields)
}

func (f Field) Normalize() Expr {
	r := f.Record.Normalize()
	if rl, ok := r.(RecordLit); ok {
		val := rl[f.FieldName]
		return val.Normalize()
	}
	return Field{Record: r, FieldName: f.FieldName}
}

func (u UnionType) Normalize() Expr {
	fields := make(map[string]Expr, len(u))
	for name, val := range u {
		if val == nil {
			// empty alternative
			fields[name] = nil
			continue
		}
		fields[name] = val.Normalize()
	}
	return UnionType(fields)
}

func (m Merge) Normalize() Expr {
	handlerN := m.Handler.Normalize()
	unionN := m.Union.Normalize()
	if handlers, ok := handlerN.(RecordLit); ok {
		if unionVal, ok := unionN.(*App); ok {
			// do we have a union alternative with a value, or a
			// subexpression?
			if field, ok := unionVal.Fn.(Field); ok {
				return Apply(
					handlers[field.FieldName],
					unionVal.Arg,
				).Normalize()
			}
		}
		if unionVal, ok := unionN.(Field); ok {
			// we have an empty union alternative
			return handlers[unionVal.FieldName].Normalize()
		}
	}
	output := Merge{
		Handler: handlerN,
		Union:   unionN,
	}
	if m.Annotation != nil {
		output.Annotation = m.Annotation.Normalize()
	}
	return output
}

func (e Embed) Normalize() Expr {
	panic("Can't normalize an expression with unresolved imports")
}

func (c Const) AlphaNormalize() Expr { return c }
func (v Var) AlphaNormalize() Expr   { return v }

func (lam *LambdaExpr) AlphaNormalize() Expr {
	if lam.Label == "_" {
		return &LambdaExpr{
			Label: "_",
			Type:  lam.Type.AlphaNormalize(),
			Body:  lam.Body.AlphaNormalize(),
		}
	} else {
		b1 := Shift(1, Var{"_", 0}, lam.Body)
		b2 := Subst(Var{lam.Label, 0}, Var{"_", 0}, b1)
		b3 := Shift(-1, Var{"x", 0}, b2)
		return &LambdaExpr{
			Label: "_",
			Type:  lam.Type.AlphaNormalize(),
			Body:  b3.AlphaNormalize(),
		}
	}
}
func (pi *Pi) AlphaNormalize() Expr {
	if pi.Label == "_" {
		return &Pi{
			Label: "_",
			Type:  pi.Type.AlphaNormalize(),
			Body:  pi.Body.AlphaNormalize(),
		}
	} else {
		B1 := Shift(1, Var{"_", 0}, pi.Body)
		B2 := Subst(Var{pi.Label, 0}, Var{"_", 0}, B1)
		B3 := Shift(-1, Var{pi.Label, 0}, B2)
		return &Pi{
			Label: "_",
			Type:  pi.Type.AlphaNormalize(),
			Body:  B3.AlphaNormalize(),
		}
	}
}
func (app *App) AlphaNormalize() Expr {
	return Apply(
		app.Fn.AlphaNormalize(),
		app.Arg.AlphaNormalize(),
	)
}

func (l Let) AlphaNormalize() Expr {
	binding := l.Bindings[0]
	if binding.Annotation != nil {
		binding.Annotation = binding.Annotation.AlphaNormalize()
	}
	x := binding.Variable
	if x == "_" {
		binding.Value = binding.Value.AlphaNormalize()
	} else {
		v1 := Shift(1, Var{"_", 0}, binding.Value)
		v2 := Subst(Var{x, 0}, Var{"_", 0}, v1)
		v3 := Shift(-1, Var{x, 0}, v2)
		binding.Value = v3.AlphaNormalize()
		binding.Variable = "_"
	}

	rest := l.Body
	if len(l.Bindings) > 1 {
		rest = Let{Bindings: l.Bindings[1:], Body: l.Body}
	}
	rest = rest.AlphaNormalize()

	b0 := Shift(1, Var{"_", 0}, rest)
	b1 := Subst(Var{x, 0}, Var{"_", 0}, b0)
	b2 := Shift(-1, Var{x, 0}, b1)
	b3 := b2.AlphaNormalize()
	b4, ok := b3.(Let)
	if ok {
		return Let{
			Bindings: append([]Binding{binding}, b4.Bindings...),
			Body:     b4.Body,
		}
	} else {
		return Let{
			Bindings: []Binding{binding},
			Body:     b3,
		}
	}
}

func (a Annot) AlphaNormalize() Expr { return a.Expr.AlphaNormalize() }

func (t Builtin) AlphaNormalize() Expr { return t }

func (d DoubleLit) AlphaNormalize() Expr { return d }

func (t TextLit) AlphaNormalize() Expr {
	newTextLit := TextLit{make(Chunks, len(t.Chunks)), t.Suffix}
	for i, chunk := range t.Chunks {
		newTextLit.Chunks[i].Prefix = chunk.Prefix
		newTextLit.Chunks[i].Expr = chunk.Expr.AlphaNormalize()
	}
	return newTextLit
}

func (n BoolLit) AlphaNormalize() Expr { return n }
func (b BoolIf) AlphaNormalize() Expr {
	return BoolIf{
		Cond: b.Cond.AlphaNormalize(),
		T:    b.T.AlphaNormalize(),
		F:    b.F.AlphaNormalize(),
	}
}

func (n NaturalLit) AlphaNormalize() Expr { return n }
func (op Operator) AlphaNormalize() Expr {
	L := op.L.AlphaNormalize()
	R := op.R.AlphaNormalize()

	return Operator{OpCode: op.OpCode, L: L, R: R}
}
func (i IntegerLit) AlphaNormalize() Expr { return i }

func (l EmptyList) AlphaNormalize() Expr { return EmptyList{l.Type.AlphaNormalize()} }
func (l NonEmptyList) AlphaNormalize() Expr {
	exprs := []Expr(l)
	vals := make([]Expr, len(exprs))
	for i, expr := range exprs {
		vals[i] = expr.AlphaNormalize()
	}
	return NonEmptyList(vals)
}

func (s Some) AlphaNormalize() Expr {
	return Some{s.Val.AlphaNormalize()}
}

func (r Record) AlphaNormalize() Expr {
	fields := make(map[string]Expr, len(r))
	for name, val := range r {
		fields[name] = val.AlphaNormalize()
	}
	return Record(fields)
}

func (r RecordLit) AlphaNormalize() Expr {
	fields := make(map[string]Expr, len(r))
	for name, val := range r {
		fields[name] = val.AlphaNormalize()
	}
	return RecordLit(fields)
}

func (f Field) AlphaNormalize() Expr {
	return Field{
		Record:    f.Record.AlphaNormalize(),
		FieldName: f.FieldName,
	}
}

func (u UnionType) AlphaNormalize() Expr {
	fields := make(map[string]Expr, len(u))
	for name, val := range u {
		if val == nil {
			// empty alternative
			fields[name] = nil
			continue
		}
		fields[name] = val.AlphaNormalize()
	}
	return UnionType(fields)
}

func (m Merge) AlphaNormalize() Expr {
	output := Merge{
		Handler: m.Handler.AlphaNormalize(),
		Union:   m.Union.AlphaNormalize(),
	}
	if m.Annotation != nil {
		output.Annotation = m.Annotation.Normalize()
	}
	return output
}

func (e Embed) AlphaNormalize() Expr {
	panic("Can't normalize an expression with unresolved imports")
}

func NewLambdaExpr(arg string, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
