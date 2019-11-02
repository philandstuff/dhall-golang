package core

import (
	"errors"
	"fmt"
)

type context map[string][]Term

func (ctx context) extend(name string, t Term) context {
	newctx := context{}
	for k, v := range ctx {
		newctx[k] = v
	}
	newctx[name] = append(newctx[name], t)
	return newctx
}

func (ctx context) freshLocal(name string) LocalVar {
	return LocalVar{Name: name, Index: len(ctx[name])}
}

// assert that a type is exactly expectedType (no judgmentallyEqual
// here)
func assertSimpleType(ctx context, expr, expectedType Term) error {
	actualType, err := normalizedTypeWith(ctx, expr)
	if err != nil {
		return err
	}
	if actualType != expectedType {
		return mkTypeError(wrongOperandType(expectedType, actualType))
	}
	return nil
}

func normalizedTypeWith(ctx context, t Term) (Term, error) {
	typ, err := typeWith(ctx, t)
	if err != nil {
		return nil, err
	}
	return Quote(Eval(typ)), nil
}

func assertNormalizedTypeIs(ctx context, expr Term, expectedType Term, msg error) error {
	t, err := normalizedTypeWith(ctx, expr)
	if err != nil {
		return err
	}
	if t != expectedType {
		return msg
	}
	return nil
}

