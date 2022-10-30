package atomic

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// Query is a wrapper around pgx Query method. Accepts querier interface,
// for example pgxpool.Pool or pgx.Conn implements it.
func Query(ctx context.Context, q querier, sql string, args ...any) (pgx.Rows, error) {
	if tx := txFromContext(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}
	return q.Query(ctx, sql, args...)
}

type executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// Exec is a wrapper around pgx Exec method. Accepts executor interface,
// for example pgxpool.Pool or pgx.Conn implements it.
func Exec(ctx context.Context, e executor, sql string, args ...any) (pgconn.CommandTag, error) {
	if tx := txFromContext(ctx); tx != nil {
		return tx.Exec(ctx, sql, args...)
	}
	return e.Exec(ctx, sql, args...)
}

type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// QueryRow is a wrapper around pgx QueryRow method. Accepts queryRower interface,
// for example pgxpool.Pool or pgx.Conn implements it.
func QueryRow(ctx context.Context, q queryRower, sql string, args ...any) pgx.Row {
	if tx := txFromContext(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}
	return q.QueryRow(ctx, sql, args...)
}
