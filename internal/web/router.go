package web

import (
	"encoding/json"
	"net/http"
	"time"

	"ajirascan/internal/auth"
	"ajirascan/internal/database"
)

func NewRouter(db *database.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/analyze", HomeHandler)
	mux.Handle("/signin", PageHandler("templates/signin.html"))
	mux.Handle("/signup", PageHandler("templates/signup.html"))
	service := auth.Service{}
	if db != nil {
		service.DB = db.SQL
	}
	mux.Handle("/api/v1/auth/register", AuthHandler(service))
	mux.Handle("/api/v1/auth/login", AuthHandler(service))
	mux.Handle("/api/v1/auth/logout", AuthHandler(service))
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