// This returns
//  Term: the element type of a list type
//  Bool: whether it succeeded
func listElementType(e Term) (Term, bool) {
	app, ok := e.(AppTerm)
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

func TypeOf(t Term) (Term, error) {
	return typeWith(context{}, t)
}

func typeWith(ctx context, t Term) (Term, error) {
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
			return FnType(Double, Text), nil
		case IntegerShow:
			return FnType(Integer, Text), nil
		case IntegerToDouble:
			return FnType(Integer, Double), nil
		case List, Optional:
			return FnType(Type, Type), nil
		case ListBuild:
			return MkΠ("a", Type,
				FnType(MkΠ("list", Type,
					MkΠ("cons", FnType(Bound("a"), FnType(Bound("list"), Bound("list"))),
						MkΠ("nil", Bound("list"),
							Bound("list")))),
					Apply(List, Bound("a")))), nil
		case ListFold:
			return MkΠ("a", Type,
				FnType(Apply(List, Bound("a")),
					MkΠ("list", Type,
						MkΠ("cons", FnType(Bound("a"), FnType(Bound("list"), Bound("list"))),
							MkΠ("nil", Bound("list"),
								Bound("list")))))), nil
		case ListLength:
			return MkΠ("a", Type,
				FnType(Apply(List, Bound("a")),
					Natural)), nil
		case ListHead, ListLast:
			return MkΠ("a", Type,
				FnType(Apply(List, Bound("a")),
					Apply(Optional, Bound("a")))), nil
		case ListReverse:
			return MkΠ("a", Type,
				FnType(Apply(List, Bound("a")),
					Apply(List, Bound("a")))), nil
		case ListIndexed:
			return MkΠ("a", Type,
				FnType(Apply(List, Bound("a")),
					Apply(List, RecordType{"index": Natural, "value": Bound("a")}))), nil
		case NaturalBuild:
			return FnType(MkΠ("natural", Type,
				MkΠ("succ", FnType(Bound("natural"), Bound("natural")),
					MkΠ("zero", Bound("natural"),
						Bound("natural")))),
				Natural), nil
		case NaturalFold:
			return FnType(
				Natural,
				MkΠ("natural", Type,
					MkΠ("succ", FnType(Bound("natural"), Bound("natural")),
						MkΠ("zero", Bound("natural"),
							Bound("natural"))))), nil
		case NaturalIsZero, NaturalOdd, NaturalEven:
			return FnType(Natural, Bool), nil
		case NaturalShow:
			return FnType(Natural, Text), nil
		case NaturalToInteger:
			return FnType(Natural, Integer), nil
		case NaturalSubtract:
			return FnType(Natural, FnType(Natural, Natural)), nil
		case None:
			return MkΠ("A", Type, Apply(Optional, Bound("A"))), nil
		case OptionalBuild:
			return MkΠ("a", Type,
				FnType(MkΠ("optional", Type,
					MkΠ("just", FnType(Bound("a"), Bound("optional")),
						MkΠ("nothing", Bound("optional"),
							Bound("optional")))),
					Apply(Optional, Bound("a")))), nil
		case OptionalFold:
			return MkΠ("a", Type,
				FnType(Apply(Optional, Bound("a")),
					MkΠ("optional", Type,
						MkΠ("just", FnType(Bound("a"), Bound("optional")),
							MkΠ("nothing", Bound("optional"),
								Bound("optional")))))), nil
		case TextShow:
			return FnType(Text, Text), nil
		default:
			return nil, mkTypeError(unhandledTypeCase)
		}
	case BoundVar:
		return nil, mkTypeError(typeCheckBoundVar(t))
	case LocalVar:
		if vals, ok := ctx[t.Name]; ok {
			if t.Index < len(vals) {
				return vals[t.Index], nil
			}
			return nil, mkTypeError(unboundVariable(t))
		}
		return nil, fmt.Errorf("Unknown variable %s", t.Name)
	case FreeVar:
		return nil, errors.New("typecheck freevar unimp")
	case AppTerm:
		fnType, err := typeWith(ctx, t.Fn)
		if err != nil {
			return nil, err
		}
		argType, err := typeWith(ctx, t.Arg)
		if err != nil {
			return nil, err
		}
		piType, ok := fnType.(PiTerm)
		if !ok {
			return nil, mkTypeError(notAFunction)
		}
		expectedType := piType.Type
		actualType := argType
		if !judgmentallyEqual(expectedType, actualType) {
			return nil, mkTypeError(typeMismatch(expectedType, actualType))
		}
		bodyTypeVal := Eval(piType).(PiValue).Range(Eval(t.Arg))
		return Quote(bodyTypeVal), nil
	case LambdaTerm:
		pi := PiTerm{Label: t.Label, Type: t.Type}
		freshLocal := ctx.freshLocal(t.Label)
		bt, err := typeWith(
			ctx.extend(t.Label, t.Type),
			subst(t.Label, freshLocal, t.Body))
		if err != nil {
			return nil, err
		}
		pi.Body = rebindLocal(freshLocal, bt)
		_, err = typeWith(ctx, pi)
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
			ctx.extend(t.Label, t.Type),
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
		body := t.Body
		binding := t.Bindings[0]
		rest := t.Bindings[1:]
		for {
			value := Quote(Eval(binding.Value))
			body = subst(binding.Variable, value, body)
			for i, b := range rest {
				rest[i].Value = subst(binding.Variable, value, b.Value)
				if b.Annotation != nil {
					rest[i].Annotation = subst(binding.Variable, value, b.Annotation)
				}
			}
			bindingType, err := typeWith(ctx, binding.Value)
			if err != nil {
				return nil, err
			}
			if binding.Annotation != nil && !judgmentallyEqual(bindingType, binding.Annotation) {
				return nil, mkTypeError(annotMismatch(binding.Annotation, bindingType))
			}
			ctx = ctx.extend(binding.Variable, bindingType)
			if len(rest) == 0 {
				break
			}
			binding = rest[0]
			rest = rest[1:]
		}
		return typeWith(ctx, body)
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
		if !judgmentallyEqual(t.Annotation, actualType) {
			return nil, mkTypeError(annotMismatch(t.Annotation, actualType))
		}
		// ─────────────────
		// Γ ⊢ (t : T₀) : T₀
		return t.Annotation, nil
	case DoubleLit:
		return Double, nil
	case TextLitTerm:
		for _, chunk := range t.Chunks {
			err := assertNormalizedTypeIs(ctx, chunk.Expr, Text,
				errors.New("Interpolated expression is not Text"))
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
		if t, _ := typeWith(ctx, L); t != Type {
			return nil, mkTypeError(ifBranchMustBeTerm)
		}
		R, err := typeWith(ctx, t.F)
		if err != nil {
			return nil, err
		}
		if t, _ := typeWith(ctx, R); t != Type {
			return nil, mkTypeError(ifBranchMustBeTerm)
		}
		if !judgmentallyEqual(L, R) {
			return nil, mkTypeError(ifBranchMismatch)
		}
		return L, nil
	case IntegerLit:
		return Integer, nil
	case OpTerm:
		switch t.OpCode {
		case OrOp, AndOp, EqOp, NeOp:
			err := assertSimpleType(ctx, t.L, Bool)
			if err != nil {
				return nil, err
			}
			err = assertSimpleType(ctx, t.R, Bool)
			if err != nil {
				return nil, err
			}
			return Bool, nil
		case PlusOp, TimesOp:
			err := assertSimpleType(ctx, t.L, Natural)
			if err != nil {
				return nil, err
			}
			err = assertSimpleType(ctx, t.R, Natural)
			if err != nil {
				return nil, err
			}
			return Natural, nil
		case TextAppendOp:
			err := assertSimpleType(ctx, t.L, Text)
			if err != nil {
				return nil, err
			}
			err = assertSimpleType(ctx, t.R, Text)
			if err != nil {
				return nil, err
			}
			return Text, nil
		case ListAppendOp:
			lt, err := normalizedTypeWith(ctx, t.L)
			if err != nil {
				return nil, err
			}
			rt, err := normalizedTypeWith(ctx, t.R)
			if err != nil {
				return nil, err
			}

			lElemT, ok := listElementType(lt)
			if !ok {
				return nil, fmt.Errorf("Can't use list concatenate operator on a %s", lt)
			}
			rElemT, ok := listElementType(rt)
			if !ok {
				return nil, fmt.Errorf("Can't use list concatenate operator on a %s", rt)
			}
			if !judgmentallyEqual(lElemT, rElemT) {
				return nil, fmt.Errorf("Can't append a %s to a %s", lt, rt)
			}
			return lt, nil
		default:
			return nil, fmt.Errorf("Internal error: unknown opcode %v", t.OpCode)
		}
	case EmptyList:
		_, err := typeWith(ctx, t.Type)
		if err != nil {
			return nil, err
		}
		listType := Eval(t.Type)
		app, ok := listType.(AppValue)
		if !ok {
			return nil, mkTypeError(invalidListType)
		}
		if app.Fn != List {
			return nil, mkTypeError(invalidListType)
		}
		return Quote(listType), nil
	case NonEmptyList:
		T0, err := typeWith(ctx, t[0])
		if err != nil {
			return nil, err
		}
		T0type, err := typeWith(ctx, T0)
		if err != nil {
			return nil, err
		}
		if T0type != Type {
			return nil, mkTypeError(invalidListType)
		}
		for _, e := range t[1:] {
			T1, err := typeWith(ctx, e)
			if err != nil {
				return nil, err
			}
			if !judgmentallyEqual(T0, T1) {
				return nil, mkTypeError(mismatchedListElements(T0, T1))
			}
		}
		return Apply(List, T0), nil
	case Some:
		A, err := typeWith(ctx, t.Val)
		if err != nil {
			return nil, err
		}
		Atype, err := typeWith(ctx, A)
		if err != nil {
			return nil, err
		}
		if Atype != Type {
			return nil, mkTypeError(invalidSome)
		}
		return Apply(Optional, A), nil
	case RecordType:
		recordUniverse := Type
		for _, v := range t {
			fieldUniverse, err := normalizedTypeWith(ctx, v)
			if err != nil {
				return nil, err
			}
			if u, ok := fieldUniverse.(Universe); ok {
				if recordUniverse < u {
					recordUniverse = u
				}
			} else {
				return nil, mkTypeError(invalidFieldType)
			}
		}
		return recordUniverse, nil
	case RecordLit:
		recordType := RecordType{}
		for k, v := range t {
			fieldType, err := normalizedTypeWith(ctx, v)
			if err != nil {
				return nil, err
			}
			recordType[k] = fieldType
		}
		if _, err := typeWith(ctx, recordType); err != nil {
			return nil, err
		}
		return recordType, nil
	case ToMap:
		recordTypeTerm, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeTerm.(RecordType)
		if !ok {
			return nil, mkTypeError(cantAccess)
		}

		if len(recordType) == 0 {
			if t.Type == nil {
				return nil, mkTypeError(missingToMapType)
			}
			tt, err := typeWith(ctx, t.Type)
			if err != nil {
				return nil, err
			}
			if tt != Type {
				return nil, mkTypeError(invalidToMapRecordKind)
			}
			tTerm := Quote(Eval(t.Type))
			t, ok := tTerm.(AppTerm)
			if !ok || t.Fn != List {
				return nil, mkTypeError(invalidToMapType(tTerm))
			}
			rt, ok := t.Arg.(RecordType)
			if !ok || len(rt) != 2 || rt["mapKey"] != Text || rt["mapValue"] == nil {
				return nil, mkTypeError(invalidToMapType(tTerm))
			}
			return t, nil
		}

		var elemType Term
		for _, v := range recordType {
			if elemType == nil {
				elemType = v
			} else {
				if !judgmentallyEqual(elemType, v) {
					return nil, mkTypeError(heterogenousRecordToMap)
				}
			}
		}
		inferred := AppTerm{List, RecordType{"mapKey": Text, "mapValue": elemType}}
		if t.Type == nil {
			return inferred, nil
		}
		if !judgmentallyEqual(inferred, t.Type) {
			return nil, mkTypeError(mapTypeMismatch(inferred, t.Type))
		}
		return t.Type, nil
	case Field:
		recordTypeTerm, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeTerm.(RecordType)
		if !ok {
			return nil, mkTypeError(cantAccess)
		}
		fieldType, ok := recordType[t.FieldName]
		if !ok {
			return nil, mkTypeError(missingField)
		}
		return fieldType, nil
	case Project:
		recordTypeTerm, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeTerm.(RecordType)
		if !ok {
			return nil, mkTypeError(cantProject)
		}
		result := make(RecordType, len(t.FieldNames))
		for _, name := range t.FieldNames {
			var ok bool
			result[name], ok = recordType[name]
			if !ok {
				return nil, mkTypeError(missingField)
			}
		}
		return result, nil
	case ProjectType:
		recordTypeTerm, err := typeWith(ctx, t.Record)
		if err != nil {
			return nil, err
		}
		recordType, ok := recordTypeTerm.(RecordType)
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
		result := make(RecordType, len(selector))
		for name, typ := range Quote(selector).(RecordType) {
			fieldType, ok := recordType[name]
			if !ok {
				return nil, mkTypeError(missingField)
			}
			if !judgmentallyEqual(fieldType, typ) {
				return nil, mkTypeError(projectionTypeMismatch(typ, fieldType))
			}
			result[name] = typ
		}
		return result, nil
	case UnionType:
		return nil, errors.New("UnionType type unimplemented")
	case Merge:
		return nil, errors.New("Merge type unimplemented")
	case Assert:
		return nil, errors.New("Assert type unimplemented")
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

func annotMismatch(expectedType, actualType Term) typeMessage {
	return twoArgTypeMessage{
		format: "Expression doesn't match annotation\n" +
			"\n" +
			"expected %v but got %v",
		expr0: expectedType,
		expr1: actualType,
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
		format: "❰toMap❱ result type doesn't match annotation" +
			"\n" +
			"map had type %v but was annotated %v",
		expr0: inferred,
		expr1: annotated,
	}
}

func invalidToMapType(expr Term) typeMessage {
	return oneArgTypeMessage{
		format: "An empty ❰toMap❱ was annotated with an invalid type" +
			"\n" +
			"%v",
		expr: expr,
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

func typeCheckBoundVar(boundVar Term) typeMessage {
	return oneArgTypeMessage{
		format: "Internal error: shouldn't ever see BoundVar in TypeOf(), but saw %s",
		expr:   boundVar,
	}
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
	notAFunction       = staticTypeMessage{"Not a function"}
	untyped            = staticTypeMessage{"❰Sort❱ has no type, kind, or sort"}

	invalidToMapRecordKind  = staticTypeMessage{"❰toMap❱ expects a record of kind ❰Type❱"}
	heterogenousRecordToMap = staticTypeMessage{"❰toMap❱ expects a homogenous record"}
	missingToMapType        = staticTypeMessage{"An empty ❰toMap❱ requires a type annotation"}

	cantAccess              = staticTypeMessage{"Not a record or a union"}
	cantProject             = staticTypeMessage{"Not a record"}
	cantProjectByExpression = staticTypeMessage{"Selector is not a record type"}
	missingField            = staticTypeMessage{"Missing record field"}

	unhandledTypeCase = staticTypeMessage{"Internal error: unhandled case in TypeOf()"}
)
