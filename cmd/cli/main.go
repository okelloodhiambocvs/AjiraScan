package main

import (
	"flag"
	"fmt"
	"os"

	"ajirascan/internal/ats"
)

func read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
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

	// CLEAN REPORT OUTPUT (NEW FORMAT)
	fmt.Println("\n==================================")
	fmt.Println("AJIRASCAN ANALYSIS COMPLETE")
	fmt.Println("==================================\n")

	fmt.Printf("JOB TYPE: %s\n", result.JobType)
	fmt.Printf("ATS SCORE: %d/100\n", result.Score)

	switch {
	case result.Score >= 80:
		fmt.Println("RATING: Excellent Match")
	case result.Score >= 60:
		fmt.Println("RATING: Strong Match")
	case result.Score >= 40:
		fmt.Println("RATING: Moderate Match")
	default:
		fmt.Println("RATING: Weak Match")
	}

	fmt.Println("\n----------------------------------")
	fmt.Println("CORE INSIGHTS")
	fmt.Println("----------------------------------")

	// Top 5 matched keywords only
	for i, m := range result.Matched {
		if i >= 5 {
			break
		}
		fmt.Println("✓", m)
	}

	// Top 5 missing keywords only
	for i, m := range result.Missing {
		if i >= 5 {
			break
		}
		fmt.Println("✕", m)
	}

	fmt.Println("\n----------------------------------")
	fmt.Println("SECTION COVERAGE")
	fmt.Println("----------------------------------")

	for _, s := range result.FoundSections {
		fmt.Println("✓", s)
	}

	for _, s := range result.MissingSections {
		fmt.Println("✕", s)
	}

	fmt.Println("\n----------------------------------")
	fmt.Println("TOP RECOMMENDATIONS")
	fmt.Println("----------------------------------")

	// Merge + limit suggestions (avoid noise explosion)
	suggestions := append(result.Suggestions,
		ats.SectionSuggestions(result.MissingSections)...,
	)

	frequencySuggestions := ats.FrequencySuggestions(result.KeywordFrequency)
	suggestions = append(suggestions, frequencySuggestions...)

	for i, s := range suggestions {
		if i >= 6 {
			break
		}
		fmt.Println("•", s)
	}

	// PDF EXPORT OPTION
	if *output != "" {
		err := ats.ExportCVToPDF(read(*cv), result.Improvements, *output)
		if err != nil {
			fmt.Println("\nFailed to export PDF:", err)
			return
		}

		fmt.Println("\n----------------------------------")
		fmt.Println("PDF EXPORT")
		fmt.Println("----------------------------------")
		fmt.Println("✓ CV exported successfully to:", *output)
	}

	fmt.Println("\n==================================")
	fmt.Println("END OF REPORT")
	fmt.Println("==================================")
}