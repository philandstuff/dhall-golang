package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/philandstuff/dhall-golang/binary"
	"github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/eval"
	"github.com/philandstuff/dhall-golang/imports"
	"github.com/philandstuff/dhall-golang/parser"
)

func main() {
	expr, err := parser.ParseReader("-", os.Stdin)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	resolvedExpr, err := imports.Load(expr.(core.Term))
	if err != nil {
		log.Fatalf("Import resolve error: %v", err)
	}
	inferredType, err := eval.TypeOf(resolvedExpr)
	if err != nil {
		log.Fatalf("Type error: %v", err)
	}
	fmt.Fprint(os.Stderr, inferredType)
	fmt.Fprintln(os.Stderr)
	fmt.Println(eval.AlphaBetaEval(resolvedExpr))

	var buf = new(bytes.Buffer)
	binary.EncodeAsCbor(buf, eval.Quote(eval.AlphaBetaEval(resolvedExpr)))
	final, err := binary.DecodeAsCbor(buf)
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}
	fmt.Printf("decoded as %+v\n", final)
}
