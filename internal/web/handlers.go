package web

import (
	"html/template"
	"net/http"
	"strings"

	"ajirascan/internal/ats"
)

/*
Template helper functions
*/
var funcs = template.FuncMap{

	"add": func(a, b int) int {
		return a + b
	},

	"verdict": func(score int) string {

		switch {

		case score >= 80:
			return "Excellent Match"

		case score >= 60:
			return "Strong Match"

		case score >= 40:
			return "Moderate Match"

		default:
			return "Weak Match"
		}
	},
}

/*
Load templates
*/
var tmpl = template.Must(

	template.New(
		"index.html",
	).
		Funcs(funcs).
		ParseFiles(
			"templates/index.html",
		),
)

func HomeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {
	case http.MethodGet:
		renderHome(w, nil)
		return
	case http.MethodPost:
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form submission", http.StatusBadRequest)
			return
		}

		cv := r.PostForm.Get("cv")
		job := r.PostForm.Get("job")

		if strings.TrimSpace(cv) == "" || strings.TrimSpace(job) == "" {
			http.Error(w, "CV and job description are required", http.StatusBadRequest)
			return
		}
		if len(cv) > 256<<10 || len(job) > 256<<10 {
			http.Error(w, "CV and job description must each be at most 256 KiB", http.StatusRequestEntityTooLarge)
			return
		}

		result := ats.Analyze(
			cv,
			job,
		)

		renderHome(w, result)
		return
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func renderHome(w http.ResponseWriter, data any) {
	err := tmpl.Execute(
		w,
		data,
	)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
