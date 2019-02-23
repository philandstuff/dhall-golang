package main

import (
	"reflect"

	"github.com/philandstuff/dhall-golang/ast"
)

func Unmarshal(e ast.Expr, out interface{}) {
	v := reflect.ValueOf(out)
	unmarshal(e, v.Elem())
}

func unmarshal(e ast.Expr, v reflect.Value) {
	switch e := e.(type) {
	case ast.DoubleLit:
		v.SetFloat(float64(e))
	case ast.BoolLit:
		v.SetBool(bool(e))
	case ast.NaturalLit:
		v.SetInt(int64(e))
	case ast.IntegerLit:
		v.SetInt(int64(e))
	case ast.NonEmptyList:
		slice := reflect.MakeSlice(v.Type(), len(e), len(e))
		for i, expr := range e {
			unmarshal(expr, slice.Index(i))
		}
		v.Set(slice)
	}
}
