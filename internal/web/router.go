package web

import (
	"encoding/json"
	"net/http"
	"time"

	"ajirascan/internal/auth"
	"ajirascan/internal/database"
)

type RouterOptions struct {
	SessionTTL   time.Duration
	CookieSecure bool
}

func NewRouter(db *database.DB, options ...RouterOptions) http.Handler {
	settings := RouterOptions{SessionTTL: 24 * time.Hour}
	if len(options) > 0 {
		settings = options[0]
	}
	if settings.SessionTTL <= 0 {
		settings.SessionTTL = 24 * time.Hour
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/analyze", HomeHandler)
	mux.Handle("/signin", PageHandler(templatePath("signin.html")))
	mux.Handle("/signup", PageHandler(templatePath("signup.html")))
	mux.Handle("/features", PageHandler(templatePath("features.html")))
	mux.Handle("/pricing", PageHandler(templatePath("pricing.html")))
	service := auth.Service{SessionTTL: settings.SessionTTL}
	if db != nil {
		service.DB = db.SQL
	}
	authHandler := AuthHandler(service, CookieSettings{Secure: settings.CookieSecure, MaxAge: int(settings.SessionTTL.Seconds())})
	mux.Handle("/api/v1/auth/register", authHandler)
	mux.Handle("/api/v1/auth/login", authHandler)
	mux.Handle("/api/v1/auth/logout", authHandler)
	mux.Handle("/dashboard", auth.RequireSession(service, DashboardHandler(service)))
	mux.Handle("/api/v1/me", auth.RequireSession(service, MeHandler(service)))
	for path, page := range legalPages() {
		mux.Handle(path, LegalHandler(page))
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if db == nil || !db.Ready(r.Context()) {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not ready"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func utcNow() time.Time {
	return time.Now().UTC()
}
