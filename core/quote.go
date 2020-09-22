package core

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/v5/term"
)

// Quote takes the Value v and turns it back into a Term.
func Quote(v Value) term.Term {
	return quoteWith(quoteContext{}, false, v)
}

// QuoteAlphaNormal takes the Value v and turns it back into a Term,
// in alpha-normal form: ie all labels are changed to `_`.
func QuoteAlphaNormal(v Value) term.Term {
	return quoteWith(quoteContext{}, true, v)
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

func quoteWith(ctx quoteContext, shouldAlphaNormalize bool, v Value) term.Term {
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
		result = term.App{result, quoteWith(ctx, shouldAlphaNormalize, v.n)}
		if v.typ == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, shouldAlphaNormalize, v.typ)}
		if v.succ == nil {
			return result
		}
		return term.App{result, quoteWith(ctx, shouldAlphaNormalize, v.succ)}
	case naturalIsZero:
		return term.NaturalIsZero
	case naturalOdd:
		return term.NaturalOdd
	case naturalShow:
		return term.NaturalShow
	case naturalSubtract:
		if v.a != nil {
			return term.App{term.NaturalSubtract, quoteWith(ctx, shouldAlphaNormalize, v.a)}
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
	case none:
		return term.None
	case textShow:
		return term.TextShow
	case list:
		return term.List
	case listBuild:
		if v.typ != nil {
			return term.App{term.ListBuild, quoteWith(ctx, shouldAlphaNormalize, v.typ)}
		}
		return term.ListBuild
	case listFold:
		var result term.Term = term.ListFold
		if v.typ1 == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, shouldAlphaNormalize, v.typ1)}
		if v.list == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, shouldAlphaNormalize, v.list)}
		if v.typ2 == nil {
			return result
		}
		result = term.App{result, quoteWith(ctx, shouldAlphaNormalize, v.typ2)}
		if v.cons == nil {
			return result
		}
		return term.App{result, quoteWith(ctx, shouldAlphaNormalize, v.cons)}
	case listHead:
		if v.typ != nil {
			return term.App{term.ListHead, quoteWith(ctx, shouldAlphaNormalize, v.typ)}
		}
		return term.ListHead
	case listIndexed:
		if v.typ != nil {
			return term.App{term.ListIndexed, quoteWith(ctx, shouldAlphaNormalize, v.typ)}
		}
		return term.ListIndexed
	case listLength:
		if v.typ != nil {
			return term.App{term.ListLength, quoteWith(ctx, shouldAlphaNormalize, v.typ)}
		}
		return term.ListLength
	case listLast:
		if v.typ != nil {
			return term.App{term.ListLast, quoteWith(ctx, shouldAlphaNormalize, v.typ)}
		}
		return term.ListLast
	case listReverse:
		if v.typ != nil {
			return term.App{term.ListReverse, quoteWith(ctx, shouldAlphaNormalize, v.typ)}
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
		label := v.Label
		if shouldAlphaNormalize {
			label = "_"
		}
		bodyVal := v.Call(quoteVar{Name: label, Index: ctx[label]})
		return term.Lambda{
			Label: label,
			Type:  quoteWith(ctx, shouldAlphaNormalize, v.Domain),
			Body:  quoteWith(ctx.extend(label), shouldAlphaNormalize, bodyVal),
		}
	case Pi:
		label := v.Label
		if shouldAlphaNormalize {
			label = "_"
		}
		bodyVal := v.Codomain(quoteVar{Name: label, Index: ctx[label]})
		return term.Pi{
			Label: label,
			Type:  quoteWith(ctx, shouldAlphaNormalize, v.Domain),
			Body:  quoteWith(ctx.extend(label), shouldAlphaNormalize, bodyVal),
		}
	case app:
		return term.App{
			Fn:  quoteWith(ctx, shouldAlphaNormalize, v.Fn),
			Arg: quoteWith(ctx, shouldAlphaNormalize, v.Arg),
		}
	case oper:
		return term.Op{
			OpCode: v.OpCode,
			L:      quoteWith(ctx, shouldAlphaNormalize, v.L),
			R:      quoteWith(ctx, shouldAlphaNormalize, v.R),
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
		return term.Apply(term.List, quoteWith(ctx, shouldAlphaNormalize, v.Type))
	case EmptyList:
		return term.EmptyList{Type: quoteWith(ctx, shouldAlphaNormalize, v.Type)}
	case NonEmptyList:
		l := term.NonEmptyList{}
		for _, e := range v {
			l = append(l, quoteWith(ctx, shouldAlphaNormalize, e))
		}
		return l
	case PlainTextLit:
		return term.PlainText(string(v))
	case interpolatedText:
		var newChunks term.Chunks
		for _, chunk := range v.Chunks {
			newChunks = append(newChunks, term.Chunk{
				Prefix: chunk.Prefix,
				Expr:   quoteWith(ctx, shouldAlphaNormalize, chunk.Expr),
			})
		}
		return term.TextLit{
			Chunks: newChunks,
			Suffix: v.Suffix,
		}
	case ifVal:
		return term.If{
			Cond: quoteWith(ctx, shouldAlphaNormalize, v.Cond),
			T:    quoteWith(ctx, shouldAlphaNormalize, v.T),
			F:    quoteWith(ctx, shouldAlphaNormalize, v.F),
		}
	case OptionalOf:
		return term.Apply(term.Optional, quoteWith(ctx, shouldAlphaNormalize, v.Type))
	case Some:
		return term.Some{Val: quoteWith(ctx, shouldAlphaNormalize, v.Val)}
	case NoneOf:
		return term.Apply(term.None, quoteWith(ctx, shouldAlphaNormalize, v.Type))
	case RecordType:
		rt := term.RecordType{}
		for k, v := range v {
			rt[k] = quoteWith(ctx, shouldAlphaNormalize, v)
		}
		return rt
	case RecordLit:
		rt := term.RecordLit{}
		for k, v := range v {
			rt[k] = quoteWith(ctx, shouldAlphaNormalize, v)
		}
		return rt
	case toMap:
		result := term.ToMap{Record: quoteWith(ctx, shouldAlphaNormalize, v.Record)}
		if v.Type != nil {
			result.Type = quoteWith(ctx, shouldAlphaNormalize, v.Type)
		}
		return result
	case field:
		return term.Field{
			Record:    quoteWith(ctx, shouldAlphaNormalize, v.Record),
			FieldName: v.FieldName,
		}
	case project:
		return term.Project{
			Record:     quoteWith(ctx, shouldAlphaNormalize, v.Record),
			FieldNames: v.FieldNames,
		}
	case with:
		return term.With{
			Record: quoteWith(ctx, shouldAlphaNormalize, v.Record),
			Path:   v.Path,
			Value:  quoteWith(ctx, shouldAlphaNormalize, v.Value),
		}
	case UnionType:
		result := term.UnionType{}
		for k, v := range v {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = quoteWith(ctx, shouldAlphaNormalize, v)
		}
		return result
	case unionConstructor:
		return term.Field{
			Record:    quoteWith(ctx, shouldAlphaNormalize, v.Type),
			FieldName: v.Alternative,
		}
	case unionVal:
		var result term.Term = term.Field{
			Record:    quoteWith(ctx, shouldAlphaNormalize, v.Type),
			FieldName: v.Alternative,
		}
		if v.Val != nil {
			result = term.App{
				Fn:  result,
				Arg: quoteWith(ctx, shouldAlphaNormalize, v.Val),
			}
		}
		return result
	case merge:
		result := term.Merge{
			Handler: quoteWith(ctx, shouldAlphaNormalize, v.Handler),
			Union:   quoteWith(ctx, shouldAlphaNormalize, v.Union),
		}
		if v.Annotation != nil {
			result.Annotation = quoteWith(ctx, shouldAlphaNormalize, v.Annotation)
		}
		return result
	case assert:
		return term.Assert{Annotation: quoteWith(ctx, shouldAlphaNormalize, v.Annotation)}
	}
	panic(fmt.Sprintf("unknown Value type %#v", v))
}
