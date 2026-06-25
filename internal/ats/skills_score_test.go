package ats

import "testing"

func TestCalculateSkillsScore(t *testing.T) {

	matched := []string{
		"excel",
		"sql",
		"python",
	}

	missing := []string{
		"powerbi",
	}

	score := CalculateSkillsScore(
		matched,
		missing,
	)

	if score != 75 {
		t.Errorf("expected 75, got %d", score)
	}
}
