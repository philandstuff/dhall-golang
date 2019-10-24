package core_test

import (
	. "github.com/philandstuff/dhall-golang/core"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Quote",
	func(v Value, expected Term) {
		Expect(Quote(v)).
			To(Equal(expected))
	},
	Entry("Type", Type, Type),
	Entry("Kind", Kind, Kind),
	Entry("Sort", Sort, Sort),
	Entry(`λ(x : Natural) → x`,
		LambdaValue{Label: "x", Domain: Natural, Fn: func(x Value) Value {
			return x
		}},
		LambdaTerm{Label: "x", Type: Natural, Body: BoundVar{"x", 0}},
	),
	Entry(`λ(x : Natural) → λ(x : Natural) → x`,
		LambdaValue{Label: "x", Domain: Natural, Fn: func(x Value) Value {
			return LambdaValue{
				Label:  "x",
				Domain: Natural,
				Fn:     func(x Value) Value { return x },
			}
		}},
		LambdaTerm{Label: "x", Type: Natural, Body: LambdaTerm{
			Label: "x", Type: Natural,
			Body: BoundVar{"x", 0}}},
	),
	Entry(`λ(x : Natural) → λ(x : Natural) → x@1`,
		LambdaValue{Label: "x", Domain: Natural, Fn: func(x1 Value) Value {
			return LambdaValue{
				Label:  "x",
				Domain: Natural,
				Fn:     func(x Value) Value { return x1 },
			}
		}},
		LambdaTerm{Label: "x", Type: Natural, Body: LambdaTerm{
			Label: "x", Type: Natural,
			Body: BoundVar{"x", 1}}},
	),
	Entry(`λ(x : Natural) → λ(y : Natural) → x`,
		LambdaValue{Label: "x", Domain: Natural, Fn: func(x Value) Value {
			return LambdaValue{
				Label:  "y",
				Domain: Natural,
				Fn:     func(y Value) Value { return x },
			}
		}},
		LambdaTerm{Label: "x", Type: Natural, Body: LambdaTerm{
			Label: "y", Type: Natural,
			Body: BoundVar{"x", 0}}},
	),
	Entry(`Natural → Natural`,
		PiValue{Label: "_", Domain: Natural, Range: func(x Value) Value {
			return Natural
		}},
		PiTerm{Label: "_", Type: Natural, Body: Natural},
	),
	Entry(`∀(a : Type) → List a`,
		PiValue{Label: "a", Domain: Type, Range: func(x Value) Value {
			return AppValue{List, x}
		}},
		PiTerm{Label: "a", Type: Type, Body: AppTerm{List, BoundVar{"a", 0}}},
	),
	Entry(`[] : List Natural`,
		EmptyListVal{Type: AppValue{Fn: List, Arg: Natural}},
		EmptyList{Type: AppTerm{Fn: List, Arg: Natural}}),
)
