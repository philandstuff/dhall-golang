package ast

import (
	"errors"
	"fmt"
	"io"
)

type TypeContext map[string][]Expr

func (ctx *TypeContext) Insert(name string, val Expr) *TypeContext {
	newctx := make(TypeContext)
	for k, v := range *ctx {
		newctx[k] = v
	}
	newctx[name] = append(newctx[name], val)
	return &newctx
}

func (ctx *TypeContext) Lookup(name string, i int) (Expr, bool) {
	slice := (*ctx)[name]
	if i >= len(slice) {
		return nil, false
	}
	// we index from the right of the slices
	return slice[len(slice)-1-i], true
}

//TODO: make this lazy
func (ctx *TypeContext) Map(f func(Expr) Expr) *TypeContext {
	newctx := make(TypeContext)
	for k, vs := range *ctx {
		a := make([]Expr, len(vs))
		for i, v := range vs {
			a[i] = f(v)
		}
		newctx[k] = a
	}
	return &newctx
}

func EmptyContext() *TypeContext {
	return &TypeContext{}
}

type (
	Expr interface {
		io.WriterTo
		Normalize() Expr
		TypeWith(*TypeContext) (Expr, error)
	}

	Const int

	Var struct {
		Name  string
		Index int
	}

	LambdaExpr struct {
		Label string
		Type  Expr
		Body  Expr
	}

	Pi struct {
		Label string
		Type  Expr
		Body  Expr
	}

	App struct {
		Fn  Expr
		Arg Expr
	}

	Annot struct {
		Expr       Expr
		Annotation Expr
	}

	double    struct{}
	DoubleLit float64

	// 'boolean' to avoid clashing with go's bool type
	boolean struct{}
	BoolLit bool

	natural     struct{}
	NaturalLit  int
	NaturalPlus struct {
		L Expr
		R Expr
	}

	integer    struct{}
	IntegerLit int

	list struct{}
	// `[] : List Natural` == EmptyList{Natural}
	EmptyList struct{ Type Expr }
	// `[2,3,4]` == NonEmptyList(2,3,4)
	NonEmptyList []Expr
)

const (
	Type Const = Const(iota)
	Kind Const = Const(iota)
	Sort Const = Const(iota)
)

// FIXME placeholder before we actually implement it
func Shift(d int, v Var, e Expr) Expr {
	switch e := e.(type) {
	case Const:
		return e
	case Var:
		if v.Name == e.Name && v.Index <= e.Index {
			return Var{Name: e.Name, Index: e.Index + d}
		} else {
			return e
		}
	case *LambdaExpr:
		var body Expr
		if e.Label == v.Name {
			body = Shift(d, Var{Name: v.Name, Index: v.Index + 1}, e.Body)
		} else {
			body = Shift(d, v, e.Body)
		}
		return &LambdaExpr{
			Label: e.Label,
			Type:  Shift(d, v, e.Type),
			Body:  body,
		}
	case *Pi:
		var body Expr
		if e.Label == v.Name {
			body = Shift(d, Var{Name: v.Name, Index: v.Index + 1}, e.Body)
		} else {
			body = Shift(d, v, e.Body)
		}
		return &Pi{
			Label: e.Label,
			Type:  Shift(d, v, e.Type),
			Body:  body,
		}
	case *App:
		return &App{
			Fn:  Shift(d, v, e.Fn),
			Arg: Shift(d, v, e.Arg),
		}
	case Annot:
		return Annot{Shift(d, v, e.Expr), Shift(d, v, e.Annotation)}
	case double:
		return e
	case DoubleLit:
		return e
	case boolean:
		return e
	case BoolLit:
		return e
	case natural:
		return e
	case NaturalLit:
		return e
	case NaturalPlus:
		return NaturalPlus{L: Shift(d, v, e.L), R: Shift(d, v, e.R)}
	case integer:
		return e
	case IntegerLit:
		return e
	case list:
		return e
	case EmptyList:
		return e
	case NonEmptyList:
		exprs := make([]Expr, len([]Expr(e)))
		for i, expr := range []Expr(e) {
			exprs[i] = Shift(d, v, expr)
		}
		return NonEmptyList(exprs)
	}
	panic("missing switch case in Shift()")
}

