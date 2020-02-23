package dhall

import (
	"fmt"
	"reflect"

	"github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/imports"
	"github.com/philandstuff/dhall-golang/parser"
)

func isMapEntryType(recordType core.RecordType) bool {
	if _, ok := recordType["mapKey"]; ok {
		if _, ok := recordType["mapValue"]; ok {
			return len(recordType) == 2
		}
	}
	return false
}

// Unmarshal takes dhall input as a byte array and parses it, resolves
// imports, typechecks, evaluates, and unmarshals it into the given
// variable.
func Unmarshal(b []byte, out interface{}) error {
	term, err := parser.Parse("-", b)
	if err != nil {
		return err
	}
	resolved, err := imports.Load(term)
	if err != nil {
		return err
	}
	_, err = core.TypeOf(resolved)
	if err != nil {
		return err
	}
	Decode(core.Eval(resolved), out)
	return nil
}

// Decode takes a core.Value and unmarshals it into the given
// variable.
func Decode(e core.Value, out interface{}) {
	v := reflect.ValueOf(out)
	decode(e, v.Elem())
}

// encode converts a reflect.Value to a core.Value with the given
// Dhall type
func encode(val reflect.Value, typ core.Value) core.Value {
	switch e := typ.(type) {
	case core.Builtin:
		switch typ {
		case core.Double:
			return core.DoubleLit(val.Float())
		case core.Bool:
			return core.BoolLit(val.Bool())
		case core.Natural:
			switch val.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return core.NaturalLit(val.Uint())
			default:
				return core.NaturalLit(val.Int())
			}
		case core.Integer:
			return core.IntegerLit(val.Int())
		case core.Text:
			return core.PlainTextLit(val.String())
		default:
			panic("unknown Builtin")
		}
	case core.ListOf:
		if val.Kind() == reflect.Map {
			mapEntryType := e.Type.(core.RecordType)
			if !isMapEntryType(mapEntryType) {
				panic("Can't unmarshal golang map into given Dhall type")
			}
			if val.Len() == 0 {
				return core.EmptyList{Type: e.Type}
			}
			l := make(core.NonEmptyList, val.Len())
			iter := val.MapRange()
			i := 0
			for iter.Next() {
				l[i] = core.RecordLit{
					"mapKey":   encode(iter.Key(), mapEntryType["mapKey"]),
					"mapValue": encode(iter.Value(), mapEntryType["mapValue"]),
				}
				i++
			}
			return l
		}
		if val.Len() == 0 {
			return core.EmptyList{Type: e.Type}
		}
		l := make(core.NonEmptyList, val.Len())
		for i := 0; i < val.Len(); i++ {
			l[i] = encode(val.Index(i), e.Type)
		}
		return l
	case core.RecordType:
		rec := core.RecordLit{}
	fields:
		for key, typ := range e {
			structType := val.Type()
			for i := 0; i < structType.NumField(); i++ {
				tag := structType.Field(i).Tag.Get("dhall")
				if key == tag {
					rec[key] = encode(val.Field(i), typ)
					continue fields
				}
			}
			rec[key] = encode(val.FieldByName(key), typ)
		}
		return rec
	default:
		panic("Don't know what to do with val")
	}
}

// dhallShim takes a Callable and wraps it so that it can be passed
// to reflect.MakeFunc().  This means it converts reflect.Value inputs
// to core.Value inputs, and converts core.Value outputs to
// reflect.Value outputs.
func dhallShim(out reflect.Type, dhallFunc core.Callable) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		var expr core.Value = dhallFunc
		for _, arg := range args {
			fn := expr.(core.Callable)
			dhallArg := encode(arg, fn.ArgType())
			expr = fn.Call(dhallArg)
		}
		ptr := reflect.New(out)
		decode(expr, ptr.Elem())
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
	if some, ok := e.(core.Some); ok {
		return flattenOptional(some.Val)
	}
	if _, ok := e.(core.NoneOf); ok {
		return nil
	}
	return e
}

func decode(e core.Value, v reflect.Value) {
	e = flattenOptional(e)
	if e == nil {
		return
	}
	switch v.Kind() {
	case reflect.Map:
		// initialise with new (non-nil) value
		v.Set(reflect.MakeMap(v.Type()))
		if _, ok := e.(core.EmptyList); ok {
			return
		}
		e := e.(core.NonEmptyList)
		recordLit := e[0].(core.RecordLit)
		if len(recordLit) != 2 {
			panic("can only unmarshal `List {mapKey : T, mapValue : U}` into go map")
		}
		for _, r := range e {
			entry := r.(core.RecordLit)
			key := reflect.New(v.Type().Key()).Elem()
			val := reflect.New(v.Type().Elem()).Elem()
			decode(entry["mapKey"], key)
			decode(entry["mapValue"], val)
			v.SetMapIndex(key, val)
		}
	case reflect.Struct:
		e := e.(core.RecordLit)
		structType := v.Type()
		for i := 0; i < structType.NumField(); i++ {
			// FIXME ignores fields in RecordLit not in Struct
			tag := structType.Field(i).Tag.Get("dhall")
			if tag != "" {
				decode(e[tag], v.Field(i))
			} else {
				decode(e[structType.Field(i).Name], v.Field(i))
			}
		}
	case reflect.Func:
		e := e.(core.Callable)
		fnType := v.Type()
		returnType := fnType.Out(0)
		fn := reflect.MakeFunc(fnType, dhallShim(returnType, e))
		v.Set(fn)
	case reflect.Slice:
		if _, ok := e.(core.EmptyList); ok {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			return
		}
		e := e.(core.NonEmptyList)
		slice := reflect.MakeSlice(v.Type(), len(e), len(e))
		for i, expr := range e {
			decode(expr, slice.Index(i))
		}
		v.Set(slice)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v.SetUint(uint64(e.(core.NaturalLit)))
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
		case core.PlainTextLit:
			v.SetString(string(e))
		default:
			panic(fmt.Sprintf("Don't know how to decode %v into %v", e, v.Kind()))
		}
	}
}
