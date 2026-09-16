package ats

func WeightedScore(matched, missing []string) int {
	matchedScore := 0
	totalScore := 0

	for _, keyword := range matched {
		weight := getWeight(keyword)
		matchedScore += weight
		totalScore += weight
	}

	for _, keyword := range missing {
		totalScore += getWeight(keyword)
	}

	if totalScore == 0 {
		return 0
	}

	if len(matched)+len(missing) == 0 {
		return 0
	}

	return (matchedScore * 100) / totalScore
}

func getWeight(keyword string) int {
	weight, exists := KeywordWeights[keyword]

	if exists {
		return weight
	}

	return 1
}

func ApplyJobContextBoost(score int, jobType JobType, matched []string) int {
	boost := 0

	switch jobType {

	case TechJob:
		for _, k := range matched {
			if KeywordWeights[k] >= 5 {
				boost += 2
			}
		}

	case NGOJob:
		for _, k := range matched {
			if k == "fundraising" || k == "advocacy" {
				boost += 3
			}
		}

	case CommunicationJob:
		for _, k := range matched {
			if k == "communication" || k == "media" {
				boost += 2
			}
		}
	}

	if score+boost > 100 {
		return 100
	}
	return score + boost
}
