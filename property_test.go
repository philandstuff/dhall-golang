package dhall_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"
)

var (
	NaturalLit = gopter.DeriveGen(
		func(n int) ast.Expr { var x ast.Expr = ast.NaturalLit(n); return x },
		func(n ast.Expr) int { return int(n.(ast.NaturalLit)) },
		gen.IntRange(0, math.MaxInt32).WithLabel("NaturalLit"),
	)
	IntegerLit = gopter.DeriveGen(
		func(i int) ast.Expr { return ast.IntegerLit(i) },
		func(i ast.Expr) int { return int(i.(ast.IntegerLit)) },
		gen.Int()).WithLabel("IntegerLit")
	LeafExpr = gen.OneGenOf(NaturalLit, IntegerLit, BasicType).WithLabel("LeafExpr")
)

func PlusOf(inner gopter.Gen) gopter.Gen {
	return gen.Struct(reflect.TypeOf(ast.Operator{}), map[string]gopter.Gen{
		"OpCode": gen.Const(ast.PlusOp),
		"L":      inner,
		"R":      inner,
	}).WithLabel("NaturalPlus")
}

func TimesOf(inner gopter.Gen) gopter.Gen {
	return gen.Struct(reflect.TypeOf(ast.Operator{}), map[string]gopter.Gen{
		"OpCode": gen.Const(ast.TimesOp),
		"L":      inner,
		"R":      inner,
	}).WithLabel("NaturalTimes")
}

func ExprOf(inner gopter.Gen) gopter.Gen {
	return gen.OneGenOf(
		LeafExpr,
		PlusOf(inner),
		TimesOf(inner),
		ListOf(inner),
	).WithLabel("ExprOf")
}

var BasicType = gen.OneConstOf(ast.Double, ast.Text, ast.Bool, ast.Natural, ast.Integer).WithLabel("BasicType")

func ListOf(inner gopter.Gen) gopter.Gen {
	return gopter.DeriveGen(
		func(expr ast.Expr) ast.Expr { return &ast.App{Fn: ast.List, Arg: expr} },
		func(listType ast.Expr) ast.Expr { return listType.(*ast.App).Arg },
		inner,
	)
}

// func AnyType(in *gopter.GenParameters) *gopter.GenResult {
// 	return gen.Weighted([]gen.WeightedGen{
// 		gen.WeightedGen{Weight: 4,
// 			Gen: BasicType,
// 		},
// 		gen.WeightedGen{Weight: 1,
// 			Gen: ListType,
// 		},
// 	})(in)
// }

func TestParseWhatYouWrite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	properties := gopter.NewProperties(nil)

	properties.Property("written expressions parse back as themselves",
		prop.ForAll(
			func(e ast.Expr) bool {
				expr, err := parser.Parse("-", []byte(fmt.Sprint(e)))
				if err != nil {
					return false
				}
				return reflect.DeepEqual(e, expr)
			},
			ExprOf(ExprOf(LeafExpr)),
		))

	properties.TestingRun(t)
}
