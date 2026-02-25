package plugin

import (
	"fmt"
	"strings"

	"github.com/impeaone/go-log-linter/pkg/analyzer"
)

func parseConfig(conf any) (analyzer.Config, error) {
	cfg := analyzer.DefaultConfig()
	if conf == nil {
		return cfg, nil
	}

	asMap, ok := conf.(map[string]any)
	if !ok {
		return cfg, fmt.Errorf("settings must be a map, got %T", conf)
	}

	if inner, ok := asMap["settings"].(map[string]any); ok {
		asMap = inner
	}

	if inner, ok := asMap["loglinter"].(map[string]any); ok {
		if s, ok := inner["settings"].(map[string]any); ok {
			asMap = s
		}
	}

	m := make(map[string]any, len(asMap))
	for k, v := range asMap {
		m[strings.ToLower(k)] = v
	}

	// bool
	if v, exists := m["requirelowercasestart"]; exists {
		b, ok := v.(bool)
		if !ok {
			return cfg, fmt.Errorf("requireLowercaseStart must be bool, got %T", v)
		}
		cfg.RequireLowercaseStart = b
	}
	if v, exists := m["forbidspecialchars"]; exists {
		b, ok := v.(bool)
		if !ok {
			return cfg, fmt.Errorf("forbidSpecialChars must be bool, got %T", v)
		}
		cfg.ForbidSpecialChars = b
	}
	if v, exists := m["forbidsensitive"]; exists {
		b, ok := v.(bool)
		if !ok {
			return cfg, fmt.Errorf("forbidSensitive must be bool, got %T", v)
		}
		cfg.ForbidSensitive = b
	}

	// string
	if v, exists := m["englishmode"]; exists {
		s, ok := v.(string)
		if !ok {
			return cfg, fmt.Errorf("englishMode must be string, got %T", v)
		}
		cfg.EnglishMode = analyzer.EnglishMode(s)
	}
	if v, exists := m["allowedcharsregex"]; exists {
		s, ok := v.(string)
		if !ok {
			return cfg, fmt.Errorf("allowedCharsRegex must be string, got %T", v)
		}
		cfg.AllowedCharsRegex = s
	}

	// []string
	if v, exists := m["sensitivekeywords"]; exists {
		ss, err := toStringSlice(v)
		if err != nil {
			return cfg, fmt.Errorf("sensitiveKeywords: %w", err)
		}
		cfg.SensitiveKeywords = ss
	}

	return cfg, nil
}

func toStringSlice(v any) ([]string, error) {
	switch x := v.(type) {
	case []string:
		return x, nil
	case []any:
		out := make([]string, 0, len(x))
		for i, it := range x {
			s, ok := it.(string)
			if !ok {
				return nil, fmt.Errorf("element %d must be string, got %T", i, it)
			}
			out = append(out, s)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("must be array of strings, got %T", v)
	}
}
