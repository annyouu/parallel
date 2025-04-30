package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/format"
	"golang.org/x/tools/go/ast/astutil"
	"log"
)

func main() {
	expr , err := parser.ParseExpr(`v + 1`)
	if err != nil {
		log.Fatal(err)
	}

	// ASTを書き換える
	newExpr := astutil.Apply(expr, func(c *astutil.Cursor) bool {
		if ident, ok := c.Node().(*ast.Ident); ok && ident.Name == "v" {
			ident.Name = "x"
		}
		return true
	}, nil).(ast.Expr)

	// ASTからソースコードに変換する
	fset := token.NewFileSet()
	var buf bytes.Buffer
	err = format.Node(&buf, fset, newExpr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())
}