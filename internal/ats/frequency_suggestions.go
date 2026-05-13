package ats

import "fmt"

func FrequencySuggestions(report []FrequencyReport) []string {
	var suggestions []string

	for _, item := range report {
		if item.Count == 1 {
			suggestions = append(
				suggestions,
				fmt.Sprintf(
					"Keyword '%s' appears only once; consider reinforcing it naturally across your CV",
					item.Keyword,
				),
			)
		}
	}

	return suggestions
}