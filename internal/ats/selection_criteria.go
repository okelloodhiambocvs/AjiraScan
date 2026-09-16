package ats

var prohibitedSelectionCriteria = map[string]bool{
	"age":         true,
	"gender":      true,
	"male":        true,
	"female":      true,
	"sex":         true,
	"race":        true,
	"ethnicity":   true,
	"religion":    true,
	"nationality": true,
	"tribe":       true,
	"marital":     true,
	"pregnant":    true,
	"pregnancy":   true,
	"disability":  true,
}

func FilterSelectionCriteria(tokens []string) (permitted []string, excluded []string) {
	for _, token := range tokens {
		if prohibitedSelectionCriteria[token] {
			excluded = append(excluded, token)
			continue
		}
		permitted = append(permitted, token)
	}
	return permitted, excluded
}
