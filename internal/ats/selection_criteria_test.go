package ats

import "testing"

func TestAnalyzeExcludesProtectedSelectionTermsFromScore(t *testing.T) {
	result := Analyze("Go developer", "Go developer female nationality")
	if result.Score != 100 {
		t.Fatalf("expected protected terms not to reduce score, got %d", result.Score)
	}
	if len(result.ExcludedSelectionCriteria) != 2 {
		t.Fatalf("expected excluded criteria to be disclosed, got %#v", result.ExcludedSelectionCriteria)
	}
}
