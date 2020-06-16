package dhall_test

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/v4"
)

// TaggedMessage is the struct we want to unmarshal from Dhall
type TaggedMessage struct {
	Name string `dhall:"name"`
	Body string `dhall:"entity"`
	Time int64  `dhall:"instant"`
}

// dhallTaggedMessage is the Dhall source we want to unmarshal
const dhallTaggedMessage = `
{ name = "Alice", entity = "Hello", instant = 1294706395881547000 }
`

func Example_tagged() {
	var m TaggedMessage
	err := dhall.Unmarshal([]byte(dhallTaggedMessage), &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", m)
	// Output:
	// {Name:Alice Body:Hello Time:1294706395881547000}
}
