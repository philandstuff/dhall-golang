package ast

import (
	"errors"
	"fmt"
	"io"
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
		io.WriterTo
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

	BuiltinType int

	DoubleLit float64

	BoolLit bool
	BoolIf  struct {
		Cond Expr
		T    Expr
		F    Expr
	}

	NaturalLit  int
	NaturalPlus struct {
		L Expr
		R Expr
	}
	NaturalTimes struct {
		L Expr
		R Expr
	}

	IntegerLit int

	// `[] : List Natural` == EmptyList{Natural}
	EmptyList struct{ Type Expr }
	// `[2,3,4]` == NonEmptyList(2,3,4)
	NonEmptyList []Expr

	Record    map[string]Expr
	RecordLit map[string]Expr
)

const (
	Type Const = iota
	Kind
	Sort
)

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
	case BuiltinType:
		return e
	case DoubleLit:
		return e
	case BoolLit:
		return e
	case BoolIf:
		return BoolIf{Cond: Shift(d, v, e.Cond), T: Shift(d, v, e.T), F: Shift(d, v, e.F)}
	case NaturalLit:
		return e
	case NaturalPlus:
		return NaturalPlus{L: Shift(d, v, e.L), R: Shift(d, v, e.R)}
	case NaturalTimes:
		return NaturalTimes{L: Shift(d, v, e.L), R: Shift(d, v, e.R)}
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
	case Record:
		fields := make(map[string]Expr, len(map[string]Expr(e)))
		for name, val := range map[string]Expr(e) {
			fields[name] = Shift(d, v, val)
		}
		return Record(fields)
	case RecordLit:
		fields := make(map[string]Expr, len(map[string]Expr(e)))
		for name, val := range map[string]Expr(e) {
			fields[name] = Shift(d, v, val)
		}
		return RecordLit(fields)
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
	case BuiltinType:
		return e
	case DoubleLit:
		return e
	case BoolLit:
		return e
	case BoolIf:
		return BoolIf{Cond: Subst(v, c, e.Cond), T: Subst(v, c, e.T), F: Subst(v, c, e.F)}
	case NaturalLit:
		return e
	case NaturalPlus:
		return NaturalPlus{L: Subst(v, c, e.L), R: Subst(v, c, e.R)}
	case NaturalTimes:
		return NaturalTimes{L: Subst(v, c, e.L), R: Subst(v, c, e.R)}
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
	case Record:
		fields := make(map[string]Expr, len(map[string]Expr(e)))
		for name, val := range map[string]Expr(e) {
			fields[name] = Subst(v, c, val)
		}
		return Record(fields)
	case RecordLit:
		fields := make(map[string]Expr, len(map[string]Expr(e)))
		for name, val := range map[string]Expr(e) {
			fields[name] = Subst(v, c, val)
		}
		return RecordLit(fields)
	}
	panic("missing switch case in Subst()")
}

func Rule(a Const, b Const) (Const, error) {
	if b == Type {
		return Type, nil
	}
	if a == Kind && b == Kind {
		return Kind, nil
	}
	if a == Sort && (b == Kind || b == Sort) {
		return Sort, nil
	}
	return Const(0), errors.New("Dependent types are not allowed")
}

const (
	Double BuiltinType = iota
	Bool
	Natural
	Integer
	List
)

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
	_ Expr = Bool
	_ Expr = BoolLit(true)
	_ Expr = BoolIf{}
	_ Expr = Natural
	_ Expr = NaturalLit(3)
	_ Expr = NaturalPlus{}
	_ Expr = NaturalTimes{}
	_ Expr = Integer
	_ Expr = IntegerLit(3)
	_ Expr = List
	_ Expr = EmptyList{Natural}
	_ Expr = NonEmptyList([]Expr{NaturalLit(3)})
	_ Expr = Record(map[string]Expr{})
	_ Expr = RecordLit(map[string]Expr{})
)

func (c Const) WriteTo(out io.Writer) (int64, error) {
	var n int
	var err error
	if c == Type {
		n, err = fmt.Fprint(out, "Type")
	} else if c == Kind {
		n, err = fmt.Fprint(out, "Kind")
	} else {
		n, err = fmt.Fprint(out, "Sort")
	}
	return int64(n), err
}

func (v Var) WriteTo(out io.Writer) (int64, error) {
	var n int
	var err error
	if v.Index == 0 {
		n, err = fmt.Fprint(out, v.Name)
	} else {
		n, err = fmt.Fprintf(out, "%s@%d", v.Name, v.Index)
	}
	return int64(n), err
}

