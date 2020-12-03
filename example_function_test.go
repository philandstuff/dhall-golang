package dhall_test

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/v6"
)

// Config is the struct we want to unmarshal from Dhall
type Config struct {
	Name string
	// Dhall lets you unmarshal functions as well as data
	Greet func(string) string
}

// configMessage is the Dhall source we want to unmarshal
const configMessage = `
{ Name = "Alice", Greet = λ(name : Text) → "Howdy, ${name}!" }
`

func Example_function() {
	// you can also unmarshal Dhall functions to Go functions
	var greet func(string) string
	err := dhall.Unmarshal([]byte(`λ(name : Text) → "Howdy, ${name}!"`), &greet)
	if err != nil {
		panic(err)
	}
	fmt.Println(greet("Alice"))
	// Output:
	// Howdy, Alice!
}
