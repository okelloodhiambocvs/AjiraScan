package config

import "testing"

func TestLoadRejectsInvalidProductionConfiguration(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("TRUSTED_ORIGINS", "http://app.example.test")
	if _, err := Load(); err == nil {
		t.Fatal("expected non-HTTPS production origin rejection")
	}
}

func TestLoadRejectsMalformedConfiguredValues(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("SESSION_TTL", "not-a-duration")
	if _, err := Load(); err == nil {
		t.Fatal("expected malformed duration rejection")
	}
}
