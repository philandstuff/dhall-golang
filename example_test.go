package dhall_test

import (
	"fmt"

	"github.com/philandstuff/dhall-golang"
)

// Message is the struct we want to unmarshal from Dhall
type Message struct {
	Name string
	Body string
	Time int64
}

func Example() {
	var m Message
	err := dhall.Unmarshal([]byte(
		`{Name = "Alice",Body = "Hello", Time = 1294706395881547000}`,
	), &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", m)
	// Output:
	// {Name:Alice Body:Hello Time:1294706395881547000}
}