// Subst(x, C, B) == B[x := C]
func Subst(v Var, c Expr, b Expr) Expr {
	switch e := b.(type) {
	case Const:
		return e
	case Var:
		if e == v {
			return c
		} else {
			return b
		}
	case *LambdaExpr:
		substType := Subst(v, c, e.Type)
		v2 := v
		if v.Name == e.Label {
			v2.Index++
		}
		substBody := Subst(v2, Shift(1, Var{Name: e.Label}, c), e.Body)
		return &LambdaExpr{
			Label: e.Label,
			Type:  substType,
			Body:  substBody,
		}
	case *Pi:
		substType := Subst(v, c, e.Type)
		v2 := v
		if v.Name == e.Label {
			v2.Index++
		}
		substBody := Subst(v2, Shift(1, Var{Name: e.Label}, c), e.Body)
		return &Pi{
			Label: e.Label,
			Type:  substType,
			Body:  substBody,
		}
	case *App:
		return &App{
			Fn:  Subst(v, c, e.Fn),
			Arg: Subst(v, c, e.Arg),
		}
	case Annot:
		return Annot{Subst(v, c, e.Expr), Subst(v, c, e.Annotation)}
	case double:
		return e
	case DoubleLit:
		return e
	case boolean:
		return e
	case BoolLit:
		return e
	case natural:
		return e
	case NaturalLit:
		return e
	case NaturalPlus:
		return NaturalPlus{L: Subst(v, c, e.L), R: Subst(v, c, e.R)}
	case integer:
		return e
	case IntegerLit:
		return e
	case list:
		return e
	case EmptyList:
		return e
	case NonEmptyList:
		exprs := make([]Expr, len([]Expr(e)))
		for i, expr := range []Expr(e) {
			exprs[i] = Subst(v, c, expr)
		}
		return NonEmptyList(exprs)
	}
	panic("missing switch case in Subst()")
}

func Rule(a Const, b Const) (Const, error) {
	if b == Type {
		return Type, nil
	}
	if a == Kind && b == Kind {
		return Kind, nil
	}
	if a == Sort && (b == Kind || b == Sort) {
		return Sort, nil
	}
	return Const(0), errors.New("Dependent types are not allowed")
}

var (
	Double  double  = double(struct{}{})
	Bool    boolean = boolean(struct{}{})
	Natural natural = natural(struct{}{})
	Integer integer = integer(struct{}{})
	List    list    = list(struct{}{})
)

const (
	True  = BoolLit(true)
	False = BoolLit(false)
)

func MakeList(first Expr, rest ...Expr) NonEmptyList {
	return NonEmptyList(append([]Expr{first}, rest...))
}

var (
	_ Expr = Type
	_ Expr = &Var{}
	_ Expr = &LambdaExpr{}
	_ Expr = &Pi{}
	_ Expr = &App{}
	_ Expr = Annot{}
	_ Expr = Double
	_ Expr = DoubleLit(3.0)
	_ Expr = Bool
	_ Expr = BoolLit(true)
	_ Expr = Natural
	_ Expr = NaturalLit(3)
	_ Expr = NaturalPlus{}
	_ Expr = Integer
	_ Expr = IntegerLit(3)
	_ Expr = List
	_ Expr = EmptyList{Natural}
	_ Expr = NonEmptyList([]Expr{NaturalLit(3)})
)

func (c Const) WriteTo(out io.Writer) (int64, error) {
	var n int
	var err error
	if c == Type {
		n, err = fmt.Fprint(out, "Type")
	} else if c == Kind {
		n, err = fmt.Fprint(out, "Kind")
	} else {
		n, err = fmt.Fprint(out, "Sort")
	}
	return int64(n), err
}

func (v Var) WriteTo(out io.Writer) (int64, error) {
	var n int
	var err error
	if v.Index == 0 {
		n, err = fmt.Fprint(out, v.Name)
	}
	n, err = fmt.Fprintf(out, "%s@%d", v.Name, v.Index)
	return int64(n), err
}

func (lam *LambdaExpr) WriteTo(out io.Writer) (int64, error) {
	w1, err := fmt.Fprintf(out, "λ(%s : ", lam.Label)
	if err != nil {
		return int64(w1), err
	}
	w2, err := lam.Type.WriteTo(out)
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := fmt.Fprint(out, ") → ")
	if err != nil {
		return int64(w1) + int64(w2) + int64(w3), err
	}
	w4, err := lam.Body.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3) + int64(w4), err
}

func (pi *Pi) WriteTo(out io.Writer) (int64, error) {
	w1, err := fmt.Fprintf(out, "∀(%s : ", pi.Label)
	if err != nil {
		return int64(w1), err
	}
	w2, err := pi.Type.WriteTo(out)
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := fmt.Fprint(out, ") → ")
	if err != nil {
		return int64(w1) + int64(w2) + int64(w3), err
	}
	w4, err := pi.Body.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3) + int64(w4), err
}

func (app *App) WriteTo(out io.Writer) (int64, error) {
	w1, err := app.Fn.WriteTo(out)
	if err != nil {
		return int64(w1), err
	}
	w2, err := fmt.Fprint(out, " ")
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := app.Arg.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3), err
}

