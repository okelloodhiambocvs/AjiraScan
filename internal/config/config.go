package config

import (
	"errors"
	"fmt"
	"net/url"
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
	sessionTTL, err := requiredDuration("SESSION_TTL", 24*time.Hour)
	if err != nil {
		return Config{}, err
	}
	maxRequestBytes, err := requiredInt64("MAX_REQUEST_BYTES", 1<<20)
	if err != nil {
		return Config{}, err
	}
	rateLimit, err := requiredInt("RATE_LIMIT_PER_MINUTE", 60)
	if err != nil {
		return Config{}, err
	}
	uploadMaxBytes, err := requiredInt64("MAX_UPLOAD_BYTES", 10<<20)
	if err != nil {
		return Config{}, err
	}
	cookieSecure, err := requiredBool("COOKIE_SECURE", false)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		Environment:     value("APP_ENV", "development"),
		Address:         value("APP_ADDR", "127.0.0.1:8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		TrustedOrigins:  csv(os.Getenv("TRUSTED_ORIGINS")),
		SessionTTL:      sessionTTL,
		MaxRequestBytes: maxRequestBytes,
		RateLimitPerMin: rateLimit,
		UploadMaxBytes:  uploadMaxBytes,
		StorageDriver:   value("STORAGE_DRIVER", "local"),
		AIProvider:      value("AI_PROVIDER", "disabled"),
		AIModel:         os.Getenv("AI_MODEL"),
		PaymentProvider: value("PAYMENT_PROVIDER", "disabled"),
	}
	cfg.RequireHTTPS = cfg.Environment == "production"
	cfg.CookieSecure = cfg.RequireHTTPS || cookieSecure
	switch cfg.Environment {
	case "development", "test", "production":
	default:
		return Config{}, errors.New("APP_ENV must be development, test, or production")
	}
	if cfg.Environment == "production" {
		if cfg.DatabaseURL == "" {
			return Config{}, errors.New("DATABASE_URL is required in production")
		}
		if len(cfg.TrustedOrigins) == 0 {
			return Config{}, errors.New("TRUSTED_ORIGINS is required in production")
		}
		for _, origin := range cfg.TrustedOrigins {
			parsed, err := url.ParseRequestURI(origin)
			if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.Path != "" {
				return Config{}, fmt.Errorf("TRUSTED_ORIGINS must contain HTTPS origins without paths")
			}
		}
	}
	if cfg.SessionTTL <= 0 || cfg.MaxRequestBytes <= 0 || cfg.UploadMaxBytes <= 0 || cfg.RateLimitPerMin <= 0 {
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

func requiredDuration(key string, fallback time.Duration) (time.Duration, error) {
	if raw := os.Getenv(key); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return 0, fmt.Errorf("%s must be a duration: %w", key, err)
		}
		return parsed, nil
	}
	return fallback, nil
}

func requiredInt(key string, fallback int) (int, error) {
	if raw := os.Getenv(key); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return 0, fmt.Errorf("%s must be an integer: %w", key, err)
		}
		return parsed, nil
	}
	return fallback, nil
}

func requiredInt64(key string, fallback int64) (int64, error) {
	if raw := os.Getenv(key); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%s must be an integer: %w", key, err)
		}
		return parsed, nil
	}
	return fallback, nil
}

func requiredBool(key string, fallback bool) (bool, error) {
	if raw := os.Getenv(key); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return false, fmt.Errorf("%s must be a boolean: %w", key, err)
		}
		return parsed, nil
	}
	return fallback, nil
}
