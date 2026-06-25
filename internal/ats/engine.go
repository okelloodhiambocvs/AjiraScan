package ats

import "ajirascan/internal/text"

type Result struct {
	Score int

	Breakdown ScoreBreakdown

	Matched []string
	Missing []string

	Suggestions []string

	FoundSections   []string
	MissingSections []string

	KeywordFrequency []FrequencyReport
	CategoryAnalysis []CategoryReport

	JobType JobType

	Improvements []CVImprovement

	// PROFESSIONAL REPORT DATA
	IssueCount int

	TailoringScore int
	ContentScore   int
	SectionScore   int
	ATSScore       int
}

func Analyze(cv, job string) Result {

	jobType := DetectJobType(job)

	cvTokens := text.Tokenize(text.Normalize(cv))
	jobTokens := text.Tokenize(text.Normalize(job))

	matched, missing := MatchKeywords(
		cvTokens,
		jobTokens,
	)

	// PHRASE MATCHING
	matchedPhrases, missingPhrases := MatchPhrases(
		cv,
		job,
	)

	sections := DetectSections(cv)

	frequency := AnalyzeKeywordFrequency(cvTokens)

	categories := AnalyzeCategories(cvTokens)

	score := WeightedScore(
		matched,
		missing,
	)

	score = ApplyJobContextBoost(
		score,
		jobType,
		matched,
	)

	// STANDARD KEYWORD SUGGESTIONS
	suggestions := Suggestions(missing)

	// PHRASE-BASED SUGGESTIONS (HIGHER VALUE)
	phraseSuggestions := PhraseSuggestions(
		missingPhrases,
	)

	// PRIORITIZE PHRASE SUGGESTIONS
	suggestions = append(
		phraseSuggestions,
		suggestions...,
	)

	// LIMIT TOTAL RECOMMENDATIONS
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	improvements := ImproveCV(cv)

	// PROFESSIONAL SCORING

	tailoring := 18
	content := 24
	atsEssentials := 11

	if len(matched) > 10 {
		tailoring += 5
	}

	if score >= 80 {
		content += 4
	}

	// DYNAMIC SECTION SCORE
	sectionScore := CalculateSectionScore(sections)

	// PHRASE-AWARE SKILLS SCORE
	skillsScore := CalculateSkillsScore(
		matched,
		missing,
		matchedPhrases,
		missingPhrases,
	)

	// ACTIONABLE ISSUE COUNT
	keywordIssues := len(missing)

	if keywordIssues > 5 {
		keywordIssues = 5
	}

	issues :=
		keywordIssues +
			len(sections.Missing)

	// SCORE BREAKDOWN
	breakdown := ScoreBreakdown{
		OverallScore: score,
		KeywordScore: score,
		SectionScore: sectionScore,
		SkillsScore:  skillsScore,
	}

	return Result{
		Score: score,

		Breakdown: breakdown,

		Matched: matched,
		Missing: missing,

		Suggestions: suggestions,

		FoundSections:   sections.Found,
		MissingSections: sections.Missing,

		KeywordFrequency: frequency,
		CategoryAnalysis: categories,

		JobType: jobType,

		Improvements: improvements,

		IssueCount: issues,

		TailoringScore: tailoring,
		ContentScore:   content,
		SectionScore:   sectionScore,
		ATSScore:       atsEssentials,
	}
}

func GeneratePDFReport(
	cv string,
	result Result,
	outputPath string,
) error {

	return ExportCVToPDF(
		cv,
		result.Improvements,
		outputPath,
	)
}
