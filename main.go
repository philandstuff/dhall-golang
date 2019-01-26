package main

import (
	"fmt"
	"log"
	"os"

	"github.com/philandstuff/dhall-golang/ast"
	"github.com/philandstuff/dhall-golang/parser"
)

func load(e ast.Expr) ast.Expr { return e }

func main() {
	text := []byte("Sort")
	expr := parser.ParseExpression(text)
	resolvedExpr := load(expr)
	inferredType, err := resolvedExpr.TypeWith(ast.EmptyContext())
	if err != nil {
		log.Fatalf("Type error: %v", err)
	}
	inferredType.WriteTo(os.Stderr)
	fmt.Fprintln(os.Stderr)
	resolvedExpr.Normalize().WriteTo(os.Stdout)
}
