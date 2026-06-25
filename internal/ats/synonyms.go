package ats

var Synonyms = map[string]string{
	"aml": "anti_money_laundering",
	"anti_money_laundering": "anti_money_laundering",

	"kyc": "know_your_customer",
	"know_your_customer": "know_your_customer",

	"hr": "human_resources",
	"human_resources": "human_resources",

	"it": "information_technology",
	"information_technology": "information_technology",
}

func NormalizeSynonym(token string) string {

	if normalized, exists := Synonyms[token]; exists {
		return normalized
	}

	return token
}
