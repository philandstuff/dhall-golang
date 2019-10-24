package eval

import . "github.com/philandstuff/dhall-golang/core"

// Quote(v) takes the Value v and turns it back into a Term.  The `i` is the
// first fresh variable index named `quote`.  Normally this will be 0 if there
// are no variables called `quote` in the context.
func Quote(v Value) Term {
	return quoteWith(quoteContext{}, v)
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

func quoteWith(ctx quoteContext, v Value) Term {
	switch v := v.(type) {
	case Universe:
		return v
	case Builtin:
		return v
	case NaturalEvenVal:
		return NaturalEven
	case NaturalIsZeroVal:
		return NaturalIsZero
	case NaturalOddVal:
		return NaturalOdd
	case NaturalShowVal:
		return NaturalShow
	case NaturalToIntegerVal:
		return NaturalToInteger
	case FreeVar:
		return v
	case LocalVar:
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
			Type:  quoteWith(ctx, v.Domain),
			Body:  quoteWith(ctx.extend(v.Label), bodyVal),
		}
	case PiValue:
		bodyVal := v.Range(QuoteVar{Name: v.Label, Index: ctx[v.Label]})
		return PiTerm{
			Label: v.Label,
			Type:  quoteWith(ctx, v.Domain),
			Body:  quoteWith(ctx.extend(v.Label), bodyVal),
		}
	case AppValue:
		return AppTerm{
			Fn:  quoteWith(ctx, v.Fn),
			Arg: quoteWith(ctx, v.Arg),
		}
	case NaturalLit:
		return v
	case DoubleLit:
		return v
	case IntegerLit:
		return v
	case BoolLit:
		return v
	case EmptyListVal:
		return EmptyList{Type: quoteWith(ctx, v.Type)}
	case NonEmptyListVal:
		l := NonEmptyList{}
		for _, e := range v {
			l = append(l, quoteWith(ctx, e))
		}
		return l
	case TextLitVal:
		return TextLitTerm{Suffix: "quote TextLitVal unimplemented but suffix:" + v.Suffix}
	case IfVal:
		return IfTerm{
			Cond: quoteWith(ctx, v.Cond),
			T:    quoteWith(ctx, v.T),
			F:    quoteWith(ctx, v.F),
		}
	case SomeVal:
		return Some{Val: quoteWith(ctx, v.Val)}
	case RecordTypeVal:
		rt := RecordType{}
		for k, v := range v {
			rt[k] = quoteWith(ctx, v)
		}
		return rt
	case RecordLitVal:
		rt := RecordLit{}
		for k, v := range v {
			rt[k] = quoteWith(ctx, v)
		}
		return rt
	case ToMapVal:
		return TextLitTerm{Suffix: "quote ToMapVal unimplemented"}
	case FieldVal:
		return TextLitTerm{Suffix: "FieldVal unimplemented"}
	case ProjectVal:
		return TextLitTerm{Suffix: "ProjectVal unimplemented"}
	case UnionTypeVal:
		return TextLitTerm{Suffix: "UnionTypeVal unimplemented"}
	case MergeVal:
		return TextLitTerm{Suffix: "MergeVal unimplemented"}
	case AssertVal:
		return Assert{Annotation: quoteWith(ctx, v.Annotation)}
	}
	panic("unknown Value type")
}
