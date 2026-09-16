package privacy

import (
	"context"
	"errors"
)

type Exporter interface {
	ExportUser(context.Context, string) ([]byte, error)
}

type DeletionRepository interface {
	CreateDeletionRequest(context.Context, string) error
}

type Service struct {
	Exporter Exporter
	Deletion DeletionRepository
}

func (s Service) Export(ctx context.Context, userID string) ([]byte, error) {
	if userID == "" || s.Exporter == nil {
		return nil, errors.New("privacy export is unavailable")
	}
	return s.Exporter.ExportUser(ctx, userID)
}

func (s Service) RequestDeletion(ctx context.Context, userID string) error {
	if userID == "" || s.Deletion == nil {
		return errors.New("deletion requests are unavailable")
	}
	return s.Deletion.CreateDeletionRequest(ctx, userID)
}
