package web

import (
	"html/template"
	"net/http"
)

func PageHandler(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func AuthUnavailableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error": "authentication persistence is not configured",
	})
}
