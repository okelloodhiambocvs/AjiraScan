package tests

import (
	"testing"

	"smart-cv-checker/internal/ats"
)

func TestAnalyzeBasic(t *testing.T) {
	cv := "Go developer with Docker"
	job := "Go Docker Kubernetes"

	result := ats.Analyze(cv, job)

	if result.Score != 66 {
		t.Errorf("Expected ~66, got %d", result.Score)
	}
}

func TestAnalyzeFullMatch(t *testing.T) {
	cv := "Go Docker Kubernetes"
	job := "Go Docker Kubernetes"

	result := ats.Analyze(cv, job)

	if result.Score != 100 {
		t.Errorf("Expected 100, got %d", result.Score)
	}
}

func TestAnalyzeNoMatch(t *testing.T) {
	cv := "Graphic Designer"
	job := "Go Kubernetes"

	result := ats.Analyze(cv, job)

	if result.Score != 0 {
		t.Errorf("Expected 0, got %d", result.Score)
	}
}

func TestAnalyzeEmpty(t *testing.T) {
	result := ats.Analyze("", "")

	if result.Score != 0 {
		t.Errorf("Expected 0 for empty input")
	}
}