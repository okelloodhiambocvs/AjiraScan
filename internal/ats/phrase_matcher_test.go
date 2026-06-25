package ats

import "testing"

func TestPhraseMatching(t *testing.T) {

	cv := `
	Experienced in risk management,
	project management,
	and regulatory compliance.
	`

	job := `
	Looking for candidates with
	risk management,
	project management,
	and customer due diligence.
	`

	matched, missing := MatchPhrases(
		cv,
		job,
	)

	if len(matched) != 2 {
		t.Errorf("expected 2 matched phrases, got %d", len(matched))
	}

	if len(missing) != 1 {
		t.Errorf("expected 1 missing phrase, got %d", len(missing))
	}
}
