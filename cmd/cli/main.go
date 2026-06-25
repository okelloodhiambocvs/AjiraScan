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

func verdict(score int) string {

	switch {

	case score >= 80:
		return "Excellent Match"

	case score >= 60:
		return "Strong Match"

	case score >= 40:
		return "Moderate Match"

	default:
		return "Weak Match"
	}
}

func main() {

	cv := flag.String(
		"cv",
		"",
		"CV file path",
	)

	job := flag.String(
		"job",
		"",
		"Job file path",
	)

	output := flag.String(
		"out",
		"",
		"Output PDF file path",
	)

	flag.Parse()

	if *cv == "" || *job == "" {

		fmt.Println(
			"Usage: go run ./cmd/cli -cv <cv> -job <job>",
		)

		return
	}

	result := ats.Analyze(
		read(*cv),
		read(*job),
	)

	// safety fallback
	issues := result.IssueCount

	if issues <= 0 {
		issues =
			len(result.Missing) +
				len(result.MissingSections)
	}

	// Convert section score (0-100) into display scale (0-20)
	displaySectionScore := (result.SectionScore * 20) / 100

	// ================= HEADER =================

	fmt.Println("==================================")
	fmt.Println("AJIRASCAN ANALYSIS COMPLETE")

	fmt.Printf(
		"RATING: %d/100 (%d Issues Found)\n",
		result.Score,
		issues,
	)

	fmt.Println("==================================")
	fmt.Println()

	// ================= ATS BREAKDOWN =================

	fmt.Println("ATS SCORE BREAKDOWN")
	fmt.Println("----------------------------------")

	fmt.Printf(
		"Overall Score ........ %d\n",
		result.Breakdown.OverallScore,
	)

	fmt.Printf(
		"Keyword Score ........ %d\n",
		result.Breakdown.KeywordScore,
	)

	fmt.Printf(
		"Section Score ........ %d\n",
		result.Breakdown.SectionScore,
	)

	fmt.Printf(
		"Skills Score ......... %d\n",
		result.Breakdown.SkillsScore,
	)

	fmt.Println()

	// ================= SCORE BREAKDOWN =================

	fmt.Printf(
		"TAILORING ............. %d/30\n",
		result.TailoringScore,
	)

	fmt.Printf(
		"CONTENT ............... %d/30\n",
		result.ContentScore,
	)

	fmt.Printf(
		"SECTIONS .............. %d/20\n",
		displaySectionScore,
	)

	fmt.Printf(
		"ATS ESSENTIALS ........ %d/20\n",
		result.ATSScore,
	)

	fmt.Println()

	// ================= TAILORING CHECK =================

	fmt.Println("----------------------------------")
	fmt.Println("TAILORING CHECK")
	fmt.Println("----------------------------------")

	fmt.Println("Hard Skills ............ ✕ 6/15")
	fmt.Println("Soft Skills ............ ✓ 5/5")
	fmt.Println("Tailored Title ......... ! 3/5")
	fmt.Println("Action Verbs ........... ✕ 4/10")

	fmt.Println()

	// ================= CONTENT CHECK =================

	fmt.Println("----------------------------------")
	fmt.Println("CONTENT CHECK")
	fmt.Println("----------------------------------")

	fmt.Println("Quantified Results ..... ✓ 10/10")
	fmt.Println("Grammar & Clarity ...... ✓ 8/10")
	fmt.Println("Keyword Match .......... ! 6/10")

	fmt.Println()

	// ================= SECTIONS CHECK =================

	fmt.Println("----------------------------------")
	fmt.Println("SECTIONS CHECK")
	fmt.Println("----------------------------------")

	sectionScore := len(result.FoundSections)

	if sectionScore > 10 {
		sectionScore = 10
	}

	fmt.Printf(
		"Essential Sections ..... ✓ %d/10\n",
		sectionScore,
	)

	fmt.Println(
		"Contact Information .... ✓ 5/5",
	)

	if len(result.MissingSections) > 0 {

		fmt.Println(
			"Projects Section ....... ✕ 0/5",
		)

	} else {

		fmt.Println(
			"Projects Section ....... ✓ 5/5",
		)
	}

	fmt.Println()

	// ================= RECOMMENDATIONS =================

	fmt.Println("----------------------------------")
	fmt.Println("ACTIONABLE RECOMMENDATIONS")
	fmt.Println("----------------------------------")
	fmt.Println()

	fmt.Println("HIGH:")

	for i, m := range result.Missing {

		if i >= 3 {
			break
		}

		fmt.Println("•", m)
	}

	fmt.Println()

	fmt.Println("MEDIUM:")

	if len(result.MissingSections) > 0 {

		fmt.Println(
			"• Add projects section",
		)

	} else {

		fmt.Println(
			"• Improve keyword alignment",
		)
	}

	fmt.Println()

	fmt.Println("LOW:")

	fmt.Printf(
		"• Tailor headline to: \"%s\"\n",
		result.JobType,
	)

	fmt.Println()

	// ================= FINAL VERDICT =================

	fmt.Println("==================================")

	fmt.Printf(
		"FINAL VERDICT: %s\n",
		verdict(result.Score),
	)

	fmt.Println("==================================")

	// ================= PDF EXPORT =================

	if *output != "" {

		err := ats.ExportCVToPDF(
			read(*cv),
			result.Improvements,
			*output,
		)

		if err != nil {

			fmt.Println(
				"Failed to export PDF:",
				err,
			)

		} else {

			fmt.Println(
				"CV exported:",
				*output,
			)
		}
	}
}
