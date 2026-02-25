package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	linter, _ := NewAnalyzer(DefaultConfig())
	analysistest.Run(t, analysistest.TestData(), linter, "./...")
}
