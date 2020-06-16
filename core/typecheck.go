package core

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/v4/term"
)

type context map[string][]Value

func (ctx context) extend(name string, t Value) context {
	newctx := context{}
	for k, v := range ctx {
		newctx[k] = v
	}
	newctx[name] = append(newctx[name], t)
	return newctx
}

func (ctx context) freshLocal(name string) term.LocalVar {
	return term.LocalVar{Name: name, Index: len(ctx[name])}
}

func assertTypeIs(ctx context, expr term.Term, expectedType Value, msg typeMessage) error {
	actualType, err := typeWith(ctx, expr)
	if err != nil {
		return err
	}
	if !AlphaEquivalent(expectedType, actualType) {
		return mkTypeError(msg)
	}
	return nil
}

// This returns
//  Value: the element type of a list type
//  Bool: whether it succeeded
func listElementType(e Value) (Value, bool) {
	list, ok := e.(ListOf)
	if !ok {
		return nil, false
	}
	return list.Type, true
}

func functionCheck(input Universe, output Universe) Universe {
	switch {
	case output == Type:
		return Type
	case input < output:
		return output
	default:
		return input
	}
}

// TypeOf typechecks a Term, returning the type in normal form.  If
// typechecking fails, an error is returned.
func TypeOf(t term.Term) (Value, error) {
	v, err := typeWith(context{}, t)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func typeWith(ctx context, t term.Term) (Value, error) {
	switch t := t.(type) {
	case term.Universe:
		switch t {
		case term.Type:
			return Kind, nil
		case term.Kind:
			return Sort, nil
		case term.Sort:
			return nil, mkTypeError(untyped)
		default:
			return nil, mkTypeError(unhandledTypeCase)
		}
	case term.Builtin:
		switch t {
		case term.Bool, term.Double, term.Integer, term.Natural, term.Text:
			return Type, nil
		case term.DoubleShow:
			return NewFnType("_", Double, Text), nil
		case term.IntegerClamp:
			return NewFnType("_", Integer, Natural), nil
		case term.IntegerNegate:
			return NewFnType("_", Integer, Integer), nil
		case term.IntegerShow:
			return NewFnType("_", Integer, Text), nil
		case term.IntegerToDouble:
			return NewFnType("_", Integer, Double), nil
		case term.List, term.Optional:
			return NewFnType("_", Type, Type), nil
		case term.ListBuild:
			return NewPi("a", Type, func(a Value) Value {
				return NewFnType("_",
					NewPi("list", Type, func(list Value) Value {
						return NewFnType("cons",
							NewFnType("_", a, NewFnType("_", list, list)),
							NewFnType("nil", list, list))
					}),
					ListOf{a})
			}), nil
		case term.ListFold:
			return NewPi("a", Type, func(a Value) Value {
				return NewFnType("_",
					ListOf{a},
					NewPi("list", Type, func(list Value) Value {
						return NewFnType("cons",
							NewFnType("_", a, NewFnType("_", list, list)),
							NewFnType("nil", list, list))
					}))
			}), nil
		case term.ListLength:
			return NewPi("a", Type, func(a Value) Value {
				return NewFnType("_", ListOf{a}, Natural)
			}), nil
		case term.ListHead, term.ListLast:
			return NewPi("a", Type, func(a Value) Value {
				return NewFnType("_", ListOf{a},
					OptionalOf{a})
			}), nil
		case term.ListReverse:
			return NewPi("a", Type, func(a Value) Value {
				return NewFnType("_", ListOf{a},
					ListOf{a})
			}), nil
		case term.ListIndexed:
			return NewPi("a", Type, func(a Value) Value {
				return NewFnType("_", ListOf{a},
					ListOf{RecordType{"index": Natural, "value": a}})
			}), nil
		case term.NaturalBuild:
			return NewFnType("_",
				NewPi("natural", Type, func(natural Value) Value {
					return NewFnType("succ",
						NewFnType("_", natural, natural),
						NewFnType("zero", natural, natural))
				}),
				Natural), nil
		case term.NaturalFold:
			return NewFnType("_",
				Natural,
				NewPi("natural", Type, func(natural Value) Value {
					return NewFnType("succ",
						NewFnType("_", natural, natural),
						NewFnType("zero", natural, natural))
				})), nil
		case term.NaturalIsZero, term.NaturalOdd, term.NaturalEven:
			return NewFnType("_", Natural, Bool), nil
		case term.NaturalShow:
			return NewFnType("_", Natural, Text), nil
		case term.NaturalToInteger:
			return NewFnType("_", Natural, Integer), nil
		case term.NaturalSubtract:
			return NewFnType("_", Natural, NewFnType("_", Natural, Natural)), nil
		case term.None:
			return NewPi("A", Type, func(A Value) Value { return OptionalOf{A} }), nil
		case term.TextShow:
			return NewFnType("_", Text, Text), nil
		default:
			return nil, mkTypeError(unhandledTypeCase)
		}
	case term.Var:
		return nil, mkTypeError(typeCheckVar(t))
	case term.LocalVar:
		if vals, ok := ctx[t.Name]; ok {
			if t.Index < len(vals) {
				return vals[t.Index], nil
			}
			return nil, mkTypeError(unboundVariable(t))
		}
		return nil, fmt.Errorf("Unknown variable %s", t.Name)
	case term.App:
		fnType, err := typeWith(ctx, t.Fn)
		if err != nil {
			return nil, err
		}
		argType, err := typeWith(ctx, t.Arg)
		if err != nil {
			return nil, err
		}
		piType, ok := fnType.(Pi)
		if !ok {
			return nil, mkTypeError(notAFunction)
		}
		expectedType := piType.Domain
		actualType := argType
		if !AlphaEquivalent(expectedType, actualType) {
			return nil, mkTypeError(typeMismatch(Quote(expectedType), Quote(actualType)))
		}
		bodyTypeVal := piType.Codomain(Eval(t.Arg))
		return bodyTypeVal, nil
	case term.Lambda:
		_, err := typeWith(ctx, t.Type)
		if err != nil {
			return nil, err
		}
		argType := Eval(t.Type)
		pi := Pi{Label: t.Label, Domain: argType}
		freshLocal := ctx.freshLocal(t.Label)
		bt, err := typeWith(
			ctx.extend(t.Label, argType),
			term.Subst(t.Label, freshLocal, t.Body))
		if err != nil {
			return nil, err
		}
		pi.Codomain = func(x Value) Value {
			rebound := term.RebindLocal(freshLocal, Quote(bt))
			return evalWith(rebound, env{
				t.Label: []Value{x},
			})
		}
		_, err = typeWith(ctx, Quote(pi))
		if err != nil {
			return nil, err
		}
		return pi, nil
	case term.Pi:
		inUniv, err := typeWith(ctx, t.Type)
		if err != nil {
			return nil, err
		}
		i, ok := inUniv.(Universe)
		if !ok {
			return nil, mkTypeError(invalidInputType)
		}
		freshLocal := ctx.freshLocal(t.Label)
		outUniv, err := typeWith(
			ctx.extend(t.Label, Eval(t.Type)),
			term.Subst(t.Label, freshLocal, t.Body))
		if err != nil {
			return nil, err
		}
		o, ok := outUniv.(Universe)
		if !ok {
			return nil, mkTypeError(invalidOutputType)
		}
		return functionCheck(i, o), nil
	case term.NaturalLit:
		return Natural, nil
	case term.Let:
		let := t
		for len(let.Bindings) > 0 {
			binding := let.Bindings[0]
			let.Bindings = let.Bindings[1:]

			bindingType, err := typeWith(ctx, binding.Value)
			if err != nil {
				return nil, err
			}

			if binding.Annotation != nil {
				_, err := typeWith(ctx, binding.Annotation)
				if err != nil {
					return nil, err
				}
				if !AlphaEquivalent(bindingType, Eval(binding.Annotation)) {
					return nil, mkTypeError(annotMismatch(binding.Annotation, Quote(bindingType)))
				}
			}

			value := Quote(Eval(binding.Value))
			let = term.Subst(binding.Variable, value, let).(term.Let)
			ctx = ctx.extend(binding.Variable, bindingType)
		}
		return typeWith(ctx, let.Body)
	case term.Annot:
		if t.Annotation != term.Sort {
			// Γ ⊢ T₀ : i
			if _, err := typeWith(ctx, t.Annotation); err != nil {
				return nil, err
			}
		}
		// Γ ⊢ t : T₁
		actualType, err := typeWith(ctx, t.Expr)
		if err != nil {
			return nil, err
		}
		// T₀ ≡ T₁
		if !AlphaEquivalent(Eval(t.Annotation), actualType) {
			return nil, mkTypeError(annotMismatch(t.Annotation, Quote(actualType)))
		}
		// ─────────────────
		// Γ ⊢ (t : T₀) : T₀
		return actualType, nil
	case term.DoubleLit:
		return Double, nil
	case term.TextLit:
		for _, chunk := range t.Chunks {
			err := assertTypeIs(ctx, chunk.Expr, Text,
				cantInterpolate)
			if err != nil {
				return nil, err
			}
		}
		return Text, nil
	case term.BoolLit:
		return Bool, nil
	case term.If:
		condType, err := typeWith(ctx, t.Cond)
		if err != nil {
			return nil, err
		}
		if condType != Bool {
			return nil, mkTypeError(invalidPredicate)
		}
		L, err := typeWith(ctx, t.T)
		if err != nil {
			return nil, err
		}
		// no need to check for err here
		if t, _ := typeWith(ctx, Quote(L)); t != Type {
			return nil, mkTypeError(ifBranchMustBeTerm)
		}
		R, err := typeWith(ctx, t.F)
		if err != nil {
			return nil, err
		}
		if t, _ := typeWith(ctx, Quote(R)); t != Type {
			return nil, mkTypeError(ifBranchMustBeTerm)
		}
		if !AlphaEquivalent(L, R) {
			return nil, mkTypeError(ifBranchMismatch)
		}
		return L, nil
	case term.IntegerLit:
		return Integer, nil
	case term.Op:
		switch t.OpCode {
		case term.OrOp, term.AndOp, term.EqOp, term.NeOp:
			err := assertTypeIs(ctx, t.L, Bool, cantBoolOp(t.OpCode))
			if err != nil {
				return nil, err
			}
			err = assertTypeIs(ctx, t.R, Bool, cantBoolOp(t.OpCode))
			if err != nil {
				return nil, err
			}
			return Bool, nil
		case term.PlusOp, term.TimesOp:
			err := assertTypeIs(ctx, t.L, Natural, cantNaturalOp(t.OpCode))
			if err != nil {
				return nil, err
			}
			err = assertTypeIs(ctx, t.R, Natural, cantNaturalOp(t.OpCode))
			if err != nil {
				return nil, err
			}
			return Natural, nil
		case term.TextAppendOp:
			err := assertTypeIs(ctx, t.L, Text, cantTextAppend)
			if err != nil {
				return nil, err
			}
			err = assertTypeIs(ctx, t.R, Text, cantTextAppend)
			if err != nil {
				return nil, err
			}
			return Text, nil
		case term.ListAppendOp:
			lt, err := typeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rt, err := typeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}

			lElemT, ok := listElementType(lt)
			if !ok {
				return nil, mkTypeError(cantListAppend)
			}
			rElemT, ok := listElementType(rt)
			if !ok {
				return nil, mkTypeError(cantListAppend)
			}
			if !AlphaEquivalent(lElemT, rElemT) {
				return nil, mkTypeError(listAppendMismatch)
			}
			return lt, nil
		case term.RecordMergeOp:
			lType, err := typeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rType, err := typeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}
			recordType := term.Op{L: Quote(lType), R: Quote(rType), OpCode: term.RecordTypeMergeOp}
			if _, err = typeWith(ctx, recordType); err != nil {
				return nil, err
			}
			return Eval(recordType), nil
		case term.RecordTypeMergeOp:
			lKind, err := typeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rKind, err := typeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}
			lt, ok := Eval(t.L).(RecordType)
			if !ok {
				return nil, mkTypeError(combineTypesRequiresRecordType)
			}
			rt, ok := Eval(t.R).(RecordType)
			if !ok {
				return nil, mkTypeError(combineTypesRequiresRecordType)
			}
			// ensure that the records are safe to merge
			if _, err := mergeRecordTypes(lt, rt); err != nil {
				return nil, err
			}
			if lKind.(Universe) > rKind.(Universe) {
				return lKind, nil
			}
			return rKind, nil
		case term.RightBiasedRecordMergeOp:
			lType, err := typeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rType, err := typeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}
			lRecord, ok := lType.(RecordType)
			if !ok {
				return nil, mkTypeError(mustCombineARecord)
			}
			rRecord, ok := rType.(RecordType)
			if !ok {
				return nil, mkTypeError(mustCombineARecord)
			}
			result := RecordType{}
			for k, v := range lRecord {
				result[k] = v
			}
			for k, v := range rRecord {
				result[k] = v
			}
			return result, nil
		case term.ImportAltOp:
			return typeWith(ctx, t.L)
		case term.EquivOp:
			lType, err := typeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rType, err := typeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}
			err = assertTypeIs(ctx, Quote(lType), Type, incomparableExpression)
			if err != nil {
				return nil, err
			}
			err = assertTypeIs(ctx, Quote(rType), Type, incomparableExpression)
			if err != nil {
				return nil, err
			}
			if !AlphaEquivalent(lType, rType) {
				return nil, mkTypeError(equivalenceTypeMismatch)
			}
			return Type, nil
		case term.CompleteOp:
			return typeWith(ctx,
				term.Annot{
					Expr: term.Op{OpCode: term.RightBiasedRecordMergeOp,
						L: term.Field{Record: t.L, FieldName: "default"},
						R: t.R,
					},
					Annotation: term.Field{Record: t.L, FieldName: "Type"},
				},
			)
		default:
			return nil, fmt.Errorf("Internal error: unknown opcode %v", t.OpCode)
		}
	case term.EmptyList:
		_, err := typeWith(ctx, t.Type)
		if err != nil {
			return nil, err
		}
		listType := Eval(t.Type)
		_, ok := listElementType(listType)
		if !ok {
			return nil, mkTypeError(invalidListType)
		}
		return listType, nil
	case term.NonEmptyList:
		T0, err := typeWith(ctx, t[0])
		if err != nil {
			return nil, err
		}
		err = assertTypeIs(ctx, Quote(T0), Type, invalidListType)
		if err != nil {
			return nil, err
		}
		for _, e := range t[1:] {
			T1, err := typeWith(ctx, e)
			if err != nil {
				return nil, err
			}
			if !AlphaEquivalent(T0, T1) {
				return nil, mkTypeError(mismatchedListElements(Quote(T0), Quote(T1)))
			}
		}
		return ListOf{T0}, nil
	case term.Some:
		A, err := typeWith(ctx, t.Val)
		if err != nil {
			return nil, err
		}
		if err = assertTypeIs(ctx, Quote(A), Type, invalidSome); err != nil {
			return nil, err
		}
		return OptionalOf{A}, nil
	case term.RecordType:
		recordUniverse := Type
		for _, v := range t {
			fieldUniverse, err := typeWith(ctx, v)
			if err != nil {
				return nil, err
			}
			u, ok := fieldUniverse.(Universe)
			if !ok {
				return nil, mkTypeError(invalidFieldType)
			}
			if recordUniverse < u {
				recordUniverse = u
			}
		}
		return recordUniverse, nil
	case term.RecordLit:
		recordType := RecordType{}
		for k, v := range t {
			fieldType, err := typeWith(ctx, v)
			if err != nil {
				return nil, err
			}
			recordType[k] = fieldType
		}
		if _, err := typeWith(ctx, Quote(recordType)); err != nil {
			return nil, err
		}
		return recordType, nil
	case term.ToMap:
		recordTypeVal, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeVal.(RecordType)
		if !ok {
			return nil, mkTypeError(cantAccess)
		}

		if len(recordType) == 0 {
			if t.Type == nil {
				return nil, mkTypeError(missingToMapType)
			}
			err = assertTypeIs(ctx, t.Type, Type, invalidToMapRecordKind)
			if err != nil {
				return nil, err
			}
			tVal := Eval(t.Type)
			t, ok := listElementType(tVal)
			if !ok {
				return nil, mkTypeError(invalidToMapType(Quote(tVal)))
			}
			rt, ok := t.(RecordType)
			if !ok || len(rt) != 2 || rt["mapKey"] != Text || rt["mapValue"] == nil {
				return nil, mkTypeError(invalidToMapType(Quote(tVal)))
			}
			return tVal, nil
		}

		var elemType Value
		for _, v := range recordType {
			if elemType == nil {
				elemType = v
			} else {
				if !AlphaEquivalent(elemType, v) {
					return nil, mkTypeError(heterogenousRecordToMap)
				}
			}
		}
		if k, _ := typeWith(ctx, Quote(elemType)); k != Type {
			return nil, mkTypeError(invalidToMapRecordKind)
		}
		inferred := ListOf{RecordType{"mapKey": Text, "mapValue": elemType}}
		if t.Type == nil {
			return inferred, nil
		}
		if _, err = typeWith(ctx, t.Type); err != nil {
			return nil, err
		}
		annot := Eval(t.Type)
		if !AlphaEquivalent(inferred, annot) {
			return nil, mkTypeError(mapTypeMismatch(Quote(inferred), t.Type))
		}
		return inferred, nil
	case term.Field:
		recordTypeVal, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeVal.(RecordType)
		if ok {
			fieldType, ok := recordType[t.FieldName]
			if !ok {
				return nil, mkTypeError(missingField)
			}
			return fieldType, nil
		}
		unionTypeV := Eval(t.Record)
		unionType, ok := unionTypeV.(UnionType)
		if !ok {
			return nil, mkTypeError(cantAccess)
		}
		alternativeType, ok := unionType[t.FieldName]
		if !ok {
			return nil, mkTypeError(missingConstructor)
		}
		if alternativeType == nil {
			return unionType, nil
		}
		return Pi{
			Label:    t.FieldName,
			Domain:   alternativeType,
			Codomain: func(Value) Value { return unionType },
		}, nil
	case term.Project:
		recordTypeVal, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeVal.(RecordType)
		if !ok {
			return nil, mkTypeError(cantProject)
		}
		result := make(RecordType, len(t.FieldNames))
		seen := make(map[string]bool)
		for _, name := range t.FieldNames {
			if seen[name] {
				return nil, mkTypeError(duplicateProjectedField(name))
			}
			seen[name] = true

			var ok bool
			result[name], ok = recordType[name]
			if !ok {
				return nil, mkTypeError(missingField)
			}
		}
		return result, nil
	case term.ProjectType:
		recordTypeVal, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeVal.(RecordType)
		if !ok {
			return nil, mkTypeError(cantProject)
		}
		_, err = typeWith(ctx, t.Selector)
		if err != nil {
			return nil, err
		}
		selectorVal := Eval(t.Selector)
		selector, ok := selectorVal.(RecordType)
		if !ok {
			return nil, mkTypeError(cantProjectByExpression)
		}
		result := make(RecordType, len(selector))
		for name, typ := range selector {
			fieldType, ok := recordType[name]
			if !ok {
				return nil, mkTypeError(missingField)
			}
			if !AlphaEquivalent(fieldType, typ) {
				return nil, mkTypeError(projectionTypeMismatch(Quote(typ), Quote(fieldType)))
			}
			result[name] = typ
		}
		return result, nil
	case term.UnionType:
		universe := Type
		for _, typ := range t {
			if typ == nil {
				// empty alternative
				continue
			}
			k, err := typeWith(ctx, typ)
			if err != nil {
				return nil, err
			}
			u, ok := k.(Universe)
			if !ok {
				return nil, mkTypeError(invalidAlternativeType)
			}
			if u > universe {
				universe = u
			}
		}
		return universe, nil
	case term.Merge:
		handlerTypeVal, err := typeWith(ctx, t.Handler)
		if err != nil {
			return nil, err
		}
		unionTypeV, err := typeWith(ctx, t.Union)
		if err != nil {
			return nil, err
		}
		handlerType, ok := handlerTypeVal.(RecordType)
		if !ok {
			return nil, mkTypeError(mustMergeARecord)
		}
		unionType, ok := unionTypeV.(UnionType)
		if !ok {
			opt, ok := unionTypeV.(OptionalOf)
			if !ok {
				return nil, mkTypeError(mustMergeUnion)
			}
			unionType = UnionType{"Some": opt.Type, "None": nil}
		}
		if len(handlerType) > len(unionType) {
			return nil, mkTypeError(unusedHandler)
		}

		if len(handlerType) == 0 {
			if t.Annotation == nil {
				return nil, mkTypeError(missingMergeType)
			}
			if _, err := typeWith(ctx, t.Annotation); err != nil {
				return nil, err
			}
			return Eval(t.Annotation), nil
		}

		var result Value
		for altName, altType := range unionType {
			fieldType, ok := handlerType[altName]
			if !ok {
				return nil, mkTypeError(missingHandler)
			}
			if altType == nil {
				if result == nil {
					result = fieldType
				} else {
					if !AlphaEquivalent(result, fieldType) {
						return nil, mkTypeError(handlerOutputTypeMismatch(Quote(result), Quote(fieldType)))
					}
				}
			} else {
				pi, ok := fieldType.(Pi)
				if !ok {
					return nil, mkTypeError(handlerNotAFunction)
				}
				if !AlphaEquivalent(altType, pi.Domain) {
					return nil, mkTypeError(handlerInputTypeMismatch(Quote(altType), Quote(pi.Domain)))
				}
				outputType := pi.Codomain(NaturalLit(1))
				outputType2 := pi.Codomain(NaturalLit(2))
				if !AlphaEquivalent(outputType, outputType2) {
					// hacky way of detecting output type depending on input
					return nil, mkTypeError(disallowedHandlerType)
				}
				if result == nil {
					result = outputType
				} else {
					if !AlphaEquivalent(result, outputType) {
						return nil, mkTypeError(handlerOutputTypeMismatch(Quote(result), Quote(outputType)))
					}
				}
			}
		}
		if t.Annotation != nil {
			if _, err := typeWith(ctx, t.Annotation); err != nil {
				return nil, err
			}
			if !AlphaEquivalent(result, Eval(t.Annotation)) {
				return nil, mkTypeError(annotMismatch(t.Annotation, Quote(result)))
			}
		}
		return result, nil
	case term.Assert:
		err := assertTypeIs(ctx, t.Annotation, Type, notAnEquivalence)
		if err != nil {
			return nil, err
		}
		oper, ok := Eval(t.Annotation).(oper)
		if !ok || oper.OpCode != term.EquivOp {
			return nil, mkTypeError(notAnEquivalence)
		}
		if !AlphaEquivalent(oper.L, oper.R) {
			return nil, mkTypeError(assertionFailed(Quote(oper.L), Quote(oper.R)))
		}
		return oper, nil
	}
	return nil, mkTypeError(unhandledTypeCase)
}

