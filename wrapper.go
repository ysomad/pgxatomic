package pgxatomic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type txStarter interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

type wrapper struct {
	db   txStarter
	opts pgx.TxOptions
}

func NewWrapper(db txStarter, opts pgx.TxOptions) (*wrapper, error) {
	if db == nil {
		return nil, errors.New("pgxatomic: db cannot be nil")
	}

	return &wrapper{
		db:   db,
		opts: opts,
	}, nil
}

// Wrap wraps txFunc in pgx.BeginTxFunc with injected pgx.Tx into context.
func (w *wrapper) Wrap(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return pgx.BeginTxFunc(ctx, w.db, w.opts, func(tx pgx.Tx) error {
		return txFunc(withTx(ctx, tx))
	})
}
