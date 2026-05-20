package web

import (
	"html/template"
	"net/http"

	"ajirascan/internal/ats"
)

// FIXED PATH: templates (NOT template)
var tmpl = template.Must(
	template.ParseFiles("templates/index.html"),
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		cv := r.FormValue("cv")
		job := r.FormValue("job")

		if cv == "" || job == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		result := ats.Analyze(cv, job)

		_ = tmpl.Execute(w, result)
		return
	}

	_ = tmpl.Execute(w, nil)
}