package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// Ensure that alphaMatcher is a valid GomegaMatcher
var _ types.GomegaMatcher = &alphaMatcher{}

var _ = Describe("AlphaEquivalent", func() {
	DescribeTable("returns false for Terms that aren't equivalent",
		func(l, r Term) {
			Expect(l).ToNot(BeAlphaEquivalentTo(r))
		},
		Entry("Bool and Natural", Bool, Natural),
		Entry("λ(a : Natural) → a` and `λ(b : Natural) → 3",
			NewLambda("a", Natural, NewVar("a")),
			NewLambda("b", Natural, NaturalLit(3))),
	)
	DescribeTable("a Term is AlphaEquivalent to itself",
		func(t Term) {
			Expect(t).To(BeAlphaEquivalentTo(t))
		},
		Entry(`λ(a : Natural) → a`,
			NewLambda("a", Natural, NewVar("a"))),
		Entry(`∀(a : Type) → List a`,
			NewPi("a", Type, Apply(List, NewVar("a")))),
		Entry(`{a = True, c = _}`,
			RecordLit{"a": True, "c": NewVar("_")}),
		Entry(`toMap {a = True, c = _}`,
			ToMap{Record: RecordLit{"a": True, "c": NewVar("_")}}),
		Entry(`toMap {a = True, c = _} ≡ toMap {a = True, c = _}`,
			Equivalent(ToMap{Record: RecordLit{"a": True, "c": NewVar("_")}}, ToMap{Record: RecordLit{"a": True, "c": NewVar("_")}})),
		Entry(`[{mapKey = "a", mapValue = True}, {mapKey = "c", mapValue = _}] ≡ [{mapKey = "a", mapValue = True}, {mapKey = "c", mapValue = _}]`,
			Equivalent(
				NewList(RecordLit{"mapKey": PlainText("a"), "mapValue": True}, RecordLit{"mapKey": PlainText("a"), "mapValue": NewVar("_")}),
				NewList(RecordLit{"mapKey": PlainText("a"), "mapValue": True}, RecordLit{"mapKey": PlainText("a"), "mapValue": NewVar("_")}))),
		Entry(`∀(_ : Bool) → [{mapKey = "a", mapValue = True}, {mapKey = "c", mapValue = _}] ≡ [{mapKey = "a", mapValue = True}, {mapKey = "c", mapValue = _}]`,
			NewPi("_", Bool, Equivalent(
				NewList(RecordLit{"mapKey": PlainText("a"), "mapValue": True}, RecordLit{"mapKey": PlainText("a"), "mapValue": NewVar("_")}),
				NewList(RecordLit{"mapKey": PlainText("a"), "mapValue": True}, RecordLit{"mapKey": PlainText("a"), "mapValue": NewVar("_")})))),
		Entry("`∀(_ : Type) → toMap {a = True, c = _} ≡ toMap {a = True, c = _}`",
			NewPi("_", Type, Equivalent(ToMap{Record: RecordLit{"a": True, "c": NewVar("_")}}, ToMap{Record: RecordLit{"a": True, "c": NewVar("_")}}))),
	)
	DescribeTable("returns true for nonidentical but alpha-equivalent Terms",
		func(l, r Term) {
			Expect(l).To(BeAlphaEquivalentTo(r))
		},
		Entry("`λ(a : Natural) → a` and `λ(b : Natural) → b`",
			NewLambda("a", Natural, NewVar("a")),
			NewLambda("b", Natural, NewVar("b"))),
		Entry("`∀(a : Type) → List a` and ∀(b : Type) → List b`",
			NewPi("a", Type, Apply(List, NewVar("a"))),
			NewPi("b", Type, Apply(List, NewVar("b")))),
	)
})
