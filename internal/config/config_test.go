package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRejectsInvalidProductionConfiguration(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("TRUSTED_ORIGINS", "http://app.example.test")
	if _, err := Load(); err == nil {
		t.Fatal("expected non-HTTPS production origin rejection")
	}
}

func TestLoadDotEnvDoesNotOverrideProcessEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.env")
	if err := os.WriteFile(path, []byte("AJIRASCAN_CONFIG_TEST=dotenv-value\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AJIRASCAN_CONFIG_TEST", "process-value")
	if err := loadDotEnv(path); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv("AJIRASCAN_CONFIG_TEST"); got != "process-value" {
		t.Fatal("process environment must remain authoritative")
	}
}

func TestLoadRejectsMalformedConfiguredValues(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("SESSION_TTL", "not-a-duration")
	if _, err := Load(); err == nil {
		t.Fatal("expected malformed duration rejection")
	}
}
