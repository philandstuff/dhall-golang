package ast

import (
	"errors"
	"fmt"
	"io"
	"log"
)

type TypeContext *map[string][]Expr

func EmptyContext() TypeContext {
	return &map[string][]Expr{}
}

type (
	Expr interface {
		Normalize() Expr
		TypeWith(TypeContext) (Expr, error)
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

	natural struct{}

	NaturalLit int
)

const (
	Type Const = Const(iota)
	Kind Const = Const(iota)
	Sort Const = Const(iota)
)

// FIXME placeholder before we actually implement it
func shift(d int, v Var, e Expr) Expr {
	return e
}

// subst(x, C, B) == B[x := C]
func subst(v Var, c Expr, b Expr) Expr {
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
		substType := subst(v, c, e.Type)
		v2 := v
		if v.Name == e.Label {
			v2.Index++
		}
		substBody := subst(v2, shift(1, Var{Name: e.Label}, c), e.Body)
		return &LambdaExpr{
			Label: e.Label,
			Type:  substType,
			Body:  substBody,
		}
	case *Pi:
		substType := subst(v, c, e.Type)
		v2 := v
		if v.Name == e.Label {
			v2.Index++
		}
		substBody := subst(v2, shift(1, Var{Name: e.Label}, c), e.Body)
		return &Pi{
			Label: e.Label,
			Type:  substType,
			Body:  substBody,
		}
	case *App:
		return &App{
			Fn:  subst(v, c, e.Fn),
			Arg: subst(v, c, e.Arg),
		}
	case natural:
		return e
	case NaturalLit:
		return e
	}
	panic("missing switch case in subst()")
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
	Natural natural = natural(struct{}{})
)

var (
	_ Expr = Type
	_ Expr = &Var{}
	_ Expr = &LambdaExpr{}
	_ Expr = &Pi{}
	_ Expr = &App{}
	_ Expr = Natural
	_ Expr = NaturalLit(3)
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

func (natural) WriteTo(out io.Writer) (int, error) { return fmt.Fprint(out, "Natural") }

func (n NaturalLit) WriteTo(out io.Writer) (int, error) { return fmt.Fprintf(out, "%d", n) }

func (c Const) TypeWith(TypeContext) (Expr, error) {
	if c == Type {
		return Kind, nil
	}
	if c == Kind {
		return Sort, nil
	}
	return nil, errors.New("Sort has no type")
}

func (v Var) TypeWith(ctx TypeContext) (Expr, error) {
	if t, ok := (*ctx)[v.Name]; ok {
		return t[0], nil
	}
	return nil, fmt.Errorf("Unbound variable %s", v.Name)
}

func (lam *LambdaExpr) TypeWith(ctx TypeContext) (Expr, error) {
	// FIXME: proper de bruijn indices to avoid variable capture
	// FIXME: modifying context in place is.. icky
	argType := lam.Type.Normalize()
	(*ctx)[lam.Label] = append([]Expr{argType}, (*ctx)[lam.Label]...)
	bodyType, err := lam.Body.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	// Restore ctx to how it was before
	(*ctx)[lam.Label] = (*ctx)[lam.Label][1:len((*ctx)[lam.Label])]

	p := &Pi{Label: lam.Label, Type: argType, Body: bodyType}
	_, err2 := p.TypeWith(ctx)
	if err2 != nil {
		return nil, err2
	}

	return p, nil
}

func (pi *Pi) TypeWith(ctx TypeContext) (Expr, error) {
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

func (app *App) TypeWith(ctx TypeContext) (Expr, error) {
	fnType, err := app.Fn.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	tF := fnType.Normalize()
	// FIXME return error rather than panic if tF isn't a
	// *Pi
	pF := tF.(*Pi)

	argType, err := app.Arg.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	// FIXME replace == with a JudgmentallyEqual() fn here
	if pF.Type == argType {
		a := shift(1, Var{Name: pF.Label}, app.Arg)
		b := subst(Var{Name: pF.Label}, a, pF.Body)
		return shift(-1, Var{Name: pF.Label}, b), nil
	} else {
		return nil, errors.New("type mismatch between lambda and applied value")
	}
}

func (natural) TypeWith(TypeContext) (Expr, error) { return Type, nil }

func (n NaturalLit) TypeWith(TypeContext) (Expr, error) { return Natural, nil }

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
		a2 := shift(1, v, a)
		b1 := subst(v, a2, l.Body)
		b2 := shift(-1, v, b1)
		return b2.Normalize()
	}
	panic("got stuck in (*App).Normalize()")
}

func (n natural) Normalize() Expr    { return n }
func (n NaturalLit) Normalize() Expr { return n }

func NewLambdaExpr(arg string, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
