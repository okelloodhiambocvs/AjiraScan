package tests

import (
	"testing"

	"smart-cv-checker/library"
)

func TestAnalyzeCV_BasicMatch(t *testing.T) {
	cv := "I developed APIs using Go and Docker"
	job := "Looking for Go developer with Docker experience"

	result := library.AnalyzeCV(cv, job)

	if result.MatchScore == 0 {
		t.Errorf("Expected non-zero match score")
	}

	if len(result.MatchedKeywords) == 0 {
		t.Errorf("Expected matched keywords")
	}
}

func TestAnalyzeCV_MissingKeywords(t *testing.T) {
	cv := "I am a graphic designer"
	job := "Looking for Go developer with Kubernetes"

	result := library.AnalyzeCV(cv, job)

	if len(result.MissingKeywords) == 0 {
		t.Errorf("Expected missing keywords")
	}
}

func TestAnalyzeCV_CopyPastePenalty(t *testing.T) {
	text := "We are looking for a Go developer with Docker experience"

	result := library.AnalyzeCV(text, text)

	if result.AuthenticityScore >= 100 {
		t.Errorf("Expected lower authenticity score for copied content")
	}
}

func TestAnalyzeCV_ActionVerbBoost(t *testing.T) {
	cv := "I built and developed backend systems using Go"
	job := "Go backend developer"

	result := library.AnalyzeCV(cv, job)

	if result.MatchScore <= 0 {
		t.Errorf("Expected boosted score from action verbs")
	}
}