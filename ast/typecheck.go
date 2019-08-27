package ast

import (
	"errors"
	"fmt"
	"reflect"
)

func judgmentallyEqual(e1 Expr, e2 Expr) bool {
	ne1 := e1.Normalize().AlphaNormalize()
	ne2 := e2.Normalize().AlphaNormalize()
	return reflect.DeepEqual(ne1, ne2)
}

// assert that a type is exactly expectedType (no judgmentallyEqual
// here)
func assertSimpleType(ctx *TypeContext, expr, expectedType Expr) error {
	actualType, err := expr.TypeWith(ctx)
	if err != nil {
		return err
	}
	actualType = actualType.Normalize()
	if actualType != expectedType {
		return fmt.Errorf("Expecting a %v, got %v", expectedType, actualType)
	}
	return nil
}

func NormalizedTypeWith(e Expr, ctx *TypeContext) (Expr, error) {
	t, err := e.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	return t.Normalize(), nil
}

func assertNormalizedTypeIs(expr Expr, ctx *TypeContext, expectedType Expr, msg error) error {
	t, err := NormalizedTypeWith(expr, ctx)
	if err != nil {
		return err
	}
	if t != expectedType {
		return msg
	}
	return nil
}

func (c Const) TypeWith(ctx *TypeContext) (Expr, error) {
	switch c {
	case Type:
		return Kind, nil
	case Kind:
		return Sort, nil
	case Sort:
		return nil, errors.New("Sort has no type")
	default:
		return nil, fmt.Errorf("Unknown Const %d", c)
	}
}

// common variable names in types
var (
	a        = MkVar("a")
	A        = MkVar("A")
	list     = MkVar("list")
	natural  = MkVar("natural")
	optional = MkVar("optional")
)

func (b Builtin) TypeWith(ctx *TypeContext) (Expr, error) {
	switch b {
	case Double, Text, Bool, Natural, Integer:
		return Type, nil
	case List, Optional:
		return FnType(Type, Type), nil
	case None:
		return &Pi{"A", Type, Apply(Optional, A)}, nil
	case NaturalBuild:
		return FnType(&Pi{"natural", Type,
			&Pi{"succ", FnType(natural, natural),
				&Pi{"zero", natural,
					natural}}},
			Natural), nil
	case NaturalFold:
		return FnType(Natural,
			&Pi{"natural", Type,
				&Pi{"succ", FnType(natural, natural),
					&Pi{"zero", natural,
						natural}}}), nil
	case NaturalIsZero, NaturalEven, NaturalOdd:
		return FnType(Natural, Bool), nil
	case NaturalToInteger:
		return FnType(Natural, Integer), nil
	case NaturalShow:
		return FnType(Natural, Text), nil
	case NaturalSubtract:
		return FnType(Natural, FnType(Natural, Natural)), nil
	case IntegerToDouble:
		return FnType(Integer, Double), nil
	case IntegerShow:
		return FnType(Integer, Text), nil
	case DoubleShow:
		return FnType(Double, Text), nil
	case TextShow:
		return FnType(Text, Text), nil
	case ListBuild:
		return &Pi{"a", Type,
			FnType(
				&Pi{"list", Type,
					&Pi{"cons", FnType(a, FnType(list, list)),
						&Pi{"nil", list, list}}},
				Apply(List, a),
			)}, nil
	case ListFold:
		list := MkVar("list")
		return &Pi{"a", Type,
			FnType(Apply(List, a),
				&Pi{"list", Type,
					&Pi{"cons", FnType(a, FnType(list, list)),
						&Pi{"nil", list,
							list}}})}, nil
	case ListLength:
		return &Pi{"a", Type,
			FnType(Apply(List, a), Natural),
		}, nil
	case ListHead, ListLast:
		return &Pi{"a", Type,
			FnType(
				Apply(List, a),
				Apply(Optional, a),
			)}, nil
	case ListIndexed:
		return &Pi{"a", Type, FnType(
			Apply(List, a),
			Apply(List, Record{"index": Natural, "value": a}),
		)}, nil
	case ListReverse:
		return &Pi{"a", Type,
			FnType(
				Apply(List, a),
				Apply(List, a),
			)}, nil
	case OptionalBuild:
		return &Pi{"a", Type,
			FnType(
				&Pi{"optional", Type,
					&Pi{"just", FnType(a, optional),
						&Pi{"nothing", optional,
							optional}}},
				Apply(Optional, a))}, nil
	case OptionalFold:
		return &Pi{"a", Type,
			FnType(Apply(Optional, a),
				&Pi{"optional", Type,
					&Pi{"just", FnType(a, optional),
						&Pi{"nothing", optional,
							optional}}},
			)}, nil
	default:
		return nil, fmt.Errorf("Unknown Builtin %s", b)
	}
}

