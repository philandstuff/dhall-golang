package dhall

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/philandstuff/dhall-golang/v5/core"
	"github.com/philandstuff/dhall-golang/v5/imports"
	"github.com/philandstuff/dhall-golang/v5/parser"
	"github.com/philandstuff/dhall-golang/v5/term"
)

func isMapEntryType(recordType map[string]core.Value) bool {
	if _, ok := recordType["mapKey"]; ok {
		if _, ok := recordType["mapValue"]; ok {
			return len(recordType) == 2
		}
	}
	return false
}

// A DecodeIfaceFunc is a function for decoding into an interface type
type DecodeIfaceFunc func(*Decoder, core.Value) (interface{}, error)

// A Decoder is responsible for decoding Dhall source into Go data.
type Decoder struct {
	interfaces map[string]map[string]DecodeIfaceFunc
}

// NewDecoder returns a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{
		interfaces: make(map[string]map[string]DecodeIfaceFunc),
	}
}

// RegisterIface registers the given interface decoding function with
// the Decoder.
func (d *Decoder) RegisterIface(target interface{}, decodeFunc DecodeIfaceFunc) {
	ifaceType := reflect.ValueOf(target).Elem().Type()
	if _, ok := d.interfaces[ifaceType.PkgPath()]; !ok {
		d.interfaces[ifaceType.PkgPath()] = make(map[string]DecodeIfaceFunc)
	}
	d.interfaces[ifaceType.PkgPath()][ifaceType.Name()] = decodeFunc
}

// Unmarshal takes dhall input as a byte array and parses it, resolves
// imports, typechecks, evaluates, and unmarshals it into the given
// variable.
func Unmarshal(b []byte, out interface{}) error {
	return NewDecoder().Unmarshal(b, out)
}

// Unmarshal takes dhall input as a byte array and parses it, resolves
// imports, typechecks, evaluates, and unmarshals it into the given
// variable.
func (d *Decoder) Unmarshal(b []byte, out interface{}) error {
	term, err := parser.Parse("-", b)
	if err != nil {
		return err
	}
	return d.unmarshalTerm(term, out)
}

// UnmarshalFile takes dhall input from a file and parses it, resolves
// imports, typechecks, evaluates, and unmarshals it into the given
// variable.
func UnmarshalFile(filename string, out interface{}) error {
	return NewDecoder().UnmarshalFile(filename, out)
}

// UnmarshalFile takes dhall input from a file and parses it, resolves
// imports, typechecks, evaluates, and unmarshals it into the given
// variable.
func (d *Decoder) UnmarshalFile(filename string, out interface{}) error {
	term, err := parser.ParseFile(filename)
	if err != nil {
		return err
	}
	return d.unmarshalTerm(term, out)
}

func (d *Decoder) unmarshalTerm(term term.Term, out interface{}) error {
	resolved, err := imports.Load(term)
	if err != nil {
		return err
	}
	_, err = core.TypeOf(resolved)
	if err != nil {
		return err
	}
	return d.Decode(core.Eval(resolved), out)
}

// Decode takes a core.Value and unmarshals it into the given
// variable.
func Decode(e core.Value, out interface{}) error {
	return NewDecoder().Decode(e, out)
}

// Decode takes a core.Value and unmarshals it into the given
// variable.
func (d *Decoder) Decode(e core.Value, out interface{}) error {
	v := reflect.ValueOf(out)
	return d.decode(e, v.Elem())
}

