package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"go/constant"
	"golang.org/x/tools/go/analysis"
)

type StrParts struct {
	ConstPrefix string // что слева
	HasDynamic  bool   // переменные/вызовы/не-константы
	Pos         token.Pos
}

// extractStringParts пытается извлечь константный префикс из выражения строки.
// Поддерживает: литералы, const-идентификаторы, "a"+x, x+"b", цепочки сложения.
func extractStringParts(pass *analysis.Pass, expr ast.Expr) StrParts {
	out := StrParts{Pos: expr.Pos()}

	// Если это compile-time константа строка — возвращаем её как префикс без динамики
	if tv, ok := pass.TypesInfo.Types[expr]; ok && tv.Value != nil && tv.Value.Kind() == constant.String {
		out.ConstPrefix = constant.StringVal(tv.Value)
		return out
	}

	switch e := expr.(type) {
	case *ast.BasicLit:
		out.HasDynamic = true
		return out

	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			out.HasDynamic = true
			return out
		}
		left := extractStringParts(pass, e.X)
		right := extractStringParts(pass, e.Y)

		// префикс = левый префикс + (если левый полностью константный) правый префикс
		out.ConstPrefix = left.ConstPrefix
		out.HasDynamic = left.HasDynamic || right.HasDynamic

		// Если слева не было динамики, то к префиксу можно добавить правый константный префикс
		if !left.HasDynamic {
			out.ConstPrefix += right.ConstPrefix
		}
		return out

	case *ast.ParenExpr:
		return extractStringParts(pass, e.X)

	case *ast.Ident:
		out.HasDynamic = true
		return out

	default:
		out.HasDynamic = true
		return out
	}
}

// containsSensitivePrefix — проверка "опасного" префикса.
func containsSensitivePrefix(prefix string, sensSet map[string]struct{}) bool {
	p := strings.ToLower(prefix)
	for kw := range sensSet {
		if strings.Contains(p, kw) {
			return true
		}
	}
	return false
}
