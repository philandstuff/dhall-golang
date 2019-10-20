package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/philandstuff/dhall-golang/core"
	"github.com/philandstuff/dhall-golang/eval"
	"github.com/philandstuff/dhall-golang/imports"
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/ugorji/go/codec"
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
	fmt.Println(eval.Eval(resolvedExpr))
	var ch codec.CborHandle
	var buf = new(bytes.Buffer)
	enc := codec.NewEncoder(buf, &ch)
	dec := codec.NewDecoder(buf, &ch)
	enc.Encode(eval.Eval(resolvedExpr))
	var final interface{}
	dec.Decode(&final)
	fmt.Printf("%+v\n", final)
}
