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

	// =========================
	// CLEAN REPORT OUTPUT
	// =========================

	fmt.Println("==================================")
	fmt.Println("AJIRASCAN ANALYSIS COMPLETE")
	fmt.Println("==================================")
	fmt.Println()

	// Score header
	ratingText := ""
	switch {
	case result.Score >= 80:
		ratingText = "Excellent Match"
	case result.Score >= 60:
		ratingText = "Strong Match"
	case result.Score >= 40:
		ratingText = "Moderate Match"
	default:
		ratingText = "Weak Match"
	}

	fmt.Printf("RATING: %d/100 (%s)\n", result.Score, ratingText)
	fmt.Println()

	// Breakdown (simple weighted display placeholder)
	fmt.Println("BREAKDOWN")
	fmt.Println("- Tailoring:", result.Score*30/100, "/30")
	fmt.Println("- Content:", result.Score*30/100, "/30")
	fmt.Println("- Sections:", result.Score*20/100, "/20")
	fmt.Println("- ATS Essentials:", result.Score*20/100, "/20")
	fmt.Println()

	// Key Insights (LIMIT noise)
	fmt.Println("KEY INSIGHTS")

	insightLimit := 5
	count := 0
	for _, m := range result.Matched {
		if count >= insightLimit {
			break
		}
		fmt.Println("✓", m)
		count++
	}

	for _, m := range result.Missing {
		if count >= insightLimit {
			break
		}
		fmt.Println("✕", m)
		count++
	}
	fmt.Println()

	// Sections coverage (clean)
	fmt.Println("SECTIONS COVERAGE")

	for _, s := range result.FoundSections {
		fmt.Println("✓", s)
	}

	for _, s := range result.MissingSections {
		fmt.Println("✕", s)
	}
	fmt.Println()

	// Top recommendations ONLY (no spam)
	fmt.Println("TOP RECOMMENDATIONS")

	recCount := 0
	for _, s := range result.Suggestions {
		if recCount >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", recCount+1, s)
		recCount++
	}

	sectionSuggestions := ats.SectionSuggestions(result.MissingSections)
	for _, s := range sectionSuggestions {
		if recCount >= 5 {
			break
		}
		recCount++
		fmt.Printf("%d. %s\n", recCount, s)
	}

	fmt.Println()
	fmt.Println("==================================")
	fmt.Println("END OF REPORT")
	fmt.Println("==================================")

	// PDF export
	if *output != "" {
		err := ats.ExportCVToPDF(read(*cv), result.Improvements, *output)
		if err != nil {
			fmt.Println("Failed to export PDF:", err)
		} else {
			fmt.Println()
			fmt.Println("CV exported successfully to:", *output)
		}
	}
}