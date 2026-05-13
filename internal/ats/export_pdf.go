package ats

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func ExportCVToPDF(cv string, improvements []CVImprovement, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "AjiraScan - ATS Optimized CV")
	pdf.Ln(12)

	// Original CV Section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Original CV")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5, cv, "", "", false)
	pdf.Ln(5)

	// Improvements Section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "CV Improvements")
	pdf.Ln(8)

	for _, imp := range improvements {
		pdf.SetFont("Arial", "B", 10)
		pdf.MultiCell(0, 5, fmt.Sprintf("Original: %s", imp.Original), "", "", false)

		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, fmt.Sprintf("Improved: %s", imp.Improved), "", "", false)

		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, fmt.Sprintf("Reason: %s", imp.Reason), "", "", false)

		pdf.Ln(3)
	}

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return err
	}

	return nil
}