package ats

import "ajirascan/internal/text"

type FrequencyReport struct {
	Keyword string
	Count   int
}

func AnalyzeKeywordFrequency(tokens []string) []FrequencyReport {
	freqMap := make(map[string]int)

	for _, t := range tokens {
		if text.IsRelevantToken(t) && !text.StopWords[t] {
			freqMap[t]++
		}
	}

	var report []FrequencyReport

	for keyword, count := range freqMap {
		report = append(report, FrequencyReport{
			Keyword: keyword,
			Count:   count,
		})
	}

	return report
}
