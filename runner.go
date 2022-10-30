package pgxatomic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type txStarter interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

type runner struct {
	db   txStarter
	opts pgx.TxOptions
}

func NewRunner(db txStarter, opts pgx.TxOptions) (*runner, error) {
	if db == nil {
		return nil, errors.New("pgxatomic: db cannot be nil")
	}

	return &runner{
		db:   db,
		opts: opts,
	}, nil
}

// Run executes txFunc within shared transaction.
func (r *runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return pgx.BeginTxFunc(ctx, r.db, r.opts, func(tx pgx.Tx) error {
		return txFunc(withTx(ctx, tx))
	})
}
