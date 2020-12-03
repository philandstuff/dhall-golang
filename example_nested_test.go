package dhall_test

import (
	"fmt"

	"github.com/philandstuff/dhall-golang/v6"
)

// NestedConfig is the struct we want to unmarshal from Dhall
type NestedConfig struct {
	Name     string
	DBConfig struct {
		Username string
		Password string
	}
}

const nestedDhallMessage = `
{ Name = "Alice", DBConfig = { Username = "foo", Password = "bar" } }
`

func Example_nested() {
	var m NestedConfig
	err := dhall.Unmarshal([]byte(nestedDhallMessage), &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", m)
	// Output:
	// {Name:Alice DBConfig:{Username:foo Password:bar}}
}
