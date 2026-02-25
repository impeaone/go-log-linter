package analyzer

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

// isSupportedLoggerCall определяет, является ли вызов логгером slog/zap
// и возвращает индекс аргумента с сообщением.
func isSupportedLoggerCall(pass *analysis.Pass, call *ast.CallExpr) (msgArgIdx int, ok bool) {
	fn := typeutil.StaticCallee(pass.TypesInfo, call)
	if fn == nil {
		// не статический вызов (может быть интерфейс, функция‑значение)
		return -1, false
	}

	obj := fn
	sig, _ := obj.Type().(*types.Signature)
	if sig == nil {
		return -1, false
	}

	recv := sig.Recv()
	// Методы логгеров (logger.Info, sugar.Infof ...)
	if recv != nil {
		if isSlogLogger(recv.Type()) && isSlogMethod(obj.Name()) {
			// slog.Logger.Info(msg string, ...)
			return 0, true
		}
		if isZapLogger(recv.Type()) && isZapLikeMethod(obj.Name()) {
			// *zap.Logger.Info(msg string, ...)
			return 0, true
		}
		if isZapSugaredLogger(recv.Type()) && isZapLikeMethod(obj.Name()) {
			// *zap.SugaredLogger.Infof(msg string, ...)
			return 0, true
		}
		return -1, false
	}

	// Функции пакета slog: slog.Info, slog.Error, …
	if pkg := obj.Pkg(); pkg != nil && pkg.Path() == "log/slog" {
		if isSlogTopLevelFunc(obj.Name()) {
			// slog.Info(msg string, ...)
			return 0, true
		}
	}

	return -1, false
}

func isSlogLogger(t types.Type) bool {
	// ожидается *log/slog.Logger или log/slog.Logger
	const slogLoggerPath = "log/slog"
	const slogLoggerName = "Logger"

	named, ok := derefNamed(t)
	if !ok {
		return false
	}
	if named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == slogLoggerPath && named.Obj().Name() == slogLoggerName
}

func isSlogMethod(name string) bool {
	switch name {
	case "Info", "Error", "Warn", "Debug", "Log", "LogAttrs":
		return true
	default:
		return false
	}
}

func isSlogTopLevelFunc(name string) bool {
	// функции пакета slog, которые логируют сообщение
	switch name {
	case "Info", "Error", "Warn", "Debug":
		return true
	default:
		return false
	}
}

func isZapLogger(t types.Type) bool {
	// go.uber.org/zap.Logger
	const zapPkgPath = "go.uber.org/zap"
	const zapLoggerName = "Logger"

	named, ok := derefNamed(t)
	if !ok {
		return false
	}
	if named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == zapPkgPath && named.Obj().Name() == zapLoggerName
}

func isZapSugaredLogger(t types.Type) bool {
	// go.uber.org/zap.SugaredLogger
	const zapPkgPath = "go.uber.org/zap"
	const zapSugarName = "SugaredLogger"

	named, ok := derefNamed(t)
	if !ok {
		return false
	}
	if named.Obj().Pkg() == nil {
		return false
	}
	return named.Obj().Pkg().Path() == zapPkgPath && named.Obj().Name() == zapSugarName
}

func isZapLikeMethod(name string) bool {
	switch name {
	case "Debug", "Debugf", "Debugw",
		"Info", "Infof", "Infow",
		"Warn", "Warnf", "Warnw",
		"Error", "Errorf", "Errorw",
		"Panic", "Panicf", "Panicw",
		"Fatal", "Fatalf", "Fatalw":
		return true
	default:
		return false
	}
}

// derefNamed снимает все * и возвращает Named‑тип, если он есть.
func derefNamed(t types.Type) (*types.Named, bool) {
	for {
		switch tt := t.(type) {
		case *types.Pointer:
			t = tt.Elem()
		case *types.Named:
			return tt, true
		default:
			return nil, false
		}
	}
}
