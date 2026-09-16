package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

func Open(ctx context.Context, databaseURL string) (*DB, error) {
	if databaseURL == "" {
		return nil, errors.New("database URL is empty")
	}
	config, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	db := stdlib.OpenDB(*config)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, err
	}
	return &DB{SQL: db}, nil
}

func (d *DB) Close() error {
	if d == nil || d.SQL == nil {
		return nil
	}
	return d.SQL.Close()
}

func (d *DB) Ready(ctx context.Context) bool {
	if d == nil || d.SQL == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return d.SQL.PingContext(ctx) == nil
}
