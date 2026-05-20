package ats

import (
	"fmt"
	"strings"
)

func FormatReport(r Result) string {

	var b strings.Builder

	b.WriteString("==================================\n")
	b.WriteString("AJIRASCAN ANALYSIS COMPLETE\n")
	b.WriteString(fmt.Sprintf("RATING: %d/100\n", r.Score))
	b.WriteString("==================================\n\n")

	b.WriteString(fmt.Sprintf("JOB TYPE: %s\n\n", r.JobType))

	// CORE SUMMARY
	b.WriteString("----------------------------------\n")
	b.WriteString("CORE SUMMARY\n")
	b.WriteString("----------------------------------\n")

	for _, s := range r.Suggestions[:min(3, len(r.Suggestions))] {
		b.WriteString("✓ " + s + "\n")
	}

	// SKILL ALIGNMENT (SIMPLIFIED)
	b.WriteString("\n----------------------------------\n")
	b.WriteString("SKILL ALIGNMENT\n")
	b.WriteString("----------------------------------\n")

	for i := 0; i < min(6, len(r.Matched)); i++ {
		b.WriteString("✓ " + r.Matched[i] + "\n")
	}

	for i := 0; i < min(6, len(r.Missing)); i++ {
		b.WriteString("✕ " + r.Missing[i] + "\n")
	}

	// SECTION COVERAGE
	b.WriteString("\n----------------------------------\n")
	b.WriteString("SECTION COVERAGE\n")
	b.WriteString("----------------------------------\n")

	for _, s := range r.FoundSections {
		b.WriteString("✓ " + s + "\n")
	}

	for _, s := range r.MissingSections {
		b.WriteString("✕ " + s + "\n")
	}

	// RECOMMENDATIONS
	b.WriteString("\n----------------------------------\n")
	b.WriteString("RECOMMENDATIONS\n")
	b.WriteString("----------------------------------\n")

	for _, s := range r.Suggestions {
		b.WriteString("• " + s + "\n")
	}

	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}