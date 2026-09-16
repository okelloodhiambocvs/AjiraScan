package tenancy

import (
	"context"
	"database/sql"
	"errors"
)

type Scope struct {
	UserID         string
	OrganizationID string
}

func (s Scope) Valid() bool {
	return s.UserID != "" && s.OrganizationID != ""
}

func Begin(ctx context.Context, db *sql.DB, scope Scope) (*sql.Tx, error) {
	if db == nil {
		return nil, errors.New("database unavailable")
	}
	if !scope.Valid() {
		return nil, errors.New("tenant scope is required")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.organization_id', $1, true)", scope.OrganizationID); err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, "SELECT set_config('app.user_id', $1, true)", scope.UserID); err != nil {
		tx.Rollback()
		return nil, err
	}
	return tx, nil
}
