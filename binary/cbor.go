package binary

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"path"

	. "github.com/philandstuff/dhall-golang/core"
	"github.com/ugorji/go/codec"
)

var cbor = newCborHandle()

var nameToBuiltin = map[string]Term{
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
	"Natural/subtract":  NaturalSubtract,

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

func decodeMap(i interface{}) (map[string]Term, error) {
	if val, ok := i.(map[interface{}]interface{}); ok {
		decodedM := make(map[string]Term, len(val))
		for n, t := range val {
			name, err := unwrapString(n)
			if err != nil {
				return nil, err
			}
			if t == nil {
				decodedM[name] = nil
			} else {
				decodedM[name], err = decode(t)
				if err != nil {
					return nil, err
				}
			}
		}
		return decodedM, nil
	}
	return nil, fmt.Errorf("couldn't interpret %v as map[string]interface{}", i)
}

func decode(decodedCbor interface{}) (Term, error) {
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
		return BoolLit(val), nil
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
				fn, err := decode(val[1])
				if err != nil {
					return nil, err
				}
				args := make([]Term, len(val)-2)
				for i, arg := range val[2:] {
					args[i], err = decode(arg)
					if err != nil {
						return nil, err
					}
				}
				return Apply(fn, args...), nil
			case 1, 2: // function or pi
				name := "_"
				var t, body Term
				var err error
				switch len(val) {
				case 3: // implicit _ name
					t, err = decode(val[1])
					if err != nil {
						return nil, err
					}
					body, err = decode(val[2])
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
					t, err = decode(val[2])
					if err != nil {
						return nil, err
					}
					body, err = decode(val[3])
					if err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("CBOR decode error: malformed function expression: %v", val)
				}
				if label == 1 {
					return LambdaTerm{name, t, body}, nil
				} else {
					return PiTerm{name, t, body}, nil
				}
			case 3: // operator
				if len(val) != 4 {
					return nil, fmt.Errorf("CBOR decode error: an operator takes exactly two arguments")
				}
				l, err := decode(val[2])
				if err != nil {
					return nil, err
				}
				r, err := decode(val[3])
				if err != nil {
					return nil, err
				}
				opcode, err := unwrapUint(val[1])
				if err != nil {
					return nil, err
				}
				if opcode > 13 {
					return nil, fmt.Errorf("CBOR decode error: unknown operator code %d", opcode)
				}
				return OpTerm{OpCode: int(opcode), L: l, R: r}, nil
			case 4: // list
				if val[1] != nil {
					if len(val) > 2 {
						return nil, fmt.Errorf("CBOR decode error: nonempty lists must not have an annotation in %v", val)
					}
					t, err := decode(val[1])
					if err != nil {
						return nil, err
					}
					return EmptyList{Type: Apply(List, t)}, nil
				}
				items := make(NonEmptyList, len(val)-2)
				for i, rawItem := range val[2:] {
					var err error
					items[i], err = decode(rawItem)
					if err != nil {
						return nil, err
					}
				}
				return items, nil
			case 5: // Some
				if len(val) != 3 || val[1] != nil {
					return nil, fmt.Errorf("CBOR decode error: malformed Some expression: %v", val)
				}
				val, err := decode(val[2])
				if err != nil {
					return nil, err
				}
				return Some{val}, nil
			case 6: // merge
				var annotation Term
				var err error
				switch len(val) {
				case 4:
					annotation, err = decode(val[3])
					if err != nil {
						return nil, err
					}
					fallthrough
				case 3:
					l, err := decode(val[1])
					if err != nil {
						return nil, err
					}
					r, err := decode(val[2])
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
					return RecordType(m), nil
				} else {
					return RecordLit(m), nil
				}
			case 9: // field access (r.x or u.x)
				recordOrUnionType, err := decode(val[1])
				if err != nil {
					return nil, err
				}
				label, err := unwrapString(val[2])
				if err != nil {
					return nil, err
				}
				return Field{Record: recordOrUnionType, FieldName: label}, nil
			case 10: // projection
				record, err := decode(val[1])
				if err != nil {
					return nil, err
				}
				switch val[2].(type) {
				case string: // r.{x,y}
					fieldNames := make([]string, len(val)-2)
					for i, fieldNameWrapped := range val[2:] {
						fieldName, err := unwrapString(fieldNameWrapped)
						if err != nil {
							return nil, err
						}
						fieldNames[i] = fieldName
					}
					return Project{
						Record:     record,
						FieldNames: fieldNames,
					}, nil
				case []interface{}: // r.(t)
					projectType, err := decode(val[2].([]interface{})[0])
					if err != nil {
						return nil, err
					}
					return ProjectType{
						Record:   record,
						Selector: projectType,
					}, nil
				}
			case 11: // union type
				m, err := decodeMap(val[1])
				if err != nil {
					return nil, err
				}
				return UnionType(m), nil
				// case 12: // union literal (deprecated)
				// case 13: // constructors (now removed)
			case 14: // if
				cond, err := decode(val[1])
				if err != nil {
					return nil, err
				}
				tBranch, err := decode(val[2])
				if err != nil {
					return nil, err
				}
				fBranch, err := decode(val[3])
				if err != nil {
					return nil, err
				}
				return IfTerm{Cond: cond, T: tBranch, F: fBranch}, nil
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
					expr, err := decode(val[i+1])
					if err != nil {
						return nil, err
					}
					chunks = append(chunks, Chunk{Prefix: prefix, Expr: expr})
				}
				s, err := unwrapString(val[i])
				if err != nil {
					return nil, err
				}
				return TextLitTerm{Chunks: chunks, Suffix: s}, nil
			case 19: // assert
				annot, err := decode(val[1])
				if err != nil {
					return nil, err
				}
				return Assert{Annotation: annot}, nil
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
					urlPath := "/"
					for i := 6; i < len(val)-1; i++ {
						pathComponent, err := unwrapString(val[i])
						if err != nil {
							return nil, err
						}

						urlPath = path.Join(urlPath, pathComponent)
					}
					query := ""
					if val[len(val)-1] != nil {
						content, err := unwrapString(val[len(val)-1])
						if err != nil {
							return nil, err
						}
						query = "?" + content
					}
					u, err := url.Parse(scheme + "://" + authority + urlPath + query)
					if err != nil {
						return nil, err
					}
					f = NewRemote(u)
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
				return Import{ImportHashed: ImportHashed{Fetchable: f}}, nil
			case 25: // let
				if len(val)%3 != 2 {
					return nil, fmt.Errorf("CBOR decode error: unexpected array length %d when decoding let", len(val))
				}
				body, err := decode(val[len(val)-1])
				if err != nil {
					return nil, err
				}
				var bindings []Binding
				for i := 1; i < len(val)-2; i = i + 3 {
					name, err := unwrapString(val[i])
					if err != nil {
						return nil, err
					}
					var annotation Term
					if val[i+1] != nil {
						annotation, err = decode(val[i+1])
						if err != nil {
							return nil, err
						}
					}
					value, err := decode(val[i+2])
					if err != nil {
						return nil, err
					}
					bindings = append(bindings, Binding{name, annotation, value})
				}
				return NewLet(body, bindings...), nil
			case 26: // annotated expression
				expr, err := decode(val[1])
				if err != nil {
					return nil, err
				}
				annotation, err := decode(val[2])
				if err != nil {
					return nil, err
				}
				return Annot{expr, annotation}, nil
			case 27: // toMap
				record, err := decode(val[1])
				if err != nil {
					return nil, err
				}
				output := ToMap{Record: record}
				if len(val) > 2 {
					typ, err := decode(val[2])
					if err != nil {
						return nil, err
					}
					output.Type = typ
				}
				return output, nil
			case 28: // [] : T -- but not in form [] : List T
				t, err := decode(val[1])
				if err != nil {
					return nil, err
				}
				return EmptyList{Type: t}, nil
			}
		}
	}
	return nil, fmt.Errorf("unimplemented while decoding %+v", decodedCbor)
}

