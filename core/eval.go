package core

import (
	"fmt"
	"strings"
)

type Env map[string][]Value

// Eval normalizes Term to a Value.
func Eval(t Term) Value {
	return evalWith(t, Env{}, false)
}

// AlphaBetaEval alpha-beta-normalizes Term to a Value.
func AlphaBetaEval(t Term) Value {
	return evalWith(t, Env{}, true)
}

func evalWith(t Term, e Env, shouldAlphaNormalize bool) Value {
	switch t := t.(type) {
	case Universe:
		return t
	case Builtin:
		switch t {
		case NaturalBuild:
			return NaturalBuildVal
		case NaturalEven:
			return NaturalEvenVal
		case NaturalFold:
			return NaturalFoldVal
		case NaturalIsZero:
			return NaturalIsZeroVal
		case NaturalOdd:
			return NaturalOddVal
		case NaturalShow:
			return NaturalShowVal
		case NaturalSubtract:
			return NaturalSubtractVal
		case NaturalToInteger:
			return NaturalToIntegerVal
		case IntegerShow:
			return IntegerShowVal
		case IntegerToDouble:
			return IntegerToDoubleVal
		case DoubleShow:
			return DoubleShowVal
		case OptionalBuild:
			return OptionalBuildVal
		case OptionalFold:
			return OptionalFoldVal
		case TextShow:
			return TextShowVal
		case ListBuild:
			return ListBuildVal
		case ListFold:
			return ListFoldVal
		case ListHead:
			return ListHeadVal
		case ListIndexed:
			return ListIndexedVal
		case ListLength:
			return ListLengthVal
		case ListLast:
			return ListLastVal
		case ListReverse:
			return ListReverseVal
		default:
			return t
		}
	case BoundVar:
		if t.Index >= len(e[t.Name]) {
			panic(fmt.Sprintf("Eval: unbound variable %s", t))
		}
		return e[t.Name][t.Index]
	case LocalVar:
		return t
	case FreeVar:
		return t
	case LambdaTerm:
		v := LambdaValue{
			Label:  t.Label,
			Domain: evalWith(t.Type, e, shouldAlphaNormalize),
			hasCall1: func(x Value) Value {
				newEnv := Env{}
				for k, v := range e {
					newEnv[k] = v
				}
				newEnv[t.Label] = append([]Value{x}, newEnv[t.Label]...)
				return evalWith(t.Body, newEnv, shouldAlphaNormalize)
			},
		}
		if shouldAlphaNormalize {
			v.Label = "_"
		}
		return v
	case PiTerm:
		v := PiValue{
			Label:  t.Label,
			Domain: evalWith(t.Type, e, shouldAlphaNormalize),
			Range: func(x Value) Value {
				newEnv := Env{}
				for k, v := range e {
					newEnv[k] = v
				}
				newEnv[t.Label] = append([]Value{x}, newEnv[t.Label]...)
				return evalWith(t.Body, newEnv, shouldAlphaNormalize)
			}}
		if shouldAlphaNormalize {
			v.Label = "_"
		}
		return v
	case AppTerm:
		fn := evalWith(t.Fn, e, shouldAlphaNormalize)
		arg := evalWith(t.Arg, e, shouldAlphaNormalize)
		return applyVal1(fn, arg)
	case Let:
		newEnv := Env{}
		for k, v := range e {
			newEnv[k] = v
		}

		for _, b := range t.Bindings {
			val := evalWith(b.Value, newEnv, shouldAlphaNormalize)
			newEnv[b.Variable] = append([]Value{val}, newEnv[b.Variable]...)
		}
		return evalWith(t.Body, newEnv, shouldAlphaNormalize)
	case Annot:
		return evalWith(t.Expr, e, shouldAlphaNormalize)
	case DoubleLit:
		return t
	case TextLitTerm:
		var str strings.Builder
		var newChunks ChunkVals
		for _, chunk := range t.Chunks {
			str.WriteString(chunk.Prefix)
			normExpr := evalWith(chunk.Expr, e, shouldAlphaNormalize)
			if text, ok := normExpr.(TextLitVal); ok {
				if len(text.Chunks) != 0 {
					// first chunk gets the rest of str
					str.WriteString(text.Chunks[0].Prefix)
					newChunks = append(newChunks,
						ChunkVal{Prefix: str.String(), Expr: text.Chunks[0].Expr})
					newChunks = append(newChunks,
						text.Chunks[1:]...)
					str.Reset()
				}
				str.WriteString(text.Suffix)

			} else {
				newChunks = append(newChunks, ChunkVal{Prefix: str.String(), Expr: normExpr})
				str.Reset()
			}
		}
		str.WriteString(t.Suffix)
		newSuffix := str.String()

		// Special case: "${<expr>}" → <expr>
		if len(newChunks) == 1 && newChunks[0].Prefix == "" && newSuffix == "" {
			return newChunks[0].Expr
		}

		return TextLitVal{Chunks: newChunks, Suffix: newSuffix}
	case BoolLit:
		return t
	case IfTerm:
		condVal := evalWith(t.Cond, e, shouldAlphaNormalize)
		if condVal == True {
			return evalWith(t.T, e, shouldAlphaNormalize)
		}
		if condVal == False {
			return evalWith(t.F, e, shouldAlphaNormalize)
		}
		tVal := evalWith(t.T, e, shouldAlphaNormalize)
		fVal := evalWith(t.F, e, shouldAlphaNormalize)
		if tVal == True && fVal == False {
			return condVal
		}
		if judgmentallyEqualVals(tVal, fVal) {
			return tVal
		}
		return IfVal{
			Cond: condVal,
			T:    evalWith(t.T, e, shouldAlphaNormalize),
			F:    evalWith(t.F, e, shouldAlphaNormalize),
		}
	case NaturalLit:
		return t
	case IntegerLit:
		return t
	case OpTerm:
		l := evalWith(t.L, e, shouldAlphaNormalize)
		r := evalWith(t.R, e, shouldAlphaNormalize)
		switch t.OpCode {
		case OrOp, AndOp, EqOp, NeOp, TextAppendOp,
			ListAppendOp, RecordMergeOp, RightBiasedRecordMergeOp,
			RecordTypeMergeOp, ImportAltOp, EquivOp:
			return TextLitVal{Suffix: "OpTerm unimplemented"}
		case PlusOp:
			ln, lok := l.(NaturalLit)
			rn, rok := r.(NaturalLit)
			if lok && rok {
				return NaturalLit(ln + rn)
			}
			if l == NaturalLit(0) {
				return r
			}
			if r == NaturalLit(0) {
				return l
			}
		case TimesOp:
			ln, lok := l.(NaturalLit)
			rn, rok := r.(NaturalLit)
			if lok && rok {
				return NaturalLit(ln * rn)
			}
			if l == NaturalLit(0) {
				return NaturalLit(0)
			}
			if r == NaturalLit(0) {
				return NaturalLit(0)
			}
			if l == NaturalLit(1) {
				return r
			}
			if r == NaturalLit(1) {
				return l
			}
			return OpValue{OpCode: t.OpCode, L: l, R: r}
		}
		return OpValue{OpCode: t.OpCode, L: l, R: r}
	case EmptyList:
		return EmptyListVal{Type: evalWith(t.Type, e, shouldAlphaNormalize)}
	case NonEmptyList:
		result := make([]Value, len(t))
		for i, t := range t {
			result[i] = evalWith(t, e, shouldAlphaNormalize)
		}
		return NonEmptyListVal(result)
	case Some:
		return SomeVal{evalWith(t.Val, e, shouldAlphaNormalize)}
	case RecordType:
		newRT := RecordTypeVal{}
		for k, v := range t {
			newRT[k] = evalWith(v, e, shouldAlphaNormalize)
		}
		return newRT
	case RecordLit:
		newRT := RecordLitVal{}
		for k, v := range t {
			newRT[k] = evalWith(v, e, shouldAlphaNormalize)
		}
		return newRT
	case ToMap:
		return TextLitVal{Suffix: "ToMap unimplemented"}
	case Field:
		return FieldVal{
			Record:    evalWith(t.Record, e, shouldAlphaNormalize),
			FieldName: t.FieldName,
		}
	case Project:
		return TextLitVal{Suffix: "Project unimplemented"}
	case ProjectType:
		return TextLitVal{Suffix: "ProjectType unimplemented"}
	case UnionType:
		return TextLitVal{Suffix: "UnionType unimplemented"}
	case Merge:
		return TextLitVal{Suffix: "Merge unimplemented"}
	case Assert:
		return AssertVal{Annotation: evalWith(t.Annotation, e, shouldAlphaNormalize)}
	default:
		panic(fmt.Sprint("unknown term type", t))
	}
}

