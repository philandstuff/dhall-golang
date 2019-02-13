package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"
	"github.com/ugorji/go/codec"
)

//go:generate pigeon -o parser/dhall.go parser/dhall.peg

func load(e ast.Expr) ast.Expr { return e }

func main() {
	expr, err := parser.ParseReader("-", os.Stdin)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}
	resolvedExpr := load(expr.(ast.Expr))
	inferredType, err := resolvedExpr.TypeWith(ast.EmptyContext())
	if err != nil {
		log.Fatalf("Type error: %v", err)
	}
	inferredType.WriteTo(os.Stderr)
	fmt.Fprintln(os.Stderr)
	//resolvedExpr.Normalize().WriteTo(os.Stdout)
	var ch codec.CborHandle
	var buf = new(bytes.Buffer)
	enc := codec.NewEncoder(buf, &ch)
	dec := codec.NewDecoder(buf, &ch)
	enc.Encode(resolvedExpr.Normalize())
	var final interface{}
	dec.Decode(&final)
	fmt.Printf("%+v\n", final)
}
