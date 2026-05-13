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

	fmt.Println("====== ATS RESULT ======")
	fmt.Println("Score:", result.Score)
	fmt.Println("Matched:", result.Matched)
	fmt.Println("Missing:", result.Missing)
}