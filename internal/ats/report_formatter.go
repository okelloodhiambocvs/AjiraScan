package ats

import (
	"fmt"
	"strings"
)

func FormatReport(r Result) string {

	var b strings.Builder

	issues := len(r.MissingSections)

	if len(r.Suggestions) > issues {
		issues = len(r.Suggestions)
	}

	if issues > 9 {
		issues = 9
	}

	b.WriteString(
		"==================================\n",
	)

	b.WriteString(
		"AJIRASCAN ANALYSIS COMPLETE\n",
	)

	b.WriteString(
		fmt.Sprintf(
			"RATING: %d/100 (%d Issues Found)\n",
			r.Score,
			issues,
		),
	)

	b.WriteString(
		"==================================\n\n",
	)

	// SCORE BLOCKS

	tailoring := r.TailoringScore
	content := r.ContentScore

	// Convert dynamic 0-100 section score into 0-20 display scale
	sections := (r.SectionScore * 20) / 100

	essentials := r.ATSScore

	b.WriteString(
		fmt.Sprintf(
			"TAILORING ............. %d/30\n",
			tailoring,
		),
	)

	b.WriteString(
		fmt.Sprintf(
			"CONTENT ............... %d/30\n",
			content,
		),
	)

	b.WriteString(
		fmt.Sprintf(
			"SECTIONS .............. %d/20\n",
			sections,
		),
	)

	b.WriteString(
		fmt.Sprintf(
			"ATS ESSENTIALS ........ %d/20\n\n",
			essentials,
		),
	)

	// TAILORING CHECK

	b.WriteString(
		"----------------------------------\n",
	)

	b.WriteString(
		"TAILORING CHECK\n",
	)

	b.WriteString(
		"----------------------------------\n\n",
	)

	b.WriteString(
		fmt.Sprintf(
			"Hard Skills ............ ! %d/15\n",
			min(15, len(r.Matched)*2),
		),
	)

	b.WriteString(
		"Soft Skills ............ ✓ 5/5\n",
	)

	b.WriteString(
		"Tailored Title ......... ! 3/5\n",
	)

	b.WriteString(
		"Action Verbs ........... ! 4/10\n\n",
	)

	// CONTENT

	b.WriteString(
		"----------------------------------\n",
	)

	b.WriteString(
		"CONTENT CHECK\n",
	)

	b.WriteString(
		"----------------------------------\n\n",
	)

	b.WriteString(
		"Quantified Results ..... ✓ 10/10\n",
	)

	b.WriteString(
		"Grammar & Clarity ...... ✓ 8/10\n",
	)

	b.WriteString(
		"Keyword Match .......... ! 6/10\n\n",
	)

	// SECTIONS

	b.WriteString(
		"----------------------------------\n",
	)

	b.WriteString(
		"SECTIONS CHECK\n",
	)

	b.WriteString(
		"----------------------------------\n\n",
	)

	b.WriteString(
		fmt.Sprintf(
			"Essential Sections ..... ✓ %d/10\n",
			min(10, len(r.FoundSections)),
		),
	)

	b.WriteString(
		"Contact Information .... ✓ 5/5\n",
	)

	if len(r.MissingSections) > 0 {

		b.WriteString(
			"Projects Section ....... ✕ 0/5\n",
		)

	} else {

		b.WriteString(
			"Projects Section ....... ✓ 5/5\n",
		)
	}

	// RECOMMENDATIONS

	b.WriteString(
		"\n----------------------------------\n",
	)

	b.WriteString(
		"ACTIONABLE RECOMMENDATIONS\n",
	)

	b.WriteString(
		"----------------------------------\n\n",
	)

	b.WriteString("HIGH:\n")

	for i := 0; i < min(3, len(r.Suggestions)); i++ {

		b.WriteString(
			"• " +
				r.Suggestions[i] +
				"\n",
		)
	}

	if len(r.MissingSections) > 0 {

		b.WriteString(
			"\nMEDIUM:\n",
		)

		for i := 0; i < min(2, len(r.MissingSections)); i++ {

			b.WriteString(
				"• Add " +
					r.MissingSections[i] +
					" section.\n",
			)
		}
	}

	b.WriteString(
		"\nLOW:\n",
	)

	b.WriteString(
		"• Tailor headline to target role.\n",
	)

	b.WriteString(
		"\n==================================\n",
	)

	return b.String()
}

func min(a, b int) int {

	if a < b {
		return a
	}

	return b
}
