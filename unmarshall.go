package main

import (
	"fmt"
	"reflect"

	"github.com/philandstuff/dhall-golang/ast"
)

func Unmarshal(e ast.Expr, out interface{}) {
	v := reflect.ValueOf(out)
	unmarshal(e, v.Elem())
}

func dhallTypeToReflectType(e ast.Expr) (reflect.Type, error) {
	e = e.Normalize()
	switch e := e.(type) {
	case ast.BuiltinType:
		switch e {
		case ast.Double:
			return reflect.TypeOf(0.0), nil
		case ast.Bool:
			return reflect.TypeOf(true), nil
		case ast.Natural:
			fallthrough
		case ast.Integer:
			return reflect.TypeOf(int(0)), nil
		}
	case *ast.App:
		switch e.Fn.(type) {
		case ast.BuiltinType:
			switch e.Fn {
			case ast.List:
				innerType, err := dhallTypeToReflectType(e.Arg)
				if err != nil {
					return nil, err
				}
				return reflect.SliceOf(innerType), nil
			}
		}
	}
	return nil, fmt.Errorf("Couldn't pick type for %+v", e)
}

func unmarshal(e ast.Expr, v reflect.Value) {
	t, err := e.TypeWith(ast.EmptyContext())
	if err != nil {
		// FIXME
		return
	}
	switch e := e.(type) {
	case ast.BoolLit:
		v.SetBool(bool(e))
	case ast.NaturalLit:
		v.SetInt(int64(e))
	case ast.IntegerLit:
		v.SetInt(int64(e))
	case ast.NonEmptyList:
		// FIXME ignored error
		innerType, _ := dhallTypeToReflectType(t.(*ast.App).Arg)
		slice := reflect.MakeSlice(reflect.SliceOf(innerType), len(e), len(e))
		for i, expr := range e {
			unmarshal(expr, slice.Index(i))
		}
		v.Set(slice)
	}
}
