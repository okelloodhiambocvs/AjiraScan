package ats

import "fmt"

func Suggestions(missing []string) []string {

	var suggestions []string

	for _, keyword := range missing {

		suggestions = append(
			suggestions,
			fmt.Sprintf(
				"Add skill or experience related to: %s",
				keyword,
			),
		)

		if len(suggestions) >= 5 {
			break
		}
	}

	return suggestions
}

func SectionSuggestions(
	missingSections []string,
) []string {

	var suggestions []string

	for _, section := range missingSections {

		suggestions = append(
			suggestions,
			fmt.Sprintf(
				"Add '%s' section",
				section,
			),
		)

		if len(suggestions) >= 5 {
			break
		}
	}

	return suggestions
}