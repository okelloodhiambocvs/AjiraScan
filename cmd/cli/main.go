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
	flag.Parse()

	if *cv == "" || *job == "" {
		fmt.Println("Usage: go run ./cmd/cli -cv <cv_file> -job <job_file>")
		return
	}

	result := ats.Analyze(read(*cv), read(*job))

	fmt.Println("====== AJIRASCAN ATS REPORT ======")
	fmt.Println()

	fmt.Println("Detected Job Type:", result.JobType)
	fmt.Println()

	fmt.Println("ATS Score:", result.Score)

	switch {
	case result.Score >= 80:
		fmt.Println("ATS Rating: Excellent Match")
	case result.Score >= 60:
		fmt.Println("ATS Rating: Strong Match")
	case result.Score >= 40:
		fmt.Println("ATS Rating: Moderate Match")
	default:
		fmt.Println("ATS Rating: Weak Match")
	}

	fmt.Println()

	fmt.Println("Matched Keywords:")
	for _, m := range result.Matched {
		fmt.Println("-", m)
	}

	fmt.Println()

	fmt.Println("Missing Keywords:")
	for _, m := range result.Missing {
		fmt.Println("-", m)
	}

	fmt.Println()

	fmt.Println("Keyword Frequency:")
	for _, item := range result.KeywordFrequency {
		fmt.Printf("- %s: %d\n", item.Keyword, item.Count)
	}

	fmt.Println()

	fmt.Println("Detected Sections:")
	for _, s := range result.FoundSections {
		fmt.Println("✓", s)
	}

	fmt.Println()

	fmt.Println("Missing Sections:")
	for _, s := range result.MissingSections {
		fmt.Println("✗", s)
	}

	fmt.Println()

	fmt.Println("Skill Category Analysis:")
	for _, c := range result.CategoryAnalysis {
		fmt.Printf("- %s: %d%% match\n", c.Category, c.Score)
	}

	fmt.Println()

	fmt.Println("Suggestions:")
	for _, s := range result.Suggestions {
		fmt.Println("-", s)
	}

	sectionSuggestions := ats.SectionSuggestions(result.MissingSections)
	for _, s := range sectionSuggestions {
		fmt.Println("-", s)
	}

	frequencySuggestions := ats.FrequencySuggestions(result.KeywordFrequency)
	for _, s := range frequencySuggestions {
		fmt.Println("-", s)
	}

	fmt.Println()
	
	fmt.Println("CV Improvement Suggestions:")

	for _, imp := range result.Improvements {
		fmt.Println("Original:", imp.Original)
		fmt.Println("Improved:", imp.Improved)
		fmt.Println("Reason:", imp.Reason)
		fmt.Println()
	}
}