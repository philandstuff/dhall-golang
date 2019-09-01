package core

// Terms in the dhall internal language
type Term interface {
	isTerm()
}

// values
type Value interface {
	isValue()
}

type Universe int

const (
	Type Universe = iota
	Kind
	Sort
)

type Builtin string

const (
	Natural Builtin = "Natural"
	List    Builtin = "List"
)

type (
	BoundVar struct {
		Name  string
		Index int
	}

	FreeVar struct {
		Name  string
		Index int
	}

	// a LocalVar is an internal sentinel value used by TypeOf() in the process
	// of typechecking the body of lambdas and pis
	// FIXME: should this be unexported?
	LocalVar struct {
		Name  string
		Index int
	}

	// a QuoteVar is an internal sentinel value used by Quote() in the process
	// of converting Values back to Terms
	// FIXME: should this be unexported?
	QuoteVar struct {
		Name  string
		Index int
	}
)

func (Universe) isTerm()  {}
func (Universe) isValue() {}

func (Builtin) isTerm()  {}
func (Builtin) isValue() {}

func (BoundVar) isTerm() {}

func (FreeVar) isTerm()  {}
func (FreeVar) isValue() {}

func (LocalVar) isTerm()  {}
func (LocalVar) isValue() {}

func (QuoteVar) isValue() {}

func Bound(name string) BoundVar {
	return BoundVar{Name: name}
}

type (
	LambdaTerm struct {
		Label string
		Type  Term
		Body  Term
	}

	PiTerm struct {
		Label string
		Type  Term
		Body  Term
	}

	AppTerm struct {
		Fn  Term
		Arg Term
	}
)

func (LambdaTerm) isTerm() {}
func (PiTerm) isTerm()     {}
func (AppTerm) isTerm()    {}

func Lambda(label string, t Term, body Term) LambdaTerm {
	return LambdaTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

func Mkλ(label string, t Term, body Term) LambdaTerm {
	return LambdaTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

func MkΠ(label string, t Term, body Term) PiTerm {
	return PiTerm{
		Label: label,
		Type:  t,
		Body:  body,
	}
}

// FnType returns a non-dependent function type
func FnType(domain Term, codomain Term) PiTerm {
	return PiTerm{
		Label: "_",
		Type:  domain,
		Body:  codomain,
	}
}

func Apply(fn Term, arg Term) AppTerm {
	return AppTerm{
		Fn:  fn,
		Arg: arg,
	}
}

type (
	// a lambda value is a go function
	LambdaValue struct {
		Label  string
		Domain Value
		Fn     func(Value) Value
	}

	// A PiValue is: the type of the domain, and a go function capturing the
	// range
	PiValue struct {
		Label  string
		Domain Value
		Range  func(Value) Value
	}

	AppValue struct {
		Fn  Value
		Arg Value
	}
)

func (LambdaValue) isValue() {}
func (PiValue) isValue()     {}

func (AppValue) isValue() {}

func MkΠval(label string, d Value, r func(Value) Value) PiValue {
	return PiValue{
		Label:  label,
		Domain: d,
		Range:  r,
	}
}

// FnTypeVal returns a non-dependent function type Value
func FnTypeVal(d Value, r Value) PiValue {
	return PiValue{
		Label:  "_",
		Domain: d,
		Range:  func(Value) Value { return r },
	}
}

type (
	// 0, 1, 2, 3...
	NaturalLit uint

	// [] : List a
	EmptyList    struct{ Type Term }
	EmptyListVal struct{ Type Value }
)

func (NaturalLit) isTerm()  {}
func (NaturalLit) isValue() {}

func (EmptyList) isTerm()     {}
func (EmptyListVal) isValue() {}