func (v Var) TypeWith(ctx *TypeContext) (Expr, error) {
	if t, ok := ctx.Lookup(v.Name, v.Index); ok {
		return t, nil
	}
	return nil, fmt.Errorf("Unbound variable %s, context was %+v", v.Name, ctx)
}

// Γ₀ ⊢ λ(x : A) → b : ∀(x : A) → B
func (lam *LambdaExpr) TypeWith(ctx *TypeContext) (Expr, error) {
	if _, err := lam.Type.TypeWith(ctx); err != nil {
		return nil, err
	}
	argType := lam.Type.Normalize()
	newctx := ctx.Insert(lam.Label, argType).Map(func(e Expr) Expr { return Shift(1, Var{Name: lam.Label}, e) })
	bodyType, err := lam.Body.TypeWith(newctx)
	if err != nil {
		return nil, err
	}

	p := &Pi{Label: lam.Label, Type: argType, Body: bodyType}
	_, err2 := p.TypeWith(ctx)
	if err2 != nil {
		return nil, err2
	}

	return p, nil
}

func Rule(a Const, b Const) (Const, error) {
	if b == Type {
		return Type, nil
	}
	if a > b {
		return a, nil
	} else {
		return b, nil
	}
}

// Γ₀ ⊢ ∀(x : A) → B : c
func (pi *Pi) TypeWith(ctx *TypeContext) (Expr, error) {
	// Γ₀ ⊢ A :⇥ i
	tA, err := NormalizedTypeWith(pi.Type, ctx)
	if err != nil {
		return nil, err
	}
	kA, ok := tA.(Const)
	if !ok {
		return nil, fmt.Errorf("Expected %v to be a Const", tA)
	}
	newctx := ctx.Insert(pi.Label, pi.Type.Normalize()).Map(func(e Expr) Expr { return Shift(1, Var{Name: pi.Label}, e) })
	tB, err := NormalizedTypeWith(pi.Body, newctx)
	if err != nil {
		return nil, err
	}
	kB, ok := tB.(Const)
	if !ok {
		return nil, errors.New("Wrong kind for body of pi type")
	}

	return Rule(kA, kB)
}

// Γ ⊢ f a₀ : B₂
func (app *App) TypeWith(ctx *TypeContext) (Expr, error) {
	// Γ ⊢ f :⇥ ∀(x : A₀) → B₀
	tF, err := NormalizedTypeWith(app.Fn, ctx)
	if err != nil {
		return nil, err
	}
	pF, ok := tF.(*Pi)
	if !ok {
		return nil, fmt.Errorf("Expected %s to be a function type", tF)
	}

	// Γ ⊢ a₀ : A₁
	A1, err := app.Arg.TypeWith(ctx)
	if err != nil {
		return nil, err
	}

	// A₀ ≡ A₁
	if judgmentallyEqual(pF.Type, A1) {
		// ↑(1, x, 0, a₀) = a₁
		a1 := Shift(1, Var{Name: pF.Label}, app.Arg)
		// B₀[x ≔ a₁] = B₁
		B1 := Subst(Var{Name: pF.Label}, a1, pF.Body)
		// ↑(-1, x, 0, B₁) = B₂
		return Shift(-1, Var{Name: pF.Label}, B1), nil
	} else {
		return nil, fmt.Errorf("type mismatch between function and applied value: `%v` `%v`", pF, A1)
	}
}

func (l Let) TypeWith(ctx *TypeContext) (Expr, error) {
	binding := l.Bindings[0]
	x := binding.Variable
	// Γ ⊢ a₀ : A₁
	valueType, err := binding.Value.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	// a₀ ⇥ a₁
	a1 := binding.Value.Normalize()
	if binding.Annotation != nil {
		// Γ ⊢ A₀ : i
		_, err := binding.Annotation.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		// A₀ ≡ A₁
		if !judgmentallyEqual(binding.Annotation, valueType) {
			return nil, fmt.Errorf("in let binding, was expecting type %v but saw %v", binding.Annotation, valueType)
		}
	}

	// TODO: optimization where binding.Value is a term

	// ↑(1, x, 0, a₁) = a₂
	a2 := Shift(1, Var{x, 0}, a1)

	rest := l.Body
	if len(l.Bindings) > 1 {
		rest = Let{Bindings: l.Bindings[1:], Body: l.Body}
	}

	// b₀[x ≔ a₂] = b₁
	b1 := Subst(Var{x, 0}, a2, rest)
	// ↑(-1, x, 0, b₁) = b₂
	b2 := Shift(-1, Var{x, 0}, b1)
	// Γ ⊢ b₂ : B
	return b2.TypeWith(ctx)
}

