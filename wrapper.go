package pgxatomic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type wrapper struct {
	p    *pgxpool.Pool
	opts pgx.TxOptions
}

func NewWrapper(p *pgxpool.Pool, opts pgx.TxOptions) (*wrapper, error) {
	if p == nil {
		return nil, errors.New("pgxatomic: pool cannot be nil")
	}

	return &wrapper{
		p:    p,
		opts: opts,
	}, nil
}

// Wrap wraps txFunc in pgx.BeginTxFunc with injected pgx.Tx into context.
func (w *wrapper) Wrap(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return pgx.BeginTxFunc(ctx, w.p, w.opts, func(tx pgx.Tx) error {
		return txFunc(withTx(ctx, tx))
	})
}
