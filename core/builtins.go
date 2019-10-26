package core

import (
	"fmt"
	"strings"
)

func naturalBuild(x Value) Value {
	var succ Value = LambdaValue{
		Label:  "x",
		Domain: Natural,
		hasCall1: func(x Value) Value {
			if n, ok := x.(NaturalLit); ok {
				return NaturalLit(n + 1)
			}
			return OpValue{OpCode: PlusOp, L: x, R: NaturalLit(1)}
		},
	}
	if app, ok := x.(AppValue); ok {
		if _, ok := app.Fn.(naturalFoldVal); ok {
			return app.Arg
		}
	}
	return applyVal3(x, Natural, succ, NaturalLit(0))
}

func naturalEven(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 0)
	}
	return nil
}

func naturalFold(n, _, s, z Value) Value {
	if n, ok := n.(NaturalLit); ok {
		result := z
		for i := 0; i < int(n); i++ {
			if succ, ok := s.(Callable1); ok {
				result = succ.Call1(result)
			} else {
				result = AppValue{s, result}
			}
		}
		return result
	}
	return nil
}

func naturalIsZero(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n == 0)
	}
	return nil
}

func naturalOdd(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 1)
	}
	return nil
}

func naturalShow(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return TextLitVal{Suffix: fmt.Sprintf("%d", n)}
	}
	return nil
}

func naturalSubtract(a, b Value) Value {
	m, mok := a.(NaturalLit)
	n, nok := b.(NaturalLit)
	if mok && nok {
		if n >= m {
			return NaturalLit(n - m)
		}
		return NaturalLit(0)
	}
	if a == NaturalLit(0) {
		return b
	}
	if b == NaturalLit(0) {
		return NaturalLit(0)
	}
	if judgmentallyEqualVals(a, b) {
		return NaturalLit(0)
	}
	return nil
}

func naturalToInteger(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return IntegerLit(n)
	}
	return nil
}

func integerShow(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		return TextLitVal{Suffix: fmt.Sprintf("%+d", i)}
	}
	return nil
}

func integerToDouble(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		return DoubleLit(i)
	}
	return nil
}

func doubleShow(x Value) Value {
	if d, ok := x.(DoubleLit); ok {
		return TextLitVal{Suffix: d.String()}
	}
	return nil
}

func optionalBuild(A0, g Value) Value {
	var some Value = LambdaValue{
		Label:  "a",
		Domain: A0,
		hasCall1: func(a Value) Value {
			return SomeVal{a}
		},
	}
	if app, ok := g.(AppValue); ok {
		if app2, ok := app.Fn.(AppValue); ok {
			if _, ok := app2.Fn.(optionalFoldVal); ok {
				return app.Arg
			}
		}
	}
	return applyVal3(g, AppValue{Optional, A0}, some, AppValue{None, A0})
}

func optionalFold(_, opt, _, some, none Value) Value {
	if s, ok := opt.(SomeVal); ok {
		if some, ok := some.(Callable1); ok {
			return some.Call1(s.Val)
		}
		return AppValue{some, s.Val}
	}
	if app, ok := opt.(AppValue); ok {
		if app.Fn == None {
			return none
		}
	}
	return nil
}

