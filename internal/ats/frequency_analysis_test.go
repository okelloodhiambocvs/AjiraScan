package ats

import "testing"

func TestAnalyzeKeywordFrequency(t *testing.T) {
	tokens := []string{
		"go",
		"go",
		"docker",
	}

	report := AnalyzeKeywordFrequency(tokens)

	if len(report) == 0 {
		t.Errorf("expected frequency report")
	}
}