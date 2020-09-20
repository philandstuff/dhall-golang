package term

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/fxamacker/cbor/v2"
)

func encMode() cbor.EncMode {
	em, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		panic(err)
	}
	return em
}

var em = encMode()

// MarshalCBOR implements cbor.Marshaler.
func (u Universe) MarshalCBOR() ([]byte, error) {
	switch u {
	case Type:
		return cbor.Marshal("Type")
	case Kind:
		return cbor.Marshal("Kind")
	case Sort:
		return cbor.Marshal("Sort")
	}
	return nil, fmt.Errorf("Invalid Universe value %d", u)
}

// MarshalCBOR implements cbor.Marshaler.
func (v Var) MarshalCBOR() ([]byte, error) {
	if v.Name == "_" {
		return cbor.Marshal(v.Index)
	}
	return cbor.Marshal([]interface{}{v.Name, v.Index})
}

// MarshalCBOR implements cbor.Marshaler.
func (app App) MarshalCBOR() ([]byte, error) {
	fn := app.Fn
	args := []interface{}{app.Arg}
	for true {
		childapp, ok := fn.(App)
		if !ok {
			break
		}
		fn = childapp.Fn
		args = append([]interface{}{childapp.Arg}, args...)
	}
	return em.Marshal(append([]interface{}{0, fn}, args...))
}

// MarshalCBOR implements cbor.Marshaler.
func (l Lambda) MarshalCBOR() ([]byte, error) {
	if l.Label == "_" {
		return em.Marshal([]interface{}{1, l.Type, l.Body})
	}
	return em.Marshal([]interface{}{1, l.Label, l.Type, l.Body})
}

// MarshalCBOR implements cbor.Marshaler.
func (p Pi) MarshalCBOR() ([]byte, error) {
	if p.Label == "_" {
		return em.Marshal([]interface{}{2, p.Type, p.Body})
	}
	return em.Marshal([]interface{}{2, p.Label, p.Type, p.Body})
}

// MarshalCBOR implements cbor.Marshaler.
func (op Op) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]interface{}{3, op.OpCode, op.L, op.R})
}

// MarshalCBOR implements cbor.Marshaler.
func (l EmptyList) MarshalCBOR() ([]byte, error) {
	if app, ok := l.Type.(App); ok {
		if app.Fn == List {
			return em.Marshal([]interface{}{4, app.Arg})
		}
	}
	return em.Marshal([]interface{}{28, l.Type})
}

// MarshalCBOR implements cbor.Marshaler.
func (l NonEmptyList) MarshalCBOR() ([]byte, error) {
	output := make([]interface{}, len(l)+2)
	output[0] = 4
	output[1] = nil
	for i, item := range l {
		output[i+2] = item
	}
	return em.Marshal(output)
}

// MarshalCBOR implements cbor.Marshaler.
func (s Some) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]interface{}{5, nil, s.Val})
}

// MarshalCBOR implements cbor.Marshaler.
func (m Merge) MarshalCBOR() ([]byte, error) {
	if m.Annotation == nil {
		return em.Marshal([]interface{}{6, m.Handler, m.Union})
	}
	return em.Marshal([]interface{}{6, m.Handler, m.Union, m.Annotation})
}

// MarshalCBOR implements cbor.Marshaler.
func (r RecordType) MarshalCBOR() ([]byte, error) {
	mapCbor, err := marshalMapInOrder(r)
	if err != nil {
		return nil, err
	}
	return em.Marshal([]interface{}{7, mapCbor})
}

// marshalMapInOrder is a bit of a hack to allow us to marshal a map
// as CBOR in Dhall order, not RFC7049-canonical order.  The key
// difference is that RFC7049 order sorts by the encoded bytes, which
// means in practice that shorter strings precede longer strings
// regardless of contents.
func marshalMapInOrder(m map[string]Term) (cbor.RawMessage, error) {
	// we will store up our CBOR in out
	out := new(bytes.Buffer)

	// first, the length marker for the map
	lenCbor, err := em.Marshal(len(m))
	if err != nil {
		return nil, err
	}
	// hack to convert major type 0 (unsigned int) to major type 5
	// (map, starting with length)
	lenCbor[0] = lenCbor[0] | 0xa0
	out.Write(lenCbor)

	// get keys in sorted order
	var sortedKeys []string
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	e := em.NewEncoder(out)
	for _, k := range sortedKeys {
		err := e.Encode(k)
		if err != nil {
			return nil, err
		}
		err = e.Encode(m[k])
		if err != nil {
			return nil, err
		}
	}
	return out.Bytes(), nil
}

// MarshalCBOR implements cbor.Marshaler.
func (r RecordLit) MarshalCBOR() ([]byte, error) {
	mapCbor, err := marshalMapInOrder(r)
	if err != nil {
		return nil, err
	}
	return em.Marshal([]interface{}{8, mapCbor})
}

// MarshalCBOR implements cbor.Marshaler.
func (field Field) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]interface{}{9, field.Record, field.FieldName})
}