func (a Annot) WriteTo(out io.Writer) (int64, error) {
	w1, err := a.Expr.WriteTo(out)
	if err != nil {
		return int64(w1), err
	}
	w2, err := fmt.Fprint(out, " : ")
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := a.Annotation.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3), err
}

func (double) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "Double")
	return int64(n), err
}

func (d DoubleLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%f", d)
	return int64(n), err
}

func (boolean) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "Bool")
	return int64(n), err
}

func (bl BoolLit) WriteTo(out io.Writer) (int64, error) {
	if bool(bl) {
		n, err := fmt.Fprint(out, "True")
		return int64(n), err
	} else {
		n, err := fmt.Fprint(out, "True")
		return int64(n), err
	}
}

func (natural) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "Natural")
	return int64(n), err
}

func (nl NaturalLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%d", nl)
	return int64(n), err
}

func (p NaturalPlus) WriteTo(out io.Writer) (int64, error) {
	w1, err := p.L.WriteTo(out)
	if err != nil {
		return int64(w1), err
	}
	w2, err := fmt.Fprint(out, " + ")
	if err != nil {
		return int64(w1) + int64(w2), err
	}
	w3, err := p.R.WriteTo(out)
	return int64(w1) + int64(w2) + int64(w3), err
}

func (integer) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "Integer")
	return int64(n), err
}

func (i IntegerLit) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprintf(out, "%d", i)
	return int64(n), err
}

func (list) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "List")
	return int64(n), err
}

func (l EmptyList) WriteTo(out io.Writer) (int64, error) {
	n, err := fmt.Fprint(out, "[] : ")
	if err != nil {
		return int64(n), err
	}
	n2, err := l.Type.WriteTo(out)
	return int64(n) + int64(n2), err
}

func (l NonEmptyList) WriteTo(out io.Writer) (int64, error) {
	var written int64
	i, err := fmt.Fprint(out, "[ ")
	written += int64(i)
	if err != nil {
		return written, err
	}
	exprs := []Expr(l)
	i6, err := exprs[0].WriteTo(out)
	written += i6
	if err != nil {
		return written, err
	}
	for _, expr := range exprs[1:] {
		i, err = fmt.Fprint(out, ", ")
		written += int64(i)
		if err != nil {
			return written, err
		}
		i6, err := expr.WriteTo(out)
		written += i6
		if err != nil {
			return written, err
		}
	}
	i, err = fmt.Fprint(out, " ]")
	written += int64(i)
	return written, err
}

func (c Const) Normalize() Expr { return c }
func (v Var) Normalize() Expr   { return v }

func (lam *LambdaExpr) Normalize() Expr {
	return &LambdaExpr{
		Label: lam.Label,
		Type:  lam.Type.Normalize(),
		Body:  lam.Body.Normalize(),
	}
}
func (pi *Pi) Normalize() Expr {
	return &Pi{
		Label: pi.Label,
		Type:  pi.Type.Normalize(),
		Body:  pi.Body.Normalize(),
	}
}
func (app *App) Normalize() Expr {
	f := app.Fn.Normalize()
	a := app.Arg.Normalize()
	if l, ok := f.(*LambdaExpr); ok {
		v := Var{Name: l.Label}
		a2 := Shift(1, v, a)
		b1 := Subst(v, a2, l.Body)
		b2 := Shift(-1, v, b1)
		return b2.Normalize()
	}
	return app
}

func (a Annot) Normalize() Expr { return a.Expr.Normalize() }

func (d double) Normalize() Expr    { return d }
func (d DoubleLit) Normalize() Expr { return d }

func (n boolean) Normalize() Expr    { return n }
func (n BoolLit) Normalize() Expr    { return n }
func (n natural) Normalize() Expr    { return n }
func (n NaturalLit) Normalize() Expr { return n }
func (p NaturalPlus) Normalize() Expr {
	L := p.L.Normalize().(NaturalLit)
	R := p.R.Normalize().(NaturalLit)
	return NaturalLit(int(L) + int(R))
}

func (i integer) Normalize() Expr    { return i }
func (i IntegerLit) Normalize() Expr { return i }

func (l list) Normalize() Expr      { return l }
func (l EmptyList) Normalize() Expr { return EmptyList{l.Type.Normalize()} }
func (l NonEmptyList) Normalize() Expr {
	exprs := []Expr(l)
	vals := make([]Expr, len(exprs))
	for i, expr := range exprs {
		vals[i] = expr.Normalize()
	}
	return NonEmptyList(vals)
}

func NewLambdaExpr(arg string, argType Expr, body Expr) *LambdaExpr {
	return &LambdaExpr{
		Label: arg,
		Type:  argType,
		Body:  body,
	}
}
