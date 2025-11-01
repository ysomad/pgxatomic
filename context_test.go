package pgxatomic

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testKey string

func TestWithTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx1 := NewMockTx(ctrl)
	tx2 := NewMockTx(ctrl)

	tests := []struct {
		name string
		ctx  context.Context
		tx   pgx.Tx
	}{
		{
			name: "nil tx",
			ctx:  context.Background(),
			tx:   nil,
		},
		{
			name: "tx1",
			ctx:  context.Background(),
			tx:   tx1,
		},
		{
			name: "tx2",
			ctx:  context.TODO(),
			tx:   tx2,
		},
		{
			name: "override",
			ctx:  context.WithValue(context.Background(), txKey{}, "not-a-tx"),
			tx:   tx1,
		},
		{
			name: "with other key",
			ctx:  context.WithValue(context.Background(), testKey("other"), "value"),
			tx:   tx2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCtx := WithTx(tt.ctx, tt.tx)
			assert.NotNil(t, gotCtx)

			if tt.tx != nil {
				gotTx, ok := gotCtx.Value(txKey{}).(pgx.Tx)
				assert.True(t, ok)
				assert.Equal(t, tt.tx, gotTx)
			}
		})
	}
}

func TestTxFromContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tx1 := NewMockTx(ctrl)
	tx2 := NewMockTx(ctrl)

	tests := []struct {
		name string
		ctx  context.Context
		want pgx.Tx
	}{
		{
			name: "nil tx",
			ctx:  context.WithValue(context.Background(), txKey{}, nil),
			want: nil,
		},
		{
			name: "tx1",
			ctx:  context.WithValue(context.Background(), txKey{}, tx1),
			want: tx1,
		},
		{
			name: "tx2",
			ctx:  context.WithValue(context.TODO(), txKey{}, tx2),
			want: tx2,
		},
		{
			name: "not found",
			ctx:  context.Background(),
			want: nil,
		},
		{
			name: "other key",
			ctx:  context.WithValue(context.Background(), testKey("other"), "value"),
			want: nil,
		},
		{
			name: "not tx",
			ctx:  context.WithValue(context.Background(), txKey{}, "not-a-tx"),
			want: nil,
		},
		{
			name: "nil context",
			ctx:  nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TxFromContext(tt.ctx)
			assert.Equal(t, tt.want, got)
		})
	}
}
