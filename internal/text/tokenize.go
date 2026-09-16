package text

import "strings"

func Tokenize(input string) []string {
	words := strings.Fields(input)

	var cleaned []string

	for _, w := range words {
		if IsRelevantToken(w) && !StopWords[w] {
			cleaned = append(cleaned, w)
		}
	}

	return cleaned
}
