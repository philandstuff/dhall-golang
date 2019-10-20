package eval

import (
	"fmt"

	. "github.com/philandstuff/dhall-golang/core"
)

func subst(name string, replacement, t Term) Term {
	return substAtLevel(0, name, replacement, t)
}

func substAtLevel(i int, name string, replacement, t Term) Term {
	switch t := t.(type) {
	case Universe:
		return t
	case Builtin:
		return t
	case BoundVar:
		if t.Name == name && t.Index == i {
			return replacement
		}
		return t
	case FreeVar:
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
		return TextLitTerm{Suffix: "let unimplemented"}
	case Annot:
		return substAtLevel(i, name, replacement, t.Expr)
	case DoubleLit:
		return t
	case TextLitTerm:
		return TextLitTerm{Suffix: "TextLit unimplemented but here's the suffix: " + t.Suffix}
	case BoolLit:
		return t
	case IfTerm:
		return TextLitTerm{Suffix: "If unimplemented"}
	case IntegerLit:
		return t
	case OpTerm:
		return TextLitTerm{Suffix: "OpTerm unimplemented"}
	case EmptyList:
		return EmptyList{Type: substAtLevel(i, name, replacement, t.Type)}
	case NonEmptyList:
		return TextLitTerm{Suffix: "NonEmptyList unimplemented"}
	case Some:
		return Some{substAtLevel(i, name, replacement, t.Val)}
	case RecordType:
		return TextLitTerm{Suffix: "RecordType unimplemented"}
	case RecordLit:
		return TextLitTerm{Suffix: "RecordLit unimplemented"}
	case ToMap:
		return TextLitTerm{Suffix: "ToMap unimplemented"}
	case Field:
		return TextLitTerm{Suffix: "Field unimplemented"}
	case Project:
		return TextLitTerm{Suffix: "Project unimplemented"}
	case ProjectType:
		return TextLitTerm{Suffix: "ProjectType unimplemented"}
	case UnionType:
		return TextLitTerm{Suffix: "UnionType unimplemented"}
	case Merge:
		return TextLitTerm{Suffix: "Merge unimplemented"}
	case Assert:
		return Assert{Annotation: substAtLevel(i, name, replacement, t.Annotation)}
	case Import:
		return t
	default:
		panic(fmt.Sprintf("unknown term type %v", t))
	}
}

func rebindLocal(local LocalVar, t Term) Term {
	return rebindAtLevel(0, local, t)
}

func rebindAtLevel(i int, local LocalVar, t Term) Term {
	switch t := t.(type) {
	case Universe:
		return t
	case Builtin:
		return t
	case BoundVar:
		return t
	case LocalVar:
		if t == local {
			return BoundVar{
				Name:  t.Name,
				Index: i,
			}
		}
		return t
	case FreeVar:
		return t
	case LambdaTerm:
		return LambdaTerm{
			Label: t.Label,
			Type:  rebindAtLevel(i, local, t.Type),
			Body:  rebindAtLevel(i+1, local, t.Body),
		}
	case PiTerm:
		return PiTerm{
			Label: t.Label,
			Type:  rebindAtLevel(i, local, t.Type),
			Body:  rebindAtLevel(i+1, local, t.Body),
		}
	case AppTerm:
		return AppTerm{
			Fn:  rebindAtLevel(i, local, t.Fn),
			Arg: rebindAtLevel(i, local, t.Arg),
		}
	case NaturalLit:
		return t
	case EmptyList:
		return EmptyList{
			Type: rebindAtLevel(i, local, t.Type),
		}
	default:
		panic(fmt.Sprintf("unknown term type %v", t))
	}
}
