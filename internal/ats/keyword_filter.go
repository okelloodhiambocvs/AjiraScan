package ats

var IgnoredKeywords = map[string]bool{
	// Generic hiring words
	"candidate":  true,
	"applicant":  true,
	"position":   true,
	"role":       true,
	"vacancy":    true,
	"job":        true,
	"employment": true,

	// Generic requirement words
	"required":         true,
	"preferred":        true,
	"requirements":     true,
	"qualification":    true,
	"qualifications":   true,
	"responsibility":   true,
	"responsibilities": true,
	"accountability":   true,
	"accountabilities": true,
	"duty":             true,
	"duties":           true,

	// Generic work verbs
	"work":     true,
	"working":  true,
	"manage":   true,
	"managed":  true,
	"support":  true,
	"provide":  true,
	"provided": true,
	"assist":   true,
	"assisting": true,
	"ensure":   true,

	// Geographic / organization noise
	"kenya":  true,
	"africa": true,

	// Common document noise
	"page":         true,
	"pages":        true,
	"resource":     true,
	"resources":    true,
	"reference":    true,
	"references":   true,
	"appendix":     true,
	"annex":        true,
	"draft":        true,
	"directorate":  true,
	"application":  true,
	"applications": true,

	// Low ATS value words
	"practical": true,
	"action":    true,
	"general":   true,
	"various":   true,
	"related":   true,
	"including": true,
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
