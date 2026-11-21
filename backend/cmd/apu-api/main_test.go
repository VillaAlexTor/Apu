package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("APU_DATABASE_URL", "postgres://apu:local@localhost:5433/apu")
	t.Setenv("APU_SESSION_TTL", "24h")
	t.Setenv("APU_COOKIE_SECURE", "true")

	config, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if config.sessionTTL.String() != "24h0m0s" {
		t.Fatalf("unexpected session TTL: %s", config.sessionTTL)
	}
	if !config.cookieSecure {
		t.Fatal("expected secure cookies")
	}
}

func TestLoadConfigRequiresDatabaseURL(t *testing.T) {
	_ = os.Unsetenv("APU_DATABASE_URL")
	if _, err := loadConfig(); err == nil {
		t.Fatal("expected missing database URL to fail")
	}
}

func TestProductionRequiresMFAKey(t *testing.T) {
	t.Setenv("APU_ENV", "production")
	t.Setenv("APU_DATABASE_URL", "postgres://apu:local@localhost:5433/apu")
	t.Setenv("APU_COOKIE_SECURE", "true")
	t.Setenv("APU_ALLOWED_ORIGIN", "https://panel.apu.example")
	t.Setenv("APU_SMTP_ADDRESS", "smtp.example:587")
	t.Setenv("APU_SMTP_FROM", "APU <no-reply@apu.example>")
	_ = os.Unsetenv("APU_MFA_ENCRYPTION_KEY")
	if _, err := loadConfig(); err == nil {
		t.Fatal("expected production without MFA encryption key to fail")
	}
}

func TestProductionRejectsExposedDevelopmentTokens(t *testing.T) {
	t.Setenv("APU_ENV", "production")
	t.Setenv("APU_DATABASE_URL", "postgres://apu:local@localhost:5433/apu")
	t.Setenv("APU_COOKIE_SECURE", "true")
	t.Setenv("APU_ALLOWED_ORIGIN", "https://panel.apu.example")
	t.Setenv("APU_MFA_ENCRYPTION_KEY", "AQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQE")
	t.Setenv("APU_EXPOSE_DEVELOPMENT_TOKENS", "true")
	t.Setenv("APU_SMTP_ADDRESS", "smtp.example:587")
	t.Setenv("APU_SMTP_FROM", "APU <no-reply@apu.example>")
	if _, err := loadConfig(); err == nil {
		t.Fatal("expected production to reject exposed development tokens")
	}
}
