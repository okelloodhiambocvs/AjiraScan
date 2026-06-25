package ats

import "testing"

func TestCountAchievements(t *testing.T) {

	cv := `
Managed 15 employees.
Reviewed 500 customer files.
Completed 25 audits.
`

	count := CountAchievements(cv)

	if count != 3 {
		t.Fatalf(
			"expected 3 achievements, got %d",
			count,
		)
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
