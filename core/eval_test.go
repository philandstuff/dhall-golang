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
		Expect(evalWith(Var{Name: "foo"}, Env{"foo": []Value{Type}}, false)).
			To(Equal(Type))
	})
	It("Free variable", func() {
		Expect(evalWith(Var{Name: "foo"}, Env{}, false)).
			To(Equal(Var{Name: "foo"}))
	})
	It("Lambda id function", func() {
		f := Eval(Lambda("x", Type, Var{Name: "x"})).(LambdaValue)
		Expect(f.Call(Type)).
			To(Equal(Type))
		Expect(f.Label).
			To(Equal("x"))
	})
	It("Lambda id function with alpha normalization", func() {
		f := AlphaBetaEval(Lambda("x", Type, Var{Name: "x"})).(LambdaValue)
		Expect(f.Call(Type)).
			To(Equal(Type))
		Expect(f.Label).
			To(Equal("_"))
	})
	Describe("application", func() {
		It("To neutral", func() {
			Expect(Eval(Apply(Var{Name: "f"}, Var{Name: "x"}))).
				To(Equal(AppValue{Fn: Var{Name: "f"}, Arg: Var{Name: "x"}}))
		})
		It("To lambda", func() {
			Expect(Eval(Apply(
				Lambda("x", Kind, Var{Name: "x"}),
				Type,
			))).
				To(Equal(Type))
		})
	})
})
