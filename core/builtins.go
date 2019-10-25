package core

import (
	"fmt"
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
	if g, ok := x.(Callable3); ok {
		if result := g.Call3(Natural, succ, NaturalLit(0)); result != nil {
			return result
		}
	}
	return applyVal(x, Natural, succ, NaturalLit(0))
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
	if g, ok := g.(Callable3); ok {
		if result := g.Call3(AppValue{Optional, A0}, some, AppValue{None, A0}); result != nil {
			return result
		}
	}
	return applyVal(g, AppValue{Optional, A0}, some, AppValue{None, A0})
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
)
