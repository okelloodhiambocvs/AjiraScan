package text

import "testing"

func TestNormalize(t *testing.T) {
	out := Normalize("Go Developer!!!\nDocker")
	expected := "go developer docker"

	if out != expected {
		t.Errorf("expected '%s', got '%s'", expected, out)
	}
}

func TestTokenize(t *testing.T) {
	tokens := Tokenize("go docker backend")

	if len(tokens) != 3 {
		t.Errorf("expected 3 tokens")
	}
}

func TestFrequency(t *testing.T) {
	freq := Frequency([]string{"go", "go", "docker"})

	if freq["go"] != 2 {
		t.Errorf("expected go count = 2")
	}
}
