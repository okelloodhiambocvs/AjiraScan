package text

func Frequency(tokens []string) map[string]int {
	freq := make(map[string]int)

	for _, t := range tokens {
		freq[t]++
	}

	return freq
}
