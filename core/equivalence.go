package core

import (
	"math"
)

// AlphaEquivalent reports if two Values are equal after
// alpha-normalization, as defined by the standard.  Broadly, two
// values are alpha-equivalent if they are structurally identical,
// ignoring label names.
func AlphaEquivalent(v1 Value, v2 Value) bool {
	return alphaEquivalentWith(0, v1, v2)
}

func alphaEquivalentWith(level int, v1 Value, v2 Value) bool {
	switch v1 := v1.(type) {
	case Universe, Builtin,
		naturalBuild, naturalEven, naturalFold,
		naturalIsZero, naturalOdd, naturalShow,
		naturalSubtract, naturalToInteger,
		integerShow, integerClamp, integerNegate, integerToDouble,
		doubleShow,
		optional, optionalBuild, optionalFold, none,
		textShow,
		list, listBuild, listFold, listHead, listIndexed,
		listLength, listLast, listReverse,
		freeVar, localVar, quoteVar,
		NaturalLit, IntegerLit, BoolLit, PlainTextLit:
		return v1 == v2
	case DoubleLit:
		v2, ok := v2.(DoubleLit)
		return ok && v1 == v2 && math.Signbit(float64(v1)) == math.Signbit(float64(v2))
	case lambda:
		v2, ok := v2.(lambda)
		if !ok {
			return false
		}
		// we deliberately ignore the Labels here
		return alphaEquivalentWith(level, v1.Domain, v2.Domain) &&
			alphaEquivalentWith(
				level+1,
				v1.Call(quoteVar{Name: "_", Index: level}),
				v2.Call(quoteVar{Name: "_", Index: level}),
			)
	case Pi:
		v2, ok := v2.(Pi)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Domain, v2.Domain) &&
			alphaEquivalentWith(
				level+1,
				v1.Codomain(quoteVar{Name: "_", Index: level}),
				v2.Codomain(quoteVar{Name: "_", Index: level}),
			)
	case app:
		v2, ok := v2.(app)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Fn, v2.Fn) &&
			alphaEquivalentWith(level, v1.Arg, v2.Arg)
	case oper:
		v2, ok := v2.(oper)
		if !ok {
			return false
		}
		return v1.OpCode == v2.OpCode &&
			alphaEquivalentWith(level, v1.L, v2.L) &&
			alphaEquivalentWith(level, v1.R, v2.R)
	case ListOf:
		v2, ok := v2.(ListOf)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Type, v2.Type)
	case EmptyList:
		v2, ok := v2.(EmptyList)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Type, v2.Type)
	case NonEmptyList:
		v2, ok := v2.(NonEmptyList)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for i := range v1 {
			if !alphaEquivalentWith(level, v1[i], v2[i]) {
				return false
			}
		}
		return true
	case interpolatedText:
		v2, ok := v2.(interpolatedText)
		if !ok {
			return false
		}
		if v1.Suffix != v2.Suffix ||
			len(v1.Chunks) != len(v2.Chunks) {
			return false
		}
		for i, c1 := range v1.Chunks {
			c2 := v2.Chunks[i]
			if c1.Prefix != c2.Prefix {
				return false
			}
			if !alphaEquivalentWith(level, c1.Expr, c2.Expr) {
				return false
			}
		}
		return true
	case ifVal:
		v2, ok := v2.(ifVal)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Cond, v2.Cond) &&
			alphaEquivalentWith(level, v1.T, v2.T) &&
			alphaEquivalentWith(level, v1.F, v2.F)
	case OptionalOf:
		v2, ok := v2.(OptionalOf)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Type, v2.Type)
	case Some:
		v2, ok := v2.(Some)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Val, v2.Val)
	case NoneOf:
		v2, ok := v2.(NoneOf)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Type, v2.Type)
	case RecordType:
		v2, ok := v2.(RecordType)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for k := range v1 {
			if v2[k] == nil ||
				!alphaEquivalentWith(level, v1[k], v2[k]) {
				return false
			}
		}
		return true
	case RecordLit:
		v2, ok := v2.(RecordLit)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for k := range v1 {
			if v2[k] == nil ||
				!alphaEquivalentWith(level, v1[k], v2[k]) {
				return false
			}
		}
		return true
	case toMap:
		v2, ok := v2.(toMap)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Record, v2.Record) &&
			alphaEquivalentWith(level, v1.Type, v2.Type)
	case field:
		v2, ok := v2.(field)
		if !ok {
			return false
		}
		return v1.FieldName == v2.FieldName &&
			alphaEquivalentWith(level, v1.Record, v2.Record)
	case project:
		v2, ok := v2.(project)
		if !ok {
			return false
		}
		if len(v1.FieldNames) != len(v2.FieldNames) {
			return false
		}
		for i := range v1.FieldNames {
			if v1.FieldNames[i] != v2.FieldNames[i] {
				return false
			}
		}
		return alphaEquivalentWith(level, v1.Record, v2.Record)
	case UnionType:
		v2, ok := v2.(UnionType)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for k := range v1 {
			if v1[k] == nil {
				if v2[k] != nil {
					return false
				}
				continue
			}
			if !alphaEquivalentWith(level, v1[k], v2[k]) {
				return false
			}
		}
		return true
	case merge:
		v2, ok := v2.(merge)
		if !ok {
			return false
		}
		if v1.Annotation != nil {
			if v2.Annotation == nil {
				return false
			}
			if !alphaEquivalentWith(level, v1.Annotation, v2.Annotation) {
				return false
			}
		}
		return alphaEquivalentWith(level, v1.Handler, v2.Handler) &&
			alphaEquivalentWith(level, v1.Union, v2.Union)
	case assert:
		v2, ok := v2.(assert)
		if !ok {
			return false
		}
		return alphaEquivalentWith(level, v1.Annotation, v2.Annotation)
	}
	panic("unknown Value type")
}
