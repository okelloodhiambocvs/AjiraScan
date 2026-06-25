package ats

import "testing"

func TestCalculateSectionScore(t *testing.T) {

	report := SectionReport{
		Found: []string{
			"summary",
			"skills",
			"experience",
			"education",
		},
	}

	score := CalculateSectionScore(report)

	if score != 50 {
		t.Errorf("expected 50, got %d", score)
	}
}
