package ats

import (
	"sort"

	"ajirascan/internal/text"
)

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

	keywords := make([]string, 0, len(freqMap))
	for keyword := range freqMap {
		keywords = append(keywords, keyword)
	}
	sort.Strings(keywords)

	report := make([]FrequencyReport, 0, len(keywords))
	for _, keyword := range keywords {
		count := freqMap[keyword]
		report = append(report, FrequencyReport{
			Keyword: keyword,
			Count:   count,
		})
	}

	return report
}
