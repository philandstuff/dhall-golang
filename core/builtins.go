package core

import (
	"fmt"
	"strings"
)

func (naturalBuildVal) Call(x Value) Value {
	var succ Value = lambdaValue{
		Label:  "x",
		Domain: Natural,
		Fn: func(x Value) Value {
			if n, ok := x.(NaturalLit); ok {
				return NaturalLit(n + 1)
			}
			return opValue{OpCode: PlusOp, L: x, R: NaturalLit(1)}
		},
	}
	return applyVal(x, Natural, succ, NaturalLit(0))
}

func (naturalBuildVal) ArgType() Value {
	return NewPiVal("natural", Type, func(natural Value) Value {
		return NewFnTypeVal("succ", NewFnTypeVal("_", natural, natural),
			NewFnTypeVal("zero", natural,
				natural))
	})
}

func (naturalEvenVal) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 0)
	}
	return nil
}

func (naturalEvenVal) ArgType() Value { return Natural }

func (fold naturalFoldVal) Call(x Value) Value {
	if fold.n == nil {
		return naturalFoldVal{n: x}
	}
	if fold.typ == nil {
		return naturalFoldVal{
			n:   fold.n,
			typ: x,
		}
	}
	if fold.succ == nil {
		return naturalFoldVal{
			n:    fold.n,
			typ:  fold.typ,
			succ: x,
		}
	}
	zero := x
	if n, ok := fold.n.(NaturalLit); ok {
		result := zero
		for i := 0; i < int(n); i++ {
			result = applyVal(fold.succ, result)
		}
		return result
	}
	return nil
}

func (fold naturalFoldVal) ArgType() Value {
	if fold.n == nil {
		return Natural
	}
	if fold.typ == nil {
		return Type
	}
	if fold.succ == nil {
		return NewFnTypeVal("_", fold.typ, fold.typ)
	}
	// zero
	return fold.typ
}

func (naturalIsZeroVal) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n == 0)
	}
	return nil
}

func (naturalIsZeroVal) ArgType() Value { return Natural }

func (naturalOddVal) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 1)
	}
	return nil
}

func (naturalOddVal) ArgType() Value { return Natural }

func (naturalShowVal) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return TextLitVal{Suffix: fmt.Sprintf("%d", n)}
	}
	return nil
}

func (naturalShowVal) ArgType() Value { return Natural }

func (sub naturalSubtractVal) Call(x Value) Value {
	if sub.a == nil {
		return naturalSubtractVal{a: x}
	}
	m, mok := sub.a.(NaturalLit)
	n, nok := x.(NaturalLit)
	if mok && nok {
		if n >= m {
			return NaturalLit(n - m)
		}
		return NaturalLit(0)
	}
	if sub.a == NaturalLit(0) {
		return x
	}
	if x == NaturalLit(0) {
		return NaturalLit(0)
	}
	if AlphaEquivalentVals(sub.a, x) {
		return NaturalLit(0)
	}
	return nil
}

func (naturalSubtractVal) ArgType() Value { return Natural }

func (naturalToIntegerVal) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return IntegerLit(n)
	}
	return nil
}

func (naturalToIntegerVal) ArgType() Value { return Natural }

func (integerClampVal) Call(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		if i < 0 {
			return NaturalLit(0)
		}
		return NaturalLit(i)
	}
	return nil
}

func (integerClampVal) ArgType() Value { return Integer }

func (integerNegateVal) Call(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		return IntegerLit(-i)
	}
	return nil
}

func (integerNegateVal) ArgType() Value { return Integer }

func (integerShowVal) Call(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		return TextLitVal{Suffix: fmt.Sprintf("%+d", i)}
	}
	return nil
}

func (integerShowVal) ArgType() Value { return Integer }

func (integerToDoubleVal) Call(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		return DoubleLit(i)
	}
	return nil
}

func (integerToDoubleVal) ArgType() Value { return Integer }

func (doubleShowVal) Call(x Value) Value {
	if d, ok := x.(DoubleLit); ok {
		return TextLitVal{Suffix: d.String()}
	}
	return nil
}

func (doubleShowVal) ArgType() Value { return Double }

func (build optionalBuildVal) Call(x Value) Value {
	if build.typ == nil {
		return optionalBuildVal{typ: x}
	}
	var some Value = lambdaValue{
		Label:  "a",
		Domain: build.typ,
		Fn: func(a Value) Value {
			return SomeVal{a}
		},
	}
	return applyVal(x, AppValue{Optional, build.typ}, some, AppValue{None, build.typ})
}

