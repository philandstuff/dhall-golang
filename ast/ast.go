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

	double    struct{}
	DoubleLit float64

	natural     struct{}
	NaturalLit  int
	NaturalPlus struct {
		L Expr
		R Expr
	}

	integer    struct{}
	IntegerLit int

	list    struct{}
	ListLit struct {
		// Content must not be empty (use nil for no content)
		Content    []Expr
		Annotation Expr
	}
)

const (
	Type Const = Const(iota)
	Kind Const = Const(iota)
	Sort Const = Const(iota)
)

// FIXME placeholder before we actually implement it
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
	case double:
		return e
	case DoubleLit:
		return e
	case natural:
		return e
	case NaturalLit:
		return e
	case NaturalPlus:
		return NaturalPlus{L: Shift(d, v, e.L), R: Shift(d, v, e.R)}
	case integer:
		return e
	case IntegerLit:
		return e
	case list:
		return e
	case ListLit:
		if e.Content == nil {
			return e
		}
		exprs := make([]Expr, len(e.Content))
		for i, expr := range e.Content {
			exprs[i] = Shift(d, v, expr)
		}
		return MakeAnnotatedList(e.Annotation, exprs...)
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
	case double:
		return e
	case DoubleLit:
		return e
	case natural:
		return e
	case NaturalLit:
		return e
	case NaturalPlus:
		return NaturalPlus{L: Subst(v, c, e.L), R: Subst(v, c, e.R)}
	case integer:
		return e
	case IntegerLit:
		return e
	case list:
		return e
	case ListLit:
		if e.Content == nil {
			return e
		}
		exprs := make([]Expr, len(e.Content))
		for i, expr := range e.Content {
			exprs[i] = Subst(v, c, expr)
		}
		return MakeAnnotatedList(e.Annotation, exprs...)
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

var (
	Double  double  = double(struct{}{})
	Natural natural = natural(struct{}{})
	Integer integer = integer(struct{}{})
	List    list    = list(struct{}{})
)

func MakeList(content ...Expr) ListLit {
	return ListLit{Content: content}
}

func MakeAnnotatedList(annotation Expr, content ...Expr) ListLit {
	return ListLit{
		Content:    content,
		Annotation: annotation,
	}
}

var (
	_ Expr = Type
	_ Expr = &Var{}
	_ Expr = &LambdaExpr{}
	_ Expr = &Pi{}
	_ Expr = &App{}
	_ Expr = Double
	_ Expr = DoubleLit(3.0)
	_ Expr = Natural
	_ Expr = NaturalLit(3)
	_ Expr = NaturalPlus{}
	_ Expr = Integer
	_ Expr = IntegerLit(3)
	_ Expr = List
	_ Expr = ListLit{}
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
	}
	n, err = fmt.Fprintf(out, "%s@%d", v.Name, v.Index)
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

func (double) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "Double")
	return int64(n), err
}

func (d DoubleLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%f", d)
	return int64(n), err
}

func (natural) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "Natural")
	return int64(n), err
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

func (integer) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "Integer")
	return int64(n), err
}

func (i IntegerLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%d", i)
	return int64(n), err
}

func (list) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "List")
	return int64(n), err
}

func (l ListLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%d", l)
	return int64(n), err
}

func judgmentallyEqual(e1 Expr, e2 Expr) bool { return e1 == e2 }

func (c Const) TypeWith(*TypeContext) (Expr, error) {
	if c == Type {
		return Kind, nil
	}
	if c == Kind {
		return Sort, nil
	}
	return nil, errors.New("Sort has no type")
}

func (v Var) TypeWith(ctx *TypeContext) (Expr, error) {
	if t, ok := ctx.Lookup(v.Name, 0); ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unbound variable %s, context was %+v", v.Name, ctx)
}

func (lam *LambdaExpr) TypeWith(ctx *TypeContext) (Expr, error) {
	if _, err := lam.Type.TypeWith(ctx); err != nil {
		return nil, err
	}
	argType := lam.Type.Normalize()
	newctx := ctx.Insert(lam.Label, argType).Map(func(e Expr) Expr { return Shift(1, Var{Name: lam.Label}, e) })
	bodyType, err := lam.Body.TypeWith(newctx)
	if err != nil {
		return nil, err
	}

	p := &Pi{Label: lam.Label, Type: argType, Body: bodyType}
	_, err2 := p.TypeWith(ctx)
	if err2 != nil {
		return nil, err2
	}

	return p, nil
}