type typeError struct {
	ctx     context
	message typeMessage
}

func mkTypeError(message typeMessage) typeError {
	return typeError{message: message}
}

func (t typeError) Error() string {
	return t.message.String()
}

type typeMessage interface {
	String() string
}

type staticTypeMessage struct{ text string }
type oneArgTypeMessage struct {
	format string
	expr   term.Term
}
type twoArgTypeMessage struct {
	format string
	expr0  term.Term
	expr1  term.Term
}

func (m staticTypeMessage) String() string { return m.text }
func (m oneArgTypeMessage) String() string {
	return fmt.Sprintf(m.format, m.expr)
}
func (m twoArgTypeMessage) String() string {
	return fmt.Sprintf(m.format, m.expr0, m.expr1)
}

func unboundVariable(e term.Term) typeMessage {
	return oneArgTypeMessage{
		format: "Unbound variable: %v",
		expr:   e,
	}
}

func annotMismatch(annotation, actualType term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "Expression doesn't match annotation\n" +
			"\n" +
			"Expression of type %v was annotated %v",
		expr0: actualType,
		expr1: annotation,
	}
}

func wrongOperandType(expectedType, actualType term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "Expected %v but got %v",
		expr0:  expectedType,
		expr1:  actualType,
	}
}

func typeMismatch(expectedType, actualType term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "Wrong type of function argument\n" +
			"\n" +
			"expected %v but got %v",
		expr0: expectedType,
		expr1: actualType,
	}
}

