package ats

import "regexp"

func CountAchievements(cv string) int {

	re := regexp.MustCompile(`\d+`)

	matches := re.FindAllString(
		cv,
		-1,
	)

	return len(matches)
}
