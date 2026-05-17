package web

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"ajirascan/internal/ats"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmpl := template.Must(
			template.ParseFiles("templates/index.html"),
		)

		tmpl.Execute(w, nil)
		return
	}

	cv := r.FormValue("cv")

	file, header, err := r.FormFile("cvfile")

	if err == nil {

		defer file.Close()

		filename := strings.ToLower(header.Filename)

		switch {

		case strings.HasSuffix(filename, ".txt"):

			cv = ReadTXT(file)

		case strings.HasSuffix(filename, ".docx"):

			tempPath := "temp.docx"

			tempFile, _ := os.Create(tempPath)

			io.Copy(tempFile, file)

			tempFile.Close()

			parsed, parseErr := ReadDOCX(tempPath)

			if parseErr == nil {
				cv = parsed
			}

			os.Remove(tempPath)

		case strings.HasSuffix(filename, ".pdf"):

			tempPath := "temp.pdf"

			tempFile, _ := os.Create(tempPath)

			io.Copy(tempFile, file)

			tempFile.Close()

			parsed, parseErr := ReadPDF(tempPath)

			if parseErr == nil {
				cv = parsed
			}

			os.Remove(tempPath)
		}
	}

	job := r.FormValue("job")

	result := ats.Analyze(cv, job)

	tmpl := template.Must(
		template.ParseFiles("templates/index.html"),
	)

	tmpl.Execute(w, result)
}