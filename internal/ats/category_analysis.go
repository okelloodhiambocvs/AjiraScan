package ats

type CategoryReport struct {
	Category string
	Matched  []string
	Score    int
}

func AnalyzeCategories(tokens []string) []CategoryReport {
	tokenSet := make(map[string]bool)

	for _, t := range tokens {
		tokenSet[t] = true
	}

	var reports []CategoryReport

	for category, keywords := range SkillCategories {
		var matched []string

		for _, k := range keywords {
			if tokenSet[k] {
				matched = append(matched, k)
			}
		}

		score := 0
		if len(keywords) > 0 {
			score = (len(matched) * 100) / len(keywords)
		} else {
			score = 0
		}

		if len(matched) == 0 {
			score = 0
		}

		reports = append(reports, CategoryReport{
			Category: category,
			Matched:  matched,
			Score:    score,
		})
	}

	return reports
}