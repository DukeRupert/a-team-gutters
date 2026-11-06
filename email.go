package main

import (
	"fmt"
	"strings"
)

// emailConfig is the Postmark configuration for estimate-request emails.
type emailConfig struct {
	Token string
	From  string
	To    string
}

func (c emailConfig) enabled() bool { return c.Token != "" }

// loadEmailConfig reads the Postmark settings. All three are required unless
// EMAIL_DISABLED=1, which is meant for local development only. getenv is
// os.Getenv in production; tests pass a map lookup.
func loadEmailConfig(getenv func(string) string) (emailConfig, error) {
	if getenv("EMAIL_DISABLED") == "1" {
		return emailConfig{}, nil
	}

	c := emailConfig{
		Token: getenv("POSTMARK_SERVER_TOKEN"),
		From:  getenv("POSTMARK_FROM"),
		To:    getenv("POSTMARK_TO"),
	}

	var missing []string
	if c.Token == "" {
		missing = append(missing, "POSTMARK_SERVER_TOKEN")
	}
	if c.From == "" {
		missing = append(missing, "POSTMARK_FROM")
	}
	if c.To == "" {
		missing = append(missing, "POSTMARK_TO")
	}
	if len(missing) > 0 {
		return emailConfig{}, fmt.Errorf("estimate emails are not configured: %s unset (set EMAIL_DISABLED=1 to run without email in local dev)", strings.Join(missing, ", "))
	}
	return c, nil
}
