package ats

import "strings"

func MatchPhrases(cv, job string) ([]string, []string) {

	cv = strings.ToLower(cv)
	job = strings.ToLower(job)

	var matched []string
	var missing []string

	for _, phrase := range ATSPhrases {

		if !strings.Contains(job, phrase) {
			continue
		}

		if strings.Contains(cv, phrase) {
			matched = append(matched, phrase)
		} else {
			missing = append(missing, phrase)
		}
	}

	return matched, missing
}
