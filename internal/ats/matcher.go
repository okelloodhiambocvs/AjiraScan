package ats

func MatchKeywords(cvTokens, jobTokens []string) ([]string, []string) {

	cvSet := make(map[string]bool)

	for _, t := range cvTokens {
		cvSet[NormalizeSynonym(t)] = true
	}

	matched := []string{}
	missing := []string{}
	seen := make(map[string]bool)

	for _, j := range jobTokens {

		normalized := NormalizeSynonym(j)
		if seen[normalized] {
			continue
		}
		seen[normalized] = true

		if cvSet[normalized] {
			matched = append(matched, j)
		} else {
			missing = append(missing, j)
		}
	}

	return matched, missing
}
