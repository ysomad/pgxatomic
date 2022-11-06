package pgxatomic

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type db interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// runner starts transaction in Run method by wrapping txFunc using db,
// pgx.Conn and pgxpool.Pool implements db.
type runner struct {
	db   db
	opts pgx.TxOptions
}

func NewRunner(db db, o pgx.TxOptions) (runner, error) {
	if db == nil {
		return runner{}, errors.New("pgxatomic: db cannot be nil")
	}
	return runner{
		db:   db,
		opts: o,
	}, nil
}

// Run wraps txFunc in pgx.BeginTxFunc with injected pgx.Tx into context and runs it.
func (r runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return pgx.BeginTxFunc(ctx, r.db, r.opts, func(tx pgx.Tx) error {
		return txFunc(withTx(ctx, tx))
	})
}
