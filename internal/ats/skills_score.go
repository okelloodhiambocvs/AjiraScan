package ats

func CalculateSkillsScore(matched, missing []string) int {

	total := len(matched) + len(missing)

	if total == 0 {
		return 0
	}

	return (len(matched) * 100) / total
}