func textShow(a0 Value) Value {
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

func listBuild(A0, g Value) Value {
	var cons Value = LambdaValue{
		Label:  "a",
		Domain: A0,
		hasCall1: func(a Value) Value {
			return LambdaValue{
				Label:  "as",
				Domain: AppValue{List, A0},
				hasCall1: func(as Value) Value {
					if _, ok := as.(EmptyListVal); ok {
						return NonEmptyListVal{a}
					}
					if as, ok := as.(NonEmptyListVal); ok {
						return append(NonEmptyListVal{a}, as...)
					}
					return OpValue{OpCode: ListAppendOp, L: NonEmptyListVal{a}, R: as}
				},
			}
		},
	}
	if app, ok := g.(AppValue); ok {
		if app2, ok := app.Fn.(AppValue); ok {
			if _, ok := app2.Fn.(listFoldVal); ok {
				return app.Arg
			}
		}
	}
	return applyVal3(g, AppValue{List, A0}, cons, EmptyListVal{AppValue{List, A0}})
}

func listFold(_, l, _, cons, empty Value) Value {
	if _, ok := l.(EmptyListVal); ok {
		return empty
	}
	if l, ok := l.(NonEmptyListVal); ok {
		result := empty
		for i := len(l) - 1; i >= 0; i-- {
			result = applyVal2(cons, l[i], result)
		}
		return result
	}
	return nil
}

func listLength(_, l Value) Value {
	if _, ok := l.(EmptyListVal); ok {
		return NaturalLit(0)
	}
	if l, ok := l.(NonEmptyListVal); ok {
		return NaturalLit(len(l))
	}
	return nil
}

func listHead(T, l Value) Value {
	if _, ok := l.(EmptyListVal); ok {
		return AppValue{None, T}
	}
	if l, ok := l.(NonEmptyListVal); ok {
		return SomeVal{l[0]}
	}
	return nil
}

func listLast(T, l Value) Value {
	if _, ok := l.(EmptyListVal); ok {
		return AppValue{None, T}
	}
	if l, ok := l.(NonEmptyListVal); ok {
		return SomeVal{l[len(l)-1]}
	}
	return nil
}

func listIndexed(T, l Value) Value {
	if _, ok := l.(EmptyListVal); ok {
		return EmptyListVal{AppValue{
			List,
			RecordTypeVal{"index": Natural, "value": T},
		}}
	}
	if l, ok := l.(NonEmptyListVal); ok {
		var result []Value
		for i, v := range l {
			result = append(result,
				RecordLitVal{"index": NaturalLit(i), "value": v})
		}
		return NonEmptyListVal(result)
	}
	return nil
}

func listReverse(T, l Value) Value {
	if _, ok := l.(EmptyListVal); ok {
		return l
	}
	if l, ok := l.(NonEmptyListVal); ok {
		result := make([]Value, len(l))
		for i, v := range l {
			result[len(l)-i-1] = v
		}
		return NonEmptyListVal(result)
	}
	return nil
}

var (
	NaturalBuildVal     = naturalBuildVal{naturalBuild}
	NaturalEvenVal      = naturalEvenVal{naturalEven}
	NaturalFoldVal      = naturalFoldVal{naturalFold}
	NaturalIsZeroVal    = naturalIsZeroVal{naturalIsZero}
	NaturalOddVal       = naturalOddVal{naturalOdd}
	NaturalShowVal      = naturalShowVal{naturalShow}
	NaturalSubtractVal  = naturalSubtractVal{naturalSubtract}
	NaturalToIntegerVal = naturalToIntegerVal{naturalToInteger}
	IntegerShowVal      = integerShowVal{integerShow}
	IntegerToDoubleVal  = integerToDoubleVal{integerToDouble}
	DoubleShowVal       = doubleShowVal{doubleShow}

	OptionalBuildVal = optionalBuildVal{optionalBuild}
	OptionalFoldVal  = optionalFoldVal{optionalFold}

	TextShowVal = textShowVal{textShow}

	ListBuildVal   = listBuildVal{listBuild}
	ListFoldVal    = listFoldVal{listFold}
	ListLengthVal  = listLengthVal{listLength}
	ListHeadVal    = listHeadVal{listHead}
	ListLastVal    = listLastVal{listLast}
	ListIndexedVal = listIndexedVal{listIndexed}
	ListReverseVal = listReverseVal{listReverse}

	_ Callable1 = NaturalBuildVal
	_ Callable1 = NaturalEvenVal
	_ Callable1 = NaturalIsZeroVal
	_ Callable1 = NaturalOddVal
	_ Callable1 = NaturalShowVal
	_ Callable1 = NaturalToIntegerVal
	_ Callable1 = IntegerShowVal
	_ Callable1 = IntegerToDoubleVal
	_ Callable1 = DoubleShowVal
	_ Callable1 = TextShowVal

	_ Callable2 = NaturalSubtractVal
	_ Callable2 = OptionalBuildVal

	_ Callable4 = NaturalFoldVal

	_ Callable5 = OptionalFoldVal
)
