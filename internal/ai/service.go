package ai

import (
	"context"
	"errors"
	"strings"
	"time"
)

const MaxInputBytes = 64 << 10

type Request struct {
	Operation string
	CVText    string
	JobText   string
	Locale    string
}

type Output struct {
	Score               int
	MatchedRequirements []string
	MissingRequirements []string
	RelevantExperience  []string
	Skills              []string
	Explanation         string
	Limitations         []string
	ModelProvider       string
	ModelName           string
}

type Provider interface {
	Analyze(context.Context, Request) (Output, error)
	Name() string
}

type Service struct {
	Provider Provider
	Timeout  time.Duration
}

func (s Service) Analyze(ctx context.Context, request Request) (Output, error) {
	if err := validateRequest(request); err != nil {
		return Output{}, err
	}
	if s.Provider == nil {
		return Output{}, errors.New("AI provider is not configured")
	}
	timeout := s.Timeout
	if timeout == 0 {
		timeout = 20 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	output, err := s.Provider.Analyze(ctx, request)
	if err != nil {
		return Output{}, err
	}
	if err := validateOutput(output); err != nil {
		return Output{}, err
	}
	output.ModelProvider = s.Provider.Name()
	return output, nil
}

func validateRequest(request Request) error {
	if request.Operation == "" || len(request.CVText) == 0 || len(request.CVText) > MaxInputBytes || len(request.JobText) > MaxInputBytes {
		return errors.New("invalid AI request size or operation")
	}
	if strings.Contains(strings.ToLower(request.CVText), "ignore previous instructions") {
		return errors.New("untrusted document contains prompt-injection pattern")
	}
	return nil
}

func validateOutput(output Output) error {
	if output.Score < 0 || output.Score > 100 {
		return errors.New("AI output score is outside 0-100")
	}
	if len(output.Explanation) > 12000 {
		return errors.New("AI output is too large")
	}
	return nil
}

type DisabledProvider struct{}

func (DisabledProvider) Analyze(context.Context, Request) (Output, error) {
	return Output{}, errors.New("AI provider is disabled")
}

func (DisabledProvider) Name() string { return "disabled" }
