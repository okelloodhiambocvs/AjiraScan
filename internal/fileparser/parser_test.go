package fileparser

import (
	"strings"
	"testing"

	"ajirascan/internal/documents"
)

func TestReadTXTRejectsEmptyAndOversizedInput(t *testing.T) {
	if _, err := ReadTXT(strings.NewReader(" \n")); err != ErrNoExtractableText {
		t.Fatalf("expected empty-text error, got %v", err)
	}
	overLimit := strings.Repeat("a", documents.MaxExtractedTextBytes+1)
	if _, err := ReadTXT(strings.NewReader(overLimit)); err == nil {
		t.Fatal("expected oversized text rejection")
	}
}

func TestReadTXTReturnsBoundedText(t *testing.T) {
	content, err := ReadTXT(strings.NewReader("Candidate experience"))
	if err != nil || content != "Candidate experience" {
		t.Fatalf("unexpected text parsing result: %q, %v", content, err)
	}
}
