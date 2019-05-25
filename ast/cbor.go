package ast

import (
	"errors"
	"fmt"
	"net/url"
	"path"

	"github.com/ugorji/go/codec"
)

var nameToBuiltin = map[string]Expr{
	"Type": Type,
	"Kind": Kind,
	"Sort": Sort,

	"Double":   Double,
	"Text":     Text,
	"Bool":     Bool,
	"Natural":  Natural,
	"Integer":  Integer,
	"List":     List,
	"Optional": Optional,
	"None":     None,

	"Natural/build":     NaturalBuild,
	"Natural/fold":      NaturalFold,
	"Natural/isZero":    NaturalIsZero,
	"Natural/even":      NaturalEven,
	"Natural/odd":       NaturalOdd,
	"Natural/toInteger": NaturalToInteger,
	"Natural/show":      NaturalShow,

	"Integer/toDouble": IntegerToDouble,
	"Integer/show":     IntegerShow,

	"Double/show": DoubleShow,

	"Text/show": TextShow,

	"List/build":   ListBuild,
	"List/fold":    ListFold,
	"List/length":  ListLength,
	"List/head":    ListHead,
	"List/last":    ListLast,
	"List/indexed": ListIndexed,
	"List/reverse": ListReverse,

	"Optional/build": OptionalBuild,
	"Optional/fold":  OptionalFold,
}

func unwrapUint(i interface{}) (uint, error) {
	if val, ok := i.(uint64); ok {
		return uint(val), nil
	}
	return 0, fmt.Errorf("couldn't interpret %v as uint", i)
}

func unwrapInt(i interface{}) (int, error) {
	if val, ok := i.(uint64); ok {
		return int(val), nil
	}
	if val, ok := i.(int64); ok {
		return int(val), nil
	}
	return 0, fmt.Errorf("couldn't interpret %v as int", i)
}

