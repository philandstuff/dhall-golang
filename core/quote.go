package core

// Quote takes the Value v and turns it back into a Term.  The `i` is the
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
	case naturalBuildVal:
		return NaturalBuild
	case naturalEvenVal:
		return NaturalEven
	case naturalFoldVal:
		var result Term = NaturalFold
		if v.n == nil {
			return result
		}
		result = AppTerm{result, quoteWith(ctx, v.n)}
		if v.typ == nil {
			return result
		}
		result = AppTerm{result, quoteWith(ctx, v.typ)}
		if v.succ == nil {
			return result
		}
		return AppTerm{result, quoteWith(ctx, v.succ)}
	case naturalIsZeroVal:
		return NaturalIsZero
	case naturalOddVal:
		return NaturalOdd
	case naturalShowVal:
		return NaturalShow
	case naturalSubtractVal:
		if v.a != nil {
			return AppTerm{NaturalSubtract, quoteWith(ctx, v.a)}
		}
		return NaturalSubtract
	case naturalToIntegerVal:
		return NaturalToInteger
	case integerShowVal:
		return IntegerShow
	case integerToDoubleVal:
		return IntegerToDouble
	case doubleShowVal:
		return DoubleShow
	case optionalBuildVal:
		if v.typ != nil {
			return AppTerm{OptionalBuild, quoteWith(ctx, v.typ)}
		}
		return OptionalBuild
	case optionalFoldVal:
		var result Term = OptionalFold
		if v.typ1 == nil {
			return result
		}
		result = AppTerm{result, quoteWith(ctx, v.typ1)}
		if v.opt == nil {
			return result
		}
		result = AppTerm{result, quoteWith(ctx, v.opt)}
		if v.typ2 == nil {
			return result
		}
		result = AppTerm{result, quoteWith(ctx, v.typ2)}
		if v.some == nil {
			return result
		}
		return AppTerm{result, quoteWith(ctx, v.some)}
	case textShowVal:
		return TextShow
	case listBuildVal:
		if v.typ != nil {
			return AppTerm{ListBuild, quoteWith(ctx, v.typ)}
		}
		return ListBuild
	case listFoldVal:
		var result Term = ListFold
		if v.typ1 == nil {
			return result
		}
		result = AppTerm{result, quoteWith(ctx, v.typ1)}
		if v.list == nil {
			return result
		}
		result = AppTerm{result, quoteWith(ctx, v.list)}
		if v.typ2 == nil {
			return result
		}
		result = AppTerm{result, quoteWith(ctx, v.typ2)}
		if v.cons == nil {
			return result
		}
		return AppTerm{result, quoteWith(ctx, v.cons)}
	case listHeadVal:
		if v.typ != nil {
			return AppTerm{ListHead, quoteWith(ctx, v.typ)}
		}
		return ListHead
	case listIndexedVal:
		if v.typ != nil {
			return AppTerm{ListIndexed, quoteWith(ctx, v.typ)}
		}
		return ListIndexed
	case listLengthVal:
		if v.typ != nil {
			return AppTerm{ListLength, quoteWith(ctx, v.typ)}
		}
		return ListLength
	case listLastVal:
		if v.typ != nil {
			return AppTerm{ListLast, quoteWith(ctx, v.typ)}
		}
		return ListLast
	case listReverseVal:
		if v.typ != nil {
			return AppTerm{ListReverse, quoteWith(ctx, v.typ)}
		}
		return ListReverse
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
		bodyVal := v.Call(QuoteVar{Name: v.Label, Index: ctx[v.Label]})
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
	case OpValue:
		return OpTerm{
			OpCode: v.OpCode,
			L:      quoteWith(ctx, v.L),
			R:      quoteWith(ctx, v.R),
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
		var newChunks Chunks
		for _, chunk := range v.Chunks {
			newChunks = append(newChunks, Chunk{
				Prefix: chunk.Prefix,
				Expr:   quoteWith(ctx, chunk.Expr),
			})
		}
		return TextLitTerm{
			Chunks: newChunks,
			Suffix: v.Suffix,
		}
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
		result := ToMap{Record: quoteWith(ctx, v.Record)}
		if v.Type != nil {
			result.Type = quoteWith(ctx, v.Type)
		}
		return result
	case FieldVal:
		return Field{
			Record:    quoteWith(ctx, v.Record),
			FieldName: v.FieldName,
		}
	case ProjectVal:
		return Project{
			Record:     quoteWith(ctx, v.Record),
			FieldNames: v.FieldNames,
		}
	case UnionTypeVal:
		result := UnionType{}
		for k, v := range v {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = quoteWith(ctx, v)
		}
		return result
	case MergeVal:
		result := Merge{
			Handler: quoteWith(ctx, v.Handler),
			Union:   quoteWith(ctx, v.Union),
		}
		if v.Annotation != nil {
			result.Annotation = quoteWith(ctx, v.Annotation)
		}
		return result
	case AssertVal:
		return Assert{Annotation: quoteWith(ctx, v.Annotation)}
	}
	panic("unknown Value type")
}
