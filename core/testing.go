package core

import (
	"errors"
	"fmt"

	"github.com/philandstuff/dhall-golang/v6/term"
)

// GomegaMatcher is a copy of
// github.com/onsi/gomega/types.GomegaMatcher, to avoid a dependency
// on gomega for users who don't want it.
type GomegaMatcher interface {
	Match(actual interface{}) (success bool, err error)
	FailureMessage(actual interface{}) (message string)
	NegatedFailureMessage(actual interface{}) (message string)
}

// BeAlphaEquivalentTo returns a Gomega matcher which checks that a
// Term or Value is equivalent to the expected Term or Value.  If
// either the expected or actual is a Term, it is Eval()ed first; then
// the Values are compared with AlphaEquivalentVals().
func BeAlphaEquivalentTo(expected interface{}) GomegaMatcher {
	switch e := expected.(type) {
	case term.Term:
		return &alphaMatcher{
			expected: Eval(e),
		}
	case Value:
		return &alphaMatcher{
			expected: e,
		}
	}
	panic("BeAlphaEquivalentTo needs a Term or a Value")
}

type alphaMatcher struct {
	expected Value
}

func (m *alphaMatcher) Match(actual interface{}) (success bool, err error) {
	var actualValue Value
	switch a := actual.(type) {
	case term.Term:
		actualValue = Eval(a)
	case Value:
		actualValue = a
	default:
		return false, errors.New("BeAlphaEquivalentTo expects a Term or a Value")
	}
	return AlphaEquivalent(m.expected, actualValue), nil
}

func (m *alphaMatcher) FailureMessage(actual interface{}) (message string) {
	switch a := actual.(type) {
	case Value:
		return fmt.Sprintf("Expected\n\t%#v\nto be alpha-equivalent to\n\t%#v", Quote(a), Quote(m.expected))
	default:
		return fmt.Sprintf("Expected\n\t%#v\nto be alpha-equivalent to\n\t%#v", a, Quote(m.expected))
	}
}

func (m *alphaMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	switch a := actual.(type) {
	case Value:
		return fmt.Sprintf("Expected\n\t%v\nnot to be alpha-equivalent to\n\t%v", Quote(a), Quote(m.expected))
	default:
		return fmt.Sprintf("Expected\n\t%v\nnot to be alpha-equivalent to\n\t%v", a, Quote(m.expected))
	}
}
