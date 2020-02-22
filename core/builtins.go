package core

import (
	"fmt"
	"strings"

	"github.com/philandstuff/dhall-golang/term"
)

func (naturalBuild) Call(x Value) Value {
	var succ Value = lambda{
		Label:  "x",
		Domain: Natural,
		Fn: func(x Value) Value {
			if n, ok := x.(NaturalLit); ok {
				return NaturalLit(n + 1)
			}
			return oper{OpCode: term.PlusOp, L: x, R: NaturalLit(1)}
		},
	}
	return apply(x, Natural, succ, NaturalLit(0))
}

func (naturalBuild) ArgType() Value {
	return NewPi("natural", Type, func(natural Value) Value {
		return NewFnType("succ", NewFnType("_", natural, natural),
			NewFnType("zero", natural,
				natural))
	})
}

func (naturalEven) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 0)
	}
	return nil
}

func (naturalEven) ArgType() Value { return Natural }

func (fold naturalFold) Call(x Value) Value {
	if fold.n == nil {
		return naturalFold{n: x}
	}
	if fold.typ == nil {
		return naturalFold{
			n:   fold.n,
			typ: x,
		}
	}
	if fold.succ == nil {
		return naturalFold{
			n:    fold.n,
			typ:  fold.typ,
			succ: x,
		}
	}
	zero := x
	if n, ok := fold.n.(NaturalLit); ok {
		result := zero
		for i := 0; i < int(n); i++ {
			result = apply(fold.succ, result)
		}
		return result
	}
	return nil
}

func (fold naturalFold) ArgType() Value {
	if fold.n == nil {
		return Natural
	}
	if fold.typ == nil {
		return Type
	}
	if fold.succ == nil {
		return NewFnType("_", fold.typ, fold.typ)
	}
	// zero
	return fold.typ
}

func (naturalIsZero) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n == 0)
	}
	return nil
}

func (naturalIsZero) ArgType() Value { return Natural }

func (naturalOdd) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return BoolLit(n%2 == 1)
	}
	return nil
}

func (naturalOdd) ArgType() Value { return Natural }

func (naturalShow) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return PlainTextLit(fmt.Sprintf("%d", n))
	}
	return nil
}

func (naturalShow) ArgType() Value { return Natural }

func (sub naturalSubtract) Call(x Value) Value {
	if sub.a == nil {
		return naturalSubtract{a: x}
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
	if AlphaEquivalent(sub.a, x) {
		return NaturalLit(0)
	}
	return nil
}

func (naturalSubtract) ArgType() Value { return Natural }

func (naturalToInteger) Call(x Value) Value {
	if n, ok := x.(NaturalLit); ok {
		return IntegerLit(n)
	}
	return nil
}

func (naturalToInteger) ArgType() Value { return Natural }

func (integerClamp) Call(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		if i < 0 {
			return NaturalLit(0)
		}
		return NaturalLit(i)
	}
	return nil
}

func (integerClamp) ArgType() Value { return Integer }

func (integerNegate) Call(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		return IntegerLit(-i)
	}
	return nil
}

func (integerNegate) ArgType() Value { return Integer }

func (integerShow) Call(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		return PlainTextLit(fmt.Sprintf("%+d", i))
	}
	return nil
}

func (integerShow) ArgType() Value { return Integer }

func (integerToDouble) Call(x Value) Value {
	if i, ok := x.(IntegerLit); ok {
		return DoubleLit(i)
	}
	return nil
}

func (integerToDouble) ArgType() Value { return Integer }

func (doubleShow) Call(x Value) Value {
	if d, ok := x.(DoubleLit); ok {
		return PlainTextLit(d.String())
	}
	return nil
}

func (doubleShow) ArgType() Value { return Double }

func (optional) Call(x Value) Value { return OptionalOf{x} }
func (optional) ArgType() Value     { return Type }

func (build optionalBuild) Call(x Value) Value {
	if build.typ == nil {
		return optionalBuild{typ: x}
	}
	var some Value = lambda{
		Label:  "a",
		Domain: build.typ,
		Fn: func(a Value) Value {
			return Some{a}
		},
	}
	return apply(x, OptionalOf{build.typ}, some, NoneOf{build.typ})
}

