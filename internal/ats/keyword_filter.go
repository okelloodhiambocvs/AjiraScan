package ats

var IgnoredKeywords = map[string]bool{
	"kenya":            true,
	"practical":        true,
	"action":           true,
	"accountabilities": true,
	"responsibilities": true,
	"duties":           true,
	"role":             true,
	"position":         true,
	"candidate":        true,
	"applicant":        true,
	"work":             true,
	"working":          true,
	"manage":           true,
	"support":          true,
	"provide":          true,
	"required":         true,
	"preferred":        true,
}

func FilterKeywords(keywords []string) []string {

	var result []string

	for _, k := range keywords {

		if IgnoredKeywords[k] {
			continue
		}

		result = append(result, k)
	}

	return result
}
