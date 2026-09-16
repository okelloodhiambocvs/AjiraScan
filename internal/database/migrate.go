package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ApplyMigrations(ctx context.Context, db *sql.DB, directory string) error {
	if db == nil {
		return errors.New("database is nil")
	}
	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())"); err != nil {
		return err
	}
	files, err := filepath.Glob(filepath.Join(directory, "*.up.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, file := range files {
		version := strings.TrimSuffix(filepath.Base(file), ".up.sql")
		var applied bool
		if err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&applied); err != nil {
			return err
		}
		if applied {
			continue
		}
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, string(sqlBytes)); err == nil {
			_, err = tx.ExecContext(ctx, "INSERT INTO schema_migrations(version) VALUES($1)", version)
		}
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("apply %s: %w", version, err)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
