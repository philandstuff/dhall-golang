package core

import (
	"fmt"
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

func (ctx context) freshLocal(name string) localVar {
	return localVar{Name: name, Index: len(ctx[name])}
}

func assertTypeIs(ctx context, expr Term, expectedType Value, msg typeMessage) error {
	actualType, err := typeWith(ctx, expr)
	if err != nil {
		return err
	}
	if !AlphaEquivalentVals(expectedType, actualType) {
		return mkTypeError(msg)
	}
	return nil
}

// This returns
//  Value: the element type of a list type
//  Bool: whether it succeeded
func listElementType(e Value) (Value, bool) {
	app, ok := e.(AppValue)
	if !ok || app.Fn != List {
		return nil, false
	}
	return app.Arg, true
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

func TypeOf(t Term) (Value, error) {
	v, err := typeWith(context{}, t)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func typeWith(ctx context, t Term) (Value, error) {
	switch t := t.(type) {
	case Universe:
		switch t {
		case Type:
			return Kind, nil
		case Kind:
			return Sort, nil
		case Sort:
			return nil, mkTypeError(untyped)
		default:
			return nil, mkTypeError(unhandledTypeCase)
		}
	case Builtin:
		switch t {
		case Bool, Double, Integer, Natural, Text:
			return Type, nil
		case DoubleShow:
			return NewFnTypeVal("_", Double, Text), nil
		case IntegerClamp:
			return NewFnTypeVal("_", Integer, Natural), nil
		case IntegerNegate:
			return NewFnTypeVal("_", Integer, Integer), nil
		case IntegerShow:
			return NewFnTypeVal("_", Integer, Text), nil
		case IntegerToDouble:
			return NewFnTypeVal("_", Integer, Double), nil
		case List, Optional:
			return NewFnTypeVal("_", Type, Type), nil
		case ListBuild:
			return NewPiVal("a", Type, func(a Value) Value {
				return NewFnTypeVal("_",
					NewPiVal("list", Type, func(list Value) Value {
						return NewFnTypeVal("cons",
							NewFnTypeVal("_", a, NewFnTypeVal("_", list, list)),
							NewFnTypeVal("nil", list, list))
					}),
					AppValue{List, a})
			}), nil
		case ListFold:
			return NewPiVal("a", Type, func(a Value) Value {
				return NewFnTypeVal("_",
					AppValue{List, a},
					NewPiVal("list", Type, func(list Value) Value {
						return NewFnTypeVal("cons",
							NewFnTypeVal("_", a, NewFnTypeVal("_", list, list)),
							NewFnTypeVal("nil", list, list))
					}))
			}), nil
		case ListLength:
			return NewPiVal("a", Type, func(a Value) Value {
				return NewFnTypeVal("_", AppValue{List, a}, Natural)
			}), nil
		case ListHead, ListLast:
			return NewPiVal("a", Type, func(a Value) Value {
				return NewFnTypeVal("_", AppValue{List, a},
					AppValue{Optional, a})
			}), nil
		case ListReverse:
			return NewPiVal("a", Type, func(a Value) Value {
				return NewFnTypeVal("_", AppValue{List, a},
					AppValue{List, a})
			}), nil
		case ListIndexed:
			return NewPiVal("a", Type, func(a Value) Value {
				return NewFnTypeVal("_", AppValue{List, a},
					AppValue{List, RecordTypeVal{"index": Natural, "value": a}})
			}), nil
		case NaturalBuild:
			return NewFnTypeVal("_",
				NewPiVal("natural", Type, func(natural Value) Value {
					return NewFnTypeVal("succ",
						NewFnTypeVal("_", natural, natural),
						NewFnTypeVal("zero", natural, natural))
				}),
				Natural), nil
		case NaturalFold:
			return NewFnTypeVal("_",
				Natural,
				NewPiVal("natural", Type, func(natural Value) Value {
					return NewFnTypeVal("succ",
						NewFnTypeVal("_", natural, natural),
						NewFnTypeVal("zero", natural, natural))
				})), nil
		case NaturalIsZero, NaturalOdd, NaturalEven:
			return NewFnTypeVal("_", Natural, Bool), nil
		case NaturalShow:
			return NewFnTypeVal("_", Natural, Text), nil
		case NaturalToInteger:
			return NewFnTypeVal("_", Natural, Integer), nil
		case NaturalSubtract:
			return NewFnTypeVal("_", Natural, NewFnTypeVal("_", Natural, Natural)), nil
		case None:
			return NewPiVal("A", Type, func(A Value) Value { return AppValue{Optional, A} }), nil
		case OptionalBuild:
			return NewPiVal("a", Type, func(a Value) Value {
				return NewFnTypeVal("_",
					NewPiVal("optional", Type, func(optional Value) Value {
						return NewFnTypeVal("just",
							NewFnTypeVal("_", a, optional),
							NewFnTypeVal("nothing", optional, optional))
					}),
					AppValue{Optional, a})
			}), nil
		case OptionalFold:
			return NewPiVal("a", Type, func(a Value) Value {
				return NewFnTypeVal("_",
					AppValue{Optional, a},
					NewPiVal("optional", Type, func(optional Value) Value {
						return NewFnTypeVal("just",
							NewFnTypeVal("_", a, optional),
							NewFnTypeVal("nothing", optional, optional))
					}))
			}), nil
		case TextShow:
			return NewFnTypeVal("_", Text, Text), nil
		default:
			return nil, mkTypeError(unhandledTypeCase)
		}
	case Var:
		return nil, mkTypeError(typeCheckVar(t))
	case localVar:
		if vals, ok := ctx[t.Name]; ok {
			if t.Index < len(vals) {
				return vals[t.Index], nil
			}
			return nil, mkTypeError(unboundVariable(t))
		}
		return nil, fmt.Errorf("Unknown variable %s", t.Name)
	case AppTerm:
		fnType, err := typeWith(ctx, t.Fn)
		if err != nil {
			return nil, err
		}
		argType, err := typeWith(ctx, t.Arg)
		if err != nil {
			return nil, err
		}
		piType, ok := fnType.(PiValue)
		if !ok {
			return nil, mkTypeError(notAFunction)
		}
		expectedType := piType.Domain
		actualType := argType
		if !AlphaEquivalentVals(expectedType, actualType) {
			return nil, mkTypeError(typeMismatch(Quote(expectedType), Quote(actualType)))
		}
		bodyTypeVal := piType.Range(Eval(t.Arg))
		return bodyTypeVal, nil
	case LambdaTerm:
		_, err := typeWith(ctx, t.Type)
		if err != nil {
			return nil, err
		}
		argType := Eval(t.Type)
		pi := PiValue{Label: t.Label, Domain: argType}
		freshLocal := ctx.freshLocal(t.Label)
		bt, err := typeWith(
			ctx.extend(t.Label, argType),
			subst(t.Label, freshLocal, t.Body))
		if err != nil {
			return nil, err
		}
		pi.Range = func(x Value) Value {
			rebound := rebindLocal(freshLocal, Quote(bt))
			return evalWith(rebound, Env{
				t.Label: []Value{x},
			}, false)
		}
		_, err = typeWith(ctx, Quote(pi))
		if err != nil {
			return nil, err
		}
		return pi, nil
	case PiTerm:
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
			subst(t.Label, freshLocal, t.Body))
		if err != nil {
			return nil, err
		}
		o, ok := outUniv.(Universe)
		if !ok {
			return nil, mkTypeError(invalidOutputType)
		}
		return functionCheck(i, o), nil
	case NaturalLit:
		return Natural, nil
	case Let:
		let := t
		for len(let.Bindings) > 0 {
			binding := let.Bindings[0]
			let.Bindings = let.Bindings[1:]

			bindingType, err := typeWith(ctx, binding.Value)
			if err != nil {
				return nil, err
			}

			if binding.Annotation != nil {
				_, err := typeWith(ctx, binding.Value)
				if err != nil {
					return nil, err
				}
				if !AlphaEquivalentVals(bindingType, Eval(binding.Annotation)) {
					return nil, mkTypeError(annotMismatch(binding.Annotation, Quote(bindingType)))
				}
			}

			value := Quote(Eval(binding.Value))
			let = subst(binding.Variable, value, let).(Let)
			ctx = ctx.extend(binding.Variable, bindingType)
		}
		return typeWith(ctx, let.Body)
	case Annot:
		if t.Annotation != Sort {
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
		if !AlphaEquivalentVals(Eval(t.Annotation), actualType) {
			return nil, mkTypeError(annotMismatch(t.Annotation, Quote(actualType)))
		}
		// ─────────────────
		// Γ ⊢ (t : T₀) : T₀
		return actualType, nil
	case DoubleLit:
		return Double, nil
	case TextLitTerm:
		for _, chunk := range t.Chunks {
			err := assertTypeIs(ctx, chunk.Expr, Text,
				cantInterpolate)
			if err != nil {
				return nil, err
			}
		}
		return Text, nil
	case BoolLit:
		return Bool, nil
	case IfTerm:
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
		if !AlphaEquivalentVals(L, R) {
			return nil, mkTypeError(ifBranchMismatch)
		}
		return L, nil
	case IntegerLit:
		return Integer, nil
	case OpTerm:
		switch t.OpCode {
		case OrOp, AndOp, EqOp, NeOp:
			err := assertTypeIs(ctx, t.L, Bool, cantBoolOp(t.OpCode))
			if err != nil {
				return nil, err
			}
			err = assertTypeIs(ctx, t.R, Bool, cantBoolOp(t.OpCode))
			if err != nil {
				return nil, err
			}
			return Bool, nil
		case PlusOp, TimesOp:
			err := assertTypeIs(ctx, t.L, Natural, cantNaturalOp(t.OpCode))
			if err != nil {
				return nil, err
			}
			err = assertTypeIs(ctx, t.R, Natural, cantNaturalOp(t.OpCode))
			if err != nil {
				return nil, err
			}
			return Natural, nil
		case TextAppendOp:
			err := assertTypeIs(ctx, t.L, Text, cantTextAppend)
			if err != nil {
				return nil, err
			}
			err = assertTypeIs(ctx, t.R, Text, cantTextAppend)
			if err != nil {
				return nil, err
			}
			return Text, nil
		case ListAppendOp:
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
			if !AlphaEquivalentVals(lElemT, rElemT) {
				return nil, mkTypeError(listAppendMismatch)
			}
			return lt, nil
		case RecordMergeOp:
			lType, err := typeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rType, err := typeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}
			recordType := OpTerm{L: Quote(lType), R: Quote(rType), OpCode: RecordTypeMergeOp}
			if _, err = typeWith(ctx, recordType); err != nil {
				return nil, err
			}
			return Eval(recordType), nil
		case RecordTypeMergeOp:
			lKind, err := typeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rKind, err := typeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}
			lt, ok := Eval(t.L).(RecordTypeVal)
			if !ok {
				return nil, mkTypeError(combineTypesRequiresRecordType)
			}
			rt, ok := Eval(t.R).(RecordTypeVal)
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
		case RightBiasedRecordMergeOp:
			lType, err := typeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rType, err := typeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}
			lRecord, ok := lType.(RecordTypeVal)
			if !ok {
				return nil, mkTypeError(mustCombineARecord)
			}
			rRecord, ok := rType.(RecordTypeVal)
			if !ok {
				return nil, mkTypeError(mustCombineARecord)
			}
			result := RecordTypeVal{}
			for k, v := range lRecord {
				result[k] = v
			}
			for k, v := range rRecord {
				result[k] = v
			}
			return result, nil
		case ImportAltOp:
			return typeWith(ctx, t.L)
		case EquivOp:
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
			err = assertTypeIs(ctx, Quote(lType), Type, incomparableExpression)
			if err != nil {
				return nil, err
			}
			if !AlphaEquivalentVals(lType, rType) {
				return nil, mkTypeError(equivalenceTypeMismatch)
			}
			return Type, nil
		case CompleteOp:
			return typeWith(ctx,
				Annot{
					Expr: OpTerm{OpCode: RightBiasedRecordMergeOp,
						L: Field{Record: t.L, FieldName: "default"},
						R: t.R,
					},
					Annotation: Field{Record: t.L, FieldName: "Type"},
				},
			)
		default:
			return nil, fmt.Errorf("Internal error: unknown opcode %v", t.OpCode)
		}
	case EmptyList:
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
	case NonEmptyList:
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
			if !AlphaEquivalentVals(T0, T1) {
				return nil, mkTypeError(mismatchedListElements(Quote(T0), Quote(T1)))
			}
		}
		return AppValue{List, T0}, nil
	case Some:
		A, err := typeWith(ctx, t.Val)
		if err != nil {
			return nil, err
		}
		if err = assertTypeIs(ctx, Quote(A), Type, invalidSome); err != nil {
			return nil, err
		}
		return AppValue{Optional, A}, nil
	case RecordType:
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
	case RecordLit:
		recordType := RecordTypeVal{}
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
	case ToMap:
		recordTypeVal, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeVal.(RecordTypeVal)
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
			rt, ok := t.(RecordTypeVal)
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
				if !AlphaEquivalentVals(elemType, v) {
					return nil, mkTypeError(heterogenousRecordToMap)
				}
			}
		}
		if k, _ := typeWith(ctx, Quote(elemType)); k != Type {
			return nil, mkTypeError(invalidToMapRecordKind)
		}
		inferred := AppValue{List, RecordTypeVal{"mapKey": Text, "mapValue": elemType}}
		if t.Type == nil {
			return inferred, nil
		}
		if _, err = typeWith(ctx, t.Type); err != nil {
			return nil, err
		}
		annot := Eval(t.Type)
		if !AlphaEquivalentVals(inferred, annot) {
			return nil, mkTypeError(mapTypeMismatch(Quote(inferred), t.Type))
		}
		return inferred, nil
	case Field:
		recordTypeVal, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeVal.(RecordTypeVal)
		if ok {
			fieldType, ok := recordType[t.FieldName]
			if !ok {
				return nil, mkTypeError(missingField)
			}
			return fieldType, nil
		}
		unionTypeV := Eval(t.Record)
		unionType, ok := unionTypeV.(unionTypeVal)
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
		return PiValue{
			Label:  t.FieldName,
			Domain: alternativeType,
			Range:  func(Value) Value { return unionType },
		}, nil
	case Project:
		recordTypeVal, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeVal.(RecordTypeVal)
		if !ok {
			return nil, mkTypeError(cantProject)
		}
		result := make(RecordTypeVal, len(t.FieldNames))
		for _, name := range t.FieldNames {
			var ok bool
			result[name], ok = recordType[name]
			if !ok {
				return nil, mkTypeError(missingField)
			}
		}
		return result, nil
	case ProjectType:
		recordTypeVal, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeVal.(RecordTypeVal)
		if !ok {
			return nil, mkTypeError(cantProject)
		}
		_, err = typeWith(ctx, t.Selector)
		if err != nil {
			return nil, err
		}
		selectorVal := Eval(t.Selector)
		selector, ok := selectorVal.(RecordTypeVal)
		if !ok {
			return nil, mkTypeError(cantProjectByExpression)
		}
		result := make(RecordTypeVal, len(selector))
		for name, typ := range selector {
			fieldType, ok := recordType[name]
			if !ok {
				return nil, mkTypeError(missingField)
			}
			if !AlphaEquivalentVals(fieldType, typ) {
				return nil, mkTypeError(projectionTypeMismatch(Quote(typ), Quote(fieldType)))
			}
			result[name] = typ
		}
		return result, nil
	case UnionType:
		if len(t) == 0 {
			return Type, nil
		}
		var c Universe
		first := true
		for _, typ := range t {
			if typ == nil {
				// empty alternative
				continue
			}
			k, err := typeWith(ctx, typ)
			if err != nil {
				return nil, err
			}
			if first {
				var ok bool
				c, ok = k.(Universe)
				if !ok {
					return nil, mkTypeError(invalidAlternativeType)
				}
			} else {
				if !AlphaEquivalentVals(c, k) {
					return nil, mkTypeError(alternativeAnnotationMismatch)
				}
			}
			if c == Sort {
				if Eval(typ) != Kind {
					return nil, mkTypeError(invalidAlternativeType)
				}
			}
			first = false
		}
		return c, nil
	case Merge:
		handlerTypeVal, err := typeWith(ctx, t.Handler)
		if err != nil {
			return nil, err
		}
		unionTypeV, err := typeWith(ctx, t.Union)
		if err != nil {
			return nil, err
		}
		handlerType, ok := handlerTypeVal.(RecordTypeVal)
		if !ok {
			return nil, mkTypeError(mustMergeARecord)
		}
		unionType, ok := unionTypeV.(unionTypeVal)
		if !ok {
			apply, ok := unionTypeV.(AppValue)
			if !ok || apply.Fn != Optional {
				return nil, mkTypeError(mustMergeUnion)
			}
			unionType = unionTypeVal{"Some": apply.Arg, "None": nil}
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
					if !AlphaEquivalentVals(result, fieldType) {
						return nil, mkTypeError(handlerOutputTypeMismatch(Quote(result), Quote(fieldType)))
					}
				}
			} else {
				pi, ok := fieldType.(PiValue)
				if !ok {
					return nil, mkTypeError(handlerNotAFunction)
				}
				if !AlphaEquivalentVals(altType, pi.Domain) {
					return nil, mkTypeError(handlerInputTypeMismatch(Quote(altType), Quote(pi.Domain)))
				}
				outputType := pi.Range(NaturalLit(1))
				outputType2 := pi.Range(NaturalLit(2))
				if !AlphaEquivalentVals(outputType, outputType2) {
					// hacky way of detecting output type depending on input
					return nil, mkTypeError(disallowedHandlerType)
				}
				if result == nil {
					result = outputType
				} else {
					if !AlphaEquivalentVals(result, outputType) {
						return nil, mkTypeError(handlerOutputTypeMismatch(Quote(result), Quote(outputType)))
					}
				}
			}
		}
		if t.Annotation != nil {
			if _, err := typeWith(ctx, t.Annotation); err != nil {
				return nil, err
			}
			if !AlphaEquivalentVals(result, Eval(t.Annotation)) {
				return nil, mkTypeError(annotMismatch(t.Annotation, Quote(result)))
			}
		}
		return result, nil
	case Assert:
		err := assertTypeIs(ctx, t.Annotation, Type, notAnEquivalence)
		if err != nil {
			return nil, err
		}
		op, ok := Eval(t.Annotation).(opValue)
		if !ok || op.OpCode != EquivOp {
			return nil, mkTypeError(notAnEquivalence)
		}
		if !AlphaEquivalentVals(op.L, op.R) {
			return nil, mkTypeError(assertionFailed(Quote(op.L), Quote(op.R)))
		}
		return op, nil
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
	expr   Term
}
type twoArgTypeMessage struct {
	format string
	expr0  Term
	expr1  Term
}

