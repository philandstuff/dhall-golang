package term

import (
	"fmt"
	"reflect"
)

// Subst takes a Term and finds all instances of a variable called
// `name` and replaces them with the replacement.
func Subst(name string, replacement, t Term) Term {
	return substAtLevel(0, name, replacement, t)
}

func substAtLevel(i int, name string, replacement, t Term) Term {
	switch t := t.(type) {
	case Universe:
		return t
	case Builtin:
		return t
	case Var:
		if t.Name == name && t.Index == i {
			return replacement
		}
		return t
	case LocalVar:
		return t
	case LambdaTerm:
		j := i
		if t.Label == name {
			j = i + 1
		}
		return LambdaTerm{
			Label: t.Label,
			Type:  substAtLevel(i, name, replacement, t.Type),
			Body:  substAtLevel(j, name, replacement, t.Body),
		}
	case PiTerm:
		j := i
		if t.Label == name {
			j = i + 1
		}
		return PiTerm{
			Label: t.Label,
			Type:  substAtLevel(i, name, replacement, t.Type),
			Body:  substAtLevel(j, name, replacement, t.Body),
		}
	case AppTerm:
		return AppTerm{
			Fn:  substAtLevel(i, name, replacement, t.Fn),
			Arg: substAtLevel(i, name, replacement, t.Arg),
		}
	case NaturalLit:
		return t
	case Let:
		newLet := Let{}
		for _, b := range t.Bindings {
			newBinding := Binding{
				Variable: b.Variable,
				Value:    substAtLevel(i, name, replacement, b.Value),
			}
			if b.Annotation != nil {
				newBinding.Annotation = substAtLevel(i, name, replacement, b.Annotation)
			}
			newLet.Bindings = append(newLet.Bindings, newBinding)
			if b.Variable == name {
				i = i + 1
			}
		}
		newLet.Body = substAtLevel(i, name, replacement, t.Body)
		return newLet
	case Annot:
		return substAtLevel(i, name, replacement, t.Expr)
	case DoubleLit:
		return t
	case TextLitTerm:
		result := TextLitTerm{Suffix: t.Suffix}
		if t.Chunks == nil {
			return result
		}
		result.Chunks = Chunks{}
		for _, chunk := range t.Chunks {
			result.Chunks = append(result.Chunks,
				Chunk{
					Prefix: chunk.Prefix,
					Expr:   substAtLevel(i, name, replacement, chunk.Expr),
				})
		}
		return result
	case BoolLit:
		return t
	case IfTerm:
		return IfTerm{
			Cond: substAtLevel(i, name, replacement, t.Cond),
			T:    substAtLevel(i, name, replacement, t.T),
			F:    substAtLevel(i, name, replacement, t.F),
		}
	case IntegerLit:
		return t
	case OpTerm:
		return OpTerm{
			OpCode: t.OpCode,
			L:      substAtLevel(i, name, replacement, t.L),
			R:      substAtLevel(i, name, replacement, t.R),
		}
	case EmptyList:
		return EmptyList{Type: substAtLevel(i, name, replacement, t.Type)}
	case NonEmptyList:
		result := make(NonEmptyList, len(t))
		for j, e := range t {
			result[j] = substAtLevel(i, name, replacement, e)
		}
		return result
	case Some:
		return Some{substAtLevel(i, name, replacement, t.Val)}
	case RecordType:
		result := make(RecordType, len(t))
		for k, v := range t {
			result[k] = substAtLevel(i, name, replacement, v)
		}
		return result
	case RecordLit:
		result := make(RecordLit, len(t))
		for k, v := range t {
			result[k] = substAtLevel(i, name, replacement, v)
		}
		return result
	case ToMap:
		result := ToMap{Record: substAtLevel(i, name, replacement, t.Record)}
		if t.Type != nil {
			result.Type = substAtLevel(i, name, replacement, t.Type)
		}
		return result
	case Field:
		return Field{
			Record:    substAtLevel(i, name, replacement, t.Record),
			FieldName: t.FieldName,
		}
	case Project:
		return Project{
			Record:     substAtLevel(i, name, replacement, t.Record),
			FieldNames: t.FieldNames,
		}
	case ProjectType:
		return ProjectType{
			Record:   substAtLevel(i, name, replacement, t.Record),
			Selector: substAtLevel(i, name, replacement, t.Selector),
		}
	case UnionType:
		result := make(UnionType, len(t))
		for k, v := range t {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = substAtLevel(i, name, replacement, v)
		}
		return result
	case Merge:
		result := Merge{
			Handler: substAtLevel(i, name, replacement, t.Handler),
			Union:   substAtLevel(i, name, replacement, t.Union),
		}
		if t.Annotation != nil {
			result.Annotation = substAtLevel(i, name, replacement, t.Annotation)
		}
		return result
	case Assert:
		return Assert{Annotation: substAtLevel(i, name, replacement, t.Annotation)}
	case Import:
		return t
	default:
		panic(fmt.Sprintf("unknown term type %+v (%v)", t, reflect.ValueOf(t).Type()))
	}
}