func mismatchedListElements(firstType, nthType term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "List elements should all have the same type\n" +
			"\n" +
			"first element had type %v but there was an element of type %v",
		expr0: firstType,
		expr1: nthType,
	}
}

func mapTypeMismatch(inferred, annotated term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "❰toMap❱ result type doesn't match annotation\n" +
			"\n" +
			"map had type %v but was annotated %v",
		expr0: inferred,
		expr1: annotated,
	}
}

func invalidToMapType(expr term.Term) typeMessage {
	return oneArgTypeMessage{
		format: "An empty ❰toMap❱ was annotated with an invalid type\n" +
			"\n" +
			"%v",
		expr: expr,
	}
}

func handlerOutputTypeMismatch(type1, type2 term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "Handlers should have the same output type\n" +
			"\n" +
			"Saw handlers of types %v and %v",
		expr0: type1,
		expr1: type2,
	}
}

func handlerInputTypeMismatch(altType, inputType term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "Wrong handler input type\n" +
			"\n" +
			"Expected input type %v but saw %v",
		expr0: altType,
		expr1: inputType,
	}
}

func projectionTypeMismatch(firstType, secondType term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "Projection type mismatch\n" +
			"\n" +
			"tried to project a %v but the field had type %v",
		expr0: firstType,
		expr1: secondType,
	}
}