func (m staticTypeMessage) String() string { return m.text }
func (m oneArgTypeMessage) String() string {
	return fmt.Sprintf(m.format, m.expr)
}
func (m twoArgTypeMessage) String() string {
	return fmt.Sprintf(m.format, m.expr0, m.expr1)
}

func unboundVariable(e Term) typeMessage {
	return oneArgTypeMessage{
		format: "Unbound variable: %v",
		expr:   e,
	}
}

func annotMismatch(annotation, actualType Term) typeMessage {
	return twoArgTypeMessage{
		format: "Expression doesn't match annotation\n" +
			"\n" +
			"Expression of type %v was annotated %v",
		expr0: actualType,
		expr1: annotation,
	}
}

func wrongOperandType(expectedType, actualType Term) typeMessage {
	return twoArgTypeMessage{
		format: "Expected %v but got %v",
		expr0:  expectedType,
		expr1:  actualType,
	}
}

func typeMismatch(expectedType, actualType Term) typeMessage {
	return twoArgTypeMessage{
		format: "Wrong type of function argument\n" +
			"\n" +
			"expected %v but got %v",
		expr0: expectedType,
		expr1: actualType,
	}
}

func mismatchedListElements(firstType, nthType Term) typeMessage {
	return twoArgTypeMessage{
		format: "List elements should all have the same type\n" +
			"\n" +
			"first element had type %v but there was an element of type %v",
		expr0: firstType,
		expr1: nthType,
	}
}

