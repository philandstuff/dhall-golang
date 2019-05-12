package ast

import (
	"fmt"

	"github.com/ugorji/go/codec"
)

var _ codec.Selfer = Type // Const
var _ codec.Selfer = Var{}
var _ codec.Selfer = &LambdaExpr{}
var _ codec.Selfer = &Pi{}
var _ codec.Selfer = &App{}
var _ codec.Selfer = Let{}
var _ codec.Selfer = Annot{}
var _ codec.Selfer = Double // Builtin
var _ codec.Selfer = BoolIf{}
var _ codec.Selfer = EmptyList{}
var _ codec.Selfer = NonEmptyList{}
var _ codec.Selfer = Operator{}
var _ codec.Selfer = TextLit{}
var _ codec.Selfer = NaturalLit(0)
var _ codec.Selfer = IntegerLit(0)
var _ codec.Selfer = Some{}
var _ codec.Selfer = Record{}
var _ codec.Selfer = RecordLit{}
var _ codec.Selfer = Field{}
var _ codec.Selfer = UnionType{}
var _ codec.Selfer = Merge{}
var _ codec.Selfer = Embed{}

func NewCborHandle() codec.CborHandle {
	var h codec.CborHandle
	h.Canonical = true
	return h
}

func (t Const) CodecEncodeSelf(e *codec.Encoder) {
	switch t {
	case Type:
		e.Encode("Type")
	case Kind:
		e.Encode("Kind")
	case Sort:
		e.Encode("Sort")
	default:
		panic(fmt.Sprintf("unknown type %d\n", t))
	}
}

func (v Var) CodecEncodeSelf(e *codec.Encoder) {
	if v.Name == "_" {
		e.Encode(v.Index)
	} else {
		e.Encode([]interface{}{v.Name, v.Index})
	}
}

func (l *LambdaExpr) CodecEncodeSelf(e *codec.Encoder) {
	if l.Label == "_" {
		e.Encode([]interface{}{1, l.Type, l.Body})
	} else {
		e.Encode([]interface{}{1, l.Label, l.Type, l.Body})
	}
}

func (p *Pi) CodecEncodeSelf(e *codec.Encoder) {
	if p.Label == "_" {
		e.Encode([]interface{}{2, p.Type, p.Body})
	} else {
		e.Encode([]interface{}{2, p.Label, p.Type, p.Body})
	}
}

func (a *App) CodecEncodeSelf(e *codec.Encoder) {
	fn := a.Fn
	args := []interface{}{a.Arg}
	for true {
		parentapp, ok := fn.(*App)
		if !ok {
			break
		}
		fn = parentapp.Fn
		args = append([]interface{}{parentapp.Arg}, args...)
	}
	e.Encode(append([]interface{}{0, fn}, args...))
}

func (l Let) CodecEncodeSelf(e *codec.Encoder) {
	output := make([]interface{}, len(l.Bindings)*3+2)
	output[0] = 25
	for i, binding := range l.Bindings {
		output[3*i+1] = binding.Variable
		output[3*i+2] = binding.Annotation
		output[3*i+3] = binding.Value
	}
	output[len(output)-1] = l.Body
	e.Encode(output)
}

func (a Annot) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode([]interface{}{26, a.Expr, a.Annotation})
}

func (t Builtin) CodecEncodeSelf(e *codec.Encoder) {
	switch t {
	case Double:
		e.Encode("Double")
	case Text:
		e.Encode("Text")
	case Bool:
		e.Encode("Bool")
	case Natural:
		e.Encode("Natural")
	case Integer:
		e.Encode("Integer")
	case List:
		e.Encode("List")
	case Optional:
		e.Encode("Optional")
	case None:
		e.Encode("None")
	default:
		panic(fmt.Sprintf("unknown type %d\n", t))
	}
}

func (bi BoolIf) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode([]interface{}{14, bi.Cond, bi.T, bi.F})
}

func (t TextLit) CodecEncodeSelf(e *codec.Encoder) {
	output := []interface{}{18}
	for _, chunk := range t.Chunks {
		output = append(output, chunk.Prefix, chunk.Expr)
	}
	output = append(output, t.Suffix)
	e.Encode(output)
}

func (n NaturalLit) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode(append([]interface{}{15}, int(n)))
}

func (op Operator) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode([]interface{}{3, op.OpCode, op.L, op.R})
}

func (i IntegerLit) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode(append([]interface{}{16}, int(i)))
}

func (l EmptyList) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode([]interface{}{4, l.Type})
}

func (l NonEmptyList) CodecEncodeSelf(e *codec.Encoder) {
	items := []Expr(l)
	output := make([]interface{}, len(items)+2)
	output[0] = 4
	output[1] = nil
	for i, item := range items {
		output[i+2] = item
	}
	e.Encode(output)
}

