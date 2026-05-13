package ats

import "strings"

type CVImprovement struct {
	Original   string
	Improved   string
	Reason     string
}

func ImproveCV(cv string) []CVImprovement {
	lines := strings.Split(cv, "\n")

	var improvements []CVImprovement

	for _, line := range lines {
		l := strings.ToLower(line)

		switch {
		case strings.Contains(l, "experience") && strings.Contains(l, "developer"):
			improvements = append(improvements, CVImprovement{
				Original: line,
				Improved: "Developed scalable backend systems using Go and Docker to support high-performance applications.",
				Reason:   "Strengthens technical impact and adds measurable value",
			})

		case strings.Contains(l, "communication"):
			improvements = append(improvements, CVImprovement{
				Original: line,
				Improved: "Led cross-functional communication efforts to improve collaboration between technical and non-technical teams.",
				Reason:   "Makes communication experience more professional and role-relevant",
			})

		case strings.Contains(l, "team"):
			improvements = append(improvements, CVImprovement{
				Original: line,
				Improved: "Collaborated effectively within cross-functional teams to deliver project objectives on time.",
				Reason:   "Rewrites vague teamwork into impact-driven statement",
			})

		default:
			// skip irrelevant lines
		}
	}

	return improvements
}