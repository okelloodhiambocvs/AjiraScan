package ats

import "testing"

func TestAnalyzeCategories(t *testing.T) {
	tokens := []string{"go", "docker", "communication"}

	result := AnalyzeCategories(tokens)

	if len(result) == 0 {
		t.Errorf("expected category results")
	}
}