// encode converts a reflect.Value to a core.Value with the given
// Dhall type
func encode(val reflect.Value, typ core.Value) (core.Value, error) {
	if opt, ok := typ.(core.OptionalOf); ok {
		switch val.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice:
			if val.IsNil() {
				return core.NoneOf{Type: typ}, nil
			}
		}
		dhallVal, err := encode(val, opt.Type)
		if err != nil {
			return nil, err
		}
		return core.Some{Val: dhallVal}, nil
	}
	switch val.Kind() {
	case reflect.Bool:
		if typ == core.Bool {
			return core.BoolLit(val.Bool()), nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		if typ == core.Integer {
			return core.IntegerLit(val.Int()), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if typ == core.Natural {
			return core.NaturalLit(val.Uint()), nil
		}
	case reflect.Float32, reflect.Float64:
		if typ == core.Double {
			return core.DoubleLit(val.Float()), nil
		}
		// no Complex32 or Complex64
		// no Array
		// no Chan
	case reflect.Func: // not implemented
		// no Interface
	case reflect.Map:
		listOf, ok := typ.(core.ListOf)
		if !ok {
			break
		}
		mapEntryType, ok := listOf.Type.(core.RecordType)
		if !ok {
			break
		}
		if !isMapEntryType(mapEntryType) {
			break
		}
		if val.Len() == 0 {
			return core.EmptyList{Type: mapEntryType}, nil
		}
		l := make(core.NonEmptyList, val.Len())
		iter := val.MapRange()
		i := 0
		for iter.Next() {
			key, err := encode(iter.Key(), mapEntryType["mapKey"])
			if err != nil {
				return nil, err
			}
			val, err := encode(iter.Value(), mapEntryType["mapValue"])
			if err != nil {
				return nil, err
			}
			l[i] = core.RecordLit{
				"mapKey":   key,
				"mapValue": val,
			}
			i++
		}
		return l, nil
	case reflect.Ptr:
		return encode(val.Elem(), typ)
	case reflect.Slice:
		e, ok := typ.(core.ListOf)
		if !ok {
			break
		}
		if val.Len() == 0 {
			return core.EmptyList{Type: e.Type}, nil
		}
		l := make(core.NonEmptyList, val.Len())
		var err error
		for i := 0; i < val.Len() && err == nil; i++ {
			l[i], err = encode(val.Index(i), e.Type)
		}
		return l, err
	case reflect.String:
		if typ == core.Text {
			return core.PlainTextLit(val.String()), nil
		}
	case reflect.Struct:
		e, ok := typ.(core.RecordType)
		if !ok {
			break
		}
		rec := core.RecordLit{}
	fields:
		for key, typ := range e {
			var err error
			structType := val.Type()
			for i := 0; i < structType.NumField(); i++ {
				tag := structType.Field(i).Tag.Get("dhall")
				if key == tag {
					rec[key], err = encode(val.Field(i), typ)
					if err != nil {
						return nil, err
					}
					continue fields
				}
			}
			rec[key], err = encode(val.FieldByName(key), typ)
			if err != nil {
				return nil, err
			}
		}
		return rec, nil
		// no UnsafePointer
	}
	return nil, fmt.Errorf("Can't encode %v as %v", val, typ)
}

// dhallShim takes a Callable and wraps it so that it can be passed
// to reflect.MakeFunc().  This means it converts reflect.Value inputs
// to core.Value inputs, and converts core.Value outputs to
// reflect.Value outputs.
func (d *Decoder) dhallShim(out reflect.Type, dhallFunc core.Callable) func([]reflect.Value) []reflect.Value {
	return func(args []reflect.Value) []reflect.Value {
		var expr core.Value = dhallFunc
		for _, arg := range args {
			fn := expr.(core.Callable)
			dhallArg, err := encode(arg, fn.ArgType())
			if err != nil {
				// if the func was well-typed, this shouldn't happen
				panic(err)
			}
			expr = fn.Call(dhallArg)
		}
		ptr := reflect.New(out)
		err := d.decode(expr, ptr.Elem())
		if err != nil {
			// if the func was well-typed, this shouldn't happen
			panic(err)
		}
		return []reflect.Value{ptr.Elem()}
	}
}

func mkTestVal(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		return reflect.New(t.Elem())
	}
	return reflect.Zero(t)
}

// flattenSome(e) returns:
//
//  flattenSome(x) if e is Some x
//  e              otherwise
//
// note that there may be Somes buried deeper in e; we just strip any outer
// Some layers.
func flattenSome(e core.Value) core.Value {
	if some, ok := e.(core.Some); ok {
		return flattenSome(some.Val)
	}
	return e
}

