package ats

import "ajirascan/internal/text"

type FrequencyReport struct {
	Keyword string
	Count   int
}

func AnalyzeKeywordFrequency(tokens []string) []FrequencyReport {
	freqMap := text.Frequency(tokens)

	var report []FrequencyReport

	for keyword, count := range freqMap {
		report = append(report, FrequencyReport{
			Keyword: keyword,
			Count:   count,
		})
	}

	return report
}