package ats

import "ajirascan/internal/text"

type Result struct {
	Score   int
	Matched []string
	Missing []string
}

func Analyze(cv, job string) Result {
	cvTokens := text.Tokenize(text.Normalize(cv))
	jobTokens := text.Tokenize(text.Normalize(job))

	matched, missing := MatchKeywords(cvTokens, jobTokens)

	score := Score(matched, missing)

	return Result{
		Score:   score,
		Matched: matched,
		Missing: missing,
	}
}