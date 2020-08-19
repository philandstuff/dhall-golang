package term

import (
	"fmt"
	"reflect"
)

// TransformSubexprs returns a Term with f() applied to each immediate
// subexpression.
func TransformSubexprs(t Term, f func(Term) Term) Term {
	switch t := t.(type) {
	case Universe, Builtin, Var, LocalVar, NaturalLit, DoubleLit, BoolLit,
		IntegerLit:
		return t
	case Lambda:
		return Lambda{
			Label: t.Label,
			Type:  f(t.Type),
			Body:  f(t.Body),
		}
	case Pi:
		return Pi{
			Label: t.Label,
			Type:  f(t.Type),
			Body:  f(t.Body),
		}
	case App:
		return App{
			Fn:  f(t.Fn),
			Arg: f(t.Arg),
		}
	case Let:
		newLet := Let{}
		for _, b := range t.Bindings {
			newBinding := Binding{
				Variable: b.Variable,
				Value:    f(b.Value),
			}
			if b.Annotation != nil {
				newBinding.Annotation = f(b.Annotation)
			}
			newLet.Bindings = append(newLet.Bindings, newBinding)
		}
		newLet.Body = f(t.Body)
		return newLet
	case Annot:
		return Annot{
			Expr:       f(t.Expr),
			Annotation: f(t.Annotation),
		}
	case TextLit:
		result := TextLit{Suffix: t.Suffix}
		if t.Chunks == nil {
			return result
		}
		result.Chunks = Chunks{}
		for _, chunk := range t.Chunks {
			result.Chunks = append(result.Chunks,
				Chunk{
					Prefix: chunk.Prefix,
					Expr:   f(chunk.Expr),
				})
		}
		return result
	case If:
		return If{
			Cond: f(t.Cond),
			T:    f(t.T),
			F:    f(t.F),
		}
	case Op:
		return Op{
			OpCode: t.OpCode,
			L:      f(t.L),
			R:      f(t.R),
		}
	case EmptyList:
		return EmptyList{Type: f(t.Type)}
	case NonEmptyList:
		result := make(NonEmptyList, len(t))
		for j, e := range t {
			result[j] = f(e)
		}
		return result
	case Some:
		return Some{f(t.Val)}
	case RecordType:
		result := make(RecordType, len(t))
		for k, v := range t {
			result[k] = f(v)
		}
		return result
	case RecordLit:
		result := make(RecordLit, len(t))
		for k, v := range t {
			result[k] = f(v)
		}
		return result
	case ToMap:
		result := ToMap{Record: f(t.Record)}
		if t.Type != nil {
			result.Type = f(t.Type)
		}
		return result
	case Field:
		return Field{
			Record:    f(t.Record),
			FieldName: t.FieldName,
		}
	case Project:
		return Project{
			Record:     f(t.Record),
			FieldNames: t.FieldNames,
		}
	case ProjectType:
		return ProjectType{
			Record:   f(t.Record),
			Selector: f(t.Selector),
		}
	case UnionType:
		result := make(UnionType, len(t))
		for k, v := range t {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = f(v)
		}
		return result
	case Merge:
		result := Merge{
			Handler: f(t.Handler),
			Union:   f(t.Union),
		}
		if t.Annotation != nil {
			result.Annotation = f(t.Annotation)
		}
		return result
	case Assert:
		return Assert{Annotation: f(t.Annotation)}
	case Import:
		return t
	default:
		panic(fmt.Sprintf("unknown term type %+v (%v)", t, reflect.ValueOf(t).Type()))
	}
}