// a marker type for CBOR-encoding purposes
type cborBox struct{ content Term }

var _ codec.Selfer = &cborBox{}

func box(expr Term) *cborBox { return &cborBox{content: expr} }

func (b *cborBox) CodecEncodeSelf(e *codec.Encoder) {
	switch val := b.content.(type) {
	case Var:
		if val.Name == "_" {
			e.Encode(val.Index)
		} else {
			e.Encode([]interface{}{val.Name, val.Index})
		}
	case Universe:
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
	case AppTerm:
		fn := val.Fn
		args := []interface{}{box(val.Arg)}
		for true {
			parentapp, ok := fn.(AppTerm)
			if !ok {
				break
			}
			fn = parentapp.Fn
			args = append([]interface{}{box(parentapp.Arg)}, args...)
		}
		e.Encode(append([]interface{}{0, box(fn)}, args...))

	case LambdaTerm:
		if val.Label == "_" {
			e.Encode([]interface{}{1, box(val.Type), box(val.Body)})
		} else {
			e.Encode([]interface{}{1, val.Label, box(val.Type), box(val.Body)})
		}
	case PiTerm:
		if val.Label == "_" {
			e.Encode([]interface{}{2, box(val.Type), box(val.Body)})
		} else {
			e.Encode([]interface{}{2, val.Label, box(val.Type), box(val.Body)})
		}
	case OpTerm:
		e.Encode([]interface{}{3, val.OpCode, box(val.L), box(val.R)})
	case EmptyList:
		if app, ok := val.Type.(AppTerm); ok {
			if app.Fn == List {
				e.Encode([]interface{}{4, box(app.Arg)})
				break
			}
		}
		e.Encode([]interface{}{28, box(val.Type)})
	case NonEmptyList:
		output := make([]interface{}, len(val)+2)
		output[0] = 4
		output[1] = nil
		for i, item := range val {
			output[i+2] = box(item)
		}
		e.Encode(output)
	case Some:
		e.Encode([]interface{}{5, nil, box(val.Val)})
	case Merge:
		if val.Annotation != nil {
			e.Encode([]interface{}{6, box(val.Handler), box(val.Union), box(val.Annotation)})
		} else {
			e.Encode([]interface{}{6, box(val.Handler), box(val.Union)})
		}
	case RecordType:
		items := make(map[string]*cborBox)
		for k, v := range val {
			items[k] = box(v)
		}
		// we rely on the EncodeOptions having Canonical set
		// so that we get sorted keys in our map
		e.Encode([]interface{}{7, items})
	case RecordLit:
		items := make(map[string]*cborBox)
		for k, v := range val {
			items[k] = box(v)
		}
		// we rely on the EncodeOptions having Canonical set
		// so that we get sorted keys in our map
		e.Encode([]interface{}{8, items})
	case ToMap:
		if val.Type != nil {
			e.Encode([]interface{}{27, box(val.Record), box(val.Type)})
		} else {
			e.Encode([]interface{}{27, box(val.Record)})
		}
	case Field:
		e.Encode([]interface{}{9, box(val.Record), val.FieldName})
	case Project:
		output := make([]interface{}, len(val.FieldNames)+2)
		output[0] = 10
		output[1] = box(val.Record)
		for i, name := range val.FieldNames {
			output[i+2] = name
		}
		e.Encode(output)
	case ProjectType:
		e.Encode([]interface{}{
			10,
			box(val.Record),
			[]interface{}{
				box(val.Selector),
			}})
	case UnionType:
		items := make(map[string]*cborBox)
		for k, v := range val {
			items[k] = box(v)
		}
		// we rely on the EncodeOptions having Canonical set
		// so that we get sorted keys in our map
		e.Encode([]interface{}{11, items})
	case BoolLit:
		e.Encode(bool(val))
	case IfTerm:
		e.Encode([]interface{}{14, box(val.Cond), box(val.T), box(val.F)})
	case NaturalLit:
		e.Encode(append([]interface{}{15}, int(val)))
	case IntegerLit:
		e.Encode(append([]interface{}{16}, int(val)))
	case DoubleLit:
		// special-case values to encode as float16
		if float64(val) == 0.0 { // 0.0
			if math.Signbit(float64(val)) {
				e.Encode(codec.Raw([]byte{0xf9, 0x80, 0x00}))
			} else {
				e.Encode(codec.Raw([]byte{0xf9, 0x00, 0x00}))
			}
		} else if math.IsNaN(float64(val)) { // NaN
			e.Encode(codec.Raw([]byte{0xf9, 0x7e, 0x00}))
		} else if math.IsInf(float64(val), 1) { // Infinity
			e.Encode(codec.Raw([]byte{0xf9, 0x7c, 0x00}))
		} else if math.IsInf(float64(val), -1) { // -Infinity
			e.Encode(codec.Raw([]byte{0xf9, 0xfc, 0x00}))
		} else {
			single := float32(val)
			if float64(single) == float64(val) {
				e.Encode(single)
			} else {
				e.Encode(float64(val))
			}
		}
	case TextLitTerm:
		output := []interface{}{18}
		for _, chunk := range val.Chunks {
			output = append(output, chunk.Prefix, box(chunk.Expr))
		}
		output = append(output, val.Suffix)
		e.Encode(output)
	case Assert:
		e.Encode([]interface{}{19, box(val.Annotation)})
	case Import:
		r := val.Fetchable
		// we have crafted the ImportMode constants to match the expected CBOR values
		mode := val.ImportMode
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
			e.Encode([]interface{}{24, nil, mode, 7})
		default:
			panic("can't happen")
		}
	case Let:
		output := []interface{}{25}
		// flatten multiple lets into one
		for {
			for _, binding := range val.Bindings {
				output = append(output, binding.Variable)
				output = append(output, box(binding.Annotation))
				output = append(output, box(binding.Value))
			}
			// there's probably a nicer way to do this...
			nextLet, ok := val.Body.(Let)
			if !ok {
				break
			}
			val = nextLet
		}
		output = append(output, box(val.Body))
		e.Encode(output)
	case Annot:
		e.Encode([]interface{}{26, box(val.Expr), box(val.Annotation)})
	default:
		e.Encode(b.content)
	}
}

func (b *cborBox) CodecDecodeSelf(d *codec.Decoder) {
	var raw interface{}
	d.MustDecode(&raw)
	expr, err := decode(raw)
	if err != nil {
		panic(err)
	}
	b.content = expr
}

// EncodeAsCbor encodes a Term as CBOR and writes it to the io.Writer
func EncodeAsCbor(w io.Writer, e Term) error {
	enc := codec.NewEncoder(w, cbor)
	return enc.Encode(box(e))
}

// DecodeAsCbor decodes CBOR from the io.Reader and returns the resulting Expr
func DecodeAsCbor(r io.Reader) (Term, error) {
	var b cborBox
	dec := codec.NewDecoder(r, cbor)
	err := dec.Decode(&b)
	return b.content, err
}

func newCborHandle() *codec.CborHandle {
	var h codec.CborHandle
	h.Canonical = true
	h.SkipUnexpectedTags = true
	h.Raw = true
	return &h
}

const (
	HttpImport     = 0
	HttpsImport    = 1
	AbsoluteImport = 2
	HereImport     = 3
	ParentImport   = 4
	HomeImport     = 5
)
