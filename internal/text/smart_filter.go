package text

var noiseWords = map[string]bool{
	"experience": true,
	"strong":     true,
	"skills":     true,
	"skill":      true,
	"building":   true,
	"worked":     true,
	"work":       true,
	"good":       true,
	"team":       true,
	"teams":      true,
	"ability":    true,
	"responsible": true,
	"developed":  true,
	"develop":    true,
	"using":      true,
	"based":      true,
}

func IsRelevantToken(token string) bool {
	if noiseWords[token] {
		return false
	}

	if len(token) < 2 {
		return false
	}

	return true
}