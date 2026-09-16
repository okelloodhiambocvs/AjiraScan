package text

import (
	"regexp"
	"strings"
)

func Normalize(input string) string {
	input = strings.ToLower(input)
	input = strings.ReplaceAll(input, "\n", " ")

	re := regexp.MustCompile(`[^\w\s]`)
	input = re.ReplaceAllString(input, "")

	return strings.TrimSpace(input)
}
