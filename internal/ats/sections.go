package ats

import "strings"

type SectionReport struct {
	Found   []string
	Missing []string
}

var CommonSections = []string{
	"summary",
	"profile",
	"skills",
	"experience",
	"education",
	"certifications",
	"projects",
	"leadership",
}

func DetectSections(cv string) SectionReport {
	cv = strings.ToLower(cv)

	var found []string
	var missing []string

	for _, section := range CommonSections {
		if strings.Contains(cv, section) {
			found = append(found, section)
		} else {
			missing = append(missing, section)
		}
	}

	return SectionReport{
		Found:   found,
		Missing: missing,
	}
}