func (build optionalBuildVal) ArgType() Value {
	if build.typ == nil {
		return Type
	}
	return NewPiVal("optional", Type, func(optional Value) Value {
		return NewFnTypeVal("just", NewFnTypeVal("_", build.typ, optional),
			NewFnTypeVal("nothing", optional,
				optional))
	})
}

func (fold optionalFoldVal) Call(x Value) Value {
	if fold.typ1 == nil {
		return optionalFoldVal{typ1: x}
	}
	if fold.opt == nil {
		return optionalFoldVal{typ1: fold.typ1, opt: x}
	}
	if fold.typ2 == nil {
		return optionalFoldVal{
			typ1: fold.typ1,
			opt:  fold.opt,
			typ2: x,
		}
	}
	if fold.some == nil {
		return optionalFoldVal{
			typ1: fold.typ1,
			opt:  fold.opt,
			typ2: fold.typ2,
			some: x,
		}
	}
	none := x
	if s, ok := fold.opt.(SomeVal); ok {
		return applyVal(fold.some, s.Val)
	}
	if app, ok := fold.opt.(AppValue); ok {
		if app.Fn == None {
			return none
		}
	}
	return nil
}

func (fold optionalFoldVal) ArgType() Value {
	if fold.typ1 == nil {
		return Type
	}
	if fold.opt == nil {
		return applyVal(Optional, fold.typ1)
	}
	if fold.typ2 == nil {
		return Type
	}
	if fold.some == nil {
		return NewFnTypeVal("_", fold.typ1, fold.typ2)
	}
	// none
	return fold.typ2
}

func (textShowVal) Call(a0 Value) Value {
	if t, ok := a0.(TextLitVal); ok {
		if t.Chunks == nil || len(t.Chunks) == 0 {
			var out strings.Builder
			out.WriteRune('"')
			for _, r := range t.Suffix {
				switch r {
				case '"':
					out.WriteString(`\"`)
				case '$':
					out.WriteString(`\u0024`)
				case '\\':
					out.WriteString(`\\`)
				case '\b':
					out.WriteString(`\b`)
				case '\f':
					out.WriteString(`\f`)
				case '\n':
					out.WriteString(`\n`)
				case '\r':
					out.WriteString(`\r`)
				case '\t':
					out.WriteString(`\t`)
				default:
					if r < 0x1f {
						out.WriteString(fmt.Sprintf(`\u%04x`, r))
					} else {
						out.WriteRune(r)
					}
				}
			}
			out.WriteRune('"')
			return TextLitVal{Suffix: out.String()}
		}
	}
	return nil
}

func (textShowVal) ArgType() Value {
	return Text
}

func (l listBuildVal) Call(x Value) Value {
	if l.typ == nil {
		return listBuildVal{typ: x}
	}
	var cons Value = lambdaValue{
		Label:  "a",
		Domain: l.typ,
		Fn: func(a Value) Value {
			return lambdaValue{
				Label:  "as",
				Domain: AppValue{List, l.typ},
				Fn: func(as Value) Value {
					if _, ok := as.(EmptyListVal); ok {
						return NonEmptyListVal{a}
					}
					if as, ok := as.(NonEmptyListVal); ok {
						return append(NonEmptyListVal{a}, as...)
					}
					return opValue{OpCode: ListAppendOp, L: NonEmptyListVal{a}, R: as}
				},
			}
		},
	}
	return applyVal(x, AppValue{List, l.typ}, cons, EmptyListVal{AppValue{List, l.typ}})
}

func (l listBuildVal) ArgType() Value {
	if l.typ == nil {
		return Type
	}
	return NewPiVal("list", Type, func(list Value) Value {
		return NewFnTypeVal("cons", NewFnTypeVal("_", l.typ, NewFnTypeVal("_", list, list)),
			NewFnTypeVal("nil", list,
				list))
	})
}

func (l listFoldVal) Call(x Value) Value {
	if l.typ1 == nil {
		return listFoldVal{typ1: x}
	}
	if l.list == nil {
		return listFoldVal{typ1: l.typ1, list: x}
	}
	if l.typ2 == nil {
		return listFoldVal{
			typ1: l.typ1,
			list: l.list,
			typ2: x,
		}
	}
	if l.cons == nil {
		return listFoldVal{
			typ1: l.typ1,
			list: l.list,
			typ2: l.typ2,
			cons: x,
		}
	}
	empty := x
	if _, ok := l.list.(EmptyListVal); ok {
		return empty
	}
	if list, ok := l.list.(NonEmptyListVal); ok {
		result := empty
		for i := len(list) - 1; i >= 0; i-- {
			result = applyVal(l.cons, list[i], result)
		}
		return result
	}
	return nil
}

