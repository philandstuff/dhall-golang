package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/philandstuff/dhall-golang/v4/binary"
	"github.com/philandstuff/dhall-golang/v4/core"
	"github.com/philandstuff/dhall-golang/v4/imports"
	"github.com/philandstuff/dhall-golang/v4/parser"
)

func main() {
	expr, err := parser.ParseReader("-", os.Stdin)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	resolvedExpr, err := imports.Load(expr)
	if err != nil {
		log.Fatalf("Import resolve error: %v", err)
	}
	inferredType, err := core.TypeOf(resolvedExpr)
	if err != nil {
		log.Fatalf("Type error: %v", err)
	}
	fmt.Fprint(os.Stderr, inferredType)
	fmt.Fprintln(os.Stderr)
	fmt.Println(core.Eval(resolvedExpr))

	var buf = new(bytes.Buffer)
	binary.EncodeAsCbor(buf, core.QuoteAlphaNormal(core.Eval(resolvedExpr)))
	final, err := binary.DecodeAsCbor(buf)
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}
	fmt.Printf("decoded as %+v\n", final)
}