func (s Some) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode([]interface{}{5, nil, s.Val})
}

func (r Record) CodecEncodeSelf(e *codec.Encoder) {
	items := map[string]Expr(r)
	// we rely on the EncodeOptions having Canonical set
	// so that we get sorted keys in our map
	output := []interface{}{7, items}
	e.Encode(output)
}

func (r RecordLit) CodecEncodeSelf(e *codec.Encoder) {
	items := map[string]Expr(r)
	// we rely on the EncodeOptions having Canonical set
	// so that we get sorted keys in our map
	output := []interface{}{8, items}
	e.Encode(output)
}

func (f Field) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode([]interface{}{9, f.Record, f.FieldName})
}

func (u UnionType) CodecEncodeSelf(e *codec.Encoder) {
	items := map[string]Expr(u)
	// we rely on the EncodeOptions having Canonical set
	// so that we get sorted keys in our map
	output := []interface{}{11, items}
	e.Encode(output)
}

func (m Merge) CodecEncodeSelf(e *codec.Encoder) {
	if m.Annotation != nil {
		e.Encode([]interface{}{6, m.Handler, m.Union, m.Annotation})
	} else {
		e.Encode([]interface{}{6, m.Handler, m.Union})
	}
}

const (
	HttpImport     = 0
	HttpsImport    = 1
	AbsoluteImport = 2
	HereImport     = 3
	ParentImport   = 4
	HomeImport     = 5
)

func (i Embed) CodecEncodeSelf(e *codec.Encoder) {
	r := i.Fetchable
	mode := 0
	if i.ImportMode == RawText {
		mode = 1
	}
	var hash interface{} // unimplemented, leave as nil for now
	switch rr := r.(type) {
	case EnvVar:
		e.Encode([]interface{}{24, hash, mode, 6, string(rr)})
	case Local:
		if rr.IsAbs() {
			toEncode := []interface{}{24, hash, mode, AbsoluteImport}
			for _, component := range rr.PathComponents() {
				toEncode = append(toEncode, component)
			}
			e.Encode(toEncode)
		} else if rr.IsRelativeToParent() {
			toEncode := []interface{}{24, hash, mode, ParentImport}
			for _, component := range rr.PathComponents() {
				toEncode = append(toEncode, component)
			}
			e.Encode(toEncode)
		} else if rr.IsRelativeToHome() {
			toEncode := []interface{}{24, hash, mode, HomeImport}
			for _, component := range rr.PathComponents() {
				toEncode = append(toEncode, component)
			}
			e.Encode(toEncode)
		} else {
			toEncode := []interface{}{24, hash, mode, HereImport}
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
		toEncode := []interface{}{24, hash, mode, scheme, headers, rr.Authority()}
		for _, component := range rr.PathComponents() {
			toEncode = append(toEncode, component)
		}
		toEncode = append(toEncode, rr.Query())
		e.Encode(toEncode)
	case Missing:
		e.Encode("Missing unimplemented")
	default:
		panic("can't happen")
	}
}

// ignored
func (Const) CodecDecodeSelf(*codec.Decoder)        {}
func (Var) CodecDecodeSelf(*codec.Decoder)          {}
func (*LambdaExpr) CodecDecodeSelf(*codec.Decoder)  {}
func (*Pi) CodecDecodeSelf(*codec.Decoder)          {}
func (*App) CodecDecodeSelf(*codec.Decoder)         {}
func (Let) CodecDecodeSelf(*codec.Decoder)          {}
func (Annot) CodecDecodeSelf(*codec.Decoder)        {}
func (Builtin) CodecDecodeSelf(*codec.Decoder)      {}
func (BoolIf) CodecDecodeSelf(*codec.Decoder)       {}
func (TextLit) CodecDecodeSelf(*codec.Decoder)      {}
func (NaturalLit) CodecDecodeSelf(*codec.Decoder)   {}
func (Operator) CodecDecodeSelf(*codec.Decoder)     {}
func (IntegerLit) CodecDecodeSelf(*codec.Decoder)   {}
func (EmptyList) CodecDecodeSelf(*codec.Decoder)    {}
func (NonEmptyList) CodecDecodeSelf(*codec.Decoder) {}
func (Some) CodecDecodeSelf(*codec.Decoder)         {}
func (Record) CodecDecodeSelf(*codec.Decoder)       {}
func (RecordLit) CodecDecodeSelf(*codec.Decoder)    {}
func (Field) CodecDecodeSelf(*codec.Decoder)        {}
func (UnionType) CodecDecodeSelf(*codec.Decoder)    {}
func (Merge) CodecDecodeSelf(*codec.Decoder)        {}
func (Embed) CodecDecodeSelf(*codec.Decoder)        {}
