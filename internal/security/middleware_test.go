package security

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddlewareAddsSecurityHeaders(t *testing.T) {
	handler := New(100, 10, nil, slog.Default()).Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Header().Get("Content-Security-Policy") == "" || response.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("expected response security headers")
	}
}

func TestMiddlewareRejectsUnknownOrigin(t *testing.T) {
	handler := New(100, 10, []string{"https://app.example.test"}, slog.Default()).Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Origin", "https://evil.example.test")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden origin, got %d", response.Code)
	}
}

func TestMiddlewareAllowsSameOriginWithoutConfiguration(t *testing.T) {
	handler := New(100, 10, nil, slog.Default()).Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/analyze", nil)
	request.Header.Set("Origin", "http://localhost:8080")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected same-origin request, got %d", response.Code)
	}
}

func TestMiddlewareRejectsMutationWithoutOrigin(t *testing.T) {
	handler := New(100, 10, nil, slog.Default()).Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "http://localhost:8080/analyze", nil))
	if response.Code != http.StatusForbidden {
		t.Fatalf("expected mutation without Origin to be rejected, got %d", response.Code)
	}
}

func TestRateLimiterExpiresInactiveBucketsAndBoundsNewKeys(t *testing.T) {
	limiter := newRateLimiter(2)
	limiter.maxBuckets = 1
	limiter.buckets["expired"] = &bucket{start: time.Now().Add(-3 * time.Minute), count: 1}
	if !limiter.allow("active") {
		t.Fatal("expected expired bucket to be reclaimed")
	}
	if limiter.allow("another") {
		t.Fatal("expected limiter to reject a new key when bucket capacity is reached")
	}
}
