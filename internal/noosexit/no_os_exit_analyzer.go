package noosexit

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// NoOsExitAnalyzer - forbids os.Exit direct usage in main function
var Analyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "checks against os.Exit calls in main()",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// look for "main" package
		if file.Name.Name != "main" {
			continue
		}

		// search for main function in current file
		var mainFunc ast.Node
		for _, d := range file.Decls {
			f, ok := d.(*ast.FuncDecl)
			if ok && f.Name.Name == "main" {
				mainFunc = d
			}
		}

		// main function not found
		if mainFunc == nil {
			continue
		}

		ast.Inspect(mainFunc, func(node ast.Node) bool {
			c, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			s, ok := c.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			i, ok := s.X.(*ast.Ident)
			if !ok || i.Name != "os" {
				return true
			}

			if s.Sel.Name == "Exit" {
				pass.Reportf(c.Pos(), "os.Exit calls are prohibited in main()")
			}

			return true
		})
	}
	return nil, nil
}