func (pi *Pi) TypeWith(ctx *TypeContext) (Expr, error) {
	argType, err := pi.Type.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	tA := argType.Normalize()
	// FIXME return error rather than panic if tA isn't a
	// Const
	kA := tA.(Const)
	// FIXME: proper de bruijn indices to avoid variable capture
	// FIXME: modifying context in place is.. icky
	(*ctx)[pi.Label] = append([]Expr{pi.Type.Normalize()}, (*ctx)[pi.Label]...)
	bodyType, err := pi.Body.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	tB := bodyType.Normalize()
	// FIXME return error rather than panic if tA isn't a
	// Const
	kB := tB.(Const)
	// Restore ctx to how it was before
	(*ctx)[pi.Label] = (*ctx)[pi.Label][1:len((*ctx)[pi.Label])]

	return Rule(kA, kB)
}

func (app *App) TypeWith(ctx *TypeContext) (Expr, error) {
	fnType, err := app.Fn.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	tF := fnType.Normalize()
	pF, ok := tF.(*Pi)
	if !ok {
		return nil, fmt.Errorf("Expected %s to be a function type", tF)
	}

	argType, err := app.Arg.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	// FIXME replace == with a JudgmentallyEqual() fn here
	if pF.Type == argType {
		a := Shift(1, Var{Name: pF.Label}, app.Arg)
		b := Subst(Var{Name: pF.Label}, a, pF.Body)
		return Shift(-1, Var{Name: pF.Label}, b), nil
	} else {
		return nil, errors.New("type mismatch between lambda and applied value")
	}
}

func (double) TypeWith(*TypeContext) (Expr, error) { return Type, nil }

func (DoubleLit) TypeWith(*TypeContext) (Expr, error) { return Double, nil }

func (natural) TypeWith(*TypeContext) (Expr, error) { return Type, nil }

func (NaturalLit) TypeWith(*TypeContext) (Expr, error) { return Natural, nil }

func (p NaturalPlus) TypeWith(ctx *TypeContext) (Expr, error) {
	L, err := p.L.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	L = L.Normalize()
	if L != Natural {
		return nil, fmt.Errorf("Expecting a Natural, can't add %s", L)
	}
	R, err := p.R.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	R = R.Normalize()
	if R != Natural {
		return nil, fmt.Errorf("Expecting a Natural, can't add %s", R)
	}
	return Natural, nil
}

func (integer) TypeWith(*TypeContext) (Expr, error) { return Type, nil }

func (IntegerLit) TypeWith(*TypeContext) (Expr, error) { return Integer, nil }

func (list) TypeWith(*TypeContext) (Expr, error) { return &Pi{"_", Type, Type}, nil }

func (l ListLit) TypeWith(ctx *TypeContext) (Expr, error) {
	if l.Annotation != nil {
		t := l.Annotation
		k, err := t.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		if k.Normalize() != Type {
			return nil, fmt.Errorf("List annotation %s is not a Type", t)
		}
		for _, elem := range l.Content {
			t2, err := elem.TypeWith(ctx)
			if err != nil {
				return nil, err
			}
			if !judgmentallyEqual(t, t2) {
				return nil, fmt.Errorf("Types %s and %s don't match", t, t2)
			}
		}
		return &App{List, t}, nil
	}
	// Annotation is nil, we have to infer type
	if l.Content == nil {
		return nil, fmt.Errorf("Empty lists must be annotated with type")
	}
	t, err := l.Content[0].TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	k, err := t.TypeWith(ctx)
	if k.Normalize() != Type {
		return nil, fmt.Errorf("Invalid type for List elements")
	}
	for _, elem := range l.Content[1:] {
		t2, err := elem.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		if !judgmentallyEqual(t, t2) {
			return nil, fmt.Errorf("All List elements must have same type, but types %s and %s don't match", t, t2)
		}
	}
	return &App{List, t}, nil
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

func (d double) Normalize() Expr    { return d }
func (d DoubleLit) Normalize() Expr { return d }

func (n natural) Normalize() Expr    { return n }
func (n NaturalLit) Normalize() Expr { return n }
func (p NaturalPlus) Normalize() Expr {
	L := p.L.Normalize().(NaturalLit)
	R := p.R.Normalize().(NaturalLit)
	return NaturalLit(int(L) + int(R))
}

func (i integer) Normalize() Expr    { return i }
func (i IntegerLit) Normalize() Expr { return i }

func (l list) Normalize() Expr { return l }
func (l ListLit) Normalize() Expr {
	if l.Content == nil {
		return MakeAnnotatedList(l.Annotation.Normalize())
	}
	vals := make([]Expr, len(l.Content))
	for i, expr := range l.Content {
		vals[i] = expr.Normalize()
	}
	return MakeList(vals...)
}

func NewLambdaExpr(arg string, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
