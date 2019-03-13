package main_test

import (
	"bytes"
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
		func(n int) ast.NaturalLit { return ast.NaturalLit(n) },
		func(n ast.NaturalLit) int { return int(n) },
		gen.IntRange(0, math.MaxInt32)).WithLabel("NaturalLit")
	IntegerLit = gopter.DeriveGen(
		func(i int) ast.IntegerLit { return ast.IntegerLit(i) },
		func(i ast.IntegerLit) int { return int(i) },
		gen.Int()).WithLabel("IntegerLit")
	Expr = gen.OneGenOf(NaturalLit, IntegerLit, AnyType)
)
var BasicType = gen.OneConstOf(ast.Double, ast.Text, ast.Bool, ast.Natural, ast.Integer)

func ListType(in *gopter.GenParameters) *gopter.GenResult {
	return gopter.DeriveGen(
		func(e ast.Expr) *ast.App { return &ast.App{Fn: ast.List, Arg: e} },
		func(a *ast.App) ast.Expr { return a.Arg },
		BasicType,
	)(in)
}

func AnyType(in *gopter.GenParameters) *gopter.GenResult {
	return gen.Weighted([]gen.WeightedGen{
		gen.WeightedGen{Weight: 4,
			Gen: BasicType,
		},
		gen.WeightedGen{Weight: 1,
			Gen: ListType,
		},
	})(in)
}

func TestParseWhatYouWrite(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("written expressions parse back as themselves",
		prop.ForAll(
			func(e ast.Expr) bool {
				buf := new(bytes.Buffer)
				_, err := e.WriteTo(buf)
				if err != nil {
					return false
				}
				expr, err := parser.ParseReader("-", buf)
				if err != nil {
					return false
				}
				return reflect.DeepEqual(e, expr)
			},
			Expr,
		))

	properties.TestingRun(t)
}
