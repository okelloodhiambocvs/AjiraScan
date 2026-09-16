package ats

import "testing"

func TestEngine_BasicMatch(t *testing.T) {
	res := Analyze(
		"Go developer with Docker",
		"Go Docker Kubernetes Backend",
	)

	if res.Score == 0 {
		t.Errorf("expected non-zero score")
	}

	if len(res.Matched) == 0 {
		t.Errorf("expected matched keywords")
	}
}
