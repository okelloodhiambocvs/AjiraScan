package fileparser

import (
	"errors"
	"io"
	"strings"

	"ajirascan/internal/documents"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
)

const MaxPDFPages = 100

var ErrNoExtractableText = errors.New("document contains no extractable text")

func ReadTXT(reader io.Reader) (string, error) {
	data, err := documents.ReadBounded(reader, documents.MaxExtractedTextBytes)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(string(data)) == "" {
		return "", ErrNoExtractableText
	}
	return string(data), nil
}

func ReadDOCX(path string) (string, error) {

	doc, err := docx.ReadDocxFile(path)
	if err != nil {
		return "", err
	}
	defer doc.Close()

	content := doc.Editable().GetContent()
	if len(content) > documents.MaxExtractedTextBytes {
		return "", errors.New("extracted document text exceeds configured limit")
	}
	if strings.TrimSpace(content) == "" {
		return "", ErrNoExtractableText
	}
	return content, nil
}

func ReadPDF(path string) (string, error) {

	file, reader, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var builder strings.Builder
	totalPages := reader.NumPage()
	if totalPages <= 0 {
		return "", ErrNoExtractableText
	}
	if totalPages > MaxPDFPages {
		return "", errors.New("PDF exceeds configured page limit")
	}

	for i := 1; i <= totalPages; i++ {

		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		if builder.Len()+len(text) > documents.MaxExtractedTextBytes {
			return "", errors.New("extracted PDF text exceeds configured limit")
		}
		builder.WriteString(text)
	}

	content := builder.String()
	if strings.TrimSpace(content) == "" {
		return "", ErrNoExtractableText
	}
	return content, nil
}
