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

	// NEW
	AchievementCount int
}

func Analyze(cv, job string) Result {

	jobType := DetectJobType(job)

	cvTokens := text.Tokenize(
		text.Normalize(cv),
	)

	jobTokens := text.Tokenize(
		text.Normalize(job),
	)

	matched, missing := MatchKeywords(
		cvTokens,
		jobTokens,
	)

	// FILTER LOW-VALUE WORDS
	matched = FilterKeywords(matched)
	missing = FilterKeywords(missing)

	// PHRASE MATCHING
	matchedPhrases, missingPhrases := MatchPhrases(
		cv,
		job,
	)

	sections := DetectSections(cv)

	frequency := AnalyzeKeywordFrequency(
		cvTokens,
	)

	categories := AnalyzeCategories(
		cvTokens,
	)

	score := WeightedScore(
		matched,
		missing,
	)

	score = ApplyJobContextBoost(
		score,
		jobType,
		matched,
	)

	// =========================
	// RECOMMENDATIONS
	// =========================

	phraseSuggestions := PhraseSuggestions(
		missingPhrases,
	)

	var suggestions []string

	if len(phraseSuggestions) > 0 {
		suggestions = phraseSuggestions
	} else {
		suggestions = Suggestions(missing)
	}

	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	improvements := ImproveCV(cv)

	// =========================
	// ACHIEVEMENTS
	// =========================

	achievementCount := CountAchievements(cv)

	// =========================
	// TAILORING SCORE (30)
	// =========================

	tailoring := 15

	if len(matched) >= 5 {
		tailoring += 5
	}

	if len(matched) >= 10 {
		tailoring += 5
	}

	if len(matchedPhrases) >= 2 {
		tailoring += 5
	}

	if tailoring > 30 {
		tailoring = 30
	}

	// =========================
	// CONTENT SCORE (30)
	// =========================

	content := 18

	if achievementCount >= 3 {
		content += 12
	} else if achievementCount >= 1 {
		content += 6
	}

	if score >= 80 {
		content += 4
	}

	if content > 30 {
		content = 30
	}

	// =========================
	// SECTION SCORE (100)
	// =========================

	sectionScore := CalculateSectionScore(
		sections,
	)

	// =========================
	// ATS ESSENTIALS (20)
	// =========================

	atsEssentials := 10

	if len(sections.Found) >= 5 {
		atsEssentials += 4
	}

	if achievementCount >= 3 {
		atsEssentials += 3
	}

	if len(matchedPhrases) >= 2 {
		atsEssentials += 3
	}

	if atsEssentials > 20 {
		atsEssentials = 20
	}

	// =========================
	// SKILLS SCORE
	// =========================

	skillsScore := CalculateSkillsScore(
		matched,
		missing,
		matchedPhrases,
		missingPhrases,
	)

	// =========================
	// ISSUE COUNT
	// =========================

	issues := 0

	if len(missing) > 0 {
		issues += min(3, len(missing))
	}

	if len(missingPhrases) > 0 {
		issues += min(2, len(missingPhrases))
	}

	if len(sections.Missing) > 0 {
		issues += min(3, len(sections.Missing))
	}

	// =========================
	// SCORE BREAKDOWN
	// =========================

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

		AchievementCount: achievementCount,
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
