package ats

func CalculateSectionScore(report SectionReport) int {

	totalSections := len(CommonSections)

	if totalSections == 0 {
		return 0
	}

	return (len(report.Found) * 100) / totalSections
}
