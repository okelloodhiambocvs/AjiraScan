package ats

import "ajirascan/internal/text"

type Result struct {
	Score int

	Matched []string
	Missing []string

	Suggestions []string

	FoundSections   []string
	MissingSections []string

	KeywordFrequency []FrequencyReport
	CategoryAnalysis []CategoryReport

	JobType JobType

	Improvements []CVImprovement

	// NEW PROFESSIONAL REPORT DATA
	IssueCount int

	TailoringScore int
	ContentScore   int
	SectionScore   int
	ATSScore        int
}

func Analyze(cv, job string) Result {

	jobType := DetectJobType(job)

	cvTokens := text.Tokenize(text.Normalize(cv))
	jobTokens := text.Tokenize(text.Normalize(job))

	matched, missing := MatchKeywords(
		cvTokens,
		jobTokens,
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

	suggestions := Suggestions(missing)

	improvements := ImproveCV(cv)

	// PROFESSIONAL SCORING

	tailoring := 18
	content := 24
	sectionScore := 15
	atsEssentials := 11

	if len(matched) > 10 {
		tailoring += 5
	}

	if len(sections.Found) >= 6 {
		sectionScore += 3
	}

	if score >= 80 {
		content += 4
	}

	issues :=
		len(missing) +
			len(sections.Missing)

	return Result{
		Score: score,

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