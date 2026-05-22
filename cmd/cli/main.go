package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"ajirascan/internal/ats"
)

func read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// formats dotted alignment like OxfordCV style
func line(label string, value string) string {
	const width = 30
	spaces := width - len(label)
	if spaces < 1 {
		spaces = 1
	}
	return fmt.Sprintf("%s %s %s", label, strings.Repeat(".", spaces), value)
}

func main() {

	cv := flag.String("cv", "", "CV file path")
	job := flag.String("job", "", "Job file path")
	output := flag.String("out", "", "Output PDF file path")
	flag.Parse()

	if *cv == "" || *job == "" {
		fmt.Println("Usage: go run ./cmd/cli -cv <cv_file> -job <job_file> -out <output_pdf>")
		return
	}

	result := ats.Analyze(read(*cv), read(*job))

	// ================= HEADER =================
	fmt.Println("==================================")
	fmt.Println("AJIRASCAN ANALYSIS COMPLETE")
	fmt.Println("==================================")
	fmt.Println()

	ratingText := "Weak Match"
	switch {
	case result.Score >= 80:
		ratingText = "Excellent Match"
	case result.Score >= 60:
		ratingText = "Strong Match"
	case result.Score >= 40:
		ratingText = "Moderate Match"
	}

	fmt.Printf("RATING: %d/100 (%s)\n", result.Score, ratingText)
	fmt.Println()

	// ================= SCORE BREAKDOWN =================
	fmt.Println("TAILORING .............")
	fmt.Println("CONTENT ...............")
	fmt.Println("SECTIONS ..............")
	fmt.Println("ATS ESSENTIALS ........")
	fmt.Println()

	// If category scores exist, render properly
	for _, c := range result.CategoryAnalysis {
		fmt.Println(line(c.Category, fmt.Sprintf("%d/100", c.Score)))
	}
	fmt.Println()

	// ================= CORE INSIGHTS =================
	fmt.Println("----------------------------------")
	fmt.Println("CORE INSIGHTS")
	fmt.Println("----------------------------------")

	for i, m := range result.Matched {
		if i >= 8 {
			break
		}
		fmt.Println("✓", m)
	}

	for i, m := range result.Missing {
		if i >= 5 {
			break
		}
		fmt.Println("✕", m)
	}

	fmt.Println()

	// ================= SECTION COVERAGE =================
	fmt.Println("----------------------------------")
	fmt.Println("SECTION COVERAGE")
	fmt.Println("----------------------------------")

	for _, s := range result.FoundSections {
		fmt.Println("✓", s)
	}

	for _, s := range result.MissingSections {
		fmt.Println("✕", s)
	}

	fmt.Println()

	// ================= RECOMMENDATIONS =================
	fmt.Println("----------------------------------")
	fmt.Println("TOP RECOMMENDATIONS")
	fmt.Println("----------------------------------")

	for i, s := range result.Suggestions {
		if i >= 6 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, s)
	}

	fmt.Println()
	fmt.Println("==================================")
	fmt.Println("END OF REPORT")
	fmt.Println("==================================")

	// ================= PDF EXPORT =================
	if *output != "" {
		err := ats.ExportCVToPDF(read(*cv), result.Improvements, *output)
		if err != nil {
			fmt.Println("Failed to export PDF:", err)
		} else {
			fmt.Println("CV exported successfully to:", *output)
		}
	}
}