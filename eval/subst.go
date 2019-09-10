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
	case LambdaTerm:
		return LambdaTerm{
			Label: t.Label,
			Type:  substAtLevel(i, name, replacement, t.Type),
			Body:  substAtLevel(i+1, name, replacement, t.Body),
		}
	case PiTerm:
		return PiTerm{
			Label: t.Label,
			Type:  substAtLevel(i, name, replacement, t.Type),
			Body:  substAtLevel(i+1, name, replacement, t.Body),
		}
	case AppTerm:
		return AppTerm{
			Fn:  substAtLevel(i, name, replacement, t.Fn),
			Arg: substAtLevel(i, name, replacement, t.Arg),
		}
	case NaturalLit:
		return t
	case EmptyList:
		return EmptyList{
			Type: substAtLevel(i, name, replacement, t.Type),
		}
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
