package binary

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"

	"github.com/fxamacker/cbor/v2"
	. "github.com/philandstuff/dhall-golang/v4/term"
)

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

	"Integer/clamp":    IntegerClamp,
	"Integer/negate":   IntegerNegate,
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
		// Type, Double, List/fold
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
					return Lambda{name, t, body}, nil
				} else {
					return Pi{name, t, body}, nil
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
				return Op{OpCode: OpCode(opcode), L: l, R: r}, nil
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
				return If{Cond: cond, T: tBranch, F: fBranch}, nil
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
				return TextLit{Chunks: chunks, Suffix: s}, nil
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
					f = NewRemoteFile(u)
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
					f = LocalFile(file)
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

// EncodeAsCbor encodes a Term as CBOR and writes it to the io.Writer
func EncodeAsCbor(w io.Writer, e Term) error {
	em, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return err
	}
	return em.NewEncoder(w).Encode(e)
}

// DecodeAsCbor decodes CBOR from the io.Reader and returns the resulting Expr
func DecodeAsCbor(r io.Reader) (Term, error) {
	var b interface{}
	dm, err := cbor.DecOptions{MaxNestedLevels: 64}.DecMode()
	if err != nil {
		return nil, err
	}
	dec := dm.NewDecoder(r)
	err = dec.Decode(&b)
	if err != nil {
		return nil, err
	}
	return decode(b)
}