// MaybeTransformSubexprs returns a Term with f() applied to each
// immediate subexpression.  If f() returns an error at any point,
// MaybeTransformSubexprs returns that error.
func MaybeTransformSubexprs(t Term, f func(Term) (Term, error)) (Term, error) {
	switch t := t.(type) {
	case Universe, Builtin, Var, LocalVar, NaturalLit, DoubleLit, BoolLit,
		IntegerLit:
		return t, nil
	case Lambda:
		typ, err := f(t.Type)
		if err != nil {
			return nil, err
		}
		body, err := f(t.Body)
		if err != nil {
			return nil, err
		}
		return Lambda{Label: t.Label, Type: typ, Body: body}, nil
	case Pi:
		typ, err := f(t.Type)
		if err != nil {
			return nil, err
		}
		body, err := f(t.Body)
		if err != nil {
			return nil, err
		}
		return Pi{Label: t.Label, Type: typ, Body: body}, nil
	case App:
		fn, err := f(t.Fn)
		if err != nil {
			return nil, err
		}
		arg, err := f(t.Arg)
		if err != nil {
			return nil, err
		}
		return App{Fn: fn, Arg: arg}, nil
	case Let:
		var err error
		newLet := Let{}
		for _, b := range t.Bindings {
			value, err := f(b.Value)
			if err != nil {
				return nil, err
			}
			newBinding := Binding{
				Variable: b.Variable,
				Value:    value,
			}
			if b.Annotation != nil {
				newBinding.Annotation, err = f(b.Annotation)
				if err != nil {
					return nil, err
				}
			}
			newLet.Bindings = append(newLet.Bindings, newBinding)
		}
		newLet.Body, err = f(t.Body)
		return newLet, err
	case Annot:
		expr, err := f(t.Expr)
		if err != nil {
			return nil, err
		}
		annotation, err := f(t.Annotation)
		if err != nil {
			return nil, err
		}
		return Annot{Expr: expr, Annotation: annotation}, nil
	case TextLit:
		result := TextLit{Suffix: t.Suffix}
		if t.Chunks == nil {
			return result, nil
		}
		result.Chunks = Chunks{}
		for _, chunk := range t.Chunks {
			expr, err := f(chunk.Expr)
			if err != nil {
				return nil, err
			}
			result.Chunks = append(result.Chunks,
				Chunk{Prefix: chunk.Prefix, Expr: expr})
		}
		return result, nil
	case If:
		cond, err := f(t.Cond)
		if err != nil {
			return nil, err
		}
		T, err := f(t.T)
		if err != nil {
			return nil, err
		}
		F, err := f(t.F)
		if err != nil {
			return nil, err
		}
		return If{Cond: cond, T: T, F: F}, nil
	case Op:
		l, err := f(t.L)
		if err != nil {
			return nil, err
		}
		r, err := f(t.R)
		if err != nil {
			return nil, err
		}
		return Op{OpCode: t.OpCode, L: l, R: r}, nil
	case EmptyList:
		typ, err := f(t.Type)
		return EmptyList{Type: typ}, err
	case NonEmptyList:
		result := make(NonEmptyList, len(t))
		for j, e := range t {
			var err error
			result[j], err = f(e)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case Some:
		val, err := f(t.Val)
		return Some{val}, err
	case RecordType:
		result := make(RecordType, len(t))
		for k, v := range t {
			var err error
			result[k], err = f(v)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case RecordLit:
		result := make(RecordLit, len(t))
		for k, v := range t {
			var err error
			result[k], err = f(v)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case ToMap:
		record, err := f(t.Record)
		if err != nil {
			return nil, err
		}
		result := ToMap{Record: record}
		if t.Type != nil {
			result.Type, err = f(t.Type)
		}
		return result, err
	case Field:
		record, err := f(t.Record)
		return Field{
			Record:    record,
			FieldName: t.FieldName,
		}, err
	case Project:
		record, err := f(t.Record)
		return Project{
			Record:     record,
			FieldNames: t.FieldNames,
		}, err
	case ProjectType:
		record, err := f(t.Record)
		if err != nil {
			return nil, err
		}
		selector, err := f(t.Selector)
		if err != nil {
			return nil, err
		}
		return ProjectType{Record: record, Selector: selector}, nil
	case UnionType:
		result := make(UnionType, len(t))
		for k, v := range t {
			if v == nil {
				result[k] = nil
				continue
			}
			var err error
			result[k], err = f(v)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case Merge:
		handler, err := f(t.Handler)
		if err != nil {
			return nil, err
		}
		union, err := f(t.Union)
		if err != nil {
			return nil, err
		}
		result := Merge{
			Handler: handler,
			Union:   union,
		}
		if t.Annotation != nil {
			result.Annotation, err = f(t.Annotation)
		}
		return result, err
	case Assert:
		annotation, err := f(t.Annotation)
		return Assert{Annotation: annotation}, err
	case Import:
		return t, nil
	default:
		if t == nil {
			panic(fmt.Sprintf("nil term"))
		}
		panic(fmt.Sprintf("unknown term type %+v (%v)", t, reflect.ValueOf(t).Type()))
	}
}
