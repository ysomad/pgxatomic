package atomic

import (
	"context"
	"errors"
	"fmt"

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

// Run is a helper method for runWithOpts function.
func (r *runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return runWithOpts(ctx, r.db, r.opts, txFunc)
}

// Run executes txFunc within shared transaction.
func Run(ctx context.Context, db txStarter, txFunc func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("atomic: begin transaction - %w", err)
	}
	return run(ctx, tx, txFunc)
}

// runWithOpts executes txFunc withing shared transaction with pgx.TxOptions.
func runWithOpts(ctx context.Context, db txStarter, opts pgx.TxOptions, txFunc func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("atomic: begin transaction - %w", err)
	}
	return run(ctx, tx, txFunc)
}

// run executes txFunc with injected transaction in context and commits or rollback on error.
func run(ctx context.Context, tx pgx.Tx, txFunc func(ctx context.Context) error) (err error) {
	err = txFunc(withTx(ctx, tx))
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("atomic: rollback - %w", rbErr)
		}

		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("atomic: commit - %w", err)
	}

	return nil
}