func (l listFoldVal) ArgType() Value {
	if l.typ1 == nil {
		return Type
	}
	if l.list == nil {
		return applyVal(List, l.typ1)
	}
	if l.typ2 == nil {
		return Type
	}
	if l.cons == nil {
		return NewFnTypeVal("_", l.typ1, NewFnTypeVal("_", l.typ2, l.typ2))
	}
	// nil
	return l.typ2
}

func (length listLengthVal) Call(x Value) Value {
	if length.typ == nil {
		return listLengthVal{typ: x}
	}
	if _, ok := x.(EmptyListVal); ok {
		return NaturalLit(0)
	}
	if l, ok := x.(NonEmptyListVal); ok {
		return NaturalLit(len(l))
	}
	return nil
}

func (length listLengthVal) ArgType() Value {
	if length.typ == nil {
		return Type
	}
	return applyVal(List, length.typ)
}

func (head listHeadVal) Call(x Value) Value {
	if head.typ == nil {
		return listHeadVal{typ: x}
	}
	if _, ok := x.(EmptyListVal); ok {
		return AppValue{None, head.typ}
	}
	if l, ok := x.(NonEmptyListVal); ok {
		return SomeVal{l[0]}
	}
	return nil
}

func (head listHeadVal) ArgType() Value {
	if head.typ == nil {
		return Type
	}
	return applyVal(List, head.typ)
}

func (last listLastVal) Call(x Value) Value {
	if last.typ == nil {
		return listLastVal{typ: x}
	}
	if _, ok := x.(EmptyListVal); ok {
		return AppValue{None, last.typ}
	}
	if l, ok := x.(NonEmptyListVal); ok {
		return SomeVal{l[len(l)-1]}
	}
	return nil
}

func (last listLastVal) ArgType() Value {
	if last.typ == nil {
		return Type
	}
	return applyVal(List, last.typ)
}

func (indexed listIndexedVal) Call(x Value) Value {
	if indexed.typ == nil {
		return listIndexedVal{typ: x}
	}
	if _, ok := x.(EmptyListVal); ok {
		return EmptyListVal{AppValue{
			List,
			RecordTypeVal{"index": Natural, "value": indexed.typ},
		}}
	}
	if l, ok := x.(NonEmptyListVal); ok {
		var result []Value
		for i, v := range l {
			result = append(result,
				RecordLitVal{"index": NaturalLit(i), "value": v})
		}
		return NonEmptyListVal(result)
	}
	return nil
}

func (indexed listIndexedVal) ArgType() Value {
	if indexed.typ == nil {
		return Type
	}
	return applyVal(List, indexed.typ)
}

func (rev listReverseVal) Call(x Value) Value {
	if rev.typ == nil {
		return listReverseVal{typ: x}
	}
	if _, ok := x.(EmptyListVal); ok {
		return x
	}
	if l, ok := x.(NonEmptyListVal); ok {
		result := make([]Value, len(l))
		for i, v := range l {
			result[len(l)-i-1] = v
		}
		return NonEmptyListVal(result)
	}
	return nil
}

func (rev listReverseVal) ArgType() Value {
	if rev.typ == nil {
		return Type
	}
	return applyVal(List, rev.typ)
}

var (
	NaturalBuildVal     Callable = naturalBuildVal{}
	NaturalEvenVal      Callable = naturalEvenVal{}
	NaturalFoldVal      Callable = naturalFoldVal{}
	NaturalIsZeroVal    Callable = naturalIsZeroVal{}
	NaturalOddVal       Callable = naturalOddVal{}
	NaturalShowVal      Callable = naturalShowVal{}
	NaturalSubtractVal  Callable = naturalSubtractVal{}
	NaturalToIntegerVal Callable = naturalToIntegerVal{}
	IntegerClampVal     Callable = integerClampVal{}
	IntegerNegateVal    Callable = integerNegateVal{}
	IntegerShowVal      Callable = integerShowVal{}
	IntegerToDoubleVal  Callable = integerToDoubleVal{}
	DoubleShowVal       Callable = doubleShowVal{}

	OptionalBuildVal Callable = optionalBuildVal{}
	OptionalFoldVal  Callable = optionalFoldVal{}

	TextShowVal Callable = textShowVal{}

	ListBuildVal   Callable = listBuildVal{}
	ListFoldVal    Callable = listFoldVal{}
	ListLengthVal  Callable = listLengthVal{}
	ListHeadVal    Callable = listHeadVal{}
	ListLastVal    Callable = listLastVal{}
	ListIndexedVal Callable = listIndexedVal{}
	ListReverseVal Callable = listReverseVal{}
)
