package eval

import (
	"reflect"

	. "github.com/philandstuff/dhall-golang/core"
)

func judgmentallyEqual(t1 Term, t2 Term) bool {
	// TODO: alpha normalization
	ne1 := Eval(t1, Env{})
	ne2 := Eval(t2, Env{})
	return reflect.DeepEqual(ne1, ne2)
}