func (a Annot) TypeWith(ctx *TypeContext) (Expr, error) {
	if a.Annotation == Sort {
		// Γ ⊢ t : Sort
		err := assertNormalizedTypeIs(a.Expr, ctx, Sort,
			fmt.Errorf("Expected %v to have type Sort", a.Expr))
		if err != nil {
			return nil, err
		}
		// ─────────────────────
		// Γ ⊢ (t : Sort) : Sort
		return Sort, nil
	}

	// Γ ⊢ T₀ : i
	_, err := a.Annotation.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	// Γ ⊢ t : T₁
	t2, err := a.Expr.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	// T₀ ≡ T₁
	if !judgmentallyEqual(a.Annotation, t2) {
		return nil, fmt.Errorf("Annotation mismatch: inferred type %v but annotated %v", t2, a.Annotation)
	}
	// ─────────────────
	// Γ ⊢ (t : T₀) : T₀
	return t2, nil
}

func (DoubleLit) TypeWith(*TypeContext) (Expr, error) { return Double, nil }

func (t TextLit) TypeWith(ctx *TypeContext) (Expr, error) {
	for _, chunk := range t.Chunks {
		err := assertNormalizedTypeIs(chunk.Expr, ctx, Text,
			errors.New("Interpolated expression is not Text"))
		if err != nil {
			return nil, err
		}
	}
	return Text, nil
}

func (BoolLit) TypeWith(*TypeContext) (Expr, error) { return Bool, nil }

func (b BoolIf) TypeWith(ctx *TypeContext) (Expr, error) {
	// Γ ⊢ t :⇥ Bool
	err := assertNormalizedTypeIs(b.Cond, ctx, Bool,
		errors.New("if condition must be type Bool"))
	if err != nil {
		return nil, err
	}
	// Γ ⊢ l : L
	tType, err := b.T.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	// Γ ⊢ L :⇥ Type
	err = assertNormalizedTypeIs(tType, ctx, Type, errors.New("if branches must have terms, not types"))
	if err != nil {
		return nil, err
	}
	// Γ ⊢ r : R
	fType, err := b.F.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	// Γ ⊢ R :⇥ Type
	err = assertNormalizedTypeIs(fType, ctx, Type, errors.New("if branches must have terms, not types"))
	if err != nil {
		return nil, err
	}
	// L ≡ R
	if !judgmentallyEqual(tType, fType) {
		return nil, errors.New("if branches must have matching types")
	}
	// ──────────────────────────
	// Γ ⊢ if t then l else r : L
	return tType, nil
}

func (NaturalLit) TypeWith(*TypeContext) (Expr, error) { return Natural, nil }

func mergeRecords(l Record, r Record) (Record, error) {
	var err error
	output := make(Record)
	for k, v := range l {
		output[k] = v
	}
	for k, v := range r {
		if lField, ok := output[k]; ok {
			lSubrecord, Lok := lField.(Record)
			rSubrecord, Rok := v.(Record)
			if !(Lok && Rok) {
				return nil, errors.New("Record mismatch")
			}
			output[k], err = mergeRecords(lSubrecord, rSubrecord)
			if err != nil {
				return nil, err
			}
		} else {
			output[k] = v
		}
	}
	return output, nil
}

