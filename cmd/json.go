package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/philandstuff/dhall-golang/v5"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

func cmdJSON(c *cli.Context) error {
	var data interface{}
	err := dhall.UnmarshalReader("-", os.Stdin, &data)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}
