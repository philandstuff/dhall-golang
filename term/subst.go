package term

// Subst takes a Term and finds all instances of a variable called
// `name` and replaces them with the replacement.
func Subst(name string, replacement, t Term) Term {
	return substAtLevel(0, name, replacement, t)
}

func substAtLevel(i int, name string, replacement, t Term) Term {
	switch t := t.(type) {
	case Var:
		if t.Name == name && t.Index == i {
			return replacement
		}
		return t
	case Lambda:
		j := i
		if t.Label == name {
			j = i + 1
		}
		return Lambda{
			Label: t.Label,
			Type:  substAtLevel(i, name, replacement, t.Type),
			Body:  substAtLevel(j, name, replacement, t.Body),
		}
	case Pi:
		j := i
		if t.Label == name {
			j = i + 1
		}
		return Pi{
			Label: t.Label,
			Type:  substAtLevel(i, name, replacement, t.Type),
			Body:  substAtLevel(j, name, replacement, t.Body),
		}
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
	default:
		return TransformSubexprs(t,
			func(t Term) Term {
				return substAtLevel(i, name, replacement, t)
			})
	}
}

// RebindLocal takes a Term and finds all instances of a LocalVar and
// replaces them with the equivalent Var.
func RebindLocal(local LocalVar, t Term) Term {
	return rebindAtLevel(0, local, t)
}

func rebindAtLevel(i int, local LocalVar, t Term) Term {
	switch t := t.(type) {
	case LocalVar:
		if t == local {
			return Var{
				Name:  t.Name,
				Index: i,
			}
		}
		return t
	case Lambda:
		j := i
		if t.Label == local.Name {
			j = i + 1
		}
		return Lambda{
			Label: t.Label,
			Type:  rebindAtLevel(i, local, t.Type),
			Body:  rebindAtLevel(j, local, t.Body),
		}
	case Pi:
		j := i
		if t.Label == local.Name {
			j = i + 1
		}
		return Pi{
			Label: t.Label,
			Type:  rebindAtLevel(i, local, t.Type),
			Body:  rebindAtLevel(j, local, t.Body),
		}
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
	default:
		return TransformSubexprs(t,
			func(t Term) Term {
				return rebindAtLevel(i, local, t)
			})
	}
}
