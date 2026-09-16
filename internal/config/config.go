package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment     string
	Address         string
	DatabaseURL     string
	TrustedOrigins  []string
	SessionTTL      time.Duration
	CookieSecure    bool
	MaxRequestBytes int64
	RateLimitPerMin int
	UploadMaxBytes  int64
	StorageDriver   string
	AIProvider      string
	AIModel         string
	PaymentProvider string
	RequireHTTPS    bool
}

func Load() (Config, error) {
	cfg := Config{
		Environment:     value("APP_ENV", "development"),
		Address:         value("APP_ADDR", "127.0.0.1:8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		TrustedOrigins:  csv(os.Getenv("TRUSTED_ORIGINS")),
		SessionTTL:      duration("SESSION_TTL", 24*time.Hour),
		MaxRequestBytes: int64Value("MAX_REQUEST_BYTES", 1<<20),
		RateLimitPerMin: intValue("RATE_LIMIT_PER_MINUTE", 60),
		UploadMaxBytes:  int64Value("MAX_UPLOAD_BYTES", 10<<20),
		StorageDriver:   value("STORAGE_DRIVER", "local"),
		AIProvider:      value("AI_PROVIDER", "disabled"),
		AIModel:         os.Getenv("AI_MODEL"),
		PaymentProvider: value("PAYMENT_PROVIDER", "disabled"),
	}
	cfg.RequireHTTPS = cfg.Environment == "production"
	cfg.CookieSecure = cfg.RequireHTTPS || boolValue("COOKIE_SECURE", false)
	if cfg.Environment == "production" && cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required in production")
	}
	if cfg.MaxRequestBytes <= 0 || cfg.UploadMaxBytes <= 0 || cfg.RateLimitPerMin <= 0 {
		return Config{}, errors.New("request, upload and rate-limit settings must be positive")
	}
	return cfg, nil
}

func value(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func csv(v string) []string {
	var result []string
	for _, part := range strings.Split(v, ",") {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	return result
}

func duration(key string, fallback time.Duration) time.Duration {
	if raw := os.Getenv(key); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil {
			return parsed
		}
	}
	return fallback
}

func intValue(key string, fallback int) int {
	if raw := os.Getenv(key); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			return parsed
		}
	}
	return fallback
}

func int64Value(key string, fallback int64) int64 {
	if raw := os.Getenv(key); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func boolValue(key string, fallback bool) bool {
	if raw := os.Getenv(key); raw != "" {
		if parsed, err := strconv.ParseBool(raw); err == nil {
			return parsed
		}
	}
	return fallback
}
