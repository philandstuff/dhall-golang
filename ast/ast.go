package ast

type Expr interface {
	expr()
}

type LabelNode struct {
	Value string
}

type LambdaExpr struct {
	Label *LabelNode
	Type  Expr
	Body  Expr
}

func (*LabelNode) expr()  {}
func (*LambdaExpr) expr() {}

func NewLabelNode(value string) *LabelNode {
	return &LabelNode{
		Value: value,
	}
}

func NewLambdaExpr(arg *LabelNode, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
