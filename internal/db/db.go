package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Tx interface with method signatures matching the sqlx.Tx struct
type Tx interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	Commit() error
	Rollback() error
}

// DB interface for general database operations
type DB interface {
	Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}

// SqlxDB implements DB interface using sqlx
type SqlxDB struct {
	*sqlx.DB
}

// Ensure that SqlxDB implements the DB interface
func (db *SqlxDB) Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.DB.QueryxContext(ctx, query, args...)
}

func (db *SqlxDB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.DB.GetContext(ctx, dest, query, args...)
}

func (db *SqlxDB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.DB.SelectContext(ctx, dest, query, args...)
}

func (db *SqlxDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.DB.ExecContext(ctx, query, args...)
}

func (db *SqlxDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := db.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &TxWrapper{tx}, nil
}

// Commit and Rollback for transactions
func (db *SqlxDB) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panicking after rollback to ensure panics are not swallowed
		}
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

// TxWrapper wraps *sqlx.Tx to implement the Tx interface
type TxWrapper struct {
	*sqlx.Tx
}

// Exec overrides the Exec method to match the Tx interface
func (tx *TxWrapper) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.Tx.ExecContext(ctx, query, args...)
}

// Query overrides the Query method to match the Tx interface
func (tx *TxWrapper) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.Tx.QueryContext(ctx, query, args...)
}

// QueryRow overrides the QueryRow method to match the Tx interface
func (tx *TxWrapper) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return tx.Tx.QueryRowContext(ctx, query, args...)
}

// Commit commits the transaction
func (tx *TxWrapper) Commit() error {
	return tx.Tx.Commit()
}

// Rollback rolls back the transaction
func (tx *TxWrapper) Rollback() error {
	return tx.Tx.Rollback()
}