func (lam *LambdaExpr) WriteTo(out io.Writer) (int64, error) {
	w1, err := fmt.Fprintf(out, "λ(%s : ", lam.Label)
	if err != nil {
		return int64(w1), err
	}
	w2, err := lam.Type.WriteTo(out)
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := fmt.Fprint(out, ") → ")
	if err != nil {
		return int64(w1) + int64(w2) + int64(w3), err
	}
	w4, err := lam.Body.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3) + int64(w4), err
}

func (pi *Pi) WriteTo(out io.Writer) (int64, error) {
	w1, err := fmt.Fprintf(out, "∀(%s : ", pi.Label)
	if err != nil {
		return int64(w1), err
	}
	w2, err := pi.Type.WriteTo(out)
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := fmt.Fprint(out, ") → ")
	if err != nil {
		return int64(w1) + int64(w2) + int64(w3), err
	}
	w4, err := pi.Body.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3) + int64(w4), err
}

func (app *App) WriteTo(out io.Writer) (int64, error) {
	w1, err := app.Fn.WriteTo(out)
	if err != nil {
		return int64(w1), err
	}
	w2, err := fmt.Fprint(out, " ")
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := app.Arg.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3), err
}

func (l Let) WriteTo(out io.Writer) (int64, error) {
	panic("unimplemented")
}

func (a Annot) WriteTo(out io.Writer) (int64, error) {
	w1, err := a.Expr.WriteTo(out)
	if err != nil {
		return int64(w1), err
	}
	w2, err := fmt.Fprint(out, " : ")
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := a.Annotation.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3), err
}

func (t BuiltinType) WriteTo(out io.Writer) (int64, error) {
	var n int
	var err error
	switch t {
	case Double:
		n, err = fmt.Fprint(out, "Double")
	case Bool:
		n, err = fmt.Fprint(out, "Bool")
	case Natural:
		n, err = fmt.Fprint(out, "Natural")
	case Integer:
		n, err = fmt.Fprint(out, "Integer")
	case List:
		n, err = fmt.Fprint(out, "List")
	default:
		panic(fmt.Sprintf("unknown type %d\n", t))
	}
	return int64(n), err
}

func (d DoubleLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%f", d)
	return int64(n), err
}

func (bl BoolLit) WriteTo(out io.Writer) (int64, error) {
	if bool(bl) {
		n, err := fmt.Fprint(out, "True")
		return int64(n), err
	} else {
		n, err := fmt.Fprint(out, "True")
		return int64(n), err
	}
}

func (b BoolIf) WriteTo(out io.Writer) (int64, error) {
	w1, err := fmt.Fprint(out, "if ")
	if err != nil {
		return int64(w1), err
	}
	w2, err := b.Cond.WriteTo(out)
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := fmt.Fprint(out, " then ")
	if err != nil {
		return int64(w1) + int64(w2) + int64(w3), err
	}
	w4, err := b.T.WriteTo(out)
	if err != nil {
		return int64(w1) + int64(w2) + int64(w3) + int64(w4), err
	}
	w5, err := fmt.Fprint(out, " else ")
	if err != nil {
		return int64(w1) + int64(w2) + int64(w3) + int64(w4) + int64(w5), err
	}
	w6, err := b.F.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3) + int64(w4) + int64(w5) + int64(w6), err
}

func (nl NaturalLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%d", nl)
	return int64(n), err
}

func (p NaturalPlus) WriteTo(out io.Writer) (int64, error) {
	w1, err := p.L.WriteTo(out)
	if err != nil {
		return int64(w1), err
	}
	w2, err := fmt.Fprint(out, " + ")
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := p.R.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3), err
}

func (p NaturalTimes) WriteTo(out io.Writer) (int64, error) {
	w1, err := p.L.WriteTo(out)
	if err != nil {
		return int64(w1), err
	}
	w2, err := fmt.Fprint(out, " * ")
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := p.R.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3), err
}

func (i IntegerLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%d", i)
	return int64(n), err
}

func (l EmptyList) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "[] : ")
	if err != nil {
		return int64(n), err
	}
	n2, err := l.Type.WriteTo(out)
	return int64(n) + int64(n2), err
}

func (l NonEmptyList) WriteTo(out io.Writer) (int64, error) {
	var written int64
	i, err := fmt.Fprint(out, "[ ")
	written += int64(i)
	if err != nil {
		return written, err
	}
	exprs := []Expr(l)
	i6, err := exprs[0].WriteTo(out)
	written += i6
	if err != nil {
		return written, err
	}
	for _, expr := range exprs[1:] {
		i, err = fmt.Fprint(out, ", ")
		written += int64(i)
		if err != nil {
			return written, err
		}
		i6, err := expr.WriteTo(out)
		written += i6
		if err != nil {
			return written, err
		}
	}
	i, err = fmt.Fprint(out, " ]")
	written += int64(i)
	return written, err
}

