package eval

import . "github.com/philandstuff/dhall-golang/core"

// Quote(v) takes the Value v and turns it back into a Term.  The `i` is the
// first fresh variable index named `quote`.  Normally this will be 0 if there
// are no variables called `quote` in the context.
func Quote(v Value) Term {
	return quote(quoteContext{}, v)
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

func quote(ctx quoteContext, v Value) Term {
	switch v := v.(type) {
	case Universe:
		return v
	case Builtin:
		return v
	case FreeVar:
		return v
	case QuoteVar:
		return BoundVar{
			Name:  v.Name,
			Index: ctx[v.Name] - v.Index - 1,
		}
	case LambdaValue:
		bodyVal := v.Fn(QuoteVar{v.Label, ctx[v.Label]})
		return LambdaTerm{
			Label: v.Label,
			Type:  quote(ctx, v.Domain),
			Body:  quote(ctx.extend(v.Label), bodyVal),
		}
	case PiValue:
		bodyVal := v.Range(QuoteVar{v.Label, ctx[v.Label]})
		return PiTerm{
			Label: v.Label,
			Type:  quote(ctx, v.Domain),
			Body:  quote(ctx.extend(v.Label), bodyVal),
		}
	case AppNeutral:
		return AppTerm{
			Fn:  quote(ctx, v.Fn),
			Arg: quote(ctx, v.Arg),
		}
	case NaturalLit:
		return v
	case EmptyListVal:
		return EmptyList{Type: quote(ctx, v.Type)}
	}
	panic("unknown Value type")
}
