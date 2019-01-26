package ast

import (
	"errors"
	"fmt"
)

type TypeContext *map[string][]Expr

func EmptyContext() TypeContext {
	return &map[string][]Expr{}
}

type (
	Expr interface {
		Normalize() Expr
		TypeWith(TypeContext) (Expr, error)
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

	natural struct{}

	NaturalLit int
)

const (
	Type Const = Const(iota)
	Kind Const = Const(iota)
	Sort Const = Const(iota)
)

var (
	Natural natural = natural(struct{}{})
)

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
	return nil, errors.New("Unimplemented")
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
func (n natural) Normalize() Expr    { return n }
func (n NaturalLit) Normalize() Expr { return n }

func NewLambdaExpr(arg string, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
