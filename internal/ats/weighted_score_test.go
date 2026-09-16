package ats

import "testing"

func TestWeightedScore(t *testing.T) {
	score := WeightedScore(
		[]string{"go", "docker"},
		[]string{"kubernetes"},
	)

	if score <= 0 {
		t.Errorf("expected positive score")
	}
}
