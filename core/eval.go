package core

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/philandstuff/dhall-golang/v5/term"
)

type env map[string][]Value

// Eval normalizes Term to a Value.
func Eval(t term.Term) Value {
	return evalWith(t, env{})
}

func evalWith(t term.Term, e env) Value {
	switch t := t.(type) {
	case term.Universe:
		return Universe(t)
	case term.Builtin:
		switch t {
		case term.Bool:
			return Bool
		case term.Natural:
			return Natural
		case term.NaturalBuild:
			return NaturalBuild
		case term.NaturalEven:
			return NaturalEven
		case term.NaturalFold:
			return NaturalFold
		case term.NaturalIsZero:
			return NaturalIsZero
		case term.NaturalOdd:
			return NaturalOdd
		case term.NaturalShow:
			return NaturalShow
		case term.NaturalSubtract:
			return NaturalSubtract
		case term.NaturalToInteger:
			return NaturalToInteger
		case term.Integer:
			return Integer
		case term.IntegerClamp:
			return IntegerClamp
		case term.IntegerNegate:
			return IntegerNegate
		case term.IntegerShow:
			return IntegerShow
		case term.IntegerToDouble:
			return IntegerToDouble
		case term.Double:
			return Double
		case term.DoubleShow:
			return DoubleShow
		case term.Optional:
			return Optional
		case term.None:
			return None
		case term.Text:
			return Text
		case term.TextShow:
			return TextShow
		case term.List:
			return List
		case term.ListBuild:
			return ListBuild
		case term.ListFold:
			return ListFold
		case term.ListHead:
			return ListHead
		case term.ListIndexed:
			return ListIndexed
		case term.ListLength:
			return ListLength
		case term.ListLast:
			return ListLast
		case term.ListReverse:
			return ListReverse
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
		return lambda{
			Label:  t.Label,
			Domain: evalWith(t.Type, e),
			Fn: func(x Value) Value {
				newEnv := env{}
				for k, v := range e {
					newEnv[k] = v
				}
				newEnv[t.Label] = append([]Value{x}, newEnv[t.Label]...)
				return evalWith(t.Body, newEnv)
			},
		}
	case term.Pi:
		return Pi{
			Label:  t.Label,
			Domain: evalWith(t.Type, e),
			Codomain: func(x Value) Value {
				newEnv := env{}
				for k, v := range e {
					newEnv[k] = v
				}
				newEnv[t.Label] = append([]Value{x}, newEnv[t.Label]...)
				return evalWith(t.Body, newEnv)
			}}
	case term.App:
		fn := evalWith(t.Fn, e)
		arg := evalWith(t.Arg, e)
		return apply(fn, arg)
	case term.Let:
		newEnv := env{}
		for k, v := range e {
			newEnv[k] = v
		}

		for _, b := range t.Bindings {
			val := evalWith(b.Value, newEnv)
			newEnv[b.Variable] = append([]Value{val}, newEnv[b.Variable]...)
		}
		return evalWith(t.Body, newEnv)
	case term.Annot:
		return evalWith(t.Expr, e)
	case term.DoubleLit:
		return DoubleLit(t)
	case term.TextLit:
		var str strings.Builder
		var newChunks chunks
		for _, chk := range t.Chunks {
			str.WriteString(chk.Prefix)
			normExpr := evalWith(chk.Expr, e)
			if text, ok := normExpr.(PlainTextLit); ok {
				str.WriteString(string(text))
			} else if text, ok := normExpr.(interpolatedText); ok {
				// first chunk gets the rest of str
				str.WriteString(text.Chunks[0].Prefix)
				newChunks = append(newChunks,
					chunk{Prefix: str.String(), Expr: text.Chunks[0].Expr})
				newChunks = append(newChunks,
					text.Chunks[1:]...)
				str.Reset()
				str.WriteString(text.Suffix)
			} else {
				newChunks = append(newChunks, chunk{Prefix: str.String(), Expr: normExpr})
				str.Reset()
			}
		}
		str.WriteString(t.Suffix)
		newSuffix := str.String()

		// Special case: "${<expr>}" → <expr>
		if len(newChunks) == 1 && newChunks[0].Prefix == "" && newSuffix == "" {
			return newChunks[0].Expr
		}

		// Special case: no chunks -> PlainTextLit
		if len(newChunks) == 0 {
			return PlainTextLit(newSuffix)
		}

		return interpolatedText{Chunks: newChunks, Suffix: newSuffix}
	case term.BoolLit:
		return BoolLit(t)
	case term.If:
		condVal := evalWith(t.Cond, e)
		if condVal == True {
			return evalWith(t.T, e)
		}
		if condVal == False {
			return evalWith(t.F, e)
		}
		tVal := evalWith(t.T, e)
		fVal := evalWith(t.F, e)
		if tVal == True && fVal == False {
			return condVal
		}
		if AlphaEquivalent(tVal, fVal) {
			return tVal
		}
		return ifVal{
			Cond: condVal,
			T:    evalWith(t.T, e),
			F:    evalWith(t.F, e),
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
				e)
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
				e)
		}
		l := evalWith(t.L, e)
		r := evalWith(t.R, e)
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
				if AlphaEquivalent(l, r) {
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
				if AlphaEquivalent(l, r) {
					return l
				}
			case term.EqOp:
				if lok && bool(lb) {
					return r
				}
				if rok && bool(rb) {
					return l
				}
				if AlphaEquivalent(l, r) {
					return True
				}
			case term.NeOp:
				if lok && !bool(lb) {
					return r
				}
				if rok && !bool(rb) {
					return l
				}
				if AlphaEquivalent(l, r) {
					return False
				}
			}
		case term.ListAppendOp:
			if _, ok := l.(EmptyList); ok {
				return r
			}
			if _, ok := r.(EmptyList); ok {
				return l
			}
			ll, lok := l.(NonEmptyList)
			rl, rok := r.(NonEmptyList)
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
			lR, lOk := l.(RecordLit)
			rR, rOk := r.(RecordLit)

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
			lRT, lOk := l.(RecordType)
			rRT, rOk := r.(RecordType)

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
			lLit, lOk := l.(RecordLit)
			rLit, rOk := r.(RecordLit)
			if lOk && len(lLit) == 0 {
				return r
			}
			if rOk && len(rLit) == 0 {
				return l
			}
			if lOk && rOk {
				result := RecordLit{}
				for k, v := range lLit {
					result[k] = v
				}
				for k, v := range rLit {
					result[k] = v
				}
				return result
			}
			if AlphaEquivalent(l, r) {
				return l
			}
		case term.ImportAltOp:
			// nothing special
		case term.EquivOp:
			// nothing special
		}
		return oper{OpCode: t.OpCode, L: l, R: r}
	case term.EmptyList:
		return EmptyList{Type: evalWith(t.Type, e)}
	case term.NonEmptyList:
		result := make([]Value, len(t))
		for i, t := range t {
			result[i] = evalWith(t, e)
		}
		return NonEmptyList(result)
	case term.Some:
		return Some{evalWith(t.Val, e)}
	case term.RecordType:
		newRT := RecordType{}
		for k, v := range t {
			newRT[k] = evalWith(v, e)
		}
		return newRT
	case term.RecordLit:
		newRT := RecordLit{}
		for k, v := range t {
			newRT[k] = evalWith(v, e)
		}
		return newRT
	case term.ToMap:
		recordVal := evalWith(t.Record, e)
		record, ok := recordVal.(RecordLit)
		if ok {
			if len(record) == 0 {
				return EmptyList{Type: evalWith(t.Type, e)}
			}
			fieldnames := []string{}
			for k := range record {
				fieldnames = append(fieldnames, k)
			}
			sort.Strings(fieldnames)
			result := make(NonEmptyList, len(fieldnames))
			for i, k := range fieldnames {
				result[i] = RecordLit{"mapKey": PlainTextLit(k), "mapValue": record[k]}
			}
			return result
		}
		toMapVal := toMap{Record: recordVal}
		if t.Type != nil {
			toMapVal.Type = evalWith(t.Type, e)
		}
		return toMapVal
	case term.Field:
		record := evalWith(t.Record, e)
		for { // simplifications
			if proj, ok := record.(project); ok {
				record = proj.Record
				continue
			}
			op, ok := record.(oper)
			if ok && op.OpCode == term.RecordMergeOp {
				if l, ok := op.L.(RecordLit); ok {
					if lField, ok := l[t.FieldName]; ok {
						return field{
							Record: oper{
								L:      RecordLit{t.FieldName: lField},
								R:      op.R,
								OpCode: term.RecordMergeOp,
							},
							FieldName: t.FieldName,
						}
					}
					record = op.R
					continue
				}
				if r, ok := op.R.(RecordLit); ok {
					if rField, ok := r[t.FieldName]; ok {
						return field{
							Record: oper{
								L:      op.L,
								R:      RecordLit{t.FieldName: rField},
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
				if l, ok := op.L.(RecordLit); ok {
					if lField, ok := l[t.FieldName]; ok {
						return field{
							Record: oper{
								L:      RecordLit{t.FieldName: lField},
								R:      op.R,
								OpCode: term.RightBiasedRecordMergeOp,
							},
							FieldName: t.FieldName,
						}
					}
					record = op.R
					continue
				}
				if r, ok := op.R.(RecordLit); ok {
					if rField, ok := r[t.FieldName]; ok {
						return rField
					}
					record = op.L
					continue
				}
			}
			break
		}
		if lit, ok := record.(RecordLit); ok {
			return lit[t.FieldName]
		}
		if union, ok := record.(UnionType); ok {
			if union[t.FieldName] == nil {
				return unionVal{
					Type:        union,
					Alternative: t.FieldName,
				}
			}
			return unionConstructor{
				Type:        union,
				Alternative: t.FieldName,
			}
		}
		return field{
			Record:    record,
			FieldName: t.FieldName,
		}
	case term.Project:
		record := evalWith(t.Record, e)
		fieldNames := t.FieldNames
		sort.Strings(fieldNames)
		// simplifications
		for {
			if proj, ok := record.(project); ok {
				record = proj.Record
				continue
			}
			op, ok := record.(oper)
			if ok && op.OpCode == term.RightBiasedRecordMergeOp {
				if r, ok := op.R.(RecordLit); ok {
					notOverridden := []string{}
					overrides := RecordLit{}
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
					return oper{
						OpCode: term.RightBiasedRecordMergeOp,
						L: project{
							Record:     op.L,
							FieldNames: notOverridden,
						},
						R: overrides,
					}
				}
			}

			break
		}
		if lit, ok := record.(RecordLit); ok {
			result := make(RecordLit)
			for _, k := range fieldNames {
				result[k] = lit[k]
			}
			return result
		}
		if len(fieldNames) == 0 {
			return RecordLit{}
		}
		return project{
			Record:     record,
			FieldNames: fieldNames,
		}
	case term.ProjectType:
		// if `t` typechecks, `t.Selector` has to eval to a
		// RecordTypeVal, so this is safe
		s := evalWith(t.Selector, e).(RecordType)
		fieldNames := make([]string, 0, len(s))
		for fieldName := range s {
			fieldNames = append(fieldNames, fieldName)
		}
		return evalWith(
			term.Project{
				Record:     t.Record,
				FieldNames: fieldNames,
			},
			e)
	case term.UnionType:
		result := make(UnionType, len(t))
		for k, v := range t {
			if v == nil {
				result[k] = nil
				continue
			}
			result[k] = evalWith(v, e)
		}
		return result
	case term.Merge:
		handlerVal := evalWith(t.Handler, e)
		union := evalWith(t.Union, e)
		if handlers, ok := handlerVal.(RecordLit); ok {
			if unionLit, ok := union.(unionVal); ok {
				if unionLit.Val == nil {
					// empty union alternative
					return handlers[unionLit.Alternative]
				}
				return apply(handlers[unionLit.Alternative], unionLit.Val)
			}
			if some, ok := union.(Some); ok {
				// Treating Optional as < Some a | None >
				return apply(
					handlers["Some"],
					some.Val,
				)
			}
			if _, ok := union.(NoneOf); ok {
				// Treating Optional as < Some a | None >
				return handlers["None"]
			}
		}
		output := merge{
			Handler: handlerVal,
			Union:   union,
		}
		if t.Annotation != nil {
			output.Annotation = evalWith(t.Annotation, e)
		}
		return output
	case term.Assert:
		return assert{Annotation: evalWith(t.Annotation, e)}
	case term.With:
		record := evalWith(t.Record, e)
		value := evalWith(t.Value, e)
		output := record
		here := record
		var recordLit RecordLit
		depth := 0
		for _, component := range t.Path {
			var ok bool
			recordLit, ok = here.(RecordLit)
			if !ok {
				break
			}
			here = recordLit[component]
			depth = depth + 1
		}
		desugared := desugarWith(here, t.Path[depth:], value)
		if depth == 0 {
			return desugared
		}
		recordLit[t.Path[depth-1]] = desugared
		return output
	default:
		panic(fmt.Sprint("unknown term type", t))
	}
}

// desugarWith converts a `r with a.b...c = v` term to the equivalent,
// defined by desugar-with() in the Dhall standard.  Note that path
// may be of length 0, in which case value is returned.
func desugarWith(abstractRecord Value, path []string, value Value) Value {
	if len(path) == 0 {
		return value
	}
	return oper{
		OpCode: term.RightBiasedRecordMergeOp,
		L:      abstractRecord,
		R: RecordLit{path[0]: desugarWith(
			field{abstractRecord, path[0]},
			path[1:],
			value,
		)},
	}
}

func apply(fn Value, args ...Value) Value {
	out := fn
	for _, arg := range args {
		if f, ok := out.(Callable); ok {
			if result := f.Call(arg); result != nil {
				out = result
				continue
			}
		}
		out = app{Fn: out, Arg: arg}
	}
	return out
}

func mergeRecordTypes(l RecordType, r RecordType) (RecordType, error) {
	var err error
	result := make(RecordType)
	for k, v := range l {
		result[k] = v
	}
	for k, v := range r {
		if lField, ok := result[k]; ok {
			lSubrecord, Lok := lField.(RecordType)
			rSubrecord, Rok := v.(RecordType)
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

func mustMergeRecordLitVals(l RecordLit, r RecordLit) RecordLit {
	output := make(RecordLit)
	for k, v := range l {
		output[k] = v
	}
	for k, v := range r {
		if lField, ok := output[k]; ok {
			lSubrecord, Lok := lField.(RecordLit)
			rSubrecord, Rok := v.(RecordLit)
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
