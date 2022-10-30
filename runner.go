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
		return nil, errors.New("atomic: db cannot be nil")
	}

	return &runner{
		db:   db,
		opts: opts,
	}, nil
}

func (r *runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return execTxFunc(ctx, r.db, r.opts, txFunc)
}

func Run(ctx context.Context, db txStarter, txFunc func(ctx context.Context) error) error {
	return execTxFunc(ctx, db, pgx.TxOptions{}, txFunc)
}

// execTxFunc executes txFunc withing shared transaction.
func execTxFunc(ctx context.Context, db txStarter, opts pgx.TxOptions, txFunc func(ctx context.Context) error) error {
	return pgx.BeginTxFunc(ctx, db, opts, func(tx pgx.Tx) error {
		return txFunc(withTx(ctx, tx))
	})
}
