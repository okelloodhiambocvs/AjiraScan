package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 2
	argonKeyLen  uint32 = 32
)

func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return fmt.Sprintf("argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonTime, argonThreads, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash)), nil
}

func VerifyPassword(encoded, password string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[0] != "argon2id" || parts[1] != "v=19" {
		return false, errors.New("invalid password hash format")
	}
	var memory, timeCost uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[2], "m=%d,t=%d,p=%d", &memory, &timeCost, &threads); err != nil {
		return false, errors.New("invalid password hash parameters")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, errors.New("invalid password hash salt")
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("invalid password hash digest")
	}
	actual := argon2.IDKey([]byte(password), salt, timeCost, memory, threads, uint32(len(expected)))
	return subtle.ConstantTimeCompare(expected, actual) == 1, nil
}

func ValidatePassword(password string) error {
	if len(password) < 12 || len(password) > 1024 {
		return errors.New("password must be between 12 and 1024 bytes")
	}
	var upper, lower, digit bool
	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			upper = true
		case r >= 'a' && r <= 'z':
			lower = true
		case r >= '0' && r <= '9':
			digit = true
		}
	}
	if !upper || !lower || !digit {
		return errors.New("password must include upper-case, lower-case and digit characters")
	}
	return nil
}
