package fileparser

import (
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
)

func ReadTXT(reader io.Reader) string {
	buf := new(strings.Builder)
	io.Copy(buf, reader)
	return buf.String()
}

func ReadDOCX(path string) (string, error) {

	doc, err := docx.ReadDocxFile(path)
	if err != nil {
		return "", err
	}
	defer doc.Close()

	return doc.Editable().GetContent(), nil
}

func ReadPDF(path string) (string, error) {

	file, reader, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var builder strings.Builder
	totalPages := reader.NumPage()

	for i := 1; i <= totalPages; i++ {

		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err == nil {
			builder.WriteString(text)
		}
	}

	return builder.String(), nil
}