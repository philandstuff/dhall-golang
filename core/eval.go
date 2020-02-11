package core

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/philandstuff/dhall-golang/term"
)

type Env map[string][]Value

// Eval normalizes Term to a Value.
func Eval(t term.Term) Value {
	return evalWith(t, Env{}, false)
}

// AlphaBetaEval alpha-beta-normalizes Term to a Value.
func AlphaBetaEval(t term.Term) Value {
	return evalWith(t, Env{}, true)
}

func evalWith(t term.Term, e Env, shouldAlphaNormalize bool) Value {
	switch t := t.(type) {
	case term.Universe:
		return Universe(t)
	case term.Builtin:
		switch t {
		case term.NaturalBuild:
			return NaturalBuildVal
		case term.NaturalEven:
			return NaturalEvenVal
		case term.NaturalFold:
			return NaturalFoldVal
		case term.NaturalIsZero:
			return NaturalIsZeroVal
		case term.NaturalOdd:
			return NaturalOddVal
		case term.NaturalShow:
			return NaturalShowVal
		case term.NaturalSubtract:
			return NaturalSubtractVal
		case term.NaturalToInteger:
			return NaturalToIntegerVal
		case term.IntegerClamp:
			return IntegerClampVal
		case term.IntegerNegate:
			return IntegerNegateVal
		case term.IntegerShow:
			return IntegerShowVal
		case term.IntegerToDouble:
			return IntegerToDoubleVal
		case term.DoubleShow:
			return DoubleShowVal
		case term.Optional:
			return OptionalVal
		case term.OptionalBuild:
			return OptionalBuildVal
		case term.OptionalFold:
			return OptionalFoldVal
		case term.None:
			return NoneVal
		case term.TextShow:
			return TextShowVal
		case term.List:
			return ListVal
		case term.ListBuild:
			return ListBuildVal
		case term.ListFold:
			return ListFoldVal
		case term.ListHead:
			return ListHeadVal
		case term.ListIndexed:
			return ListIndexedVal
		case term.ListLength:
			return ListLengthVal
		case term.ListLast:
			return ListLastVal
		case term.ListReverse:
			return ListReverseVal
		default:
			return Builtin(t)
		}
	case term.Var:
		if t.Index >= len(e[t.Name]) {
			return freeVar{t.Name, t.Index - len(e[t.Name])}
		}
		return e[t.Name][t.Index]
	case term.LocalVar:
		return localVar(t)
	case term.Lambda:
		v := lambdaValue{
			Label:  t.Label,
			Domain: evalWith(t.Type, e, shouldAlphaNormalize),
			Fn: func(x Value) Value {
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
	case term.Pi:
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
	case term.App:
		fn := evalWith(t.Fn, e, shouldAlphaNormalize)
		arg := evalWith(t.Arg, e, shouldAlphaNormalize)
		return applyVal(fn, arg)
	case term.Let:
		newEnv := Env{}
		for k, v := range e {
			newEnv[k] = v
		}

		for _, b := range t.Bindings {
			val := evalWith(b.Value, newEnv, shouldAlphaNormalize)
			newEnv[b.Variable] = append([]Value{val}, newEnv[b.Variable]...)
		}
		return evalWith(t.Body, newEnv, shouldAlphaNormalize)
	case term.Annot:
		return evalWith(t.Expr, e, shouldAlphaNormalize)
	case term.DoubleLit:
		return DoubleLit(t)
	case term.TextLit:
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
	case term.BoolLit:
		return BoolLit(t)
	case term.If:
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
		if AlphaEquivalentVals(tVal, fVal) {
			return tVal
		}
		return ifVal{
			Cond: condVal,
			T:    evalWith(t.T, e, shouldAlphaNormalize),
			F:    evalWith(t.F, e, shouldAlphaNormalize),
		}
	case term.NaturalLit:
		return NaturalLit(t)
	case term.IntegerLit:
		return IntegerLit(t)
	case term.Op:
		// these are cases where we *don't* evaluate t.L and t.R up front
		switch t.OpCode {
		case term.TextAppendOp:
			return evalWith(
				term.TextLit{Chunks: term.Chunks{{Expr: t.L}, {Expr: t.R}}},
				e, shouldAlphaNormalize)
		case term.CompleteOp:
			return evalWith(
				term.Annot{
					Expr: term.Op{
						OpCode: term.RightBiasedRecordMergeOp,
						L:      term.Field{t.L, "default"},
						R:      t.R,
					},
					Annotation: term.Field{t.L, "Type"},
				},
				e, shouldAlphaNormalize)
		}
		l := evalWith(t.L, e, shouldAlphaNormalize)
		r := evalWith(t.R, e, shouldAlphaNormalize)
		switch t.OpCode {
		case term.OrOp, term.AndOp, term.EqOp, term.NeOp:
			lb, lok := l.(BoolLit)
			rb, rok := r.(BoolLit)
			switch t.OpCode {
			case term.OrOp:
				if lok {
					if lb {
						return True
					}
					return r
				}
				if rok {
					if rb {
						return True
					}
					return l
				}
				if AlphaEquivalentVals(l, r) {
					return l
				}
			case term.AndOp:
				if lok {
					if lb {
						return r
					}
					return False
				}
				if rok {
					if rb {
						return l
					}
					return False
				}
				if AlphaEquivalentVals(l, r) {
					return l
				}
			case term.EqOp:
				if lok && bool(lb) {
					return r
				}
				if rok && bool(rb) {
					return l
				}
				if AlphaEquivalentVals(l, r) {
					return True
				}
			case term.NeOp:
				if lok && !bool(lb) {
					return r
				}
				if rok && !bool(rb) {
					return l
				}
				if AlphaEquivalentVals(l, r) {
					return False
				}
			}
		case term.ListAppendOp:
			if _, ok := l.(EmptyListVal); ok {
				return r
			}
			if _, ok := r.(EmptyListVal); ok {
				return l
			}
			ll, lok := l.(NonEmptyListVal)
			rl, rok := r.(NonEmptyListVal)
			if lok && rok {
				return append(ll, rl...)
			}
		case term.PlusOp:
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
		case term.TimesOp:
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
		case term.RecordMergeOp:
			lR, lOk := l.(RecordLitVal)
			rR, rOk := r.(RecordLitVal)

			if lOk && len(lR) == 0 {
				return r
			}
			if rOk && len(rR) == 0 {
				return l
			}
			if lOk && rOk {
				return mustMergeRecordLitVals(lR, rR)
			}
		case term.RecordTypeMergeOp:
			lRT, lOk := l.(RecordTypeVal)
			rRT, rOk := r.(RecordTypeVal)

			if lOk && len(lRT) == 0 {
				return r
			}
			if rOk && len(rRT) == 0 {
				return l
			}
			if lOk && rOk {
				result, err := mergeRecordTypes(lRT, rRT)
				if err != nil {
					panic(err) // shouldn't happen for well-typed terms
				}
				return result
			}
		case term.RightBiasedRecordMergeOp:
			lLit, lOk := l.(RecordLitVal)
			rLit, rOk := r.(RecordLitVal)
			if lOk && len(lLit) == 0 {
				return r
			}
			if rOk && len(rLit) == 0 {
				return l
			}
			if lOk && rOk {
				result := RecordLitVal{}
				for k, v := range lLit {
					result[k] = v
				}
				for k, v := range rLit {
					result[k] = v
				}
				return result
			}
			if AlphaEquivalentVals(l, r) {
				return l
			}
		case term.ImportAltOp:
			// nothing special
		case term.EquivOp:
			// nothing special
		}
		return opValue{OpCode: t.OpCode, L: l, R: r}
	case term.EmptyList:
		return EmptyListVal{Type: evalWith(t.Type, e, shouldAlphaNormalize)}
	case term.NonEmptyList:
		result := make([]Value, len(t))
		for i, t := range t {
			result[i] = evalWith(t, e, shouldAlphaNormalize)
		}
		return NonEmptyListVal(result)
	case term.Some:
		return SomeVal{evalWith(t.Val, e, shouldAlphaNormalize)}
	case term.RecordType:
		newRT := RecordTypeVal{}
		for k, v := range t {
			newRT[k] = evalWith(v, e, shouldAlphaNormalize)
		}
		return newRT
	case term.RecordLit:
		newRT := RecordLitVal{}
		for k, v := range t {
			newRT[k] = evalWith(v, e, shouldAlphaNormalize)
		}
		return newRT
	case term.ToMap:
		recordVal := evalWith(t.Record, e, shouldAlphaNormalize)
		record, ok := recordVal.(RecordLitVal)
		if ok {
			if len(record) == 0 {
				return EmptyListVal{Type: evalWith(t.Type, e, shouldAlphaNormalize)}
			}
			fieldnames := []string{}
			for k := range record {
				fieldnames = append(fieldnames, k)
			}
			sort.Strings(fieldnames)
			result := make(NonEmptyListVal, len(fieldnames))
			for i, k := range fieldnames {
				result[i] = RecordLitVal{"mapKey": TextLitVal{Suffix: k}, "mapValue": record[k]}
			}
			return result
		}
		return toMapVal{
			Record: record,
			Type:   evalWith(t.Type, e, shouldAlphaNormalize),
		}
	case term.Field:
		record := evalWith(t.Record, e, shouldAlphaNormalize)
		for { // simplifications
			if proj, ok := record.(projectVal); ok {
				record = proj.Record
				continue
			}
			op, ok := record.(opValue)
			if ok && op.OpCode == term.RecordMergeOp {
				if l, ok := op.L.(RecordLitVal); ok {
					if lField, ok := l[t.FieldName]; ok {
						return fieldVal{
							Record: opValue{
								L:      RecordLitVal{t.FieldName: lField},
								R:      op.R,
								OpCode: term.RecordMergeOp,
							},
							FieldName: t.FieldName,
						}
					}
					record = op.R
					continue
				}
				if r, ok := op.R.(RecordLitVal); ok {
					if rField, ok := r[t.FieldName]; ok {
						return fieldVal{
							Record: opValue{
								L:      op.L,
								R:      RecordLitVal{t.FieldName: rField},
								OpCode: term.RecordMergeOp,
							},
							FieldName: t.FieldName,
						}
					}
					record = op.L
					continue
				}
			}
			if ok && op.OpCode == term.RightBiasedRecordMergeOp {
				if l, ok := op.L.(RecordLitVal); ok {
					if lField, ok := l[t.FieldName]; ok {
						return fieldVal{
							Record: opValue{
								L:      RecordLitVal{t.FieldName: lField},
								R:      op.R,
								OpCode: term.RightBiasedRecordMergeOp,
							},
							FieldName: t.FieldName,
						}
					}
					record = op.R
					continue
				}
				if r, ok := op.R.(RecordLitVal); ok {
					if rField, ok := r[t.FieldName]; ok {
						return rField
					}
					record = op.L
					continue
				}
			}
			break
		}
		if lit, ok := record.(RecordLitVal); ok {
			return lit[t.FieldName]
		}
		return fieldVal{
			Record:    record,
			FieldName: t.FieldName,
		}
	case term.Project:
		record := evalWith(t.Record, e, shouldAlphaNormalize)
		fieldNames := t.FieldNames
		sort.Strings(fieldNames)
		// simplifications
		for {
			if proj, ok := record.(projectVal); ok {
				record = proj.Record
				continue
			}
			op, ok := record.(opValue)
			if ok && op.OpCode == term.RightBiasedRecordMergeOp {
				if r, ok := op.R.(RecordLitVal); ok {
					notOverridden := []string{}
					overrides := RecordLitVal{}
					for _, fieldName := range fieldNames {
						if override, ok := r[fieldName]; ok {
							overrides[fieldName] = override
						} else {
							notOverridden = append(notOverridden, fieldName)
						}
					}
					if len(notOverridden) == 0 {
						return overrides
					}
					return opValue{
						OpCode: term.RightBiasedRecordMergeOp,
						L: projectVal{
							Record:     op.L,
							FieldNames: notOverridden,
						},
						R: overrides,
					}
				}
			}

			break
		}
		if lit, ok := record.(RecordLitVal); ok {
			result := make(RecordLitVal)
			for _, k := range fieldNames {
				result[k] = lit[k]
			}
			return result
		}
		if len(fieldNames) == 0 {
			return RecordLitVal{}
		}
		return projectVal{
			Record:     record,
			FieldNames: fieldNames,
		}
	case term.ProjectType:
		// if `t` typechecks, `t.Selector` has to eval to a
		// RecordTypeVal, so this is safe
		s := evalWith(t.Selector, e, shouldAlphaNormalize).(RecordTypeVal)
		fieldNames := make([]string, 0, len(s))
		for fieldName := range s {
			fieldNames = append(fieldNames, fieldName)
		}
		return evalWith(
			term.Project{
				Record:     t.Record,
				FieldNames: fieldNames,
			},
			e, shouldAlphaNormalize)
	case term.UnionType:
		result := make(unionTypeVal, len(t))
		for k, v := range t {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = evalWith(v, e, shouldAlphaNormalize)
		}
		return result
	case term.Merge:
		handlerVal := evalWith(t.Handler, e, shouldAlphaNormalize)
		unionVal := evalWith(t.Union, e, shouldAlphaNormalize)
		if handlers, ok := handlerVal.(RecordLitVal); ok {
			// TODO: test tricky Field inputs
			if union, ok := unionVal.(appValue); ok {
				if field, ok := union.Fn.(fieldVal); ok {
					return applyVal(
						handlers[field.FieldName],
						union.Arg,
					)
				}
			}
			if union, ok := unionVal.(fieldVal); ok {
				// empty union alternative
				return handlers[union.FieldName]
			}
			if some, ok := unionVal.(SomeVal); ok {
				// Treating Optional as < Some a | None >
				return applyVal(
					handlers["Some"],
					some.Val,
				)
			}
			if _, ok := unionVal.(NoneOf); ok {
				// Treating Optional as < Some a | None >
				return handlers["None"]
			}
		}
		output := mergeVal{
			Handler: handlerVal,
			Union:   unionVal,
		}
		if t.Annotation != nil {
			output.Annotation = evalWith(t.Annotation, e, shouldAlphaNormalize)
		}
		return output
	case term.Assert:
		return assertVal{Annotation: evalWith(t.Annotation, e, shouldAlphaNormalize)}
	default:
		panic(fmt.Sprint("unknown term type", t))
	}
}

func applyVal(fn Value, args ...Value) Value {
	out := fn
	for _, arg := range args {
		if f, ok := out.(Callable); ok {
			if result := f.Call(arg); result != nil {
				out = result
				continue
			}
		}
		out = appValue{Fn: out, Arg: arg}
	}
	return out
}

func mergeRecordTypes(l RecordTypeVal, r RecordTypeVal) (RecordTypeVal, error) {
	var err error
	result := make(RecordTypeVal)
	for k, v := range l {
		result[k] = v
	}
	for k, v := range r {
		if lField, ok := result[k]; ok {
			lSubrecord, Lok := lField.(RecordTypeVal)
			rSubrecord, Rok := v.(RecordTypeVal)
			if !(Lok && Rok) {
				return nil, errors.New("Record mismatch")
			}
			result[k], err = mergeRecordTypes(lSubrecord, rSubrecord)
			if err != nil {
				return nil, err
			}
		} else {
			result[k] = v
		}
	}
	return result, nil
}

func mustMergeRecordLitVals(l RecordLitVal, r RecordLitVal) RecordLitVal {
	output := make(RecordLitVal)
	for k, v := range l {
		output[k] = v
	}
	for k, v := range r {
		if lField, ok := output[k]; ok {
			lSubrecord, Lok := lField.(RecordLitVal)
			rSubrecord, Rok := v.(RecordLitVal)
			if !(Lok && Rok) {
				// typecheck ought to have caught this
				panic("Record mismatch")
			}
			output[k] = mustMergeRecordLitVals(lSubrecord, rSubrecord)
		} else {
			output[k] = v
		}
	}
	return output
}
