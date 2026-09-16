package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateAccount   = errors.New("account already exists")
)

type Registration struct {
	AccountType      string
	FirstName        string
	LastName         string
	Email            string
	Password         string
	OrganizationName string
}

type Service struct {
	DB         *sql.DB
	SessionTTL time.Duration
}

func (s Service) Register(ctx context.Context, r Registration) (string, error) {
	if s.DB == nil {
		return "", errors.New("authentication database unavailable")
	}
	if !validRegistration(r) {
		return "", errors.New("invalid registration")
	}
	hash, err := HashPassword(r.Password)
	if err != nil {
		return "", err
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	var userID string
	err = tx.QueryRowContext(ctx, "INSERT INTO users(email,password_hash,email_verified_at) VALUES($1,$2,now()) RETURNING id", strings.ToLower(strings.TrimSpace(r.Email)), hash).Scan(&userID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return "", ErrDuplicateAccount
		}
		return "", err
	}
	fullName := strings.TrimSpace(r.FirstName + " " + r.LastName)
	if _, err = tx.ExecContext(ctx, "INSERT INTO applicant_profiles(user_id,full_name) VALUES($1,$2)", userID, fullName); err != nil {
		return "", err
	}
	if r.AccountType == "employer" {
		slug := uniqueSlug(r.OrganizationName, userID)
		var orgID string
		if err = tx.QueryRowContext(ctx, "INSERT INTO organizations(name,slug) VALUES($1,$2) RETURNING id", strings.TrimSpace(r.OrganizationName), slug).Scan(&orgID); err != nil {
			return "", err
		}
		if _, err = tx.ExecContext(ctx, "INSERT INTO organization_members(organization_id,user_id,role) VALUES($1,$2,$3)", orgID, userID, OrganizationAdmin); err != nil {
			return "", err
		}
	}
	token, tokenHash, err := NewSessionToken()
	if err != nil {
		return "", err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO sessions(user_id,token_hash,expires_at) VALUES($1,$2,$3)", userID, tokenHash, time.Now().Add(s.ttl())); err != nil {
		return "", err
	}
	if err = tx.Commit(); err != nil {
		return "", err
	}
	return token, nil
}

func (s Service) Login(ctx context.Context, email, password string) (string, error) {
	if s.DB == nil {
		return "", errors.New("authentication database unavailable")
	}
	var id, hash string
	var lockedUntil sql.NullTime
	err := s.DB.QueryRowContext(ctx, "SELECT id,password_hash,locked_until FROM users WHERE email=$1 AND deleted_at IS NULL", strings.ToLower(strings.TrimSpace(email))).Scan(&id, &hash, &lockedUntil)
	if err != nil || (lockedUntil.Valid && lockedUntil.Time.After(time.Now())) {
		return "", ErrInvalidCredentials
	}
	ok, err := VerifyPassword(hash, password)
	if err != nil || !ok {
		_, _ = s.DB.ExecContext(ctx, "UPDATE users SET failed_login_count=failed_login_count+1, locked_until=CASE WHEN failed_login_count+1 >= 5 THEN now()+interval '15 minutes' ELSE locked_until END WHERE id=$1", id)
		return "", ErrInvalidCredentials
	}
	token, tokenHash, err := NewSessionToken()
	if err != nil {
		return "", err
	}
	_, err = s.DB.ExecContext(ctx, "UPDATE users SET failed_login_count=0, locked_until=NULL WHERE id=$1", id)
	if err == nil {
		_, err = s.DB.ExecContext(ctx, "INSERT INTO sessions(user_id,token_hash,expires_at) VALUES($1,$2,$3)", id, tokenHash, time.Now().Add(s.ttl()))
	}
	return token, err
}

func (s Service) Logout(ctx context.Context, token string) error {
	if s.DB == nil {
		return errors.New("authentication database unavailable")
	}
	_, err := s.DB.ExecContext(ctx, "UPDATE sessions SET revoked_at=now() WHERE token_hash=$1 AND revoked_at IS NULL", HashSessionToken(token))
	return err
}

func (s Service) ttl() time.Duration {
	if s.SessionTTL <= 0 {
		return 24 * time.Hour
	}
	return s.SessionTTL
}

func validRegistration(r Registration) bool {
	if r.AccountType != "applicant" && r.AccountType != "employer" {
		return false
	}
	if strings.TrimSpace(r.FirstName) == "" || strings.TrimSpace(r.LastName) == "" || !strings.Contains(r.Email, "@") {
		return false
	}
	return r.AccountType != "employer" || strings.TrimSpace(r.OrganizationName) != ""
}

func uniqueSlug(name, userID string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return '-'
	}, name)
	return strings.Trim(name, "-") + "-" + strings.ReplaceAll(userID[:8], "-", "")
}