func (op Operator) TypeWith(ctx *TypeContext) (Expr, error) {
	switch op.OpCode {
	case OrOp, AndOp, EqOp, NeOp:
		err := assertSimpleType(ctx, op.L, Bool)
		if err != nil {
			return nil, err
		}
		err = assertSimpleType(ctx, op.R, Bool)
		if err != nil {
			return nil, err
		}
		return Bool, nil
	case PlusOp, TimesOp:
		err := assertSimpleType(ctx, op.L, Natural)
		if err != nil {
			return nil, err
		}
		err = assertSimpleType(ctx, op.R, Natural)
		if err != nil {
			return nil, err
		}
		return Natural, nil
	case TextAppendOp:
		err := assertSimpleType(ctx, op.L, Text)
		if err != nil {
			return nil, err
		}
		err = assertSimpleType(ctx, op.R, Text)
		if err != nil {
			return nil, err
		}
		return Text, nil
	case ListAppendOp:
		lt, err := NormalizedTypeWith(op.L, ctx)
		if err != nil {
			return nil, err
		}
		rt, err := NormalizedTypeWith(op.R, ctx)
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
	case RecordMergeOp:
		lt, err := NormalizedTypeWith(op.L, ctx)
		if err != nil {
			return nil, err
		}
		ltr, ok := lt.(Record)
		if !ok {
			return nil, fmt.Errorf("The ∧ operator operates on records, not %v", lt)
		}
		_, err = NormalizedTypeWith(lt, ctx)
		if err != nil {
			return nil, err
		}
		rt, err := NormalizedTypeWith(op.R, ctx)
		if err != nil {
			return nil, err
		}
		rtr, ok := rt.(Record)
		if !ok {
			return nil, fmt.Errorf("The ∧ operator operates on records, not %v", rt)
		}
		_, err = NormalizedTypeWith(rt, ctx)
		if err != nil {
			return nil, err
		}
		return mergeRecords(ltr, rtr)
	case RightBiasedRecordMergeOp:
		lt, err := NormalizedTypeWith(op.L, ctx)
		if err != nil {
			return nil, err
		}
		ltr, ok := lt.(Record)
		if !ok {
			return nil, fmt.Errorf("The ⫽ operator operates on records, not %v", lt)
		}
		_, err = NormalizedTypeWith(lt, ctx)
		if err != nil {
			return nil, err
		}
		rt, err := NormalizedTypeWith(op.R, ctx)
		if err != nil {
			return nil, err
		}
		rtr, ok := rt.(Record)
		if !ok {
			return nil, fmt.Errorf("The ⫽ operator operates on records, not %v", rt)
		}
		_, err = NormalizedTypeWith(rt, ctx)
		if err != nil {
			return nil, err
		}
		output := make(Record)
		for k, v := range ltr {
			output[k] = v
		}
		for k, v := range rtr {
			output[k] = v
		}
		return output, nil
	case RecordTypeMergeOp:
		lt, err := NormalizedTypeWith(op.L, ctx)
		if err != nil {
			return nil, err
		}
		rt, err := NormalizedTypeWith(op.R, ctx)
		if err != nil {
			return nil, err
		}

		lr, ok := op.L.Normalize().(Record)
		if !ok {
			return nil, fmt.Errorf("The ⩓ operator operates on records, not %v", lt)
		}
		rr, ok := op.R.Normalize().(Record)
		if !ok {
			return nil, fmt.Errorf("The ⩓ operator operates on records, not %v", rt)
		}
		// ensure that the records are safe to merge
		_, err = mergeRecords(lr, rr)
		if err != nil {
			return nil, err
		}
		// if lr and rr are Records, then lt and rt must be Consts
		if lt.(Const) > rt.(Const) {
			return lt, nil
		} else {
			return rt, nil
		}
	case EquivOp:
		A0, err := op.L.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		A1, err := op.R.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		err = assertNormalizedTypeIs(A0, ctx, Type, errors.New("The ≡ operator compares terms"))
		if err != nil {
			return nil, err
		}
		err = assertNormalizedTypeIs(A1, ctx, Type, errors.New("The ≡ operator compares terms"))
		if err != nil {
			return nil, err
		}
		if !judgmentallyEqual(A0, A1) {
			return nil, errors.New("`x === y` requires x and y to be the same type")
		}
		return Type, nil
	}
	return nil, fmt.Errorf("Unknown opcode in %+v", op)
}

func (IntegerLit) TypeWith(*TypeContext) (Expr, error) { return Integer, nil }

func (l EmptyList) TypeWith(ctx *TypeContext) (Expr, error) {
	// Γ ⊢ T₀ : c
	_, err := NormalizedTypeWith(l.Type, ctx)
	if err != nil {
		return nil, err
	}
	lt := l.Type.Normalize()
	// T₀ ⇥ List T₁
	if app, ok := lt.(*App); ok {
		if app.Fn == List {
			// T₁ :⇥ Type
			err := assertSimpleType(ctx, app.Arg, Type)
			if err != nil {
				return nil, err
			}
			return app, nil
		}
	}
	return nil, fmt.Errorf("An empty list must be annotated with a List type, but I saw %v", l.Type)
}

