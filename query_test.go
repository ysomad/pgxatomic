package pgxatomic

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func ctxWithTx(mockTx pgx.Tx) func(t *testing.T) context.Context {
	return func(t *testing.T) context.Context {
		t.Helper()
		return WithTx(context.Background(), mockTx)
	}
}

func ctxWithoutTx() func(t *testing.T) context.Context {
	return func(t *testing.T) context.Context {
		t.Helper()
		return context.Background()
	}
}

func TestQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := NewMockTx(ctrl)
	mockRows := NewMockRows(ctrl)
	queryErr := errors.New("query error")

	tests := []struct {
		name      string
		setupCtx  func(t *testing.T) context.Context
		setupMock func()
		sql       string
		args      []any
		want      pgx.Rows
		wantErr   error
	}{
		{
			name:     "tx from context",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					Query(gomock.Any(), "SELECT * FROM users WHERE id = $1", 123).
					Return(mockRows, nil)
			},
			sql:     "SELECT * FROM users WHERE id = $1",
			args:    []any{123},
			want:    mockRows,
			wantErr: nil,
		},
		{
			name:     "no tx in context",
			setupCtx: ctxWithoutTx(),
			setupMock: func() {
				mockTx.EXPECT().
					Query(gomock.Any(), "SELECT * FROM products", gomock.Any()).
					Return(mockRows, nil)
			},
			sql:     "SELECT * FROM products",
			args:    []any{},
			want:    mockRows,
			wantErr: nil,
		},
		{
			name:     "error from tx",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					Query(gomock.Any(), "SELECT * FROM invalid", gomock.Any()).
					Return(nil, queryErr)
			},
			sql:     "SELECT * FROM invalid",
			args:    []any{},
			want:    nil,
			wantErr: queryErr,
		},
		{
			name:     "tx from context multiple args",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					Query(gomock.Any(), "SELECT * FROM users WHERE id = $1 AND status = $2", 42, "active").
					Return(mockRows, nil)
			},
			sql:     "SELECT * FROM users WHERE id = $1 AND status = $2",
			args:    []any{42, "active"},
			want:    mockRows,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx(t)
			tt.setupMock()

			got, err := Query(ctx, mockTx, tt.sql, tt.args...)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := NewMockTx(ctrl)
	execErr := errors.New("exec error")
	successTag := pgconn.NewCommandTag("UPDATE 5")

	tests := []struct {
		name      string
		setupCtx  func(t *testing.T) context.Context
		setupMock func()
		sql       string
		args      []any
		want      pgconn.CommandTag
		wantErr   error
	}{
		{
			name:     "tx from context",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					Exec(gomock.Any(), "UPDATE users SET name = $1 WHERE id = $2", "John", 1).
					Return(successTag, nil)
			},
			sql:     "UPDATE users SET name = $1 WHERE id = $2",
			args:    []any{"John", 1},
			want:    successTag,
			wantErr: nil,
		},
		{
			name:     "no tx in context",
			setupCtx: ctxWithoutTx(),
			setupMock: func() {
				mockTx.EXPECT().
					Exec(gomock.Any(), "DELETE FROM users WHERE id = $1", 99).
					Return(pgconn.NewCommandTag("DELETE 1"), nil)
			},
			sql:     "DELETE FROM users WHERE id = $1",
			args:    []any{99},
			want:    pgconn.NewCommandTag("DELETE 1"),
			wantErr: nil,
		},
		{
			name:     "error from tx",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					Exec(gomock.Any(), "INVALID SQL", gomock.Any()).
					Return(pgconn.CommandTag{}, execErr)
			},
			sql:     "INVALID SQL",
			args:    []any{},
			want:    pgconn.CommandTag{},
			wantErr: execErr,
		},
		{
			name:     "tx from context insert",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					Exec(gomock.Any(), "INSERT INTO users (name, email) VALUES ($1, $2)", "Alice", "alice@example.com").
					Return(pgconn.NewCommandTag("INSERT 1"), nil)
			},
			sql:     "INSERT INTO users (name, email) VALUES ($1, $2)",
			args:    []any{"Alice", "alice@example.com"},
			want:    pgconn.NewCommandTag("INSERT 1"),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx(t)
			tt.setupMock()

			got, err := Exec(ctx, mockTx, tt.sql, tt.args...)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestQueryRow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := NewMockTx(ctrl)
	mockRow := NewMockRow(ctrl)

	tests := []struct {
		name      string
		setupCtx  func(t *testing.T) context.Context
		setupMock func()
		sql       string
		args      []any
		want      pgx.Row
	}{
		{
			name:     "tx from context",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					QueryRow(gomock.Any(), "SELECT * FROM users WHERE id = $1", 42).
					Return(mockRow)
			},
			sql:  "SELECT * FROM users WHERE id = $1",
			args: []any{42},
			want: mockRow,
		},
		{
			name:     "no tx in context",
			setupCtx: ctxWithoutTx(),
			setupMock: func() {
				mockTx.EXPECT().
					QueryRow(gomock.Any(), "SELECT COUNT(*) FROM products", gomock.Any()).
					Return(mockRow)
			},
			sql:  "SELECT COUNT(*) FROM products",
			args: []any{},
			want: mockRow,
		},
		{
			name:     "tx from context multiple args",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					QueryRow(gomock.Any(), "SELECT * FROM users WHERE email = $1 AND status = $2", "test@example.com", "active").
					Return(mockRow)
			},
			sql:  "SELECT * FROM users WHERE email = $1 AND status = $2",
			args: []any{"test@example.com", "active"},
			want: mockRow,
		},
		{
			name:     "tx from context aggregate",
			setupCtx: ctxWithTx(mockTx),
			setupMock: func() {
				mockTx.EXPECT().
					QueryRow(gomock.Any(), "SELECT MAX(id) FROM users", gomock.Any()).
					Return(mockRow)
			},
			sql:  "SELECT MAX(id) FROM users",
			args: []any{},
			want: mockRow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx(t)
			tt.setupMock()

			got := QueryRow(ctx, mockTx, tt.sql, tt.args...)

			assert.Equal(t, tt.want, got)
		})
	}
}
