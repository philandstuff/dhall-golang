package core

import (
	"math"

	"github.com/philandstuff/dhall-golang/term"
)

// AlphaEquivalent reports if two Terms are equal after
// alpha-normalization, as defined by the standard.  Broadly, two
// terms are alpha-equivalent if they are structurally identical,
// ignoring label names.
func AlphaEquivalent(t1 term.Term, t2 term.Term) bool {
	v1 := Eval(t1)
	v2 := Eval(t2)
	return AlphaEquivalentVals(v1, v2)
}

// AlphaEquivalentVals reports if two Values are equal after
// alpha-normalization, as defined by the standard.  Broadly, two
// values are alpha-equivalent if they are structurally identical,
// ignoring label names.
func AlphaEquivalentVals(v1 Value, v2 Value) bool {
	return alphaEquivalentValsWith(0, v1, v2)
}

func alphaEquivalentValsWith(level int, v1 Value, v2 Value) bool {
	switch v1 := v1.(type) {
	case Universe, Builtin,
		naturalBuildVal, naturalEvenVal, naturalFoldVal,
		naturalIsZeroVal, naturalOddVal, naturalShowVal,
		naturalSubtractVal, naturalToIntegerVal,
		integerShowVal, integerClampVal, integerNegateVal, integerToDoubleVal,
		doubleShowVal,
		optionalVal, optionalBuildVal, optionalFoldVal, noneVal,
		textShowVal,
		listVal, listBuildVal, listFoldVal, listHeadVal, listIndexedVal,
		listLengthVal, listLastVal, listReverseVal,
		freeVar, localVar, quoteVar,
		NaturalLit, IntegerLit, BoolLit:
		return v1 == v2
	case DoubleLit:
		v2, ok := v2.(DoubleLit)
		return ok && v1 == v2 && math.Signbit(float64(v1)) == math.Signbit(float64(v2))
	case lambdaValue:
		v2, ok := v2.(lambdaValue)
		if !ok {
			return false
		}
		// we deliberately ignore the Labels here
		return alphaEquivalentValsWith(level, v1.Domain, v2.Domain) &&
			alphaEquivalentValsWith(
				level+1,
				v1.Call(quoteVar{Name: "_", Index: level}),
				v2.Call(quoteVar{Name: "_", Index: level}),
			)
	case PiValue:
		v2, ok := v2.(PiValue)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Domain, v2.Domain) &&
			alphaEquivalentValsWith(
				level+1,
				v1.Range(quoteVar{Name: "_", Index: level}),
				v2.Range(quoteVar{Name: "_", Index: level}),
			)
	case appValue:
		v2, ok := v2.(appValue)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Fn, v2.Fn) &&
			alphaEquivalentValsWith(level, v1.Arg, v2.Arg)
	case opValue:
		v2, ok := v2.(opValue)
		if !ok {
			return false
		}
		return v1.OpCode == v2.OpCode &&
			alphaEquivalentValsWith(level, v1.L, v2.L) &&
			alphaEquivalentValsWith(level, v1.R, v2.R)
	case ListOf:
		v2, ok := v2.(ListOf)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Type, v2.Type)
	case EmptyListVal:
		v2, ok := v2.(EmptyListVal)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Type, v2.Type)
	case NonEmptyListVal:
		v2, ok := v2.(NonEmptyListVal)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for i := range v1 {
			if !alphaEquivalentValsWith(level, v1[i], v2[i]) {
				return false
			}
		}
		return true
	case TextLitVal:
		v2, ok := v2.(TextLitVal)
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
			if !alphaEquivalentValsWith(level, c1.Expr, c2.Expr) {
				return false
			}
		}
		return true
	case ifVal:
		v2, ok := v2.(ifVal)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Cond, v2.Cond) &&
			alphaEquivalentValsWith(level, v1.T, v2.T) &&
			alphaEquivalentValsWith(level, v1.F, v2.F)
	case OptionalOf:
		v2, ok := v2.(OptionalOf)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Type, v2.Type)
	case SomeVal:
		v2, ok := v2.(SomeVal)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Val, v2.Val)
	case NoneOf:
		v2, ok := v2.(NoneOf)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Type, v2.Type)
	case RecordTypeVal:
		v2, ok := v2.(RecordTypeVal)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for k := range v1 {
			if v2[k] == nil ||
				!alphaEquivalentValsWith(level, v1[k], v2[k]) {
				return false
			}
		}
		return true
	case RecordLitVal:
		v2, ok := v2.(RecordLitVal)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		for k := range v1 {
			if v2[k] == nil ||
				!alphaEquivalentValsWith(level, v1[k], v2[k]) {
				return false
			}
		}
		return true
	case toMapVal:
		v2, ok := v2.(toMapVal)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Record, v2.Record) &&
			alphaEquivalentValsWith(level, v1.Type, v2.Type)
	case fieldVal:
		v2, ok := v2.(fieldVal)
		if !ok {
			return false
		}
		return v1.FieldName == v2.FieldName &&
			alphaEquivalentValsWith(level, v1.Record, v2.Record)
	case projectVal:
		v2, ok := v2.(projectVal)
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
		return alphaEquivalentValsWith(level, v1.Record, v2.Record)
	case unionTypeVal:
		v2, ok := v2.(unionTypeVal)
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
			if !alphaEquivalentValsWith(level, v1[k], v2[k]) {
				return false
			}
		}
		return true
	case mergeVal:
		v2, ok := v2.(mergeVal)
		if !ok {
			return false
		}
		if v1.Annotation != nil {
			if v2.Annotation == nil {
				return false
			}
			if !alphaEquivalentValsWith(level, v1.Annotation, v2.Annotation) {
				return false
			}
		}
		return alphaEquivalentValsWith(level, v1.Handler, v2.Handler) &&
			alphaEquivalentValsWith(level, v1.Union, v2.Union)
	case assertVal:
		v2, ok := v2.(assertVal)
		if !ok {
			return false
		}
		return alphaEquivalentValsWith(level, v1.Annotation, v2.Annotation)
	}
	panic("unknown Value type")
}
