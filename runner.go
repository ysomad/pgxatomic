package pgxatomic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type txStarter interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// Runner starts transaction in Run method by wrapping txFunc using db,
// pgx.Conn and pgxpool.Pool implements db.
type Runner struct {
	db   txStarter
	opts pgx.TxOptions
}

func NewRunner(db txStarter, o pgx.TxOptions) (Runner, error) {
	if db == nil {
		return Runner{}, errors.New("pgxatomic: db cannot be nil")
	}

	return Runner{
		db:   db,
		opts: o,
	}, nil
}

// Run wraps txFunc in pgx.BeginTxFunc with injected pgx.Tx into context and runs it.
func (r Runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return pgx.BeginTxFunc(ctx, r.db, r.opts, func(tx pgx.Tx) error {
		return txFunc(WithTx(ctx, tx))
	})
}
