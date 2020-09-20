package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/philandstuff/dhall-golang/v5/term"
)

var _ = Describe("Eval", func() {
	It("Type", func() {
		Expect(Eval(term.Type)).To(Equal(Type))
	})
	It("Bound variable", func() {
		Expect(evalWith(term.Var{Name: "foo"}, env{"foo": []Value{Type}})).
			To(Equal(Type))
	})
	It("Free variable", func() {
		Expect(evalWith(term.Var{Name: "foo"}, env{})).
			To(Equal(freeVar{Name: "foo"}))
	})
	It("Lambda id function", func() {
		f := Eval(term.NewLambda("x", term.Type, term.Var{Name: "x"})).(lambda)
		Expect(f.Call(Type)).
			To(Equal(Type))
		Expect(f.Label).
			To(Equal("x"))
	})
	Describe("application", func() {
		It("To neutral", func() {
			Expect(Eval(term.Apply(term.Var{Name: "f"}, term.Var{Name: "x"}))).
				To(Equal(app{Fn: freeVar{Name: "f"}, Arg: freeVar{Name: "x"}}))
		})
		It("To lambda", func() {
			Expect(Eval(term.Apply(
				term.NewLambda("x", term.Kind, term.Var{Name: "x"}),
				term.Type,
			))).
				To(Equal(Type))
		})
	})
	Describe("toMap", func() {
		It("Evaluates with missing type and abstract value", func() {
			Expect(Eval(term.ToMap{
				Record: term.Var{Name: "x"},
			})).
				To(Equal(toMap{
					Record: freeVar{Name: "x"},
				}))
		})
	})
})
