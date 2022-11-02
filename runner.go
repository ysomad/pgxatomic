package pgxatomic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type runner struct {
	p *pgxpool.Pool
	o pgx.TxOptions
}

func NewRunner(p *pgxpool.Pool, o pgx.TxOptions) (*runner, error) {
	if p == nil {
		return nil, errors.New("pgxatomic: pool cannot be nil")
	}

	return &runner{
		p: p,
		o: o,
	}, nil
}

// Run wraps txFunc in pgx.BeginTxFunc with injected pgx.Tx into context and runs it.
func (w *runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return pgx.BeginTxFunc(ctx, w.p, w.o, func(tx pgx.Tx) error {
		return txFunc(withTx(ctx, tx))
	})
}
