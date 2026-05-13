package main

import (
	"ajirascan/internal/ats"
	"flag"
	"fmt"
	"os"
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
	fmt.Println("ATS Score:", result.Score)

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
	fmt.Println("Suggestions:")
	for _, s := range result.Suggestions {
		fmt.Println("-", s)
	}
}