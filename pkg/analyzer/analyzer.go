package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	Name = "loglinter"
	Doc  = "Checks log messages against predefined rules"
)

var LogLinter = &analysis.Analyzer{
	Name:     Name,
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	ins.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		msgArgIdx, ok := isSupportedLoggerCall(pass, call)
		if !ok || msgArgIdx < 0 || msgArgIdx >= len(call.Args) {
			return
		}

		basicLit, ok := call.Args[msgArgIdx].(*ast.BasicLit)
		if !ok || basicLit.Kind != token.STRING {
			return
		}

		msg, err := strconv.Unquote(basicLit.Value)
		if err != nil {
			return
		}

		violations := CheckMessageAll(msg)
		for _, v := range violations {
			pass.Reportf(basicLit.Pos(), "%s", string(v))
		}
	})

	return nil, nil
}
