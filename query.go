package pgxatomic

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Query is a wrapper around pgx Query method.
func Query(
	ctx context.Context,
	db interface {
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	},
	sql string,
	args ...any,
) (pgx.Rows, error) {
	if tx := TxFromContext(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}
	return db.Query(ctx, sql, args...)
}

// Exec is a wrapper around pgx Exec method.
func Exec(
	ctx context.Context,
	db interface {
		Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	},
	sql string,
	args ...any,
) (pgconn.CommandTag, error) {
	if tx := TxFromContext(ctx); tx != nil {
		return tx.Exec(ctx, sql, args...)
	}
	return db.Exec(ctx, sql, args...)
}

// QueryRow is a wrapper around pgx QueryRow method.
func QueryRow(
	ctx context.Context,
	db interface {
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	},
	sql string,
	args ...any,
) pgx.Row {
	if tx := TxFromContext(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}
	return db.QueryRow(ctx, sql, args...)
}
