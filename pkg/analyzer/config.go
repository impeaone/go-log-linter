package analyzer

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type EnglishMode string

const (
	EnglishASCII EnglishMode = "ascii" // только ASCII
	EnglishLatin EnglishMode = "latin" // ASCII + сообщения обязаны быть латиницей/цифрами/пробелами/разрешенной пунктуацией
)

type Config struct {
	RequireLowercaseStart bool        `json:"requireLowercaseStart" yaml:"requireLowercaseStart"`
	EnglishMode           EnglishMode `json:"englishMode" yaml:"englishMode"`

	ForbidSpecialChars bool   `json:"forbidSpecialChars" yaml:"forbidSpecialChars"`
	AllowedCharsRegex  string `json:"allowedCharsRegex" yaml:"allowedCharsRegex"`

	ForbidSensitive   bool     `json:"forbidSensitive" yaml:"forbidSensitive"`
	SensitiveKeywords []string `json:"sensitiveKeywords" yaml:"sensitiveKeywords"`
	SensitivePatterns []string `json:"sensitivePatterns" yaml:"sensitivePatterns"`
}

type compiledConfig struct {
	raw Config

	allowedRe *regexp.Regexp
	sensSet   map[string]struct{}
	sensRes   []*regexp.Regexp
}

func DefaultConfig() Config {
	return Config{
		RequireLowercaseStart: true,
		EnglishMode:           EnglishASCII,

		ForbidSpecialChars: true,
		AllowedCharsRegex:  `^[a-zA-Z0-9 ,.:?'_-]+$`,

		ForbidSensitive: true,
		SensitiveKeywords: []string{
			"password", "passwd",
			"secret",
			"api_key", "apikey",
			"token", "access_token", "auth_token",
			"credential",
		},
	}
}

func compileConfig(c Config) (*compiledConfig, error) {
	// defaults
	if c.EnglishMode == "" {
		c.EnglishMode = EnglishASCII
	}

	if c.EnglishMode != EnglishASCII && c.EnglishMode != EnglishLatin {
		return nil, fmt.Errorf("englishMode: unsupported value %q (supported: %q, %q)", c.EnglishMode, EnglishASCII, EnglishLatin)
	}

	cc := &compiledConfig{raw: c}

	if c.ForbidSpecialChars {
		if c.AllowedCharsRegex == "" {
			return nil, fmt.Errorf("allowedCharsRegex: must be non-empty when forbidSpecialChars=true")
		}
		re, err := regexp.Compile(c.AllowedCharsRegex)
		if err != nil {
			return nil, fmt.Errorf("allowedCharsRegex: invalid regexp: %w", err)
		}
		cc.allowedRe = re
	}

	if c.ForbidSensitive {
		// sensitive вручную прописанные
		cc.sensSet = make(map[string]struct{}, len(c.SensitiveKeywords))
		for _, kw := range c.SensitiveKeywords {
			kw = strings.TrimSpace(strings.ToLower(kw))
			if kw == "" {
				continue
			}
			cc.sensSet[kw] = struct{}{}
		}

		// sensitive по патерну
		if len(c.SensitivePatterns) > 0 {
			cc.sensRes = make([]*regexp.Regexp, 0, len(c.SensitivePatterns))
			for i, p := range c.SensitivePatterns {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				re, err := regexp.Compile(p)
				if err != nil {
					return nil, fmt.Errorf("sensitivePatterns[%d]: invalid regexp: %w", i, err)
				}
				cc.sensRes = append(cc.sensRes, re)
			}
		}
	}

	return cc, nil
}

// TODO: пересмотреть вот этот момент
var (
	cfgMu sync.RWMutex
	cfgC  *compiledConfig
)

func SetConfig(c Config) error {
	cc, err := compileConfig(c)
	if err != nil {
		return err
	}

	cfgMu.Lock()
	cfgC = cc
	cfgMu.Unlock()
	return nil
}

func getCompiledConfig() *compiledConfig {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfgC
}
