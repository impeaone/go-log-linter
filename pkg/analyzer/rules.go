package analyzer

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Violation - просто alias для string: строка, содержащая допущенные нарушения
type Violation string

// Вывод на ошибки линтинга строки
const (
	VLowercase Violation = "log message must start with a lowercase letter"
	VEnglish   Violation = "log message must be in English (non-ASCII characters are not allowed)"
	VSpecial   Violation = "log message must not contain special characters or emojis"
	VSensitive Violation = "log message must not contain potentially sensitive data"
)

// CheckMessageAllWith - проверяет строку на установленные правила
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

	if c.ForbidSensitive && ((len(cc.sensSet) > 0 && hasSensitive(msg, cc.sensSet)) ||
		(len(cc.sensRes) > 0 && hasSensitiveByRegex(msg, cc.sensRes))) {

		out = append(out, VSensitive)
	}

	return out
}

// startsWithLower - проверяет, что строка начинается со строчной буквы
func startsWithLower(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsLower(r)
}

// hasNonASCII - проверяет, что содержатся только ASCII символы
func hasNonASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

// hasSensitive - проверка на чувствительные данные (для ручно-введенных в yaml слов)
func hasSensitive(s string, sensSet map[string]struct{}) bool {
	lower := strings.ToLower(s)
	for kw := range sensSet {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// hasSensitiveByRegex - проверка на чувствительные данные (проверка по паттерну)
func hasSensitiveByRegex(s string, res []*regexp.Regexp) bool {
	for _, re := range res {
		if re != nil && re.MatchString(s) {
			return true
		}
	}
	return false
}
