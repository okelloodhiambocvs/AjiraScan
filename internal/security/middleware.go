package security

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Middleware struct {
	MaxRequestBytes int64
	RateLimit       int
	TrustedOrigins  map[string]struct{}
	Logger          *slog.Logger
	limiter         *rateLimiter
}

func New(maxBytes int64, perMinute int, origins []string, logger *slog.Logger) *Middleware {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		allowed[origin] = struct{}{}
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Middleware{MaxRequestBytes: maxBytes, RateLimit: perMinute, TrustedOrigins: allowed, Logger: logger, limiter: newRateLimiter(perMinute)}
}

func (m *Middleware) Wrap(next http.Handler) http.Handler {
	return m.recover(m.headers(m.requestID(m.limit(m.rateLimit(m.cors(next))))))
}

func (m *Middleware) headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'; style-src 'self' 'unsafe-inline'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			r.Body = http.MaxBytesReader(w, r.Body, m.MaxRequestBytes)
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if _, ok := m.TrustedOrigins[origin]; !ok {
				http.Error(w, "origin not allowed", http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		if !m.limiter.allow(host) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if len(id) > 128 || id == "" {
			bytes := make([]byte, 16)
			if _, err := rand.Read(bytes); err == nil {
				id = hex.EncodeToString(bytes)
			} else {
				id = time.Now().UTC().Format("20060102150405.000000000")
			}
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				m.Logger.Error("panic recovered", "request_id", w.Header().Get("X-Request-ID"), "error", recovered)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type rateLimiter struct {
	mu      sync.Mutex
	limit   int
	buckets map[string]*bucket
}

type bucket struct {
	start time.Time
	count int
}

func newRateLimiter(limit int) *rateLimiter {
	return &rateLimiter{limit: limit, buckets: make(map[string]*bucket)}
}

func (l *rateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	b := l.buckets[key]
	if b == nil || now.Sub(b.start) >= time.Minute {
		l.buckets[key] = &bucket{start: now, count: 1}
		return true
	}
	if b.count >= l.limit {
		return false
	}
	b.count++
	return true
}

func ClientIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return strings.TrimSpace(remoteAddr)
	}
	return host
}
