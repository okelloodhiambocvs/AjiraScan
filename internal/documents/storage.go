package documents

import (
	"context"
	"errors"
	"io"
	"time"
)

type ObjectStore interface {
	PutQuarantined(context.Context, string, io.Reader, int64, FileType) error
	SignedDownloadURL(context.Context, string, time.Duration) (string, error)
	Delete(context.Context, string) error
}

type MalwareScanner interface {
	Scan(context.Context, string) error
}

type UnconfiguredStore struct{}

func (UnconfiguredStore) PutQuarantined(context.Context, string, io.Reader, int64, FileType) error {
	return errors.New("object storage is not configured")
}

func (UnconfiguredStore) SignedDownloadURL(context.Context, string, time.Duration) (string, error) {
	return "", errors.New("object storage is not configured")
}

func (UnconfiguredStore) Delete(context.Context, string) error {
	return errors.New("object storage is not configured")
}

type UnconfiguredScanner struct{}

func (UnconfiguredScanner) Scan(context.Context, string) error {
	return errors.New("malware scanner is not configured")
}
