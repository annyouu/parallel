package callparallel

import (
	"go/ast"
	"strings"
	"fmt"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/astutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "parallelInParent",
	Doc:  "サブテストでt.Parallel()を呼び出しているのに親テストで呼び出していないケースを探す",
	Run:  run,
}

// テスト関数かどうかを判定
func isTestFunc(name string) bool {
	return strings.HasPrefix(name, "Test") && len(name) > 4 && name[4] >= 'A' && name[4] <= 'Z'
}

// 関数内で t.Parallel() を呼んでいるかをチェック
func containsParallelCall(body *ast.BlockStmt, _ *analysis.Pass) bool {
	if body == nil {
		return false
	}
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Parallel" {
					fmt.Println("t.Parallel()が見つかりました")
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

func run(pass *analysis.Pass) (interface{}, error) {
	// fmt.Println("アナライザが実行されました")

	parentHas := make(map[string]bool)

	// ① トップレベルのテスト関数で t.Parallel() を呼んでいるかチェック
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || !isTestFunc(fn.Name.Name) {
				continue
			}
			parentHas[fn.Name.Name] = containsParallelCall(fn.Body, pass)
			fmt.Printf("親テスト関数: %s, Parallel呼び出し: %v\n", fn.Name.Name, parentHas[fn.Name.Name])
		}
	}

	// ② t.Run を呼んでいる箇所を探索
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "Run" {
				return true
			}

			if len(call.Args) < 2 {
				return true
			}

			lit, ok := call.Args[1].(*ast.FuncLit)
			if !ok {
				return true
			}

			// 親関数名を探索
			path, _ := astutil.PathEnclosingInterval(f, call.Pos(), call.End())
			var parentTest string
			for _, node := range path {
				if fn, ok := node.(*ast.FuncDecl); ok && isTestFunc(fn.Name.Name) {
					parentTest = fn.Name.Name
					fmt.Println("親関数:", parentTest)
					break
				}
			}
			if parentTest == "" {
				return true
			}

			if containsParallelCall(lit.Body, pass) && !parentHas[parentTest] {
				fmt.Printf("エラー: サブテストでParallel()が呼ばれていますが、親テスト%qでParallel()が呼ばれていません\n",
				parentTest)
			} else {
				fmt.Printf("親テスト%qで Parallel()が呼ばれています\n", parentTest)
			}
			return true
		})
	}
	// fmt.Println("アナライザが終了しました")
	return nil, nil
}
