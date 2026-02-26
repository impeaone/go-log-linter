package analyzer

import "testing"

func TestSetConfig_RejectsBadRegex(t *testing.T) {
	c := DefaultConfig()
	c.AllowedCharsRegex = "[a-"
	if err := SetConfig(c); err == nil {
		t.Fatalf("expected error")
	}
}

func TestCheckMessageAll_UsesSensitiveKeywords(t *testing.T) {
	c := DefaultConfig()
	c.ForbidSensitive = true
	c.SensitiveKeywords = []string{"auth_token"}
	if err := SetConfig(c); err != nil {
		t.Fatalf("SetConfig: %v", err)
	}

	compConf, err := compileConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	vs := CheckMessageAllWith(compConf, "auth_token")

	found := false
	for _, v := range vs {
		if v == VSensitive {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected VSensitive, got %v", vs)
	}
}
