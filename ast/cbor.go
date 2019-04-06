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
var _ codec.Selfer = Double // BuiltinType
var _ codec.Selfer = BoolIf{}
var _ codec.Selfer = EmptyList{}
var _ codec.Selfer = NonEmptyList{}
var _ codec.Selfer = TextLit{}
var _ codec.Selfer = NaturalLit(0)
var _ codec.Selfer = NaturalPlus{}
var _ codec.Selfer = IntegerLit(0)
var _ codec.Selfer = Record(map[string]Expr{})
var _ codec.Selfer = RecordLit(map[string]Expr{})

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
	} else if v.Index == 0 {
		e.Encode(v.Name)
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
	args := []interface{}{a.Arg}
	// FIXME: support multi-arg application
	e.Encode(append([]interface{}{0, a.Fn}, args...))
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

func (t BuiltinType) CodecEncodeSelf(e *codec.Encoder) {
	switch t {
	case Double:
		e.Encode("Double")
	case Bool:
		e.Encode("Bool")
	case Natural:
		e.Encode("Natural")
	case Integer:
		e.Encode("Integer")
	case List:
		e.Encode("List")
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
		output = append(output, []interface{}{chunk.Prefix, chunk.Expr})
	}
	output = append(output, t.Suffix)
	e.Encode(output)
}

func (n NaturalLit) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode(append([]interface{}{15}, int(n)))
}

func (p NaturalPlus) CodecEncodeSelf(e *codec.Encoder) {
	e.Encode([]interface{}{3, 4, p.L, p.R})
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

// ignored
func (Const) CodecDecodeSelf(*codec.Decoder)        {}
func (Var) CodecDecodeSelf(*codec.Decoder)          {}
func (*LambdaExpr) CodecDecodeSelf(*codec.Decoder)  {}
func (*Pi) CodecDecodeSelf(*codec.Decoder)          {}
func (*App) CodecDecodeSelf(*codec.Decoder)         {}
func (Let) CodecDecodeSelf(*codec.Decoder)          {}
func (Annot) CodecDecodeSelf(*codec.Decoder)        {}
func (BuiltinType) CodecDecodeSelf(*codec.Decoder)  {}
func (BoolIf) CodecDecodeSelf(*codec.Decoder)       {}
func (TextLit) CodecDecodeSelf(*codec.Decoder)      {}
func (NaturalLit) CodecDecodeSelf(*codec.Decoder)   {}
func (NaturalPlus) CodecDecodeSelf(*codec.Decoder)  {}
func (IntegerLit) CodecDecodeSelf(*codec.Decoder)   {}
func (EmptyList) CodecDecodeSelf(*codec.Decoder)    {}
func (NonEmptyList) CodecDecodeSelf(*codec.Decoder) {}
func (Record) CodecDecodeSelf(*codec.Decoder)       {}
func (RecordLit) CodecDecodeSelf(*codec.Decoder)    {}
