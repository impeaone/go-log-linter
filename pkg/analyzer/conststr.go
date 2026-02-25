package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

func evalConstString(pass *analysis.Pass, expr ast.Expr) (s string, pos token.Pos, ok bool) {
	// Если types уже посчитали константу
	if tv, found := pass.TypesInfo.Types[expr]; found && tv.Value != nil {
		if tv.Value.Kind() == constant.String {
			return constant.StringVal(tv.Value), expr.Pos(), true
		}
	}

	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", 0, false
		}
		u, err := strconv.Unquote(e.Value)
		if err != nil {
			return "", 0, false
		}
		return u, e.Pos(), true

	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", 0, false
		}
		ls, _, lok := evalConstString(pass, e.X)
		if !lok {
			return "", 0, false
		}
		rs, _, rok := evalConstString(pass, e.Y)
		if !rok {
			return "", 0, false
		}
		return ls + rs, e.Pos(), true

	case *ast.ParenExpr:
		return evalConstString(pass, e.X)

	default:
		return "", 0, false
	}
}
