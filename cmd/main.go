package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/philandstuff/dhall-golang/v5/binary"
	"github.com/philandstuff/dhall-golang/v5/core"
	"github.com/philandstuff/dhall-golang/v5/imports"
	"github.com/philandstuff/dhall-golang/v5/parser"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

func main() {
	app := &cli.App{
		Name:  "dhall-golang",
		Usage: "Dhall implemented in Go",
		Commands: []*cli.Command{
			{
				Name:   "json",
				Usage:  "output Dhall code as JSON",
				Action: cmdJSON,
			},
			{
				Name:   "yaml",
				Usage:  "output Dhall code as YAML",
				Action: cmdYAML,
			},
		},
		Action: cmdDebug,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// cmdDebug is the original scrappy debug command
func cmdDebug(c *cli.Context) error {
	expr, err := parser.ParseReader("-", os.Stdin)
	if err != nil {
		return err
	}
	resolvedExpr, err := imports.Load(expr)
	if err != nil {
		return err
	}
	inferredType, err := core.TypeOf(resolvedExpr)
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stderr, inferredType)
	fmt.Fprintln(os.Stderr)
	fmt.Println(core.Eval(resolvedExpr))

	var buf = new(bytes.Buffer)
	binary.EncodeAsCbor(buf, core.QuoteAlphaNormal(core.Eval(resolvedExpr)))
	final, err := binary.DecodeAsCbor(buf)
	if err != nil {
		return err
	}
	fmt.Printf("decoded as %+v\n", final)
	return nil
}
