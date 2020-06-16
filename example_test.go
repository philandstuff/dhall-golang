package dhall_test

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/v4"
)

// Message is the struct we want to unmarshal from Dhall
type Message struct {
	Name string
	Body string
	Time int64
}

// dhallMessage is the Dhall source we want to unmarshal
const dhallMessage = `
{ Name = "Alice", Body = "Hello", Time = 1294706395881547000 }
`

func Example() {
	var m Message
	err := dhall.Unmarshal([]byte(dhallMessage), &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", m)
	// Output:
	// {Name:Alice Body:Hello Time:1294706395881547000}
}
