package main

import (
	"fmt"
	"log"
	"os"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"
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
	resolvedExpr.Normalize().WriteTo(os.Stdout)
}