// MarshalCBOR implements cbor.Marshaler.
func (p Project) MarshalCBOR() ([]byte, error) {
	out := []interface{}{10, p.Record}
	for _, name := range p.FieldNames {
		out = append(out, name)
	}
	return em.Marshal(out)
}

// MarshalCBOR implements cbor.Marshaler.
func (p ProjectType) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]interface{}{
		10,
		p.Record,
		[]interface{}{
			p.Selector,
		}})
}

// MarshalCBOR implements cbor.Marshaler.
func (u UnionType) MarshalCBOR() ([]byte, error) {
	mapCbor, err := marshalMapInOrder(u)
	if err != nil {
		return nil, err
	}
	return em.Marshal([]interface{}{11, mapCbor})
}

// MarshalCBOR implements cbor.Marshaler.
func (i If) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]interface{}{14, i.Cond, i.T, i.F})
}

// MarshalCBOR implements cbor.Marshaler.
func (n NaturalLit) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal([]interface{}{15, uint(n)})
}

// MarshalCBOR implements cbor.Marshaler.
func (i IntegerLit) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal([]interface{}{16, int(i)})
}

// MarshalCBOR implements cbor.Marshaler.
func (t TextLit) MarshalCBOR() ([]byte, error) {
	output := []interface{}{18}
	for _, chunk := range t.Chunks {
		output = append(output, chunk.Prefix, chunk.Expr)
	}
	output = append(output, t.Suffix)
	return em.Marshal(output)
}

// MarshalCBOR implements cbor.Marshaler.
func (a Assert) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]interface{}{19, a.Annotation})
}

// MarshalCBOR implements cbor.Marshaler.
func (w With) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]interface{}{29, w.Record, w.Path, w.Value})
}

const (
	httpImport     = 0
	httpsImport    = 1
	absoluteImport = 2
	hereImport     = 3
	parentImport   = 4
	homeImport     = 5
)

// MarshalCBOR implements cbor.Marshaler.
func (i Import) MarshalCBOR() ([]byte, error) {
	// TODO: push this down onto EnvVar/LocalFile/RemoteFile/Missing?
	r := i.Fetchable
	// we have crafted the ImportMode constants to match the expected CBOR values
	mode := i.ImportMode
	switch rr := r.(type) {
	case EnvVar:
		return cbor.Marshal([]interface{}{24, i.Hash, mode, 6, string(rr)})
	case LocalFile:
		if rr.IsAbs() {
			toEncode := []interface{}{24, i.Hash, mode, absoluteImport}
			for _, component := range rr.PathComponents() {
				toEncode = append(toEncode, component)
			}
			return cbor.Marshal(toEncode)
		}
		if rr.IsRelativeToParent() {
			toEncode := []interface{}{24, i.Hash, mode, parentImport}
			for _, component := range rr.PathComponents() {
				toEncode = append(toEncode, component)
			}
			return cbor.Marshal(toEncode)
		}
		if rr.IsRelativeToHome() {
			toEncode := []interface{}{24, i.Hash, mode, homeImport}
			for _, component := range rr.PathComponents() {
				toEncode = append(toEncode, component)
			}
			return cbor.Marshal(toEncode)
		}

		toEncode := []interface{}{24, i.Hash, mode, hereImport}
		for _, component := range rr.PathComponents() {
			toEncode = append(toEncode, component)
		}
		return cbor.Marshal(toEncode)
	case RemoteFile:
		var headers interface{} // unimplemented, leave as nil for now
		scheme := httpsImport
		if rr.IsPlainHTTP() {
			scheme = httpImport
		}
		toEncode := []interface{}{24, i.Hash, mode, scheme, headers, rr.Authority()}
		for _, component := range rr.PathComponents() {
			toEncode = append(toEncode, component)
		}
		toEncode = append(toEncode, rr.Query())
		return cbor.Marshal(toEncode)
	case Missing:
		return cbor.Marshal([]interface{}{24, nil, mode, 7})
	default:
		panic("can't happen")
	}
}

// MarshalCBOR implements cbor.Marshaler.
func (l Let) MarshalCBOR() ([]byte, error) {
	output := []interface{}{25}
	// flatten multiple lets into one
	for {
		for _, binding := range l.Bindings {
			output = append(output, binding.Variable)
			output = append(output, binding.Annotation)
			output = append(output, binding.Value)
		}
		// there's probably a nicer way to do this...
		nextLet, ok := l.Body.(Let)
		if !ok {
			break
		}
		l = nextLet
	}
	output = append(output, l.Body)
	return em.Marshal(output)
}

// MarshalCBOR implements cbor.Marshaler.
func (a Annot) MarshalCBOR() ([]byte, error) {
	return em.Marshal([]interface{}{26, a.Expr, a.Annotation})
}

// MarshalCBOR implements cbor.Marshaler.
func (t ToMap) MarshalCBOR() ([]byte, error) {
	if t.Type != nil {
		return em.Marshal([]interface{}{27, t.Record, t.Type})
	}
	return em.Marshal([]interface{}{27, t.Record})
}
