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

	txStarter interface {
		starter
		starterWithOpts
	}
)

type runner struct {
	tx txStarter
}

func NewRunner(tx txStarter) *runner {
	return &runner{tx: tx}
}

func (r *runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return Run(ctx, r.tx, txFunc)
}

func (r *runner) RunWithOpts(ctx context.Context, opts pgx.TxOptions, txFunc func(ctx context.Context) error) error {
	return RunWithOpts(ctx, r.tx, opts, txFunc)
}

// Run executes txFunc within shared transaction.
func Run(ctx context.Context, db txStarter, txFunc func(ctx context.Context) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("atomic: begin transaction - %w", err)
	}
	return run(ctx, tx, txFunc)
}

// RunWithOpts executes txFunc within shared transaction with transaction options.
func RunWithOpts(ctx context.Context, db starterWithOpts, opts pgx.TxOptions, txFunc func(ctx context.Context) error) error {
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
