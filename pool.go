package pgxatomic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool wraps pgxpool.Pool query methods with pgxatomic corresponding functions
// which injects pgx.Tx into context.
type Pool struct {
	p *pgxpool.Pool
}

func NewPool(p *pgxpool.Pool) (Pool, error) {
	if p == nil {
		return Pool{}, errors.New("pgxatomic: pool cannot be nil")
	}
	return Pool{p: p}, nil
}

func (p Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return Query(ctx, p.p, sql, args...)
}

func (p Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return QueryRow(ctx, p.p, sql, args...)
}

func (p Pool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return Exec(ctx, p.p, sql, args...)
}
