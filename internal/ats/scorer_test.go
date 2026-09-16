package ats

import "testing"

func TestScore(t *testing.T) {
	score := Score(
		[]string{"go"},
		[]string{"docker"},
	)

	if score != 50 {
		t.Errorf("expected 50, got %d", score)
	}
}
