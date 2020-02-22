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
	case naturalBuild:
		return term.NaturalBuild
	case naturalEven:
		return term.NaturalEven
	case naturalFold:
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
	case naturalIsZero:
		return term.NaturalIsZero
	case naturalOdd:
		return term.NaturalOdd
	case naturalShow:
		return term.NaturalShow
	case naturalSubtract:
		if v.a != nil {
			return term.App{term.NaturalSubtract, quoteWith(ctx, v.a)}
		}
		return term.NaturalSubtract
	case integerClamp:
		return term.IntegerClamp
	case integerNegate:
		return term.IntegerNegate
	case naturalToInteger:
		return term.NaturalToInteger
	case integerShow:
		return term.IntegerShow
	case integerToDouble:
		return term.IntegerToDouble
	case doubleShow:
		return term.DoubleShow
	case optional:
		return term.Optional
	case optionalBuild:
		if v.typ != nil {
			return term.App{term.OptionalBuild, quoteWith(ctx, v.typ)}
		}
		return term.OptionalBuild
	case optionalFold:
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
	case none:
		return term.None
	case textShow:
		return term.TextShow
	case list:
		return term.List
	case listBuild:
		if v.typ != nil {
			return term.App{term.ListBuild, quoteWith(ctx, v.typ)}
		}
		return term.ListBuild
	case listFold:
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
	case listHead:
		if v.typ != nil {
			return term.App{term.ListHead, quoteWith(ctx, v.typ)}
		}
		return term.ListHead
	case listIndexed:
		if v.typ != nil {
			return term.App{term.ListIndexed, quoteWith(ctx, v.typ)}
		}
		return term.ListIndexed
	case listLength:
		if v.typ != nil {
			return term.App{term.ListLength, quoteWith(ctx, v.typ)}
		}
		return term.ListLength
	case listLast:
		if v.typ != nil {
			return term.App{term.ListLast, quoteWith(ctx, v.typ)}
		}
		return term.ListLast
	case listReverse:
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
	case lambda:
		bodyVal := v.Call(quoteVar{Name: v.Label, Index: ctx[v.Label]})
		return term.Lambda{
			Label: v.Label,
			Type:  quoteWith(ctx, v.Domain),
			Body:  quoteWith(ctx.extend(v.Label), bodyVal),
		}
	case Pi:
		bodyVal := v.Range(quoteVar{Name: v.Label, Index: ctx[v.Label]})
		return term.Pi{
			Label: v.Label,
			Type:  quoteWith(ctx, v.Domain),
			Body:  quoteWith(ctx.extend(v.Label), bodyVal),
		}
	case app:
		return term.App{
			Fn:  quoteWith(ctx, v.Fn),
			Arg: quoteWith(ctx, v.Arg),
		}
	case oper:
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
	case EmptyList:
		return term.EmptyList{Type: quoteWith(ctx, v.Type)}
	case NonEmptyList:
		l := term.NonEmptyList{}
		for _, e := range v {
			l = append(l, quoteWith(ctx, e))
		}
		return l
	case PlainTextLit:
		return term.PlainText(string(v))
	case interpolatedText:
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
	case Some:
		return term.Some{Val: quoteWith(ctx, v.Val)}
	case NoneOf:
		return term.Apply(term.None, quoteWith(ctx, v.Type))
	case RecordType:
		rt := term.RecordType{}
		for k, v := range v {
			rt[k] = quoteWith(ctx, v)
		}
		return rt
	case RecordLit:
		rt := term.RecordLit{}
		for k, v := range v {
			rt[k] = quoteWith(ctx, v)
		}
		return rt
	case toMap:
		result := term.ToMap{Record: quoteWith(ctx, v.Record)}
		if v.Type != nil {
			result.Type = quoteWith(ctx, v.Type)
		}
		return result
	case field:
		return term.Field{
			Record:    quoteWith(ctx, v.Record),
			FieldName: v.FieldName,
		}
	case project:
		return term.Project{
			Record:     quoteWith(ctx, v.Record),
			FieldNames: v.FieldNames,
		}
	case UnionType:
		result := term.UnionType{}
		for k, v := range v {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = quoteWith(ctx, v)
		}
		return result
	case merge:
		result := term.Merge{
			Handler: quoteWith(ctx, v.Handler),
			Union:   quoteWith(ctx, v.Union),
		}
		if v.Annotation != nil {
			result.Annotation = quoteWith(ctx, v.Annotation)
		}
		return result
	case assert:
		return term.Assert{Annotation: quoteWith(ctx, v.Annotation)}
	}
	panic(fmt.Sprintf("unknown Value type %#v", v))
}
