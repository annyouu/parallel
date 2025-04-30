package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"go/format"
)

func main() {
	// 位置情報用のファイルセットを作成する
	fset := token.NewFileSet()

	// 式: 1 + 2
	expr := &ast.BinaryExpr{
		X: &ast.BasicLit{Kind: token.INT, Value: "1"},
		Op: token.ADD,
		Y: &ast.BasicLit{Kind: token.INT, Value: "2"},
	}

	// 変数宣言: var x = 1 + 2
	decl := &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent("x")},
				Values: []ast.Expr{expr},
			},
		},
	}
	
	// ファイルASTの作成
	file := &ast.File{
		Name: ast.NewIdent("main"),
		Decls: []ast.Decl{decl},
	}

	// ソースコード出力
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}