package ats

import "testing"

func TestNormalizeSynonym(t *testing.T) {

	result := NormalizeSynonym("aml")

	if result != "anti_money_laundering" {
		t.Errorf(
			"expected anti_money_laundering, got %s",
			result,
		)
	}
}
