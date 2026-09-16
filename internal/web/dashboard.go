package web

import (
	"context"
	"html/template"
	"net/http"

	"ajirascan/internal/auth"
)

type dashboardProvider interface {
	Dashboard(context.Context, auth.Principal) (auth.Dashboard, error)
}

func DashboardHandler(provider dashboardProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		principal, ok := auth.PrincipalFromContext(r.Context())
		if !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		dashboard, err := provider.Dashboard(r.Context(), principal)
		if err == auth.ErrUnauthenticated {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "dashboard unavailable", http.StatusServiceUnavailable)
			return
		}
		t, err := template.ParseFiles(templatePath("dashboard.html"))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if err := t.Execute(w, dashboard); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	})
}

func MeHandler(provider dashboardProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		principal, ok := auth.PrincipalFromContext(r.Context())
		if !ok {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		dashboard, err := provider.Dashboard(r.Context(), principal)
		if err == auth.ErrUnauthenticated {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "profile unavailable", http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusOK, dashboard)
	})
}
