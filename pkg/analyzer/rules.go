package analyzer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Violation string

const (
	VLowercase Violation = "log message must start with a lowercase letter"
	VEnglish   Violation = "log message must be in English (non-ASCII characters are not allowed)"
	VSpecial   Violation = "log message must not contain special characters or emojis"
	VSensitive Violation = "log message must not contain potentially sensitive data"
)

func CheckMessageAll(msg string) []Violation {
	if msg == "" {
		return nil
	}

	cc := getCompiledConfig()

	if cc == nil {
		return nil
	}

	c := cc.raw
	var out []Violation

	if c.RequireLowercaseStart && !startsWithLower(msg) {
		out = append(out, VLowercase)
	}

	switch c.EnglishMode {
	case EnglishASCII:
		if hasNonASCII(msg) {
			out = append(out, VEnglish)
		}
	case EnglishLatin:
		if hasNonASCII(msg) {
			out = append(out, VEnglish)
		}
	}

	if c.ForbidSpecialChars && cc.allowedRe != nil && !cc.allowedRe.MatchString(msg) {
		out = append(out, VSpecial)
	}

	// Sensitive
	if c.ForbidSensitive && len(cc.sensSet) > 0 && hasSensitive(msg, cc.sensSet) {
		out = append(out, VSensitive)
	}

	return out
}

func startsWithLower(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsLower(r)
}

func hasNonASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

func hasSensitive(s string, sensSet map[string]struct{}) bool {
	lower := strings.ToLower(s)
	for kw := range sensSet {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func CheckMessageAllWith(cc *compiledConfig, msg string) []Violation {
	if msg == "" || cc == nil {
		return nil
	}
	c := cc.raw
	var out []Violation

	if c.RequireLowercaseStart && !startsWithLower(msg) {
		out = append(out, VLowercase)
	}

	switch c.EnglishMode {
	case EnglishASCII, EnglishLatin:
		if hasNonASCII(msg) {
			out = append(out, VEnglish)
		}
	}

	if c.ForbidSpecialChars && cc.allowedRe != nil && !cc.allowedRe.MatchString(msg) {
		out = append(out, VSpecial)
	}

	if c.ForbidSensitive && len(cc.sensSet) > 0 && hasSensitive(msg, cc.sensSet) {
		out = append(out, VSensitive)
	}

	return out
}
