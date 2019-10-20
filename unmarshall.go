package dhall

import (
	"reflect"
	"sort"
	"strings"

	"github.com/philandstuff/dhall-golang/core"
)

func Unmarshal(e core.Value, out interface{}) {
	v := reflect.ValueOf(out)
	unmarshal(e, v.Elem())
}

func reflectValToDhallVal(val reflect.Value, typ core.Value) core.Value {
	switch e := typ.(type) {
	case core.Builtin:
		switch typ {
		case core.Double:
			return core.DoubleLit(val.Float())
		case core.Bool:
			return core.BoolLit(val.Bool())
		case core.Natural:
			return core.NaturalLit(val.Int())
		case core.Integer:
			return core.IntegerLit(val.Int())
		case core.Text:
			return core.TextLitVal{Suffix: val.String()}
		case core.List:
			panic("wrong Kind")
		default:
			panic("unknown Builtin")
		}
		// TODO: RecordLit
	case core.AppValue:
		switch e.Fn {
		case core.List:
			if val.Len() == 0 {
				return core.EmptyListVal{Type: e.Arg}
			}
			l := make(core.NonEmptyListVal, val.Len())
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

func argNType(fn core.LambdaValue, n int) core.Value {
	if n == 0 {
		return fn.Domain
	}
	return argNType(fn.Fn(core.FreeVar{}).(core.LambdaValue), n-1)
}

func dhallShim(out reflect.Type, dhallFunc core.LambdaValue) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		var expr core.Value = dhallFunc
		for i, arg := range args {
			dhallArg := reflectValToDhallVal(arg, argNType(dhallFunc, i))
			expr = expr.(core.LambdaValue).Fn(dhallArg)
		}
		ptr := reflect.New(out)
		unmarshal(expr, ptr.Elem())
		return []reflect.Value{ptr.Elem()}
	}
}

// flattenOptional(e) returns:
// nil                if e is None T
// flattenOptional(x) if e is Some x
// e                  otherwise
// note that there may be options buried deeper in e; we just strip any outer
// Optional layers.
func flattenOptional(e core.Value) core.Value {
	if some, ok := e.(core.SomeVal); ok {
		return flattenOptional(some.Val)
	}
	if app, ok := e.(core.AppValue); ok {
		if app.Fn == core.None {
			return nil
		}
	}
	return e
}

// assumes e : core.Type
func dhallTypeToReflectType(e core.Value) reflect.Type {
	switch e := e.(type) {
	case core.Builtin:
		switch e {
		case core.Double:
			return reflect.TypeOf(float64(0))
		case core.Bool:
			return reflect.TypeOf(true)
		case core.Integer:
			return reflect.TypeOf(int(0))
		case core.Natural:
			return reflect.TypeOf(uint(0))
		case core.Text:
			return reflect.TypeOf("foo")
		}
	case core.AppValue:
		switch e.Fn {
		case core.List:
			return reflect.SliceOf(dhallTypeToReflectType(e.Arg))
		case core.Optional:
			return dhallTypeToReflectType(e.Arg)
		}
	case core.RecordTypeVal:
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

func unmarshal(e core.Value, v reflect.Value) {
	e = flattenOptional(e)
	if e == nil {
		return
	}
	switch v.Kind() {
	case reflect.Interface:
		switch e := e.(type) {
		case core.DoubleLit:
			v.Set(reflect.ValueOf(float64(e)))
		case core.BoolLit:
			v.Set(reflect.ValueOf(bool(e)))
		case core.NaturalLit:
			v.Set(reflect.ValueOf(uint(e)))
		case core.IntegerLit:
			v.Set(reflect.ValueOf(int(e)))
		case core.TextLitVal:
			// FIXME ensure TextLitVal is uninterpolated
			v.Set(reflect.ValueOf(e.Suffix))
		case core.EmptyListVal:
			// check if it's a list of map entries
			if r, ok := e.Type.(core.RecordTypeVal); ok {
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
		case core.NonEmptyListVal:
			slice := reflect.MakeSlice(v.Type(), len(e), len(e))
			for i, expr := range e {
				unmarshal(expr, slice.Index(i))
			}
			v.Set(slice)
		}
	case reflect.Map:
		// initialise with new (non-nil) value
		v.Set(reflect.MakeMap(v.Type()))
		if _, ok := e.(core.EmptyListVal); ok {
			return
		}
		e := e.(core.NonEmptyListVal)
		recordLit := e[0].(core.RecordLitVal)
		if len(recordLit) != 2 {
			panic("can only unmarshal `List {mapKey : T, mapValue : U}` into go map")
		}
		for _, r := range e {
			entry := r.(core.RecordLitVal)
			key := reflect.New(v.Type().Key()).Elem()
			val := reflect.New(v.Type().Elem()).Elem()
			unmarshal(entry["mapKey"], key)
			unmarshal(entry["mapValue"], val)
			v.SetMapIndex(key, val)
		}
	case reflect.Struct:
		e := e.(core.RecordLitVal)
		structType := v.Type()
		for i := 0; i < structType.NumField(); i++ {
			// FIXME ignores fields in RecordLit not in Struct
			unmarshal(e[structType.Field(i).Name], v.Field(i))
		}
	case reflect.Func:
		e := e.(core.LambdaValue)
		fnType := v.Type()
		returnType := fnType.Out(0)
		fn := reflect.MakeFunc(fnType, dhallShim(returnType, e))
		v.Set(fn)
	case reflect.Slice:
		if _, ok := e.(core.EmptyListVal); ok {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			return
		}
		e := e.(core.NonEmptyListVal)
		slice := reflect.MakeSlice(v.Type(), len(e), len(e))
		for i, expr := range e {
			unmarshal(expr, slice.Index(i))
		}
		v.Set(slice)
	default:
		switch e := e.(type) {
		case core.DoubleLit:
			v.SetFloat(float64(e))
		case core.BoolLit:
			v.SetBool(bool(e))
		case core.NaturalLit:
			v.SetInt(int64(e))
		case core.IntegerLit:
			v.SetInt(int64(e))
		case core.TextLitVal:
			// FIXME: ensure TextLitVal doesn't have interpolations
			v.SetString(e.Suffix)
		}
	}
}
