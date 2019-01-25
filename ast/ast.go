package ast

type Expr interface{}

type LabelNode struct {
	Value string
}

type LambdaExpr struct {
	Label *LabelNode
	Type  Expr
	Body  Expr
}

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
