package ats

import "testing"

func TestMatcher(t *testing.T) {
	cv := []string{"go", "docker"}
	job := []string{"go", "kubernetes"}

	matched, missing := MatchKeywords(cv, job)

	if len(matched) != 1 {
		t.Errorf("expected 1 match, got %d", len(matched))
	}

	if len(missing) != 1 {
		t.Errorf("expected 1 missing, got %d", len(missing))
	}
}