func (l NonEmptyList) TypeWith(ctx *TypeContext) (Expr, error) {
	exprs := []Expr(l)
	t, err := exprs[0].TypeWith(ctx)
	if err != nil {
		return nil, err
	}

	err = assertSimpleType(ctx, t, Type)
	if err != nil {
		return nil, err
	}
	for _, elem := range exprs[1:] {
		t2, err := elem.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		if !judgmentallyEqual(t, t2) {
			return nil, fmt.Errorf("All List elements must have same type, but types %s and %s don't match", t, t2)
		}
	}
	return Apply(List, t), nil
}

// This returns
//  Expr: the element type of a list type
//  Bool: whether it succeeded
func listElementType(e Expr) (Expr, bool) {
	app, ok := e.(*App)
	if !ok {
		return nil, false
	}
	if app.Fn == List {
		return app.Arg, true
	}
	return nil, false
}

func (s Some) TypeWith(ctx *TypeContext) (Expr, error) {
	// Γ ⊢ a : A
	A, err := s.Val.TypeWith(ctx)
	if err != nil {
		return nil, err
	}

	// Γ ⊢ A : Type
	err = assertNormalizedTypeIs(A, ctx, Type,
		fmt.Errorf("Some must take a term, not a type like %v", s.Val))
	if err != nil {
		return nil, err
	}
	// ───────────────────────
	// Γ ⊢ Some a : Optional A
	return Apply(Optional, A), nil
}

func (r Record) TypeWith(ctx *TypeContext) (Expr, error) {
	fields := map[string]Expr(r)
	if len(fields) == 0 {
		return Type, nil
	}
	c := Type
	for _, typ := range fields {
		k, err := typ.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		kk, ok := k.(Const)
		if !ok {
			return nil, errors.New("Invalid field type")
		}
		if kk > c {
			c = kk
		}
	}
	return c, nil
}

func (r RecordLit) TypeWith(ctx *TypeContext) (Expr, error) {
	fields := map[string]Expr(r)
	if len(fields) == 0 {
		return Record(fields), nil
	}
	fieldTypes := make(map[string]Expr, len(fields))
	for name, val := range fields {
		typ, err := val.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		_, err = typ.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		fieldTypes[name] = typ
	}
	return Record(fieldTypes), nil
}

func (t ToMap) TypeWith(ctx *TypeContext) (Expr, error) {
	rt1, err := NormalizedTypeWith(t.Record, ctx)
	if err != nil {
		return nil, err
	}
	rt2, ok := rt1.(Record)
	if !ok {
		return nil, errors.New("`toMap` must be called on a record")
	}
	if len(rt2) != 0 {
		// nonempty record
		// Γ ⊢ e :⇥ { x : T₀, xs… }
		var elemType Expr
		for _, v := range rt2 {
			if elemType == nil {
				elemType = v
			} else {
				if !judgmentallyEqual(elemType, v) {
					return nil, fmt.Errorf("fields must be the same type, but saw different types %v and %v", elemType, v)
				}
			}
		}
		inferred := Apply(List, Record(map[string]Expr{"mapKey": Text, "mapValue": elemType}))
		if t.Type == nil {
			return inferred, nil
		} else {
			// Γ ⊢ toMap e : T₀   T₀ ≡ T₁
			if !judgmentallyEqual(inferred, t.Type) {
				return nil, errors.New("inferred type of toMap didn't match annotation")
			}
			// therefore
			// Γ ⊢ ( toMap e : T₁ ) : T₁
			return t.Type, nil
		}
	} else {
		// empty record
		// Γ ⊢ e :⇥ {}
		if t.Type == nil {
			return nil, errors.New("`toMap` on an empty record must have an annotation")
		}
		// Γ ⊢ T₀ :⇥ Type
		err = assertNormalizedTypeIs(t.Type, ctx, Type,
			errors.New("in `toMap {} : T`, T must be a Type"))
		if err != nil {
			return nil, err
		}

		// T₀ ⇥ List { mapKey : Text, mapValue : T₁ }
		normType := t.Type.Normalize()
		app, ok := normType.(*App)
		if !ok || app.Fn != List {
			return nil, errors.New("toMap must be annotated with `List { mapKey : Text, mapValue : T }`")
		}
		entryType, ok := app.Arg.(Record)
		if !ok || len(entryType) != 2 || entryType["mapKey"] != Text || entryType["mapValue"] == nil {
			return nil, errors.New("toMap must be annotated with `List { mapKey : Text, mapValue : T }`")
		}
		return app, nil
	}
}

