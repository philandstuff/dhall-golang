package ast

type (
	Expr interface {
		expr()
	}

	Label string

	LambdaExpr struct {
		Label Label
		Type  Expr
		Body  Expr
	}
)

func (Label) expr()       {}
func (*LambdaExpr) expr() {}

func NewLambdaExpr(arg Label, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
