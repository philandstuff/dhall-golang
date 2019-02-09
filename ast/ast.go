package ast

import (
	"errors"
	"fmt"
	"io"
	"log"
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
		Normalize() Expr
		TypeWith(*TypeContext) (Expr, error)
		WriteTo(io.Writer) (int, error)
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
)

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
)

func (c Const) WriteTo(out io.Writer) (int, error) {
	if c == Type {
		return fmt.Fprint(out, "Type")
	} else if c == Kind {
		return fmt.Fprint(out, "Kind")
	}
	return fmt.Fprint(out, "Sort")
}

func (v Var) WriteTo(out io.Writer) (int, error) {
	if v.Index == 0 {
		return fmt.Fprint(out, v.Name)
	}
	return fmt.Fprintf(out, "%s@%d", v.Name, v.Index)
}

func (lam *LambdaExpr) WriteTo(out io.Writer) (int, error) {
	w1, err := fmt.Fprintf(out, "λ(%s : ", lam.Label)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	w2, err := lam.Type.WriteTo(out)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	w3, err := fmt.Fprint(out, ") → ")
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	w4, err := lam.Body.WriteTo(out)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	return w1 + w2 + w3 + w4, nil
}

func (pi *Pi) WriteTo(out io.Writer) (int, error) {
	w1, err := fmt.Fprintf(out, "∀(%s : ", pi.Label)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	w2, err := pi.Type.WriteTo(out)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	w3, err := fmt.Fprint(out, ") → ")
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	w4, err := pi.Body.WriteTo(out)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	return w1 + w2 + w3 + w4, nil
}

func (app *App) WriteTo(out io.Writer) (int, error) {
	w1, err := app.Fn.WriteTo(out)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	w2, err := fmt.Fprint(out, " ")
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	w3, err := app.Arg.WriteTo(out)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	return w1 + w2 + w3, nil
}

func (double) WriteTo(out io.Writer) (int, error) { return fmt.Fprint(out, "Double") }

func (d DoubleLit) WriteTo(out io.Writer) (int, error) { return fmt.Fprintf(out, "%f", d) }

func (natural) WriteTo(out io.Writer) (int, error) { return fmt.Fprint(out, "Natural") }

func (n NaturalLit) WriteTo(out io.Writer) (int, error) { return fmt.Fprintf(out, "%d", n) }

func (p NaturalPlus) WriteTo(out io.Writer) (int, error) {
	_, err := p.L.WriteTo(out)
	if err != nil {
		log.Fatalf("Fatal error %v", err)
	}
	fmt.Fprint(out, " + ")
	return p.R.WriteTo(out)
}

func (integer) WriteTo(out io.Writer) (int, error) { return fmt.Fprint(out, "Integer") }

func (i IntegerLit) WriteTo(out io.Writer) (int, error) { return fmt.Fprintf(out, "%d", i) }

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
	panic("got stuck in (*App).Normalize()")
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

func NewLambdaExpr(arg string, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
