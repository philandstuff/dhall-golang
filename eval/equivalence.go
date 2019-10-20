package eval

import (
	"reflect"

	. "github.com/philandstuff/dhall-golang/core"
)

func judgmentallyEqual(t1 Term, t2 Term) bool {
	ne1 := AlphaBetaEval(t1)
	ne2 := AlphaBetaEval(t2)
	return reflect.DeepEqual(Quote(ne1), Quote(ne2))
}
