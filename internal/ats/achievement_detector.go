package ats

import "regexp"

var achievementPattern = regexp.MustCompile(`(?i)(?:\b\d+(?:\.\d+)?\s?%|\b(?:increased|reduced|saved|managed|led|delivered|grew|improved|achieved)\b[^.\n]{0,80}\b\d+)`)

func CountAchievements(cv string) int {

	matches := achievementPattern.FindAllString(cv, -1)

	return len(matches)
}
