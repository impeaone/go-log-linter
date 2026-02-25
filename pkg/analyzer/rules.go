package analyzer

import (
	"regexp"
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

var (
	// Разрешаем только этот набор символов.
	allowed = regexp.MustCompile(`^[a-zA-Z0-9 ,.\-:'?]+$`)

	sensitiveKeywords = []string{
		"password", "passwd",
		"secret",
		"api_key", "apikey",
		"token", "access_token", "auth_token",
		"credential",
	}
)

func CheckMessageAll(msg string) []Violation {
	if msg == "" {
		return nil
	}

	var out []Violation

	// lowercase
	if !startsWithLower(msg) {
		out = append(out, VLowercase)
	}

	// english (ASCII-only)
	if hasNonASCII(msg) {
		out = append(out, VEnglish)
	}

	// special chars / emojis
	if !allowed.MatchString(msg) {
		out = append(out, VSpecial)
	}

	// sensitive
	if hasSensitiveKeyword(msg) {
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

func hasSensitiveKeyword(s string) bool {
	lower := strings.ToLower(s)
	for _, kw := range sensitiveKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
