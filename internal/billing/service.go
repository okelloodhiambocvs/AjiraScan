package billing

import (
	"context"
	"errors"
)

type Event struct {
	Provider  string
	ID        string
	Payload   []byte
	Signature string
}

type VerifiedEvent struct {
	Event
	TransactionID string
	Status        string
}

type Provider interface {
	VerifyWebhook(context.Context, Event) (VerifiedEvent, error)
	Name() string
}

type EventStore interface {
	RecordAndApply(context.Context, VerifiedEvent) (alreadyProcessed bool, err error)
}

type Service struct {
	Provider Provider
	Store    EventStore
}

func (s Service) ProcessWebhook(ctx context.Context, event Event) error {
	if event.ID == "" || len(event.Payload) == 0 || s.Provider == nil || s.Store == nil {
		return errors.New("invalid payment webhook configuration")
	}
	verified, err := s.Provider.VerifyWebhook(ctx, event)
	if err != nil {
		return err
	}
	if verified.Provider != s.Provider.Name() || verified.ID != event.ID || verified.TransactionID == "" {
		return errors.New("payment webhook verification mismatch")
	}
	_, err = s.Store.RecordAndApply(ctx, verified)
	return err
}

type DisabledProvider struct{}

func (DisabledProvider) VerifyWebhook(context.Context, Event) (VerifiedEvent, error) {
	return VerifiedEvent{}, errors.New("payment provider is disabled")
}

func (DisabledProvider) Name() string { return "disabled" }
