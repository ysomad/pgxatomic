package pgxatomic

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type txKey struct{}

// withTx sets pgx.Tx into context.
func withTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// txFromContext return pgx.Tx from context or nil if not found.
func txFromContext(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}