func (build optionalBuild) ArgType() Value {
	if build.typ == nil {
		return Type
	}
	return NewPi("optional", Type, func(optional Value) Value {
		return NewFnType("just", NewFnType("_", build.typ, optional),
			NewFnType("nothing", optional,
				optional))
	})
}

func (fold optionalFold) Call(x Value) Value {
	if fold.typ1 == nil {
		return optionalFold{typ1: x}
	}
	if fold.opt == nil {
		return optionalFold{typ1: fold.typ1, opt: x}
	}
	if fold.typ2 == nil {
		return optionalFold{
			typ1: fold.typ1,
			opt:  fold.opt,
			typ2: x,
		}
	}
	if fold.some == nil {
		return optionalFold{
			typ1: fold.typ1,
			opt:  fold.opt,
			typ2: fold.typ2,
			some: x,
		}
	}
	none := x
	if s, ok := fold.opt.(Some); ok {
		return apply(fold.some, s.Val)
	}
	if _, ok := fold.opt.(NoneOf); ok {
		return none
	}
	return nil
}

func (fold optionalFold) ArgType() Value {
	if fold.typ1 == nil {
		return Type
	}
	if fold.opt == nil {
		return OptionalOf{fold.typ1}
	}
	if fold.typ2 == nil {
		return Type
	}
	if fold.some == nil {
		return NewFnType("_", fold.typ1, fold.typ2)
	}
	// none
	return fold.typ2
}

func (none) Call(a Value) Value { return NoneOf{a} }
func (none) ArgType() Value     { return Type }

