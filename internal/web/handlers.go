package web

import (
	"html/template"
	"net/http"

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

	if r.Method == http.MethodPost {

		cv := r.FormValue("cv")
		job := r.FormValue("job")

		if cv == "" || job == "" {

			http.Redirect(
				w,
				r,
				"/",
				http.StatusSeeOther,
			)

			return
		}

		result := ats.Analyze(
			cv,
			job,
		)

		err := tmpl.Execute(
			w,
			result,
		)

		if err != nil {

			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)

			return
		}

		return
	}

	err := tmpl.Execute(
		w,
		nil,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}