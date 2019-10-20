package eval

import (
	"errors"
	"fmt"

	. "github.com/philandstuff/dhall-golang/core"
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
		case List:
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
		bodyTypeVal := Eval(piType).(PiValue).Range(Eval(argType))
		return Quote(bodyTypeVal), nil
	case LambdaTerm:
		pi := PiTerm{Label: t.Label, Type: t.Type}
		freshLocal := LocalVar{Name: t.Label, Index: len(ctx[t.Label])}
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
		freshLocal := LocalVar{Name: t.Label, Index: len(ctx[t.Label])}
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
		return nil, errors.New("Let type unimplemented")
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
		return nil, errors.New("TextLitTerm type unimplemented")
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
		return nil, errors.New("OpTerm type unimplemented")
	case EmptyList:
		_, err := typeWith(ctx, t.Type)
		if err != nil {
			return nil, err
		}
		_, ok := Eval(t.Type).(AppValue)
		if !ok {
			return nil, mkTypeError(invalidListType)
		}
		return t.Type, nil
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
		return nil, errors.New("RecordType type unimplemented")
	case RecordLit:
		return nil, errors.New("RecordLit type unimplemented")
	case ToMap:
		return nil, errors.New("ToMap type unimplemented")
	case Field:
		return nil, errors.New("Field type unimplemented")
	case Project:
		return nil, errors.New("Project type unimplemented")
	case ProjectType:
		return nil, errors.New("ProjectType type unimplemented")
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

func typeCheckBoundVar(boundVar Term) typeMessage {
	return oneArgTypeMessage{
		format: "Internal error: shouldn't ever see BoundVar in TypeOf(), but saw %s",
		expr:   boundVar,
	}
}

var (
	ifBranchMismatch   = staticTypeMessage{"❰if❱ branches must have matching types"}
	ifBranchMustBeTerm = staticTypeMessage{"❰if❱ branch is not a term"}
	invalidListType    = staticTypeMessage{"Invalid type for ❰List❱"}
	invalidInputType   = staticTypeMessage{"Invalid function input"}
	invalidOutputType  = staticTypeMessage{"Invalid function output"}
	invalidPredicate   = staticTypeMessage{"Invalid predicate for ❰if❱"}
	invalidSome        = staticTypeMessage{"❰Some❱ argument has the wrong type"}
	notAFunction       = staticTypeMessage{"Not a function"}
	untyped            = staticTypeMessage{"❰Sort❱ has no type, kind, or sort"}

	unhandledTypeCase = staticTypeMessage{"Internal error: unhandled case in TypeOf()"}
)