func assertionFailed(leftTerm, rightTerm term.Term) typeMessage {
	return twoArgTypeMessage{
		format: "Assertion failed\n" +
			"\n" +
			"%v is not equivalent to %v",
		expr0: leftTerm,
		expr1: rightTerm,
	}
}

func typeCheckVar(boundVar term.Term) typeMessage {
	return oneArgTypeMessage{
		format: "Unbound variable %s",
		expr:   boundVar,
	}
}

func cantBoolOp(opCode term.OpCode) typeMessage {
	var opStr string
	switch opCode {
	case term.OrOp:
		opStr = "||"
	case term.AndOp:
		opStr = "&&"
	case term.EqOp:
		opStr = "=="
	case term.NeOp:
		opStr = "!="
	default:
		panic(fmt.Sprintf("unknown boolean opcode %d", opCode))
	}
	return staticTypeMessage{fmt.Sprintf("❰%s❱ only works on ❰Bool❱s", opStr)}
}

func cantNaturalOp(opCode term.OpCode) typeMessage {
	var opStr string
	switch opCode {
	case term.PlusOp:
		opStr = "+"
	case term.TimesOp:
		opStr = "*"
	default:
		panic(fmt.Sprintf("unknown natural opcode %d", opCode))
	}
	return staticTypeMessage{fmt.Sprintf("❰%s❱ only works on ❰Natural❱s", opStr)}
}