func (textShow) Call(a0 Value) Value {
	if t, ok := a0.(PlainTextLit); ok {
		var out strings.Builder
		out.WriteRune('"')
		for _, r := range t {
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
		return PlainTextLit(out.String())
	}
	return nil
}

func (textShow) ArgType() Value { return Text }

func (list) Call(x Value) Value { return ListOf{x} }
func (list) ArgType() Value     { return Type }

func (l listBuild) Call(x Value) Value {
	if l.typ == nil {
		return listBuild{typ: x}
	}
	var cons Value = lambda{
		Label:  "a",
		Domain: l.typ,
		Fn: func(a Value) Value {
			return lambda{
				Label:  "as",
				Domain: ListOf{l.typ},
				Fn: func(as Value) Value {
					if _, ok := as.(EmptyList); ok {
						return NonEmptyList{a}
					}
					if as, ok := as.(NonEmptyList); ok {
						return append(NonEmptyList{a}, as...)
					}
					return oper{OpCode: term.ListAppendOp, L: NonEmptyList{a}, R: as}
				},
			}
		},
	}
	return apply(x, ListOf{l.typ}, cons, EmptyList{ListOf{l.typ}})
}

func (l listBuild) ArgType() Value {
	if l.typ == nil {
		return Type
	}
	return NewPi("list", Type, func(list Value) Value {
		return NewFnType("cons", NewFnType("_", l.typ, NewFnType("_", list, list)),
			NewFnType("nil", list,
				list))
	})
}

func (l listFold) Call(x Value) Value {
	if l.typ1 == nil {
		return listFold{typ1: x}
	}
	if l.list == nil {
		return listFold{typ1: l.typ1, list: x}
	}
	if l.typ2 == nil {
		return listFold{
			typ1: l.typ1,
			list: l.list,
			typ2: x,
		}
	}
	if l.cons == nil {
		return listFold{
			typ1: l.typ1,
			list: l.list,
			typ2: l.typ2,
			cons: x,
		}
	}
	empty := x
	if _, ok := l.list.(EmptyList); ok {
		return empty
	}
	if list, ok := l.list.(NonEmptyList); ok {
		result := empty
		for i := len(list) - 1; i >= 0; i-- {
			result = apply(l.cons, list[i], result)
		}
		return result
	}
	return nil
}

func (l listFold) ArgType() Value {
	if l.typ1 == nil {
		return Type
	}
	if l.list == nil {
		return ListOf{l.typ1}
	}
	if l.typ2 == nil {
		return Type
	}
	if l.cons == nil {
		return NewFnType("_", l.typ1, NewFnType("_", l.typ2, l.typ2))
	}
	// nil
	return l.typ2
}

func (length listLength) Call(x Value) Value {
	if length.typ == nil {
		return listLength{typ: x}
	}
	if _, ok := x.(EmptyList); ok {
		return NaturalLit(0)
	}
	if l, ok := x.(NonEmptyList); ok {
		return NaturalLit(len(l))
	}
	return nil
}

func (length listLength) ArgType() Value {
	if length.typ == nil {
		return Type
	}
	return ListOf{length.typ}
}

func (head listHead) Call(x Value) Value {
	if head.typ == nil {
		return listHead{typ: x}
	}
	if _, ok := x.(EmptyList); ok {
		return NoneOf{head.typ}
	}
	if l, ok := x.(NonEmptyList); ok {
		return Some{l[0]}
	}
	return nil
}

func (head listHead) ArgType() Value {
	if head.typ == nil {
		return Type
	}
	return ListOf{head.typ}
}

func (last listLast) Call(x Value) Value {
	if last.typ == nil {
		return listLast{typ: x}
	}
	if _, ok := x.(EmptyList); ok {
		return NoneOf{last.typ}
	}
	if l, ok := x.(NonEmptyList); ok {
		return Some{l[len(l)-1]}
	}
	return nil
}

func (last listLast) ArgType() Value {
	if last.typ == nil {
		return Type
	}
	return ListOf{last.typ}
}

func (indexed listIndexed) Call(x Value) Value {
	if indexed.typ == nil {
		return listIndexed{typ: x}
	}
	if _, ok := x.(EmptyList); ok {
		return EmptyList{ListOf{
			RecordType{"index": Natural, "value": indexed.typ}}}
	}
	if l, ok := x.(NonEmptyList); ok {
		var result []Value
		for i, v := range l {
			result = append(result,
				RecordLit{"index": NaturalLit(i), "value": v})
		}
		return NonEmptyList(result)
	}
	return nil
}

func (indexed listIndexed) ArgType() Value {
	if indexed.typ == nil {
		return Type
	}
	return ListOf{indexed.typ}
}

func (rev listReverse) Call(x Value) Value {
	if rev.typ == nil {
		return listReverse{typ: x}
	}
	if _, ok := x.(EmptyList); ok {
		return x
	}
	if l, ok := x.(NonEmptyList); ok {
		result := make([]Value, len(l))
		for i, v := range l {
			result[len(l)-i-1] = v
		}
		return NonEmptyList(result)
	}
	return nil
}

func (rev listReverse) ArgType() Value {
	if rev.typ == nil {
		return Type
	}
	return ListOf{rev.typ}
}

// These are the builtin Callable Values.
var (
	NaturalBuild     Callable = naturalBuild{}
	NaturalEven      Callable = naturalEven{}
	NaturalFold      Callable = naturalFold{}
	NaturalIsZero    Callable = naturalIsZero{}
	NaturalOdd       Callable = naturalOdd{}
	NaturalShow      Callable = naturalShow{}
	NaturalSubtract  Callable = naturalSubtract{}
	NaturalToInteger Callable = naturalToInteger{}
	IntegerClamp     Callable = integerClamp{}
	IntegerNegate    Callable = integerNegate{}
	IntegerShow      Callable = integerShow{}
	IntegerToDouble  Callable = integerToDouble{}
	DoubleShow       Callable = doubleShow{}

	Optional      Callable = optional{}
	OptionalBuild Callable = optionalBuild{}
	OptionalFold  Callable = optionalFold{}
	None          Callable = none{}

	TextShow Callable = textShow{}

	List        Callable = list{}
	ListBuild   Callable = listBuild{}
	ListFold    Callable = listFold{}
	ListLength  Callable = listLength{}
	ListHead    Callable = listHead{}
	ListLast    Callable = listLast{}
	ListIndexed Callable = listIndexed{}
	ListReverse Callable = listReverse{}
)
