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

// テスト関数がどうかを判定するメソッド
func isTestFunc(name string) bool {
	return strings.HasPrefix(name, "Test") && len(name) > 4 && name[4] >= 'A' && name[4] <= 'Z'
}

// 親テスト関数内でt.Parallel()を呼んでいるかをチェックする(ネストした無名関数は飛ばす)
func parentHasParallel(body *ast.BlockStmt) bool {
	if body == nil {
		return false
	}
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncLit:
			// 無名関数内はスキップ
			return false
		case *ast.CallExpr:
			if sel, ok := n.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Parallel" {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// 関数内でt.Parallel()を呼んでいるかをチェックする(サブテストも含めて全て探索する)
func containsParallelCall(body *ast.BlockStmt, _ *analysis.Pass) bool {
	if body == nil {
		return false
	}
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Parallel" {
					// fmt.Println("t.Parallel()が見つかりました")
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
	parentHas := make(map[string]bool)

	// ① トップレベルのテスト関数でt.parallel()を呼んでいるかチェックする
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || !isTestFunc(fn.Name.Name) {
				continue
			}

			// 親テストではネストした無名関数をスキップする
			parentHas[fn.Name.Name] = parentHasParallel(fn.Body)
			// fmt.Printf("親テスト関数: %s, Parallel呼び出し: %v\n", fn.Name.Name, parentHas[fn.Name.Name])
		}
	}

	// ② t.Runを呼んでいる箇所を探索
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

			// 親関数名を探索する
			path, _ := astutil.PathEnclosingInterval(f, call.Pos(), call.End())
			var parentTest string
			for _, node := range path {
				if fn, ok := node.(*ast.FuncDecl); ok && isTestFunc(fn.Name.Name) {
					parentTest = fn.Name.Name
					// fmt.Println("親関数:", parentTest)
					break
				}
			}
			
			if parentTest == "" {
				return true
			}

			if containsParallelCall(lit.Body, pass) && !parentHas[parentTest] {
				// fmt.Printf("エラー: サブテストでParallelが呼ばれていますか、親テスト%qでParallel()が呼ばれていません\n", parentTest)
				pass.Reportf(call.Pos(), "エラー: サブテストでParallelが呼ばれていますが、親テスト%qでParallel()が呼ばれていません", parentTest)
			} else {
				fmt.Printf("親テスト%qでParallel()が呼ばれています\n", parentTest)
			}
			
			return true
		})
	}
	return nil, nil
}