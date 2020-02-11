package dhall_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/philandstuff/dhall-golang/term"
)

var (
	NaturalLit = gopter.DeriveGen(
		func(n int) term.Term { var x term.Term = term.NaturalLit(n); return x },
		func(n term.Term) int { return int(n.(term.NaturalLit)) },
		gen.IntRange(0, math.MaxInt32).WithLabel("NaturalLit"),
	)
	IntegerLit = gopter.DeriveGen(
		func(i int) term.Term { return term.IntegerLit(i) },
		func(i term.Term) int { return int(i.(term.IntegerLit)) },
		gen.Int()).WithLabel("IntegerLit")
	LeafExpr = gen.OneGenOf(NaturalLit, IntegerLit, BasicType).WithLabel("LeafExpr")
)

func PlusOf(inner gopter.Gen) gopter.Gen {
	return gen.Struct(reflect.TypeOf(term.Op{}), map[string]gopter.Gen{
		"OpCode": gen.Const(term.PlusOp),
		"L":      inner,
		"R":      inner,
	}).WithLabel("NaturalPlus")
}

func TimesOf(inner gopter.Gen) gopter.Gen {
	return gen.Struct(reflect.TypeOf(term.Op{}), map[string]gopter.Gen{
		"OpCode": gen.Const(term.TimesOp),
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

var BasicType = gen.OneConstOf(term.Double, term.Text, term.Bool, term.Natural, term.Integer).WithLabel("BasicType")

func ListOf(inner gopter.Gen) gopter.Gen {
	return gopter.DeriveGen(
		func(expr term.Term) term.Term { return term.Apply(term.List, expr) },
		func(listType term.Term) term.Term { return listType.(term.App).Arg },
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
	t.Skip("TODO: This is broken on the nbe branch until we reproduce proper printing")
	properties := gopter.NewProperties(nil)

	properties.Property("written expressions parse back as themselves",
		prop.ForAll(
			func(e term.Term) bool {
				expr, err := parser.Parse("-", []byte(fmt.Sprint(e)))
				if err != nil {
					return false
				}
				return core.AlphaEquivalent(core.Eval(e), core.Eval(expr))
			},
			ExprOf(ExprOf(LeafExpr)),
		))

	properties.TestingRun(t)
}
