package dhall

import (
	"reflect"
	"sort"
	"strings"

	"github.com/philandstuff/dhall-golang/ast"
)

func Unmarshal(e ast.Expr, out interface{}) {
	v := reflect.ValueOf(out)
	unmarshal(e, v.Elem())
}

func reflectValToDhallVal(val reflect.Value, typ ast.Expr) ast.Expr {
	typ = typ.Normalize()
	switch e := typ.(type) {
	case ast.Builtin:
		switch typ {
		case ast.Double:
			return ast.DoubleLit(val.Float())
		case ast.Bool:
			return ast.BoolLit(val.Bool())
		case ast.Natural:
			return ast.NaturalLit(val.Int())
		case ast.Integer:
			return ast.IntegerLit(val.Int())
		case ast.Text:
			return ast.TextLit{Suffix: val.String()}
		case ast.List:
			panic("wrong Kind")
		default:
			panic("unknown Builtin")
		}
		// TODO: RecordLit
	case *ast.App:
		switch e.Fn {
		case ast.List:
			if val.Len() == 0 {
				return ast.EmptyList{Type: e.Arg}
			}
			l := make(ast.NonEmptyList, val.Len())
			for i := 0; i < val.Len(); i++ {
				l[i] = reflectValToDhallVal(val.Index(i), e.Arg)
			}
			return l
		default:
			panic("unknown app")
		}
	default:
		panic("Don't know what to do with val")
	}
}

func argNType(fn *ast.LambdaExpr, n int) ast.Expr {
	if n == 0 {
		return fn.Type
	}
	return argNType(fn.Body.(*ast.LambdaExpr), n-1)
}

func dhallShim(out reflect.Type, dhallFunc *ast.LambdaExpr) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		var expr ast.Expr = dhallFunc
		for i, arg := range args {
			dhallArg := reflectValToDhallVal(arg, argNType(dhallFunc, i))
			expr = &ast.App{Fn: expr, Arg: dhallArg}
		}
		ptr := reflect.New(out)
		unmarshal(expr.Normalize(), ptr.Elem())
		return []reflect.Value{ptr.Elem()}
	}
}

// flattenOptional(e) returns:
// nil                if e is None T
// flattenOptional(x) if e is Some x
// e                  otherwise
// note that there may be options buried deeper in e; we just strip any outer
// Optional layers.
func flattenOptional(e ast.Expr) ast.Expr {
	if some, ok := e.(ast.Some); ok {
		return flattenOptional(some.Val)
	}
	if app, ok := e.(*ast.App); ok {
		if app.Fn == ast.None {
			return nil
		}
	}
	return e
}

// assumes e : ast.Type
func dhallTypeToReflectType(e ast.Expr) reflect.Type {
	switch e := e.(type) {
	case ast.Builtin:
		switch e {
		case ast.Double:
			return reflect.TypeOf(float64(0))
		case ast.Bool:
			return reflect.TypeOf(true)
		case ast.Integer:
			return reflect.TypeOf(int(0))
		case ast.Natural:
			return reflect.TypeOf(uint(0))
		case ast.Text:
			return reflect.TypeOf("foo")
		}
	case *ast.App:
		switch e.Fn {
		case ast.List:
			return reflect.SliceOf(dhallTypeToReflectType(e.Arg))
		case ast.Optional:
			return dhallTypeToReflectType(e.Arg)
		}
	case ast.Record:
		fields := make([]reflect.StructField, 0)
		// go through fields alphabetically
		fieldNames := make([]string, 0)
		for k := range e {
			fieldNames = append(fieldNames, k)
		}
		sort.Strings(fieldNames)
		for _, k := range fieldNames {
			fields = append(fields, reflect.StructField{
				// force upper case first letter
				Name: strings.Title(k),
				Type: dhallTypeToReflectType(e[k]),
			})
		}
		return reflect.StructOf(fields)
	}
	// Pi types?
	// union types
	panic("unknown type")
}

func unmarshal(e ast.Expr, v reflect.Value) {
	e = flattenOptional(e)
	if e == nil {
		return
	}
	switch v.Kind() {
	case reflect.Interface:
		switch e := e.(type) {
		case ast.DoubleLit:
			v.Set(reflect.ValueOf(float64(e)))
		case ast.BoolLit:
			v.Set(reflect.ValueOf(bool(e)))
		case ast.NaturalLit:
			v.Set(reflect.ValueOf(uint(e)))
		case ast.IntegerLit:
			v.Set(reflect.ValueOf(int(e)))
		case ast.TextLit:
			// FIXME ensure TextLit is uninterpolated
			v.Set(reflect.ValueOf(e.Suffix))
		case ast.EmptyList:
			// check if it's a list of map entries
			if r, ok := e.Type.(ast.Record); ok {
				if len(r) == 2 {
					for k := range r {
						if k != "mapKey" && k != "mapValue" {
							goto notmap
						}
					}
					// it's a map; the record has exactly 2 keys and they are
					// "mapKey" and "mapValue"
					v.Set(reflect.MakeMap(reflect.MapOf(
						dhallTypeToReflectType(r["mapKey"]),
						dhallTypeToReflectType(r["mapValue"]),
					)))
					return
				}
			}
		notmap:
			sliceType := reflect.SliceOf(dhallTypeToReflectType(e.Type))
			v.Set(reflect.MakeSlice(sliceType, 0, 0))
		case ast.NonEmptyList:
			slice := reflect.MakeSlice(v.Type(), len(e), len(e))
			for i, expr := range e {
				unmarshal(expr, slice.Index(i))
			}
			v.Set(slice)
		}
	case reflect.Map:
		// initialise with new (non-nil) value
		v.Set(reflect.MakeMap(v.Type()))
		if _, ok := e.(ast.EmptyList); ok {
			return
		}
		e := e.(ast.NonEmptyList)
		recordLit := e[0].(ast.RecordLit)
		if len(recordLit) != 2 {
			panic("can only unmarshal `List {mapKey : T, mapValue : U}` into go map")
		}
		for _, r := range e {
			entry := r.(ast.RecordLit)
			key := reflect.New(v.Type().Key()).Elem()
			val := reflect.New(v.Type().Elem()).Elem()
			unmarshal(entry["mapKey"], key)
			unmarshal(entry["mapValue"], val)
			v.SetMapIndex(key, val)
		}
	case reflect.Struct:
		e := e.(ast.RecordLit)
		structType := v.Type()
		for i := 0; i < structType.NumField(); i++ {
			// FIXME ignores fields in RecordLit not in Struct
			unmarshal(e[structType.Field(i).Name], v.Field(i))
		}
	case reflect.Func:
		e := e.(*ast.LambdaExpr)
		fnType := v.Type()
		returnType := fnType.Out(0)
		fn := reflect.MakeFunc(fnType, dhallShim(returnType, e))
		v.Set(fn)
	case reflect.Slice:
		if _, ok := e.(ast.EmptyList); ok {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			return
		}
		e := e.(ast.NonEmptyList)
		slice := reflect.MakeSlice(v.Type(), len(e), len(e))
		for i, expr := range e {
			unmarshal(expr, slice.Index(i))
		}
		v.Set(slice)
	default:
		switch e := e.(type) {
		case ast.DoubleLit:
			v.SetFloat(float64(e))
		case ast.BoolLit:
			v.SetBool(bool(e))
		case ast.NaturalLit:
			v.SetInt(int64(e))
		case ast.IntegerLit:
			v.SetInt(int64(e))
		case ast.TextLit:
			// FIXME: ensure TextLit doesn't have interpolations
			v.SetString(e.Suffix)
		}
	}
}
