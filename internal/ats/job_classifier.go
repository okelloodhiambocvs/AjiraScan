package ats

import "ajirascan/internal/text"

type JobType string

const (
	TechJob          JobType = "TECH"
	NGOJob           JobType = "NGO"
	CommunicationJob JobType = "COMMUNICATION"
	BusinessJob      JobType = "BUSINESS"
	UnknownJob       JobType = "UNKNOWN"
)

var jobIndicators = map[JobType][]string{
	TechJob: {
		"go", "golang", "docker", "kubernetes", "api", "backend", "microservices", "cloud",
	},

	NGOJob: {
		"fundraising", "advocacy", "donor", "ngo", "humanitarian", "community", "development",
	},

	CommunicationJob: {
		"communication", "media", "journalism", "content", "writing", "public relations",
	},

	BusinessJob: {
		"marketing", "sales", "strategy", "operations", "management", "business",
	},
}

func DetectJobType(jobText string) JobType {
	jobText = text.Normalize(jobText)
	tokens := text.Tokenize(jobText)

	scoreMap := make(map[JobType]int)

	for _, token := range tokens {
		for jobType, indicators := range jobIndicators {
			for _, indicator := range indicators {
				if token == indicator {
					scoreMap[jobType]++
				}
			}
		}
	}

	var bestType JobType = UnknownJob
	max := 0

	for jt, score := range scoreMap {
		if score > max {
			max = score
			bestType = jt
		}
	}

	return bestType
}
