package ats

import "fmt"

func Suggestions(missing []string) []string {
	var suggestions []string

	for _, keyword := range missing {
		suggestions = append(
			suggestions,
			fmt.Sprintf("Consider adding experience or skills related to '%s'", keyword),
		)
	}

	return suggestions
}