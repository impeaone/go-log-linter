package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/impeaone/go-log-linter/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglinter", New)
}

func New(conf any) (register.LinterPlugin, error) {
	return &plugin{}, nil
}

type plugin struct{}

var _ register.LinterPlugin = (*plugin)(nil)

func (*plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.LogLinter,
	}, nil
}

func (*plugin) GetLoadMode() string {
	// Если ты используешь TypesInfo (а ты используешь для определения slog/zap),
	// нужен LoadModeTypesInfo.
	return register.LoadModeTypesInfo
}
