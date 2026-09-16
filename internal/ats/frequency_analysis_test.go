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

func TestAnalyzeKeywordFrequencyUsesStableOrder(t *testing.T) {
	report := AnalyzeKeywordFrequency([]string{"docker", "go", "go"})
	if len(report) != 2 || report[0].Keyword != "docker" || report[1].Keyword != "go" {
		t.Fatalf("expected sorted frequency report, got %#v", report)
	}
}
