package main

import (
	"strings"
	"testing"
)

func envMap(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadEmailConfig(t *testing.T) {
	full := map[string]string{
		"POSTMARK_SERVER_TOKEN": "tok",
		"POSTMARK_FROM":         "from@example.com",
		"POSTMARK_TO":           "to@example.com",
	}

	c, err := loadEmailConfig(envMap(full))
	if err != nil {
		t.Fatalf("full config: unexpected error: %v", err)
	}
	if !c.enabled() || c.Token != "tok" || c.From != "from@example.com" || c.To != "to@example.com" {
		t.Fatalf("full config: got %+v", c)
	}

	// The 2026-09 incident: compose passed POSTMARK_TOKEN and ESTIMATE_TO,
	// so the names the app reads were unset.
	_, err = loadEmailConfig(envMap(map[string]string{
		"POSTMARK_TOKEN": "tok",
		"POSTMARK_FROM":  "from@example.com",
		"ESTIMATE_TO":    "to@example.com",
	}))
	if err == nil {
		t.Fatal("misnamed vars: expected error, got nil")
	}
	for _, name := range []string{"POSTMARK_SERVER_TOKEN", "POSTMARK_TO"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q does not name %s", err, name)
		}
	}
	if strings.Contains(err.Error(), "POSTMARK_FROM") {
		t.Errorf("error %q names POSTMARK_FROM, which was set", err)
	}

	c, err = loadEmailConfig(envMap(map[string]string{"EMAIL_DISABLED": "1"}))
	if err != nil || c.enabled() {
		t.Fatalf("EMAIL_DISABLED=1: got %+v, %v", c, err)
	}
}
