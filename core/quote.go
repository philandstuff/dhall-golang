package core

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/term"
)

// Quote takes the Value v and turns it back into a Term.
func Quote(v Value) term.Term {
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

func quoteWith(ctx quoteContext, v Value) term.Term {
	switch v := v.(type) {
	case Universe:
		return term.Universe(v)
	case Builtin:
		return term.Builtin(v)
	case naturalBuildVal:
		return term.NaturalBuild
	case naturalEvenVal:
		return term.NaturalEven
	case naturalFoldVal:
		var result term.Term = term.NaturalFold
		if v.n == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, v.n)}
		if v.typ == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, v.typ)}
		if v.succ == nil {
			return result
		}
		return term.App{result, quoteWith(ctx, v.succ)}
	case naturalIsZeroVal:
		return term.NaturalIsZero
	case naturalOddVal:
		return term.NaturalOdd
	case naturalShowVal:
		return term.NaturalShow
	case naturalSubtractVal:
		if v.a != nil {
			return term.App{term.NaturalSubtract, quoteWith(ctx, v.a)}
		}
		return term.NaturalSubtract
	case integerClampVal:
		return term.IntegerClamp
	case integerNegateVal:
		return term.IntegerNegate
	case naturalToIntegerVal:
		return term.NaturalToInteger
	case integerShowVal:
		return term.IntegerShow
	case integerToDoubleVal:
		return term.IntegerToDouble
	case doubleShowVal:
		return term.DoubleShow
	case optionalVal:
		return term.Optional
	case optionalBuildVal:
		if v.typ != nil {
			return term.App{term.OptionalBuild, quoteWith(ctx, v.typ)}
		}
		return term.OptionalBuild
	case optionalFoldVal:
		var result term.Term = term.OptionalFold
		if v.typ1 == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, v.typ1)}
		if v.opt == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, v.opt)}
		if v.typ2 == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, v.typ2)}
		if v.some == nil {
			return result
		}
		return term.App{result, quoteWith(ctx, v.some)}
	case noneVal:
		return term.None
	case textShowVal:
		return term.TextShow
	case listVal:
		return term.List
	case listBuildVal:
		if v.typ != nil {
			return term.App{term.ListBuild, quoteWith(ctx, v.typ)}
		}
		return term.ListBuild
	case listFoldVal:
		var result term.Term = term.ListFold
		if v.typ1 == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, v.typ1)}
		if v.list == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, v.list)}
		if v.typ2 == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, v.typ2)}
		if v.cons == nil {
			return result
		}
		return term.App{result, quoteWith(ctx, v.cons)}
	case listHeadVal:
		if v.typ != nil {
			return term.App{term.ListHead, quoteWith(ctx, v.typ)}
		}
		return term.ListHead
	case listIndexedVal:
		if v.typ != nil {
			return term.App{term.ListIndexed, quoteWith(ctx, v.typ)}
		}
		return term.ListIndexed
	case listLengthVal:
		if v.typ != nil {
			return term.App{term.ListLength, quoteWith(ctx, v.typ)}
		}
		return term.ListLength
	case listLastVal:
		if v.typ != nil {
			return term.App{term.ListLast, quoteWith(ctx, v.typ)}
		}
		return term.ListLast
	case listReverseVal:
		if v.typ != nil {
			return term.App{term.ListReverse, quoteWith(ctx, v.typ)}
		}
		return term.ListReverse
	case freeVar:
		return term.Var(v)
	case localVar:
		return term.LocalVar(v)
	case quoteVar:
		return term.Var{
			Name:  v.Name,
			Index: ctx[v.Name] - v.Index - 1,
		}
	case lambdaValue:
		bodyVal := v.Call(quoteVar{Name: v.Label, Index: ctx[v.Label]})
		return term.Lambda{
			Label: v.Label,
			Type:  quoteWith(ctx, v.Domain),
			Body:  quoteWith(ctx.extend(v.Label), bodyVal),
		}
	case PiValue:
		bodyVal := v.Range(quoteVar{Name: v.Label, Index: ctx[v.Label]})
		return term.Pi{
			Label: v.Label,
			Type:  quoteWith(ctx, v.Domain),
			Body:  quoteWith(ctx.extend(v.Label), bodyVal),
		}
	case appValue:
		return term.App{
			Fn:  quoteWith(ctx, v.Fn),
			Arg: quoteWith(ctx, v.Arg),
		}
	case opValue:
		return term.Op{
			OpCode: v.OpCode,
			L:      quoteWith(ctx, v.L),
			R:      quoteWith(ctx, v.R),
		}
	case NaturalLit:
		return term.NaturalLit(v)
	case DoubleLit:
		return term.DoubleLit(v)
	case IntegerLit:
		return term.IntegerLit(v)
	case BoolLit:
		return term.BoolLit(v)
	case ListOf:
		return term.Apply(term.List, quoteWith(ctx, v.Type))
	case EmptyListVal:
		return term.EmptyList{Type: quoteWith(ctx, v.Type)}
	case NonEmptyListVal:
		l := term.NonEmptyList{}
		for _, e := range v {
			l = append(l, quoteWith(ctx, e))
		}
		return l
	case TextLitVal:
		var newChunks term.Chunks
		for _, chunk := range v.Chunks {
			newChunks = append(newChunks, term.Chunk{
				Prefix: chunk.Prefix,
				Expr:   quoteWith(ctx, chunk.Expr),
			})
		}
		return term.TextLit{
			Chunks: newChunks,
			Suffix: v.Suffix,
		}
	case ifVal:
		return term.If{
			Cond: quoteWith(ctx, v.Cond),
			T:    quoteWith(ctx, v.T),
			F:    quoteWith(ctx, v.F),
		}
	case OptionalOf:
		return term.Apply(term.Optional, quoteWith(ctx, v.Type))
	case SomeVal:
		return term.Some{Val: quoteWith(ctx, v.Val)}
	case NoneOf:
		return term.Apply(term.None, quoteWith(ctx, v.Type))
	case RecordTypeVal:
		rt := term.RecordType{}
		for k, v := range v {
			rt[k] = quoteWith(ctx, v)
		}
		return rt
	case RecordLitVal:
		rt := term.RecordLit{}
		for k, v := range v {
			rt[k] = quoteWith(ctx, v)
		}
		return rt
	case toMapVal:
		result := term.ToMap{Record: quoteWith(ctx, v.Record)}
		if v.Type != nil {
			result.Type = quoteWith(ctx, v.Type)
		}
		return result
	case fieldVal:
		return term.Field{
			Record:    quoteWith(ctx, v.Record),
			FieldName: v.FieldName,
		}
	case projectVal:
		return term.Project{
			Record:     quoteWith(ctx, v.Record),
			FieldNames: v.FieldNames,
		}
	case unionTypeVal:
		result := term.UnionType{}
		for k, v := range v {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = quoteWith(ctx, v)
		}
		return result
	case mergeVal:
		result := term.Merge{
			Handler: quoteWith(ctx, v.Handler),
			Union:   quoteWith(ctx, v.Union),
		}
		if v.Annotation != nil {
			result.Annotation = quoteWith(ctx, v.Annotation)
		}
		return result
	case assertVal:
		return term.Assert{Annotation: quoteWith(ctx, v.Annotation)}
	}
	panic(fmt.Sprintf("unknown Value type %#v", v))
}
