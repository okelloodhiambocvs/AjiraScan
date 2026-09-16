package ats

import "testing"

func TestDetectSections(t *testing.T) {
	cv := `
	Summary
	Skills
	Experience
	Education
	`

	report := DetectSections(cv)

	if len(report.Found) != 4 {
		t.Errorf("expected 4 found sections")
	}
}