func (d *Decoder) decode(e core.Value, v reflect.Value) error {
	e = flattenSome(e)
	if _, ok := e.(core.NoneOf); ok {
		// TODO: should we fail if a None doesn't match the type?
		// (similar to EmptyList below)
		return nil
	}
types:
	switch v.Kind() {
	case reflect.Bool:
		if e, ok := e.(core.BoolLit); ok {
			v.SetBool(bool(e))
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		if e, ok := e.(core.IntegerLit); ok {
			v.SetInt(int64(e))
			return nil
		}
		if e, ok := e.(core.NaturalLit); ok {
			v.SetInt(int64(e))
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if e, ok := e.(core.NaturalLit); ok {
			v.SetUint(uint64(e))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if e, ok := e.(core.DoubleLit); ok {
			v.SetFloat(float64(e))
			return nil
		}
	case reflect.Func:
		fnType := v.Type()
		if fnType.NumIn() == 0 {
			return fmt.Errorf("You must decode into a function type with at least one input parameter, not %v", fnType)
		}
		if fnType.NumOut() != 1 {
			return fmt.Errorf("You must decode into a function type with exactly one output parameter, not %v", fnType)
		}
		returnType := fnType.Out(0)

		result := e
		for i := 0; i < fnType.NumIn(); i++ {
			callable, ok := result.(core.Callable)
			if !ok {
				break types
			}
			testValue := mkTestVal(fnType.In(i))
			testDhallVal, err := encode(testValue, callable.ArgType())
			if err != nil {
				return err
			}
			result = callable.Call(testDhallVal)
		}
		err := d.decode(result, reflect.New(fnType.Out(0)).Elem())
		if err != nil {
			return err
		}
		fn := reflect.MakeFunc(fnType, d.dhallShim(returnType, e.(core.Callable)))
		v.Set(fn)
		return nil
	case reflect.Map:
		// initialise with new (non-nil) value
		v.Set(reflect.MakeMap(v.Type()))
		if _, ok := e.(core.EmptyList); ok {
			// TODO: should this fail if type mismatches?
			// it should at least be a mapKey/mapValue type
			return nil
		}
		if e, ok := e.(core.NonEmptyList); ok {
			recordLit, ok := e[0].(core.RecordLit)
			if !ok {
				break
			}
			if !isMapEntryType(recordLit) {
				return errors.New("can only unmarshal `List {mapKey : T, mapValue : U}` into go map")
			}
			for _, r := range e {
				entry := r.(core.RecordLit)
				key := reflect.New(v.Type().Key()).Elem()
				val := reflect.New(v.Type().Elem()).Elem()
				err := d.decode(entry["mapKey"], key)
				if err != nil {
					return err
				}
				err = d.decode(entry["mapValue"], val)
				if err != nil {
					return err
				}
				v.SetMapIndex(key, val)
			}
			return nil
		}
	case reflect.Ptr:
		v.Set(reflect.New(v.Type().Elem()))
		return d.decode(e, v.Elem())
	case reflect.Slice:
		if _, ok := e.(core.EmptyList); ok {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			return nil
		}
		if e, ok := e.(core.NonEmptyList); ok {
			slice := reflect.MakeSlice(v.Type(), len(e), len(e))
			for i, expr := range e {
				err := d.decode(expr, slice.Index(i))
				if err != nil {
					return err
				}
			}
			v.Set(slice)
			return nil
		}
	case reflect.String:
		if e, ok := e.(core.PlainTextLit); ok {
			v.SetString(string(e))
			return nil
		}
	case reflect.Struct:
		if e, ok := e.(core.RecordLit); ok {
			structType := v.Type()
			for i := 0; i < structType.NumField(); i++ {
				// FIXME ignores fields in RecordLit not in Struct
				tag := structType.Field(i).Tag.Get("dhall")
				var err error
				if tag != "" {
					err = d.decode(e[tag], v.Field(i))
				} else {
					err = d.decode(e[structType.Field(i).Name], v.Field(i))
				}
				if err != nil {
					return err
				}
			}
			return nil
		}
	case reflect.Interface:
		ifaceType := v.Type()
		if decodeFunc, ok := d.interfaces[ifaceType.PkgPath()][ifaceType.Name()]; ok {
			decoded, err := decodeFunc(d, e)
			if err != nil {
				return err
			}
			decodedValue := reflect.ValueOf(decoded)
			if !decodedValue.Type().AssignableTo(ifaceType) {
				return fmt.Errorf("%+v is not assignable to %+v", decodedValue, ifaceType)
			}
			v.Set(decodedValue)
			return nil
		}
	}
	return fmt.Errorf("Don't know how to decode %v into %v", e, v.Kind())
}
