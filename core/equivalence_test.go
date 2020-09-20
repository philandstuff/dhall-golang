package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/philandstuff/dhall-golang/v5/term"
)

// Ensure that alphaMatcher is a valid GomegaMatcher
var _ types.GomegaMatcher = &alphaMatcher{}

var _ = Describe("AlphaEquivalent", func() {
	DescribeTable("returns false for Terms that aren't equivalent",
		func(l, r term.Term) {
			Expect(l).ToNot(BeAlphaEquivalentTo(r))
		},
		Entry("Bool and Natural", term.Bool, term.Natural),
		Entry("λ(a : Natural) → a` and `λ(b : Natural) → 3",
			term.NewLambda("a", term.Natural, term.NewVar("a")),
			term.NewLambda("b", term.Natural, term.NaturalLit(3))),
	)
	DescribeTable("a Term is AlphaEquivalent to itself",
		func(t term.Term) {
			Expect(t).To(BeAlphaEquivalentTo(t))
		},
		Entry(`λ(a : Natural) → a`,
			term.NewLambda("a", term.Natural, term.NewVar("a"))),
		Entry(`∀(a : Type) → List a`,
			term.NewPi("a", term.Type, term.Apply(term.List, term.NewVar("a")))),
		Entry(`{a = True, c = _}`,
			term.RecordLit{"a": term.True, "c": term.NewVar("_")}),
		Entry(`toMap {a = True, c = _}`,
			term.ToMap{Record: term.RecordLit{"a": term.True, "c": term.NewVar("_")}}),
		Entry(`toMap {a = True, c = _} ≡ toMap {a = True, c = _}`,
			term.Equivalent(term.ToMap{Record: term.RecordLit{"a": term.True, "c": term.NewVar("_")}}, term.ToMap{Record: term.RecordLit{"a": term.True, "c": term.NewVar("_")}})),
		Entry(`[{mapKey = "a", mapValue = True}, {mapKey = "c", mapValue = _}] ≡ [{mapKey = "a", mapValue = True}, {mapKey = "c", mapValue = _}]`,
			term.Equivalent(
				term.NewList(term.RecordLit{"mapKey": term.PlainText("a"), "mapValue": term.True}, term.RecordLit{"mapKey": term.PlainText("a"), "mapValue": term.NewVar("_")}),
				term.NewList(term.RecordLit{"mapKey": term.PlainText("a"), "mapValue": term.True}, term.RecordLit{"mapKey": term.PlainText("a"), "mapValue": term.NewVar("_")}))),
		Entry(`∀(_ : Bool) → [{mapKey = "a", mapValue = True}, {mapKey = "c", mapValue = _}] ≡ [{mapKey = "a", mapValue = True}, {mapKey = "c", mapValue = _}]`,
			term.NewPi("_", term.Bool, term.Equivalent(
				term.NewList(term.RecordLit{"mapKey": term.PlainText("a"), "mapValue": term.True}, term.RecordLit{"mapKey": term.PlainText("a"), "mapValue": term.NewVar("_")}),
				term.NewList(term.RecordLit{"mapKey": term.PlainText("a"), "mapValue": term.True}, term.RecordLit{"mapKey": term.PlainText("a"), "mapValue": term.NewVar("_")})))),
		Entry("`∀(_ : Type) → toMap {a = True, c = _} ≡ toMap {a = True, c = _}`",
			term.NewPi("_", term.Type, term.Equivalent(term.ToMap{Record: term.RecordLit{"a": term.True, "c": term.NewVar("_")}}, term.ToMap{Record: term.RecordLit{"a": term.True, "c": term.NewVar("_")}}))),
	)
	DescribeTable("returns true for nonidentical but alpha-equivalent Terms",
		func(l, r term.Term) {
			Expect(l).To(BeAlphaEquivalentTo(r))
		},
		Entry("`λ(a : Natural) → a` and `λ(b : Natural) → b`",
			term.NewLambda("a", term.Natural, term.NewVar("a")),
			term.NewLambda("b", term.Natural, term.NewVar("b"))),
		Entry("`∀(a : Type) → List a` and ∀(b : Type) → List b`",
			term.NewPi("a", term.Type, term.Apply(term.List, term.NewVar("a"))),
			term.NewPi("b", term.Type, term.Apply(term.List, term.NewVar("b")))),
	)
})
