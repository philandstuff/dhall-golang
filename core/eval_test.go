package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Eval", func() {
	It("Type", func() {
		Expect(Eval(Type)).To(Equal(Type))
	})
	It("Bound variable", func() {
		Expect(evalWith(BoundVar{Name: "foo"}, Env{"foo": []Value{Type}}, false)).
			To(Equal(Type))
	})
	It("Free variable", func() {
		Expect(evalWith(FreeVar{Name: "foo"}, Env{"foo": []Value{Type}}, false)).
			To(Equal(FreeVar{Name: "foo"}))
	})
	It("Lambda id function", func() {
		f := Eval(Lambda("x", Type, BoundVar{Name: "x"})).(LambdaValue)
		Expect(f.Call1(Type)).
			To(Equal(Type))
		Expect(f.Label).
			To(Equal("x"))
	})
	It("Lambda id function with alpha normalization", func() {
		f := AlphaBetaEval(Lambda("x", Type, BoundVar{Name: "x"})).(LambdaValue)
		Expect(f.Call1(Type)).
			To(Equal(Type))
		Expect(f.Label).
			To(Equal("_"))
	})
	Describe("application", func() {
		It("To neutral", func() {
			Expect(Eval(Apply(FreeVar{Name: "f"}, FreeVar{Name: "x"}))).
				To(Equal(AppValue{Fn: FreeVar{Name: "f"}, Arg: FreeVar{Name: "x"}}))
		})
		It("To lambda", func() {
			Expect(Eval(Apply(
				Lambda("x", Kind, BoundVar{Name: "x"}),
				Type,
			))).
				To(Equal(Type))
		})
	})
})
