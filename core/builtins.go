package core

import (
	"fmt"
)

func naturalEven(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 0)
	}
	return nil
}

func naturalFold(n, T, s, z Value) Value {
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

var (
	NaturalEvenVal      = naturalEvenVal{naturalEven}
	NaturalFoldVal      = naturalFoldVal{hasCall4{naturalFold}}
	NaturalIsZeroVal    = naturalIsZeroVal{naturalIsZero}
	NaturalOddVal       = naturalOddVal{naturalOdd}
	NaturalShowVal      = naturalShowVal{naturalShow}
	NaturalSubtractVal  = naturalSubtractVal{hasCall2{naturalSubtract}}
	NaturalToIntegerVal = naturalToIntegerVal{naturalToInteger}
	IntegerShowVal      = integerShowVal{integerShow}
	IntegerToDoubleVal  = integerToDoubleVal{integerToDouble}
	DoubleShowVal       = doubleShowVal{doubleShow}
)