// RebindLocal takes a Term and finds all instances of a LocalVar and
// replaces them with the equivalent Var.
func RebindLocal(local LocalVar, t Term) Term {
	return rebindAtLevel(0, local, t)
}

func rebindAtLevel(i int, local LocalVar, t Term) Term {
	switch t := t.(type) {
	case Universe:
		return t
	case Builtin:
		return t
	case Var:
		return t
	case LocalVar:
		if t == local {
			return Var{
				Name:  t.Name,
				Index: i,
			}
		}
		return t
	case LambdaTerm:
		j := i
		if t.Label == local.Name {
			j = i + 1
		}
		return LambdaTerm{
			Label: t.Label,
			Type:  rebindAtLevel(i, local, t.Type),
			Body:  rebindAtLevel(j, local, t.Body),
		}
	case PiTerm:
		j := i
		if t.Label == local.Name {
			j = i + 1
		}
		return PiTerm{
			Label: t.Label,
			Type:  rebindAtLevel(i, local, t.Type),
			Body:  rebindAtLevel(j, local, t.Body),
		}
	case AppTerm:
		return AppTerm{
			Fn:  rebindAtLevel(i, local, t.Fn),
			Arg: rebindAtLevel(i, local, t.Arg),
		}
	case NaturalLit:
		return t
	case Let:
		newLet := Let{}
		for _, b := range t.Bindings {
			newBinding := Binding{
				Variable: b.Variable,
				Value:    rebindAtLevel(i, local, b.Value),
			}
			if b.Annotation != nil {
				newBinding.Annotation = rebindAtLevel(i, local, b.Annotation)
			}
			newLet.Bindings = append(newLet.Bindings, newBinding)
			if b.Variable == local.Name {
				i = i + 1
			}
		}
		newLet.Body = rebindAtLevel(i, local, t.Body)
		return newLet
	case Annot:
		return rebindAtLevel(i, local, t.Expr)
	case DoubleLit:
		return t
	case TextLitTerm:
		result := TextLitTerm{Suffix: t.Suffix}
		if t.Chunks == nil {
			return result
		}
		result.Chunks = Chunks{}
		for _, chunk := range t.Chunks {
			result.Chunks = append(result.Chunks,
				Chunk{
					Prefix: chunk.Prefix,
					Expr:   rebindAtLevel(i, local, chunk.Expr),
				})
		}
		return result
	case BoolLit:
		return t
	case IfTerm:
		return IfTerm{
			Cond: rebindAtLevel(i, local, t.Cond),
			T:    rebindAtLevel(i, local, t.T),
			F:    rebindAtLevel(i, local, t.F),
		}
	case IntegerLit:
		return t
	case OpTerm:
		return OpTerm{
			OpCode: t.OpCode,
			L:      rebindAtLevel(i, local, t.L),
			R:      rebindAtLevel(i, local, t.R),
		}
	case EmptyList:
		return EmptyList{
			Type: rebindAtLevel(i, local, t.Type),
		}
	case NonEmptyList:
		result := make(NonEmptyList, len(t))
		for j, e := range t {
			result[j] = rebindAtLevel(i, local, e)
		}
		return result
	case Some:
		return Some{rebindAtLevel(i, local, t.Val)}
	case RecordType:
		result := make(RecordType, len(t))
		for k, v := range t {
			result[k] = rebindAtLevel(i, local, v)
		}
		return result
	case RecordLit:
		result := make(RecordLit, len(t))
		for k, v := range t {
			result[k] = rebindAtLevel(i, local, v)
		}
		return result
	case ToMap:
		result := ToMap{Record: rebindAtLevel(i, local, t.Record)}
		if t.Type != nil {
			result.Type = rebindAtLevel(i, local, t.Type)
		}
		return result
	case Field:
		return Field{
			Record:    rebindAtLevel(i, local, t.Record),
			FieldName: t.FieldName,
		}
	case Project:
		return Project{
			Record:     rebindAtLevel(i, local, t.Record),
			FieldNames: t.FieldNames,
		}
	case ProjectType:
		return ProjectType{
			Record:   rebindAtLevel(i, local, t.Record),
			Selector: rebindAtLevel(i, local, t.Selector),
		}
	case UnionType:
		result := make(UnionType, len(t))
		for k, v := range t {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = rebindAtLevel(i, local, v)
		}
		return result
	case Merge:
		result := Merge{
			Handler: rebindAtLevel(i, local, t.Handler),
			Union:   rebindAtLevel(i, local, t.Union),
		}
		if t.Annotation != nil {
			result.Annotation = rebindAtLevel(i, local, t.Annotation)
		}
		return result
	case Assert:
		return Assert{Annotation: rebindAtLevel(i, local, t.Annotation)}
	case Import:
		return t
	default:
		panic(fmt.Sprintf("unknown term type %+v (%v)", t, reflect.ValueOf(t).Type()))
	}
}
