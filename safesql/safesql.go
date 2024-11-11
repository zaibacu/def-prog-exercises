package safesql

import (
	"context"
	"database/sql"

	"github.com/empijei/def-prog-exercises/safesql/internal/raw"
)

func init() {
	raw.TrustedSQLCtor = func(unsafe string) TrustedSQL {
		return TrustedSQL{unsafe}
	}
}

type compileTimeConstant string

type TrustedSQL struct {
	s string
}

type DB struct {
	db *sql.DB
}

type Rows = sql.Rows
type Result = sql.Result

func Open(driverName string, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) QueryContext(ctx context.Context, query TrustedSQL, args ...any) (*Rows, error) {
	r, err := db.db.QueryContext(ctx, query.s, args...)

	return r, err
}

func (db *DB) ExecContext(ctx context.Context, query TrustedSQL, args ...any) (Result, error) {
	r, err := db.db.ExecContext(ctx, query.s, args...)

	return r, err
}

func New(text compileTimeConstant) TrustedSQL {
	return TrustedSQL{s: string(text)}
}
