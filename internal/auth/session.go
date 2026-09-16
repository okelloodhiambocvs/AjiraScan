package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	TokenHash []byte
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func NewSessionToken() (string, []byte, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, err
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	return token, sum[:], nil
}

func HashSessionToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}
