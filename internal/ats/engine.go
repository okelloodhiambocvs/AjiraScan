package ats

import "ajirascan/internal/text"

type Result struct {
	Score              int
	Matched            []string
	Missing            []string
	Suggestions        []string
	FoundSections      []string
	MissingSections    []string
	KeywordFrequency   []FrequencyReport
}

func Analyze(cv, job string) Result {
	cvTokens := text.Tokenize(text.Normalize(cv))
	jobTokens := text.Tokenize(text.Normalize(job))

	matched, missing := MatchKeywords(cvTokens, jobTokens)

	score := WeightedScore(matched, missing)

	suggestions := Suggestions(missing)
	
	sections := DetectSections(cv)
	
	frequency := AnalyzeKeywordFrequency(cvTokens)

	return Result{
	Score:            score,
	Matched:          matched,
	Missing:          missing,
	Suggestions:      suggestions,
	FoundSections:    sections.Found,
	MissingSections:  sections.Missing,
	KeywordFrequency: frequency,
}
}