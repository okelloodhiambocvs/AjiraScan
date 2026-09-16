package ats

import "testing"

func TestAnalyzeCategories(t *testing.T) {
	tokens := []string{"go", "docker", "communication"}

	result := AnalyzeCategories(tokens)

	if len(result) == 0 {
		t.Errorf("expected category results")
	}
}

func TestAnalyzeCategoriesUsesStableOrder(t *testing.T) {
	result := AnalyzeCategories([]string{"go", "docker", "communication"})
	for index := 1; index < len(result); index++ {
		if result[index-1].Category > result[index].Category {
			t.Fatal("expected categories to be sorted")
		}
	}
}
