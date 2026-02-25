package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	Name = "loglinter"
	Doc  = "Checks log messages against predefined rules"
)

func NewAnalyzer(cfg Config) (*analysis.Analyzer, error) {
	cc, err := compileConfig(cfg)
	if err != nil {
		return nil, err
	}

	a := &analysis.Analyzer{
		Name:     Name,
		Doc:      Doc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	a.Run = func(pass *analysis.Pass) (any, error) {
		return runWithConfig(pass, cc)
	}
	return a, nil
}

func runWithConfig(pass *analysis.Pass, cc *compiledConfig) (any, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	ins.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		msgArgIdx, ok := isSupportedLoggerCall(pass, call)
		if !ok || msgArgIdx < 0 || msgArgIdx >= len(call.Args) {
			return
		}

		expr := call.Args[msgArgIdx]

		if msg, pos, ok := evalConstString(pass, expr); ok {
			for _, v := range CheckMessageAllWith(cc, msg) {
				pass.Reportf(pos, "%s", string(v))
			}
			return
		}

		parts := extractStringParts(pass, expr)

		if parts.HasDynamic && cc != nil && cc.raw.ForbidSensitive && len(cc.sensSet) > 0 {
			if containsSensitivePrefix(parts.ConstPrefix, cc.sensSet) {
				pass.Reportf(parts.Pos, "%s", string(VSensitive))
			}
		}

		if parts.ConstPrefix != "" {
			for _, v := range CheckMessageAllWith(cc, parts.ConstPrefix) {
				pass.Reportf(parts.Pos, "%s", string(v))
			}
		}
	})

	return nil, nil
}
