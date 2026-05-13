package ats

import "testing"

func TestSuggestions(t *testing.T) {
	missing := []string{"docker", "leadership"}

	suggestions := Suggestions(missing)

	if len(suggestions) != 2 {
		t.Errorf("expected 2 suggestions")
	}
}