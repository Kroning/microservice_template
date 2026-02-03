package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type AbstractDB interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Transaction(ctx context.Context, t func(tx *sqlx.Tx) error) error
	Close() error
}
