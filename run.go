package atomic

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type starter interface {
	Begin(context.Context) (pgx.Tx, error)
}

// Run executes txFunc within shared transaction.
//
// Implementation example:
//
//	func (t *transactor) RunAtomic(ctx context.Context, txFunc func(ctx context.Context) error) error {
//			return atomic.Run(ctx, t.conn, txFunc)
//		}
func Run(ctx context.Context, db starter, txFunc func(ctx context.Context) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("atomic: begin transaction - %w", err)
	}
	return run(ctx, tx, txFunc)
}

type starterWithOpts interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

// RunWithOpts executes txFunc within shared transaction with pgx.TxOptions.
func RunWithOpts(ctx context.Context, db starterWithOpts, opts pgx.TxOptions, txFunc func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("atomic: begin transaction - %w", err)
	}
	return run(ctx, tx, txFunc)
}

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
