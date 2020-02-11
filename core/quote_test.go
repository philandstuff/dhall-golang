package core

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/philandstuff/dhall-golang/term"
)

var _ = DescribeTable("Quote",
	func(v Value, expected term.Term) {
		Expect(Quote(v)).
			To(Equal(expected))
	},
	Entry("Type", Type, term.Type),
	Entry("Kind", Kind, term.Kind),
	Entry("Sort", Sort, term.Sort),
	Entry(`λ(x : Natural) → x`,
		lambdaValue{Label: "x", Domain: Natural, Fn: func(x Value) Value {
			return x
		}},
		term.Lambda{Label: "x", Type: term.Natural, Body: term.Var{"x", 0}},
	),
	Entry(`λ(x : Natural) → λ(x : Natural) → x`,
		lambdaValue{Label: "x", Domain: Natural, Fn: func(x Value) Value {
			return lambdaValue{
				Label:  "x",
				Domain: Natural,
				Fn:     func(x Value) Value { return x },
			}
		}},
		term.Lambda{Label: "x", Type: term.Natural, Body: term.Lambda{
			Label: "x", Type: term.Natural,
			Body: term.Var{"x", 0}}},
	),
	Entry(`λ(x : Natural) → λ(x : Natural) → x@1`,
		lambdaValue{Label: "x", Domain: Natural, Fn: func(x1 Value) Value {
			return lambdaValue{
				Label:  "x",
				Domain: Natural,
				Fn:     func(x Value) Value { return x1 },
			}
		}},
		term.Lambda{Label: "x", Type: term.Natural, Body: term.Lambda{
			Label: "x", Type: term.Natural,
			Body: term.Var{"x", 1}}},
	),
	Entry(`λ(x : Natural) → λ(y : Natural) → x`,
		lambdaValue{Label: "x", Domain: Natural, Fn: func(x Value) Value {
			return lambdaValue{
				Label:  "y",
				Domain: Natural,
				Fn:     func(y Value) Value { return x },
			}
		}},
		term.Lambda{Label: "x", Type: term.Natural, Body: term.Lambda{
			Label: "y", Type: term.Natural,
			Body: term.Var{"x", 0}}},
	),
	Entry(`Natural → Natural`,
		PiValue{Label: "_", Domain: Natural, Range: func(x Value) Value {
			return Natural
		}},
		term.Pi{Label: "_", Type: term.Natural, Body: term.Natural},
	),
	Entry(`∀(a : Type) → List a`,
		PiValue{Label: "a", Domain: Type, Range: func(x Value) Value {
			return ListOf{x}
		}},
		term.Pi{Label: "a", Type: term.Type, Body: term.App{term.List, term.Var{"a", 0}}},
	),
	Entry(`[] : List Natural`,
		EmptyListVal{Type: ListOf{Type: Natural}},
		term.EmptyList{Type: term.App{Fn: term.List, Arg: term.Natural}}),
)