func unwrapString(i interface{}) (string, error) {
	if val, ok := i.(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("couldn't interpret %v as string", i)
}

func decodeMap(i interface{}) (map[string]Expr, error) {
	if val, ok := i.(map[interface{}]interface{}); ok {
		decodedM := make(map[string]Expr, len(val))
		for n, t := range val {
			name, err := unwrapString(n)
			if err != nil {
				return nil, err
			}
			if t == nil {
				decodedM[name] = nil
			} else {
				decodedM[name], err = Decode(t)
				if err != nil {
					return nil, err
				}
			}
		}
		return decodedM, nil
	}
	return nil, fmt.Errorf("couldn't interpret %v as map[string]interface{}", i)
}

func Decode(decodedCbor interface{}) (Expr, error) {
	switch val := decodedCbor.(type) {
	case uint64:
		// _@n
		return Var{Name: "_", Index: int(val)}, nil
	case string:
		// Type, Double, Optional/fold
		if builtin, ok := nameToBuiltin[val]; ok {
			return builtin, nil
		}
		return nil, fmt.Errorf("unrecognized builtin %s", val)
	case bool:
		if val {
			return True, nil
		} else {
			return False, nil
		}
	case float64:
		return DoubleLit(val), nil
	case []interface{}:
		switch label := val[0].(type) {
		case string:
			// x@n
			if label == "_" {
				return nil, errors.New("Invalid CBOR: variable explicitly named _")
			}
			if len(val) == 2 {
				index, err := unwrapUint(val[1])
				if err != nil {
					return nil, err
				}
				return Var{Name: label, Index: int(index)}, nil
			}
		case uint64:
			switch label {
			case 0: // application
				if len(val) <= 2 {
					return nil, errors.New("Invalid CBOR: must have at least one arg in application")
				}
				fn, err := Decode(val[1])
				if err != nil {
					return nil, err
				}
				args := make([]Expr, len(val)-2)
				for i, arg := range val[2:] {
					args[i], err = Decode(arg)
					if err != nil {
						return nil, err
					}
				}
				return Apply(fn, args...), nil
			case 1, 2: // function or pi
				name := "_"
				var t, body Expr
				var err error
				switch len(val) {
				case 3: // implicit _ name
					t, err = Decode(val[1])
					if err != nil {
						return nil, err
					}
					body, err = Decode(val[2])
					if err != nil {
						return nil, err
					}
				case 4:
					name, err = unwrapString(val[1])
					if err != nil {
						return nil, err
					}
					if name == "_" {
						return nil, errors.New("Invalid CBOR: explicit _ variable name")
					}
					t, err = Decode(val[2])
					if err != nil {
						return nil, err
					}
					body, err = Decode(val[3])
					if err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("CBOR decode error: malformed function expression: %v", val)
				}
				if label == 1 {
					return &LambdaExpr{name, t, body}, nil
				} else {
					return &Pi{name, t, body}, nil
				}
			case 3: // operator
				if len(val) != 4 {
					return nil, fmt.Errorf("CBOR decode error: an operator takes exactly two arguments")
				}
				l, err := Decode(val[2])
				if err != nil {
					return nil, err
				}
				r, err := Decode(val[3])
				if err != nil {
					return nil, err
				}
				opcode, err := unwrapUint(val[1])
				if err != nil {
					return nil, err
				}
				if opcode > 11 {
					return nil, fmt.Errorf("CBOR decode error: unknown operator code %d", opcode)
				}
				return Operator{OpCode: int(opcode), L: l, R: r}, nil
			case 4: // list
				if val[1] != nil {
					if len(val) > 2 {
						return nil, fmt.Errorf("CBOR decode error: nonempty lists must not have an annotation in %v", val)
					}
					t, err := Decode(val[1])
					if err != nil {
						return nil, err
					}
					return EmptyList{Type: t}, nil
				}
				items := make([]Expr, len(val)-2)
				for i, rawItem := range val[2:] {
					var err error
					items[i], err = Decode(rawItem)
					if err != nil {
						return nil, err
					}
				}
				return NonEmptyList(items), nil
			case 5: // Some
				if len(val) != 3 || val[1] != nil {
					return nil, fmt.Errorf("CBOR decode error: malformed Some expression: %v", val)
				}
				val, err := Decode(val[2])
				if err != nil {
					return nil, err
				}
				return Some{val}, nil
			case 6: // merge
				var annotation Expr
				var err error
				switch len(val) {
				case 4:
					annotation, err = Decode(val[3])
					if err != nil {
						return nil, err
					}
					fallthrough
				case 3:
					l, err := Decode(val[1])
					if err != nil {
						return nil, err
					}
					r, err := Decode(val[2])
					if err != nil {
						return nil, err
					}
					return Merge{Handler: l, Union: r, Annotation: annotation}, nil
				default:
					return nil, fmt.Errorf("CBOR decode error: malformed merge expression: %v", val)
				}
			case 7, 8: // record type or literal
				m, err := decodeMap(val[1])
				if err != nil {
					return nil, err
				}
				if label == 7 {
					return Record(m), nil
				} else {
					return RecordLit(m), nil
				}
			case 9: // field access (r.x or u.x)
				recordOrUnionType, err := Decode(val[1])
				if err != nil {
					return nil, err
				}
				label, err := unwrapString(val[2])
				if err != nil {
					return nil, err
				}
				return Field{Record: recordOrUnionType, FieldName: label}, nil
				// case 10: // projection
			case 11: // union type
				m, err := decodeMap(val[1])
				if err != nil {
					return nil, err
				}
				return UnionType(m), nil
				// case 12: // union literal (deprecated)
				// case 13: // constructors (now removed)
			case 14: // if
				cond, err := Decode(val[1])
				if err != nil {
					return nil, err
				}
				tBranch, err := Decode(val[2])
				if err != nil {
					return nil, err
				}
				fBranch, err := Decode(val[3])
				if err != nil {
					return nil, err
				}
				return BoolIf{Cond: cond, T: tBranch, F: fBranch}, nil
			case 15: // natural literal
				n, err := unwrapUint(val[1])
				if err != nil {
					return nil, err
				}
				return NaturalLit(n), nil
			case 16: // integer literal
				n, err := unwrapInt(val[1])
				if err != nil {
					return nil, err
				}
				return IntegerLit(n), nil
			case 18: // text literal
				i := 1
				var chunks Chunks
				for ; i+1 < len(val); i = i + 2 {
					prefix, err := unwrapString(val[i])
					if err != nil {
						return nil, err
					}
					expr, err := Decode(val[i+1])
					if err != nil {
						return nil, err
					}
					chunks = append(chunks, Chunk{Prefix: prefix, Expr: expr})
				}
				s, err := unwrapString(val[i])
				if err != nil {
					return nil, err
				}
				return TextLit{Chunks: chunks, Suffix: s}, nil
			case 24: // imports
				importLabel, err := unwrapInt(val[3])
				if err != nil {
					return nil, err
				}
				var f Fetchable
				switch importLabel {
				case 0, 1:
					scheme := "https"
					if importLabel == 0 {
						scheme = "http"
					}
					authority, err := unwrapString(val[5])
					if err != nil {
						return nil, err
					}
					u, err := url.Parse(scheme + "://" + authority)
					if err != nil {
						return nil, err
					}
					urlPath := "/"
					rawPath := "/"
					for i := 6; i < len(val)-1; i++ {
						pathComponent, err := unwrapString(val[i])
						if err != nil {
							return nil, err
						}

						rawPath = path.Join(rawPath, url.PathEscape(pathComponent))
						urlPath = path.Join(urlPath, pathComponent)
					}
					u.Path = urlPath
					if urlPath != rawPath {
						u.RawPath = rawPath
					}
					if val[len(val)-1] != nil {
						u.RawQuery, err = unwrapString(val[len(val)-1])
						if err != nil {
							return nil, err
						}
						if u.RawQuery == "" {
							u.ForceQuery = true
						}
					}
					f = Remote{url: u}
				case 2, 3, 4, 5:
					var file string
					if importLabel == 2 {
						file = "/"
					} else if importLabel == 3 {
						file = "."
					} else if importLabel == 4 {
						file = ".."
					} else {
						file = "~"
					}
					for i := 4; i < len(val); i++ {
						component, err := unwrapString(val[i])
						if err != nil {
							return nil, err
						}
						file = path.Join(file, component)
					}
					f = Local(file)
				case 6:
					name, err := unwrapString(val[4])
					if err != nil {
						return nil, err
					}
					f = EnvVar(name)
				case 7:
					f = Missing{}
				default:
					return nil, fmt.Errorf("CBOR decode error: couldn't decode %#v", val)
				}
				return Embed(Import{ImportHashed: ImportHashed{Fetchable: f}}), nil
			case 25: // let
				if len(val)%3 != 2 {
					return nil, fmt.Errorf("CBOR decode error: unexpected array length %d when decoding let", len(val))
				}
				body, err := Decode(val[len(val)-1])
				if err != nil {
					return nil, err
				}
				var bindings []Binding
				for i := 1; i < len(val)-2; i = i + 3 {
					name, err := unwrapString(val[i])
					if err != nil {
						return nil, err
					}
					var annotation Expr
					if val[i+1] != nil {
						annotation, err = Decode(val[i+1])
						if err != nil {
							return nil, err
						}
					}
					value, err := Decode(val[i+2])
					if err != nil {
						return nil, err
					}
					bindings = append(bindings, Binding{name, annotation, value})
				}
				return MakeLet(body, bindings...), nil
			case 26: // annotated expression
				expr, err := Decode(val[1])
				if err != nil {
					return nil, err
				}
				annotation, err := Decode(val[2])
				if err != nil {
					return nil, err
				}
				return Annot{expr, annotation}, nil
			}
		}
	}
	return nil, fmt.Errorf("unimplemented while decoding %+v", decodedCbor)
}

// a marker type for CBOR-encoding purposes
type CborBox struct{ Content Expr }

var _ codec.Selfer = &CborBox{}

func Box(expr Expr) *CborBox { return &CborBox{Content: expr} }

func (box *CborBox) CodecEncodeSelf(e *codec.Encoder) {
	switch val := box.Content.(type) {
	case Var:
		if val.Name == "_" {
			e.Encode(val.Index)
		} else {
			e.Encode([]interface{}{val.Name, val.Index})
		}
	case Const:
		switch val {
		case Type:
			e.Encode("Type")
		case Kind:
			e.Encode("Kind")
		case Sort:
			e.Encode("Sort")
		default:
			panic(fmt.Sprintf("unknown type %d\n", val))
		}
	case Builtin:
		e.Encode(string(val))
	case *App:
		fn := val.Fn
		args := []interface{}{Box(val.Arg)}
		for true {
			parentapp, ok := fn.(*App)
			if !ok {
				break
			}
			fn = parentapp.Fn
			args = append([]interface{}{Box(parentapp.Arg)}, args...)
		}
		e.Encode(append([]interface{}{0, Box(fn)}, args...))

	case *LambdaExpr:
		if val.Label == "_" {
			e.Encode([]interface{}{1, Box(val.Type), Box(val.Body)})
		} else {
			e.Encode([]interface{}{1, val.Label, Box(val.Type), Box(val.Body)})
		}
	case *Pi:
		if val.Label == "_" {
			e.Encode([]interface{}{2, Box(val.Type), Box(val.Body)})
		} else {
			e.Encode([]interface{}{2, val.Label, Box(val.Type), Box(val.Body)})
		}
	case Operator:
		e.Encode([]interface{}{3, val.OpCode, Box(val.L), Box(val.R)})
	case EmptyList:
		e.Encode([]interface{}{4, Box(val.Type)})
	case NonEmptyList:
		items := []Expr(val)
		output := make([]interface{}, len(items)+2)
		output[0] = 4
		output[1] = nil
		for i, item := range items {
			output[i+2] = Box(item)
		}
		e.Encode(output)
	case Some:
		e.Encode([]interface{}{5, nil, Box(val.Val)})
	case Merge:
		if val.Annotation != nil {
			e.Encode([]interface{}{6, Box(val.Handler), Box(val.Union), Box(val.Annotation)})
		} else {
			e.Encode([]interface{}{6, Box(val.Handler), Box(val.Union)})
		}
	case Record:
		items := make(map[string]*CborBox)
		for k, v := range val {
			items[k] = Box(v)
		}
		// we rely on the EncodeOptions having Canonical set
		// so that we get sorted keys in our map
		output := []interface{}{7, items}
		e.Encode(output)
	case RecordLit:
		items := make(map[string]*CborBox)
		for k, v := range val {
			items[k] = Box(v)
		}
		// we rely on the EncodeOptions having Canonical set
		// so that we get sorted keys in our map
		output := []interface{}{8, items}
		e.Encode(output)
	case Field:
		e.Encode([]interface{}{9, Box(val.Record), val.FieldName})
	case UnionType:
		items := make(map[string]*CborBox)
		for k, v := range val {
			items[k] = Box(v)
		}
		// we rely on the EncodeOptions having Canonical set
		// so that we get sorted keys in our map
		output := []interface{}{11, items}
		e.Encode(output)
	case BoolLit:
		e.Encode(bool(val))
	case BoolIf:
		e.Encode([]interface{}{14, Box(val.Cond), Box(val.T), Box(val.F)})
	case NaturalLit:
		e.Encode(append([]interface{}{15}, int(val)))
	case IntegerLit:
		e.Encode(append([]interface{}{16}, int(val)))
	case DoubleLit:
		// TODO: size appropriately
		single := float32(val)
		if float64(single) == float64(val) {
			e.Encode(single)
		} else {
			e.Encode(float64(val))
		}
	case TextLit:
		output := []interface{}{18}
		for _, chunk := range val.Chunks {
			output = append(output, chunk.Prefix, Box(chunk.Expr))
		}
		output = append(output, val.Suffix)
		e.Encode(output)
	case Embed:
		r := val.Fetchable
		mode := 0
		if val.ImportMode == RawText {
			mode = 1
		}
		switch rr := r.(type) {
		case EnvVar:
			e.Encode([]interface{}{24, val.Hash, mode, 6, string(rr)})
		case Local:
			if rr.IsAbs() {
				toEncode := []interface{}{24, val.Hash, mode, AbsoluteImport}
				for _, component := range rr.PathComponents() {
					toEncode = append(toEncode, component)
				}
				e.Encode(toEncode)
			} else if rr.IsRelativeToParent() {
				toEncode := []interface{}{24, val.Hash, mode, ParentImport}
				for _, component := range rr.PathComponents() {
					toEncode = append(toEncode, component)
				}
				e.Encode(toEncode)
			} else if rr.IsRelativeToHome() {
				toEncode := []interface{}{24, val.Hash, mode, HomeImport}
				for _, component := range rr.PathComponents() {
					toEncode = append(toEncode, component)
				}
				e.Encode(toEncode)
			} else {
				toEncode := []interface{}{24, val.Hash, mode, HereImport}
				for _, component := range rr.PathComponents() {
					toEncode = append(toEncode, component)
				}
				e.Encode(toEncode)
			}
		case Remote:
			var headers interface{} // unimplemented, leave as nil for now
			scheme := HttpsImport
			if rr.IsPlainHttp() {
				scheme = HttpImport
			}
			toEncode := []interface{}{24, val.Hash, mode, scheme, headers, rr.Authority()}
			for _, component := range rr.PathComponents() {
				toEncode = append(toEncode, component)
			}
			toEncode = append(toEncode, rr.Query())
			e.Encode(toEncode)
		case Missing:
			e.Encode([]interface{}{24, nil, 0, 7})
		default:
			panic("can't happen")
		}
	case Let:
		output := make([]interface{}, len(val.Bindings)*3+2)
		output[0] = 25
		for i, binding := range val.Bindings {
			output[3*i+1] = binding.Variable
			output[3*i+2] = Box(binding.Annotation)
			output[3*i+3] = Box(binding.Value)
		}
		output[len(output)-1] = Box(val.Body)
		e.Encode(output)
	case Annot:
		e.Encode([]interface{}{26, Box(val.Expr), Box(val.Annotation)})
	default:
		e.Encode(box.Content)
	}
}

func (box *CborBox) CodecDecodeSelf(d *codec.Decoder) {
	var raw interface{}
	d.MustDecode(&raw)
	expr, err := Decode(raw)
	if err != nil {
		panic(err)
	}
	box.Content = expr
}

func NewCborHandle() codec.CborHandle {
	var h codec.CborHandle
	h.Canonical = true
	h.SkipUnexpectedTags = true
	return h
}

const (
	HttpImport     = 0
	HttpsImport    = 1
	AbsoluteImport = 2
	HereImport     = 3
	ParentImport   = 4
	HomeImport     = 5
)