func (f Field) TypeWith(ctx *TypeContext) (Expr, error) {
	// Γ ⊢ u : c (used in the union case below)
	rt, err := NormalizedTypeWith(f.Record, ctx)
	if err != nil {
		return nil, err
	}

	rtRecord, ok := rt.(Record)
	if ok {
		// Γ ⊢ e :⇥ { x : T, xs… }
		ft, ok := rtRecord[f.FieldName]
		if !ok {
			return nil, fmt.Errorf("Tried to access nonexistent field %s", f.FieldName)
		}
		// ───────────
		// Γ ⊢ e.x : T
		return ft, nil
	}

	record := f.Record.Normalize()
	unionType, ok := record.(UnionType)
	if !ok {
		return nil, fmt.Errorf("Tried to access field from non-record non-union type: %s", record)
	}
	alternativeType, ok := unionType[f.FieldName]
	if !ok {
		return nil, fmt.Errorf("Tried to access nonexistent union alternative %s", f.FieldName)
	}
	if alternativeType == nil {
		// u ⇥ < x | ts… >
		// ─────────────────────
		// Γ ⊢ u.x : < x | ts… >
		return unionType, nil
	}
	// u ⇥ < x : T | ts… >
	// ────────────────────────────────────
	// Γ ⊢ u.x : ∀(x : T) → < x : T | ts… >
	return &Pi{f.FieldName, alternativeType, unionType}, nil
}

func (p Project) TypeWith(ctx *TypeContext) (Expr, error) {
	rtWrapped, err := p.Record.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	rt, ok := rtWrapped.(Record)
	if !ok {
		return nil, fmt.Errorf("Can't project from a %v", rtWrapped)
	}
	rtt, err := rt.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	result := make(Record, len(p.FieldNames))
	if len(p.FieldNames) == 0 {
		if rtt != Type {
			return nil, fmt.Errorf("Can't project empty list from a %v, a record of %v", rt, rtt)
		}
		return result, nil
	}
	for _, name := range p.FieldNames {
		var ok bool
		result[name], ok = rt[name]
		if !ok {
			return nil, fmt.Errorf("Tried to project nonexistent field %s", name)
		}
	}
	return result, nil
}

func (p ProjectType) TypeWith(ctx *TypeContext) (Expr, error) {
	rtWrapped, err := p.Record.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	rt, ok := rtWrapped.(Record)
	if !ok {
		return nil, fmt.Errorf("Can't project from a %v", rtWrapped)
	}
	_, err = p.Selector.TypeWith(ctx)
	if err != nil {
		return nil, err
	}
	selectorWrapped := p.Selector.Normalize()
	selector, ok := selectorWrapped.(Record)
	if !ok {
		return nil, fmt.Errorf("Can't use a %v to project from a record", rtWrapped)
	}
	result := make(Record, len(selector))
	for name, typ := range selector {
		origType, ok := rt[name]
		if !ok {
			return nil, fmt.Errorf("Tried to project nonexistent field %s", name)
		}
		if !judgmentallyEqual(typ, origType) {
			return nil, fmt.Errorf("While projecting field %s, selector type %v didn't match record type %v", name, typ, origType)
		}
		result[name] = typ
	}
	return result, nil
}

func (u UnionType) TypeWith(ctx *TypeContext) (Expr, error) {
	if len(u) == 0 {
		return Type, nil
	}
	var c Const
	first := true
	for _, typ := range u {
		if typ == nil {
			// empty alternative
			continue
		}
		k, err := typ.TypeWith(ctx)
		if err != nil {
			return nil, err
		}
		if first {
			var ok bool
			c, ok = k.(Const)
			if !ok {
				return nil, errors.New("Invalid alternative type")
			}
		} else {
			if c.Normalize() != k.Normalize() {
				return nil, fmt.Errorf("can't mix %s and %s", c, k)
			}
		}
		if c == Sort {
			if typ.Normalize() != Kind {
				return nil, errors.New("Invalid alternative type")
			}
		}
		first = false
	}
	return c, nil
}

