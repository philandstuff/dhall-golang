package core_test

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/parser"
)

var _ = DescribeTable("ArgType of builtins", func(src, typ string) {
	term, err := parser.Parse("-", []byte(src))
	Expect(err).To(Not(HaveOccurred()))
	expectedTypeTerm, err := parser.Parse("-", []byte(typ))
	Expect(err).To(Not(HaveOccurred()))
	Expect(core.Eval(term).(core.Callable).ArgType()).
		To(core.BeAlphaEquivalentTo(core.Eval(expectedTypeTerm)))
},
	Entry("Natural/build", `Natural/build`, `∀(natural : Type) → ∀(succ : natural → natural) → ∀(zero : natural) → natural`),
	Entry("Natural/fold arg 1", `Natural/fold`, `Natural`),
	Entry("Natural/fold arg 2", `Natural/fold 1`, `Type`),
	Entry("Natural/fold arg 3", `Natural/fold 1 Text`, `Text → Text`),
	Entry("Natural/fold arg 4", `Natural/fold 1 Text (λ(x : Text) → "some")`, `Text`),
	Entry("Natural/isZero", `Natural/isZero`, `Natural`),
	Entry("Natural/even", `Natural/even`, `Natural`),
	Entry("Natural/odd", `Natural/odd`, `Natural`),
	Entry("Natural/toInteger", `Natural/toInteger`, `Natural`),
	Entry("Natural/show", `Natural/show`, `Natural`),
	Entry("Natural/subtract arg 1", `Natural/subtract`, `Natural`),
	Entry("Natural/subtract arg 2", `Natural/subtract 3`, `Natural`),

	Entry("Text/show", `Text/show`, `Text`),

	Entry("List/build arg 1", `List/build`, `Type`),
	Entry("List/build arg 2", `List/build Natural`,
		`∀(list : Type) → ∀(cons : Natural → list → list) → ∀(nil : list) → list`),
	Entry("List/fold arg 1", `List/fold`, `Type`),
	Entry("List/fold arg 2", `List/fold Natural`, `List Natural`),
	Entry("List/fold arg 3", `List/fold Natural [1,2]`, `Type`),
	Entry("List/fold arg 4", `List/fold Natural [1,2] Text`, `Natural → Text → Text`),
	Entry("List/fold arg 5", `List/fold Natural [1,2] Text (λ(x : Natural) → λ(acc : Text) → "some")`, `Text`),
	Entry("List/length arg 1", `List/length`, `Type`),
	Entry("List/length arg 2", `List/length Natural`, `List Natural`),
	Entry("List/head arg 1", `List/head`, `Type`),
	Entry("List/head arg 2", `List/head Natural`, `List Natural`),
	Entry("List/last arg 1", `List/last`, `Type`),
	Entry("List/last arg 2", `List/last Natural`, `List Natural`),
	Entry("List/indexed arg 1", `List/indexed`, `Type`),
	Entry("List/indexed arg 2", `List/indexed Natural`, `List Natural`),
	Entry("List/reverse arg 1", `List/reverse`, `Type`),
	Entry("List/reverse arg 2", `List/reverse Natural`, `List Natural`),

	Entry("Optional/build arg 1", `Optional/build`, `Type`),
	Entry("Optional/build arg 2", `Optional/build Natural`,
		`∀(optional : Type) → ∀(just : Natural → optional) → ∀(nothing : optional) → optional`),
	Entry("Optional/fold arg 1", `Optional/fold`, `Type`),
	Entry("Optional/fold arg 2", `Optional/fold Natural`, `Optional Natural`),
	Entry("Optional/fold arg 3", `Optional/fold Natural (Some 2)`, `Type`),
	Entry("Optional/fold arg 4", `Optional/fold Natural (Some 2) Text`, `Natural → Text`),
	Entry("Optional/fold arg 5", `Optional/fold Natural (Some 2) Text (λ(x : Natural) → "some")`, `Text`),

	Entry("Integer/show", `Integer/show`, `Integer`),
	Entry("Integer/toDouble", `Integer/toDouble`, `Integer`),
	Entry("Integer/negate", `Integer/negate`, `Integer`),
	Entry("Integer/clamp", `Integer/clamp`, `Integer`),

	Entry("Double/show", `Double/show`, `Double`),
)
