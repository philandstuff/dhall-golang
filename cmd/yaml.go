package main

import (
	"fmt"
	"os"

	"github.com/philandstuff/dhall-golang/v5"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"gopkg.in/yaml.v2"
)

func cmdYAML(c *cli.Context) error {
	var data interface{}
	err := dhall.UnmarshalReader("-", os.Stdin, &data)
	if err != nil {
		return err
	}
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}