func (r Record) WriteTo(out io.Writer) (int64, error) {
	var written int64
	i, err := fmt.Fprint(out, "{ ")
	written += int64(i)
	if err != nil {
		return written, err
	}
	fields := map[string]Expr(r)
	first := true
	for name, expr := range fields {
		if !first {
			i, err = fmt.Fprint(out, ", ")
			written += int64(i)
			if err != nil {
				return written, err
			}
		}
		first = false
		i, err = fmt.Fprintf(out, "%s : ", name)
		written += int64(i)
		if err != nil {
			return written, err
		}
		i64, err := expr.WriteTo(out)
		written += i64
		if err != nil {
			return written, err
		}
	}
	i, err = fmt.Fprint(out, " }")
	written += int64(i)
	return written, err
}

func (r RecordLit) WriteTo(out io.Writer) (int64, error) {
	var written int64
	i, err := fmt.Fprint(out, "{ ")
	written += int64(i)
	if err != nil {
		return written, err
	}
	fields := map[string]Expr(r)
	first := true
	for name, expr := range fields {
		if !first {
			i, err = fmt.Fprint(out, ", ")
			written += int64(i)
			if err != nil {
				return written, err
			}
		}
		first = false
		i, err = fmt.Fprintf(out, "%s = ", name)
		written += int64(i)
		if err != nil {
			return written, err
		}
		i64, err := expr.WriteTo(out)
		written += i64
		if err != nil {
			return written, err
		}
	}
	i, err = fmt.Fprint(out, " }")
	written += int64(i)
	return written, err
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
	return app
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

func (t BuiltinType) Normalize() Expr { return t }

func (d DoubleLit) Normalize() Expr { return d }

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
	return b
}

func (n NaturalLit) Normalize() Expr { return n }
func (p NaturalPlus) Normalize() Expr {
	L := p.L.Normalize()
	R := p.R.Normalize()

	Ln, Lok := L.(NaturalLit)
	Rn, Rok := R.(NaturalLit)

	if Lok && Rok {
		return NaturalLit(int(Ln) + int(Rn))
	} else if Lok && int(Ln) == 0 {
		return R
	} else if Rok && int(Rn) == 0 {
		return L
	} else {
		return NaturalPlus{L: L, R: R}
	}
}

func (p NaturalTimes) Normalize() Expr {
	L := p.L.Normalize()
	R := p.R.Normalize()

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
	} else {
		return NaturalTimes{L: L, R: R}
	}
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

func (r Record) Normalize() Expr {
	fields := make(map[string]Expr, len(map[string]Expr(r)))
	for name, val := range map[string]Expr(r) {
		fields[name] = val.Normalize()
	}
	return Record(fields)
}

func (r RecordLit) Normalize() Expr {
	fields := make(map[string]Expr, len(map[string]Expr(r)))
	for name, val := range map[string]Expr(r) {
		fields[name] = val.Normalize()
	}
	return RecordLit(fields)
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

func (t BuiltinType) AlphaNormalize() Expr { return t }

func (d DoubleLit) AlphaNormalize() Expr { return d }

func (n BoolLit) AlphaNormalize() Expr { return n }
func (b BoolIf) AlphaNormalize() Expr {
	return BoolIf{
		Cond: b.Cond.AlphaNormalize(),
		T:    b.T.AlphaNormalize(),
		F:    b.F.AlphaNormalize(),
	}
}

func (n NaturalLit) AlphaNormalize() Expr { return n }
func (p NaturalPlus) AlphaNormalize() Expr {
	L := p.L.AlphaNormalize()
	R := p.R.AlphaNormalize()

	return NaturalPlus{L: L, R: R}
}

func (p NaturalTimes) AlphaNormalize() Expr {
	L := p.L.AlphaNormalize()
	R := p.R.AlphaNormalize()

	return NaturalTimes{L: L, R: R}
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

func (r Record) AlphaNormalize() Expr {
	fields := make(map[string]Expr, len(map[string]Expr(r)))
	for name, val := range map[string]Expr(r) {
		fields[name] = val.AlphaNormalize()
	}
	return Record(fields)
}

func (r RecordLit) AlphaNormalize() Expr {
	fields := make(map[string]Expr, len(map[string]Expr(r)))
	for name, val := range map[string]Expr(r) {
		fields[name] = val.AlphaNormalize()
	}
	return RecordLit(fields)
}

func NewLambdaExpr(arg string, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
