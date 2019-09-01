package eval_test

import (
	. "github.com/philandstuff/dhall-golang/core"
	. "github.com/philandstuff/dhall-golang/eval"

	. "github.com/onsi/gomega"
)

var _ = Describe("Eval", func() {
	It("Type", func() {
		Expect(Eval(Type, Env{})).To(Equal(Type))
	})
	It("Bound variable", func() {
		Expect(Eval(BoundVar{Name: "foo"}, Env{"foo": []Value{Type}})).
			To(Equal(Type))
	})
	It("Free variable", func() {
		Expect(Eval(FreeVar{Name: "foo"}, Env{"foo": []Value{Type}})).
			To(Equal(FreeVar{Name: "foo"}))
	})
	It("Lambda id function", func() {
		f := Eval(Lambda("x", Type, BoundVar{Name: "x"}), Env{}).(LambdaValue)
		Expect(f.Fn(Type)).
			To(Equal(Type))
	})
	Describe("application", func() {
		It("To neutral", func() {
			Expect(Eval(Apply(FreeVar{Name: "f"}, FreeVar{Name: "x"}), Env{})).
				To(Equal(AppValue{Fn: FreeVar{Name: "f"}, Arg: FreeVar{Name: "x"}}))
		})
		It("To lambda", func() {
			Expect(Eval(Apply(
				Lambda("x", Kind, BoundVar{Name: "x"}),
				Type,
			), Env{})).
				To(Equal(Type))
		})
	})
})
