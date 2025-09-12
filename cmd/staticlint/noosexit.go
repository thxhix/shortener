package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// noOsExitAnalyzer reports direct calls to os.Exit inside main() function of package main.
var noOsExitAnalyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "forbids direct calls to os.Exit in main.main; return an error or use log.Fatal in main instead",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: runNoOsExit,
}

func runNoOsExit(pass *analysis.Pass) (any, error) {
	// Only care about package main
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.CallExpr)(nil),
	}

	type fnFrame struct {
		inMain bool
	}

	var stack []fnFrame

	ins.Nodes(nodeFilter, func(n ast.Node, push bool) bool {
		switch node := n.(type) {

		case *ast.FuncDecl:
			if push {
				inMain := node.Name.Name == "main" && node.Recv == nil
				stack = append(stack, fnFrame{inMain: inMain})
			} else {
				stack = stack[:len(stack)-1]
			}
			return true

		case *ast.CallExpr:
			// Only check when we're inside main()
			if len(stack) == 0 || !stack[len(stack)-1].inMain {
				return true
			}
			// Is callee a selector (pkg.Func)?
			sel, ok := node.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			// Resolve to types.Func and check package path + name
			if obj, ok := pass.TypesInfo.Uses[sel.Sel].(*types.Func); ok {
				if obj.Name() == "Exit" && obj.Pkg() != nil && obj.Pkg().Path() == "os" {
					pass.Reportf(node.Lparen, "do not call os.Exit directly in main; return an error from run() and handle it in main, e.g., log.Fatal(err)")
				}
			}
		}
		return true
	})

	return nil, nil
}
