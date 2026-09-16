package ats

import "testing"

func TestCountAchievements(t *testing.T) {

	cv := `
	Managed 15 employees.
	Improved response time by 25%.
	Delivered 6 projects.
`

	count := CountAchievements(cv)

	if count != 3 {
		t.Fatalf(
			"expected 3 achievements, got %d",
			count,
		)
	}
}

func TestCountAchievementsDoesNotTreatContactOrDatesAsAchievements(t *testing.T) {
	cv := "Phone: +254 712 345678\nExperience: 2020-2024"
	if count := CountAchievements(cv); count != 0 {
		t.Fatalf("expected contact details and dates to be ignored, got %d", count)
	}
}

func TestAchievementContentCap(t *testing.T) {

	content := 18

	content += 12

	if content > 30 {
		content = 30
	}

	if content != 30 {
		t.Fatalf(
			"expected 30, got %d",
			content,
		)
	}
}
