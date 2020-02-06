package dhall_test

import (
	"fmt"

	"github.com/philandstuff/dhall-golang"
)

// TaggedMessage is the struct we want to unmarshal from Dhall
type TaggedMessage struct {
	Name string `json:name`
	Body string `json:body`
	Time int64  `json:instant`
}

// dhallTaggedMessage is the Dhall source we want to unmarshal
const dhallTaggedMessage = `
{ name = "Alice", body = "Hello", instant = 1294706395881547000 }
`

func Example_tagged() {
	var m Message
	err := dhall.Unmarshal([]byte(dhallMessage), &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", m)
	// Output:
	// {Name:Alice Body:Hello Time:1294706395881547000}
}