func applyVal1(fn Value, arg Value) Value {
	if f, ok := fn.(Callable1); ok {
		if result := f.Call1(arg); result != nil {
			return result
		}
	}
	if fn, ok := fn.(AppValue); ok {
		return applyVal2(fn.Fn, fn.Arg, arg)
	}
	return AppValue{Fn: fn, Arg: arg}
}

func applyVal2(fn Value, a, b Value) Value {
	if f, ok := fn.(Callable2); ok {
		if result := f.Call2(a, b); result != nil {
			return result
		}
	}
	if fn, ok := fn.(AppValue); ok {
		return applyVal3(fn.Fn, fn.Arg, a, b)
	}
	return AppValue{AppValue{fn, a}, b}
}

func applyVal3(fn Value, a, b, c Value) Value {
	if f, ok := fn.(Callable3); ok {
		if result := f.Call3(a, b, c); result != nil {
			return result
		}
	}
	if fn, ok := fn.(AppValue); ok {
		return applyVal4(fn.Fn, fn.Arg, a, b, c)
	}
	return AppValue{AppValue{AppValue{fn, a}, b}, c}
}

func applyVal4(fn Value, a, b, c, d Value) Value {
	if f, ok := fn.(Callable4); ok {
		if result := f.Call4(a, b, c, d); result != nil {
			return result
		}
	}
	if fn, ok := fn.(AppValue); ok {
		return applyVal5(fn.Fn, fn.Arg, a, b, c, d)
	}
	return AppValue{AppValue{AppValue{AppValue{fn, a}, b}, c}, d}
}

func applyVal5(fn Value, a, b, c, d, e Value) Value {
	if f, ok := fn.(Callable5); ok {
		if result := f.Call5(a, b, c, d, e); result != nil {
			return result
		}
	}
	return AppValue{AppValue{AppValue{AppValue{AppValue{fn, a}, b}, c}, d}, e}
}
