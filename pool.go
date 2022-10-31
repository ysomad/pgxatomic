package pgxatomic

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool wraps pgxpool.Pool query methods with pgxatomic versions,
// which is injecting pgx.Tx into context.
// For example pgxpool.Pool and pgx.Conn implements it.
type Pool struct {
	*pgxpool.Pool
}

func NewPool(p *pgxpool.Pool) *Pool {
	return &Pool{p}
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return Query(ctx, p, sql, args...)
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return QueryRow(ctx, p, sql, args...)
}

func (p *Pool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return Exec(ctx, p, sql, args...)
}