func mapTypeMismatch(inferred, annotated Term) typeMessage {
	return twoArgTypeMessage{
		format: "❰toMap❱ result type doesn't match annotation\n" +
			"\n" +
			"map had type %v but was annotated %v",
		expr0: inferred,
		expr1: annotated,
	}
}

func invalidToMapType(expr Term) typeMessage {
	return oneArgTypeMessage{
		format: "An empty ❰toMap❱ was annotated with an invalid type\n" +
			"\n" +
			"%v",
		expr: expr,
	}
}

func handlerOutputTypeMismatch(type1, type2 Term) typeMessage {
	return twoArgTypeMessage{
		format: "Handlers should have the same output type\n" +
			"\n" +
			"Saw handlers of types %v and %v",
		expr0: type1,
		expr1: type2,
	}
}

func handlerInputTypeMismatch(altType, inputType Term) typeMessage {
	return twoArgTypeMessage{
		format: "Wrong handler input type\n" +
			"\n" +
			"Expected input type %v but saw %v",
		expr0: altType,
		expr1: inputType,
	}
}

func projectionTypeMismatch(firstType, secondType Term) typeMessage {
	return twoArgTypeMessage{
		format: "Projection type mismatch\n" +
			"\n" +
			"tried to project a %v but the field had type %v",
		expr0: firstType,
		expr1: secondType,
	}
}

