package dhall_test

import (
	"fmt"
	"math"

	"github.com/philandstuff/dhall-golang/v4"
	"github.com/philandstuff/dhall-golang/v4/core"
	"github.com/pkg/errors"
)

type Shape interface {
	Area() float64
}

type Circle struct {
	Radius float64
}

type Rectangle struct {
	Width, Height float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// shapesMessage is the Dhall source we want to unmarshal
const shapesMessage = `
let Shape = < Circle : { Radius : Double } | Rectangle : { Width : Double, Height : Double } >
in [ Shape.Circle { Radius = 3.0 }, Shape.Rectangle { Width = 5.0, Height = 2.0 } ]
`

func decodeShape(d *dhall.Decoder, v core.Value) (interface{}, error) {
	if u, ok := v.(core.UnionVal); ok {
		if u.Alternative == "Circle" {
			var c Circle
			err := d.Decode(u.Val, &c)
			return c, err
		}
		if u.Alternative == "Rectangle" {
			var r Rectangle
			err := d.Decode(u.Val, &r)
			return r, err
		}
	}
	return nil, errors.New("Error decoding shape")
}

func Example_union() {
	var shapes []Shape
	decoder := dhall.NewDecoder()
	var shape Shape
	decoder.RegisterIface(&shape, decodeShape)
	err := decoder.Unmarshal([]byte(shapesMessage), &shapes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", shapes)
	// Output:
	// []dhall_test.Shape{dhall_test.Circle{Radius:3}, dhall_test.Rectangle{Width:5, Height:2}}
}
