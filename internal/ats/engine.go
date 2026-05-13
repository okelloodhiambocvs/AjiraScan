package ats

import "ajirascan/internal/text"

type Result struct {
	Score            int
	Matched          []string
	Missing          []string
	Suggestions      []string
	FoundSections    []string
	MissingSections  []string
	KeywordFrequency []FrequencyReport
	CategoryAnalysis []CategoryReport
	JobType          JobType
	Improvements     []CVImprovement
}

func Analyze(cv, job string) Result {
	// 1. Detect job type FIRST (context layer)
	jobType := DetectJobType(job)

	// 2. Normalize + tokenize inputs
	cvTokens := text.Tokenize(text.Normalize(cv))
	jobTokens := text.Tokenize(text.Normalize(job))

	// 3. Keyword matching
	matched, missing := MatchKeywords(cvTokens, jobTokens)

	// 4. Structural analysis
	sections := DetectSections(cv)

	// 5. Frequency analysis
	frequency := AnalyzeKeywordFrequency(cvTokens)

	// 6. Category analysis
	categories := AnalyzeCategories(cvTokens)

	// 7. Scoring engine
	score := WeightedScore(matched, missing)
	score = ApplyJobContextBoost(score, jobType, matched)

	// 8. Suggestions engine
	suggestions := Suggestions(missing)

	// 9. CV improvement engine (NEW STEP 50)
	improvements := ImproveCV(cv)

	// 10. Return full ATS result
	return Result{
		Score:            score,
		Matched:          matched,
		Missing:          missing,
		Suggestions:      suggestions,
		FoundSections:    sections.Found,
		MissingSections:  sections.Missing,
		KeywordFrequency: frequency,
		CategoryAnalysis: categories,
		JobType:          jobType,
		Improvements:     improvements,
	}
}