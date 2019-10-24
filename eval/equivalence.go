package eval

import (
	"reflect"

	. "github.com/philandstuff/dhall-golang/core"
)

func judgmentallyEqual(t1 Term, t2 Term) bool {
	v1 := AlphaBetaEval(t1)
	v2 := AlphaBetaEval(t2)
	return judgmentallyEqualVals(v1, v2)
}

func judgmentallyEqualVals(v1 Value, v2 Value) bool {
	return reflect.DeepEqual(Quote(v1), Quote(v2))
}