func (m Merge) TypeWith(ctx *TypeContext) (e Expr, outerr error) {
	ht, err := NormalizedTypeWith(m.Handler, ctx)
	if err != nil {
		return nil, err
	}
	ut, err := NormalizedTypeWith(m.Union, ctx)
	if err != nil {
		return nil, err
	}
	handlerRecord, ok := ht.(Record)
	if !ok {
		return nil, errors.New("merge handler must be a record")
	}
	unionType, ok := ut.(UnionType)
	if !ok {
		return nil, fmt.Errorf("merge arg must be of union type, but inferred type of %s was %s", m.Union, ut)
	}

	if len(handlerRecord) != len(unionType) {
		return nil, errors.New("handler record fields must match union type alternatives")
	}

	// Γ ⊢ t :⇥ {}   Γ ⊢ u :⇥ <>
	if len(handlerRecord) == 0 {
		if m.Annotation == nil {
			return nil, errors.New("empty merge requires an annotation")
		}
		// Γ ⊢ T :⇥ Type
		err := assertNormalizedTypeIs(m.Annotation, ctx, Type,
			errors.New("in `merge {=} <> : T`, T must be a Type"))
		if err != nil {
			return nil, err
		}
		// ───────────────────────
		// Γ ⊢ (merge t u : T) : T
		return m.Annotation, nil
	}
	var outputType Expr
	for altName, altType := range unionType {
		field, ok := handlerRecord[altName]
		if !ok {
			return nil, errors.New("handler record fields must match union type alternatives")
		}
		if altType == nil {
			// Γ ⊢ t :⇥ { y : T₀, ts… }   Γ ⊢ u :⇥ < y | us… >
			if outputType == nil {
				outputType = field
			} else {
				if !judgmentallyEqual(outputType, field) {
					return nil, fmt.Errorf("all handlers must output the same type, but %v was not the same as %v", outputType, field)
				}
			}
		} else {
			// Γ ⊢ t :⇥ { y : ∀(x : A₀) → T₀, ts… }
			// Γ ⊢ u :⇥ < y : A₁ | us… >
			pi, ok := field.(*Pi)
			if !ok {
				return nil, errors.New("handler must be a function")
			}
			// A₀ ≡ A₁
			if !judgmentallyEqual(altType, pi.Type) {
				return nil, fmt.Errorf("Handler for %s should take argument of type %s, not %s", altName, altType, pi.Type)
			}
			// ; `x` not free in `T₀`
			if IsFreeIn(pi.Body, altName) {
				// we maybe don't need this block, except for better
				// error reporting
				return nil, fmt.Errorf("Variable %s used in function type body %s", altName, pi.Body)
			}
			// the following call to Shift() may panic -- see
			// dhall-lang/tests/typecheck/unit/MergeHandlerFreeVar for an
			// example of when this may happen
			defer func() {
				if r := recover(); r == "tried to shift to negative" {
					outerr = fmt.Errorf("Unbound variable %s", pi.Label)
				}
			}()
			thisOutputType := Shift(-1, Var{Name: pi.Label, Index: 0}, pi.Body)
			if outputType == nil {
				outputType = thisOutputType
			} else {
				if !judgmentallyEqual(outputType, thisOutputType) {
					return nil, fmt.Errorf("all handlers must output the same type, but %v was not the same as %v", outputType, thisOutputType)
				}
			}
		}
	}

	if m.Annotation != nil && !judgmentallyEqual(outputType, m.Annotation) {
		return nil, fmt.Errorf("Expression was annotated as type %s but inferred type was %s", m.Annotation, outputType)
	}
	return outputType, nil
}

func (a Assert) TypeWith(ctx *TypeContext) (Expr, error) {
	err := assertNormalizedTypeIs(a.Annotation, ctx, Type, errors.New("assert must be annotated with a Type"))
	if err != nil {
		return nil, err
	}
	annot := a.Annotation.Normalize()
	equiv, ok := annot.(Operator)
	if !ok || equiv.OpCode != EquivOp {
		return nil, errors.New("assert must be annotated with an equivalence")
	}
	if !judgmentallyEqual(equiv.L, equiv.R) {
		return nil, fmt.Errorf("assertion failed: %v ≢ %v", equiv.L, equiv.R)
	}
	return annot, nil
}

func (e Embed) TypeWith(ctx *TypeContext) (Expr, error) {
	return nil, errors.New("Cannot typecheck an expression with unresolved imports")
}
