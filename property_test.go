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
	NaturalLit = gen.Int32Range(0, math.MaxInt32).
			Map(func(i int32) ast.Expr {
			return ast.NaturalLit(i)
		}).WithLabel("NaturalLit")
	IntegerLit = gen.Int32().
			Map(func(i int32) ast.Expr {
			return ast.IntegerLit(i)
		}).WithLabel("IntegerLit")
	Expr = gen.OneGenOf(NaturalLit, IntegerLit)
)

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
