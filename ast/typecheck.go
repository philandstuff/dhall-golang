package ast

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func judgmentallyEqual(e1 Expr, e2 Expr) bool {
	ne1 := e1.Normalize().AlphaNormalize()
	ne2 := e2.Normalize().AlphaNormalize()
	return reflect.DeepEqual(ne1, ne2)
}

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

func (l Let) TypeWith(ctx *TypeContext) (Expr, error) {
	binding := l.Bindings[0]
	x := binding.Variable
	a1 := binding.Value.Normalize()
	valueType, err := binding.Value.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	if binding.Annotation != nil {
		_, err := binding.Annotation.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		if !judgmentallyEqual(binding.Annotation, valueType) {
			return nil, errors.New("type doesn't match annotation in let")
		}
	}

	// TODO: optimization where binding.Value is a term
	a2 := Shift(1, Var{x, 0}, a1)

	rest := l.Body
	if len(l.Bindings) > 1 {
		rest = Let{Bindings: l.Bindings[1:], Body: l.Body}
	}

	b1 := Subst(Var{x, 0}, a2, rest)
	b2 := Shift(-1, Var{x, 0}, b1)
	retval, err := b2.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	return retval, nil
}

func (a Annot) TypeWith(ctx *TypeContext) (Expr, error) {
	_, err := a.Annotation.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	t2, err := a.Expr.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	if !judgmentallyEqual(a.Annotation, t2) {
		var b strings.Builder
		b.WriteString("Annotation mismatch: inferred type ")
		t2.WriteTo(&b)
		b.WriteString(" but annotated ")
		a.Annotation.WriteTo(&b)
		return nil, errors.New(b.String())
	}
	return t2, nil
}

func (t BuiltinType) TypeWith(*TypeContext) (Expr, error) {
	switch t {
	case Double:
		return Type, nil
	case Bool:
		return Type, nil
	case Natural:
		return Type, nil
	case Integer:
		return Type, nil
	case List:
		return &Pi{"_", Type, Type}, nil
	default:
		panic(fmt.Sprintf("unknown type %d\n", t))
	}
}

func (DoubleLit) TypeWith(*TypeContext) (Expr, error) { return Double, nil }

func (BoolLit) TypeWith(*TypeContext) (Expr, error) { return Bool, nil }

func (b BoolIf) TypeWith(ctx *TypeContext) (Expr, error) {
	condType, err := b.Cond.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	if condType != Bool {
		return nil, errors.New("if condition must be type Bool")
	}
	tType, err := b.T.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	ttType, err := tType.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	if ttType != Type {
		return nil, errors.New("if branches must have terms, not types")
	}
	fType, err := b.F.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	ftType, err := fType.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	if ftType != Type {
		return nil, errors.New("if branches must have terms, not types")
	}
	if !judgmentallyEqual(tType, fType) {
		return nil, errors.New("if branches must have matching types")
	}
	return tType, nil
}

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

func (IntegerLit) TypeWith(*TypeContext) (Expr, error) { return Integer, nil }

func (l EmptyList) TypeWith(ctx *TypeContext) (Expr, error) {
	t := l.Type
	k, err := t.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	if k.Normalize() != Type {
		return nil, fmt.Errorf("List annotation %s is not a Type", t)
	}
	return &App{List, t}, nil
}

func (l NonEmptyList) TypeWith(ctx *TypeContext) (Expr, error) {
	exprs := []Expr(l)
	t, err := exprs[0].TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	k, err := t.TypeWith(ctx)
	if k.Normalize() != Type {
		return nil, fmt.Errorf("Invalid type for List elements")
	}
	for _, elem := range exprs[1:] {
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
