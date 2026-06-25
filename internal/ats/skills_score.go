package ats

func CalculateSkillsScore(
	matched,
	missing,
	matchedPhrases,
	missingPhrases []string,
) int {

	totalMatches :=
		len(matched) +
			(len(matchedPhrases) * 3)

	totalMissing :=
		len(missing) +
			(len(missingPhrases) * 3)

	total := totalMatches + totalMissing

	if total == 0 {
		return 0
	}

	return (totalMatches * 100) / total
}
