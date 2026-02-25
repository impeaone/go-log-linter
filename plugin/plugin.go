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
	cfg, err := parseConfig(conf)
	if err != nil {
		return nil, err
	}
	return &plugin{cfg: cfg}, nil
}

type plugin struct {
	cfg analyzer.Config
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	a, err := analyzer.NewAnalyzer(p.cfg)
	if err != nil {
		return nil, err
	}
	return []*analysis.Analyzer{a}, nil
}

func (*plugin) GetLoadMode() string { return register.LoadModeTypesInfo }
