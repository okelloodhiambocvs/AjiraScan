package ats

import "testing"

func TestJobContextBoostIsCapped(t *testing.T) {
	matched := []string{"go", "docker", "kubernetes", "golang", "api", "backend", "microservices"}
	if got := ApplyJobContextBoost(99, TechJob, matched); got != 100 {
		t.Fatalf("expected capped score 100, got %d", got)
	}
}

func TestMatchKeywordsDeduplicatesJobRequirements(t *testing.T) {
	matched, missing := MatchKeywords([]string{"go"}, []string{"go", "go", "docker", "docker"})
	if len(matched) != 1 || len(missing) != 1 {
		t.Fatalf("expected unique job requirements, got matched=%v missing=%v", matched, missing)
	}
}
