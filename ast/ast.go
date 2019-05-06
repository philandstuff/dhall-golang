package ast

import (
	"fmt"
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

	Builtin int

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

	NaturalLit int

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

	Embed Import
)

const (
	Type Const = iota
	Kind
	Sort
)

const (
	Double Builtin = iota
	Text
	Bool
	Natural
	Integer
	List
	Optional
	None
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
	_ Expr = Embed(Import{})
)

type ImportHashed struct {
	Fetchable
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
		return &App{
			Fn:  Shift(d, v, e.Fn),
			Arg: Shift(d, v, e.Arg),
		}
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
		return e
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
		return Record(fields)
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
		return &App{
			Fn:  Subst(v, c, e.Fn),
			Arg: Subst(v, c, e.Arg),
		}
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
		return e
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
		return Record(fields)
	case Embed:
		return e
	}
	panic("missing switch case in Subst()")
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
	return fmt.Sprintf("λ(%s : %v) → %v", lam.Label, lam.Type, lam.Body)
}

func (pi *Pi) String() string {
	return fmt.Sprintf("∀(%s : %v) → %v", pi.Label, pi.Type, pi.Body)
}

func (app *App) String() string {
	return fmt.Sprintf("(%v %v)", app.Fn, app.Arg)
}

func (l Let) String() string {
	panic("unimplemented")
}

func (a Annot) String() string {
	return fmt.Sprintf("%v : %v", a.Expr, a.Annotation)
}

func (t Builtin) String() string {
	switch t {
	case Double:
		return "Double"
	case Text:
		return "Text"
	case Bool:
		return "Bool"
	case Natural:
		return "Natural"
	case Integer:
		return "Integer"
	case List:
		return "List"
	case Optional:
		return "Optional"
	case None:
		return "None"
	default:
		panic(fmt.Sprintf("unknown type %d\n", t))
	}
}

func (d DoubleLit) String() string {
	return fmt.Sprintf("%f", d)
}

func (t TextLit) String() string {
	var out strings.Builder
	out.WriteString(`"`)
	for _, chunk := range t.Chunks {
		// TODO: properly deserialise string here
		out.WriteString(chunk.Prefix)
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
	case RightBiasedRecordMergeOp:
	case RecordTypeMergeOp:
	case ImportAltOp:
	default:
		panic(fmt.Sprintf("unknown opcode in Operator struct %#v", op))
	}
	return fmt.Sprintf("(%v %s %v)", op.L, opStr, op.R)
}

func (i IntegerLit) String() string {
	return fmt.Sprintf("%+d", i)
}

func (l EmptyList) String() string {
	return fmt.Sprintf("[] : %v", l.Type)
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

func (e Embed) String() string {
	panic("unimplemented")
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
	f := app.Fn.Normalize()
	a := app.Arg.Normalize()
	if l, ok := f.(*LambdaExpr); ok {
		v := Var{Name: l.Label}
		a2 := Shift(1, v, a)
		b1 := Subst(v, a2, l.Body)
		b2 := Shift(-1, v, b1)
		return b2.Normalize()
	}
	return &App{Fn: f, Arg: a}
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
			} else {
				str.WriteString(text.Suffix)
			}
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
			return NaturalLit(int(Ln) + int(Rn))
		} else if Lok && int(Ln) == 0 {
			return R
		} else if Rok && int(Rn) == 0 {
			return L
		}
	case TimesOp:
		Ln, Lok := L.(NaturalLit)
		Rn, Rok := R.(NaturalLit)

		if Lok && Rok {
			return NaturalLit(int(Ln) * int(Rn))
		} else if Lok && int(Ln) == 0 {
			return NaturalLit(0)
		} else if Lok && int(Ln) == 1 {
			return R
		} else if Rok && int(Rn) == 0 {
			return NaturalLit(0)
		} else if Rok && int(Rn) == 1 {
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
		b1 := Shift(1, Var{"_", 0}, pi.Body)
		b2 := Subst(Var{pi.Label, 0}, Var{"_", 0}, b1)
		b3 := Shift(-1, Var{"x", 0}, b2)
		return &Pi{
			Label: "_",
			Type:  pi.Type.AlphaNormalize(),
			Body:  b3.AlphaNormalize(),
		}
	}
}
func (app *App) AlphaNormalize() Expr {
	return &App{
		Fn:  app.Fn.AlphaNormalize(),
		Arg: app.Arg.AlphaNormalize(),
	}
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