func assertionFailed(leftTerm, rightTerm Term) typeMessage {
	return twoArgTypeMessage{
		format: "Assertion failed\n" +
			"\n" +
			"%v is not equivalent to %v",
		expr0: leftTerm,
		expr1: rightTerm,
	}
}

func typeCheckVar(boundVar Term) typeMessage {
	return oneArgTypeMessage{
		format: "Unbound variable %s",
		expr:   boundVar,
	}
}

func cantBoolOp(opCode int) typeMessage {
	var opStr string
	switch opCode {
	case OrOp:
		opStr = "||"
	case AndOp:
		opStr = "&&"
	case EqOp:
		opStr = "=="
	case NeOp:
		opStr = "!="
	default:
		panic(fmt.Sprintf("unknown boolean opcode %d", opCode))
	}
	return staticTypeMessage{fmt.Sprintf("❰%s❱ only works on ❰Bool❱s", opStr)}
}

func cantNaturalOp(opCode int) typeMessage {
	var opStr string
	switch opCode {
	case PlusOp:
		opStr = "+"
	case TimesOp:
		opStr = "*"
	default:
		panic(fmt.Sprintf("unknown natural opcode %d", opCode))
	}
	return staticTypeMessage{fmt.Sprintf("❰%s❱ only works on ❰Natural❱s", opStr)}
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