func duplicateProjectedField(name string) typeMessage {
	return staticTypeMessage{fmt.Sprintf("Duplicate field ❰%s❱ in projection expression", name)}
}

var (
	ifBranchMismatch   = staticTypeMessage{"❰if❱ branches must have matching types"}
	ifBranchMustBeTerm = staticTypeMessage{"❰if❱ branch is not a term"}
	invalidFieldType   = staticTypeMessage{"Invalid field type"}
	invalidListType    = staticTypeMessage{"Invalid type for ❰List❱"}
	invalidInputType   = staticTypeMessage{"Invalid function input"}
	invalidOutputType  = staticTypeMessage{"Invalid function output"}
	invalidPredicate   = staticTypeMessage{"Invalid predicate for ❰if❱"}
	invalidSome        = staticTypeMessage{"❰Some❱ argument has the wrong type"}

	invalidAlternativeType        = staticTypeMessage{"Invalid alternative type"}
	alternativeAnnotationMismatch = staticTypeMessage{"Alternative annotation mismatch"}

	notAFunction = staticTypeMessage{"Not a function"}
	untyped      = staticTypeMessage{"❰Sort❱ has no type, kind, or sort"}

	incomparableExpression  = staticTypeMessage{"Incomparable expression"}
	equivalenceTypeMismatch = staticTypeMessage{"The two sides of the equivalence have different types"}

	invalidToMapRecordKind  = staticTypeMessage{"❰toMap❱ expects a record of kind ❰Type❱"}
	heterogenousRecordToMap = staticTypeMessage{"❰toMap❱ expects a homogenous record"}
	missingToMapType        = staticTypeMessage{"An empty ❰toMap❱ requires a type annotation"}

	mustMergeARecord      = staticTypeMessage{"❰merge❱ expects a record of handlers"}
	mustMergeUnion        = staticTypeMessage{"❰merge❱ expects a union or Optional"}
	missingMergeType      = staticTypeMessage{"An empty ❰merge❱ requires a type annotation"}
	unusedHandler         = staticTypeMessage{"Unused handler"}
	missingHandler        = staticTypeMessage{"Missing handler"}
	handlerNotAFunction   = staticTypeMessage{"Handler is not a function"}
	disallowedHandlerType = staticTypeMessage{"Disallowed handler type"}

	cantInterpolate = staticTypeMessage{"You can only interpolate ❰Text❱"}

	cantTextAppend     = staticTypeMessage{"❰++❱ only works on ❰Text❱"}
	cantListAppend     = staticTypeMessage{"❰#❱ only works on ❰List❱s"}
	listAppendMismatch = staticTypeMessage{"You can only append ❰List❱s with matching element types"}

	mustCombineARecord = staticTypeMessage{"You can only combine records"}

	combineTypesRequiresRecordType = staticTypeMessage{"❰⩓❱ requires arguments that are record types"}

	cantAccess              = staticTypeMessage{"Not a record or a union"}
	cantProject             = staticTypeMessage{"Not a record"}
	cantProjectByExpression = staticTypeMessage{"Selector is not a record type"}
	missingField            = staticTypeMessage{"Missing record field"}
	missingConstructor      = staticTypeMessage{"Missing constructor"}

	unhandledTypeCase = staticTypeMessage{"Internal error: unhandled case in TypeOf()"}

	notAnEquivalence = staticTypeMessage{"Not an equivalence"}
)
