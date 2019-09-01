package eval

import . "github.com/philandstuff/dhall-golang/core"

// Quote(v) takes the Value v and turns it back into a Term.  The `i` is the
// first fresh variable index named `quote`.  Normally this will be 0 if there
// are no variables called `quote` in the context.
func Quote(v Value) Term {
	return quoteWith(quoteContext{}, v, LocalVar{})
}

// quote, rebinding the given LocalVar back as a BoundVar
func quoteAndRebindLocal(v Value, l LocalVar) Term {
	return quoteWith(quoteContext{}, v, l)
}

// a quoteContext records how many binders of each variable name we have passed
type quoteContext map[string]int

func (q quoteContext) extend(name string) quoteContext {
	newCtx := make(quoteContext, len(q))
	for k, v := range q {
		newCtx[k] = v
	}
	newCtx[name]++
	return newCtx
}

func quoteWith(ctx quoteContext, v Value, l LocalVar) Term {
	switch v := v.(type) {
	case Universe:
		return v
	case Builtin:
		return v
	case FreeVar:
		return v
	case LocalVar:
		if v == l {
			return BoundVar{
				Name:  v.Name,
				Index: ctx[v.Name],
			}
		}
		return v
	case QuoteVar:
		return BoundVar{
			Name:  v.Name,
			Index: ctx[v.Name] - v.Index - 1,
		}
	case LambdaValue:
		bodyVal := v.Fn(QuoteVar{Name: v.Label, Index: ctx[v.Label]})
		return LambdaTerm{
			Label: v.Label,
			Type:  quoteWith(ctx, v.Domain, l),
			Body:  quoteWith(ctx.extend(v.Label), bodyVal, l),
		}
	case PiValue:
		bodyVal := v.Range(QuoteVar{Name: v.Label, Index: ctx[v.Label]})
		return PiTerm{
			Label: v.Label,
			Type:  quoteWith(ctx, v.Domain, l),
			Body:  quoteWith(ctx.extend(v.Label), bodyVal, l),
		}
	case AppValue:
		return AppTerm{
			Fn:  quoteWith(ctx, v.Fn, l),
			Arg: quoteWith(ctx, v.Arg, l),
		}
	case NaturalLit:
		return v
	case EmptyListVal:
		return EmptyList{Type: quoteWith(ctx, v.Type, l)}
	}
	panic("unknown Value type")
}
