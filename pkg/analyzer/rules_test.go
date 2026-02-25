package analyzer

import "testing"

func has(vs []Violation, v Violation) bool {
	for _, x := range vs {
		if x == v {
			return true
		}
	}
	return false
}

func TestCheckMessageAll_Valid(t *testing.T) {
	tests := []string{
		"starting server on port 8080",
		"failed to connect to database",
		"user authenticated successfully",
	}

	for _, msg := range tests {
		vs := CheckMessageAll(msg)
		if len(vs) != 0 {
			t.Fatalf("expected no violations for %q, got %v", msg, vs)
		}
	}
}

func TestCheckMessageAll_Invalid(t *testing.T) {
	if err := SetConfig(DefaultConfig()); err != nil {
		t.Fatalf("SetConfig: %v", err)
	}

	tests := []struct {
		msg  string
		want []Violation
	}{
		{"Starting server", []Violation{VLowercase}},
		{"запуск сервера", []Violation{VEnglish, VSpecial}},
		{"server started!🚀", []Violation{VEnglish, VSpecial}},
		{"user password: secret123", []Violation{VSensitive}},
		{"api_key=" + "123", []Violation{VSensitive, VSpecial}},
	}

	for _, tt := range tests {
		vs := CheckMessageAll(tt.msg)
		for _, w := range tt.want {
			if !has(vs, w) {
				t.Fatalf("msg=%q: expected violation %q, got %v", tt.msg, w, vs)
			}
		}
	}
}
