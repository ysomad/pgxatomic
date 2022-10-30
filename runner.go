package atomic

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type (
	starter interface {
		Begin(context.Context) (pgx.Tx, error)
	}

	starterWithOpts interface {
		BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	}
)

type runner struct {
	tx   starterWithOpts
	opts pgx.TxOptions
}

func NewRunner(tx starterWithOpts, opts pgx.TxOptions) *runner {
	return &runner{
		tx:   tx,
		opts: opts,
	}
}

// Run is a helper method for runWithOpts function.
func (r *runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return runWithOpts(ctx, r.tx, r.opts, txFunc)
}

// Run executes txFunc within shared transaction.
func Run(ctx context.Context, db starter, txFunc func(ctx context.Context) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("atomic: begin transaction - %w", err)
	}
	return run(ctx, tx, txFunc)
}

// runWithOpts executes txFunc withing shared transaction with pgx.TxOptions.
func runWithOpts(ctx context.Context, db starterWithOpts, opts pgx.TxOptions, txFunc func(ctx context.Context) error) error {
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
