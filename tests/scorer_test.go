package tests

import (
	"testing"

	"smart-cv-checker/internal/ats"
)

func TestScoreBasic(t *testing.T) {
	score := ats.Score([]string{"go"}, []string{"docker"})

	if score != 50 {
		t.Errorf("Expected 50, got %d", score)
	}
}

func TestScoreAllMatched(t *testing.T) {
	score := ats.Score([]string{"go", "docker"}, []string{})

	if score != 100 {
		t.Errorf("Expected 100, got %d", score)
	}
}

func TestScoreNoneMatched(t *testing.T) {
	score := ats.Score([]string{}, []string{"go", "docker"})

	if score != 0 {
		t.Errorf("Expected 0, got %d", score)
	}
}

func TestScoreEmpty(t *testing.T) {
	score := ats.Score([]string{}, []string{})

	if score != 0 {
		t.Errorf("Expected 0 for empty input")
	}
}