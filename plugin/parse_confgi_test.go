package plugin

import (
	"testing"

	"github.com/impeaone/go-log-linter/pkg/analyzer"
)

func TestParseConfig_DefaultOnNil(t *testing.T) {
	cfg, err := parseConfig(nil)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	def := analyzer.DefaultConfig()
	if cfg.EnglishMode != def.EnglishMode {
		t.Fatalf("englishMode: got %q want %q", cfg.EnglishMode, def.EnglishMode)
	}
}

func TestParseConfig_CustomValues(t *testing.T) {
	in := map[string]any{
		"requireLowercaseStart": false,
		"englishMode":           "latin",
		"forbidSpecialChars":    true,
		"allowedCharsRegex":     "^[a-z]+$",
		"forbidSensitive":       true,
		"sensitiveKeywords":     []any{"token", "api_key"},
		"allowUnderscore":       true,
		"allowEquals":           true,
	}

	cfg, err := parseConfig(in)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.RequireLowercaseStart != false {
		t.Fatalf("RequireLowercaseStart: got true want false")
	}
	if cfg.EnglishMode != analyzer.EnglishMode("latin") {
		t.Fatalf("EnglishMode: got %q", cfg.EnglishMode)
	}
	if cfg.AllowedCharsRegex != "^[a-z]+$" {
		t.Fatalf("AllowedCharsRegex: got %q", cfg.AllowedCharsRegex)
	}
	if len(cfg.SensitiveKeywords) != 2 {
		t.Fatalf("SensitiveKeywords len: got %d", len(cfg.SensitiveKeywords))
	}
}

func TestParseConfig_TypeErrors(t *testing.T) {
	_, err := parseConfig(map[string]any{
		"forbidSensitive": "yes",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
