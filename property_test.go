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
)

var (
	NaturalLit = gopter.DeriveGen(
		func(n int) core.Term { var x core.Term = core.NaturalLit(n); return x },
		func(n core.Term) int { return int(n.(core.NaturalLit)) },
		gen.IntRange(0, math.MaxInt32).WithLabel("NaturalLit"),
	)
	IntegerLit = gopter.DeriveGen(
		func(i int) core.Term { return core.IntegerLit(i) },
		func(i core.Term) int { return int(i.(core.IntegerLit)) },
		gen.Int()).WithLabel("IntegerLit")
	LeafExpr = gen.OneGenOf(NaturalLit, IntegerLit, BasicType).WithLabel("LeafExpr")
)

func PlusOf(inner gopter.Gen) gopter.Gen {
	return gen.Struct(reflect.TypeOf(core.OpTerm{}), map[string]gopter.Gen{
		"OpCode": gen.Const(core.PlusOp),
		"L":      inner,
		"R":      inner,
	}).WithLabel("NaturalPlus")
}

func TimesOf(inner gopter.Gen) gopter.Gen {
	return gen.Struct(reflect.TypeOf(core.OpTerm{}), map[string]gopter.Gen{
		"OpCode": gen.Const(core.TimesOp),
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

var BasicType = gen.OneConstOf(core.Double, core.Text, core.Bool, core.Natural, core.Integer).WithLabel("BasicType")

func ListOf(inner gopter.Gen) gopter.Gen {
	return gopter.DeriveGen(
		func(expr core.Term) core.Term { return core.Apply(core.List, expr) },
		func(listType core.Term) core.Term { return listType.(core.AppTerm).Arg },
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
			func(e core.Term) bool {
				expr, err := parser.Parse("-", []byte(fmt.Sprint(e)))
				if err != nil {
					return false
				}
				return core.AlphaEquivalent(e, expr)
			},
			ExprOf(ExprOf(LeafExpr)),
		))

	properties.TestingRun(t)
}
