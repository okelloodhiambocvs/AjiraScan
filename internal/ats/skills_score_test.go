package ats

import "testing"

func TestSkillsScoreWithPhrases(t *testing.T) {

	score := CalculateSkillsScore(
		[]string{"go", "docker"},
		[]string{"kubernetes"},
		[]string{"project management"},
		[]string{"risk management"},
	)

	if score <= 0 {
		t.Errorf("expected positive score")
	}
}
