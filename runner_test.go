package pgxatomic

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var errTest = errors.New("test error")

func TestNewRunner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMocktxStarter(ctrl)

	type args struct {
		db   txStarter
		opts pgx.TxOptions
	}

	tests := []struct {
		name      string
		args      args
		want      Runner
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "empty",
			args:      args{},
			want:      Runner{},
			assertion: assert.Error,
		},
		{
			name: "ok",
			args: args{
				db:   mockDB,
				opts: pgx.TxOptions{IsoLevel: pgx.ReadCommitted},
			},
			want: Runner{
				db:   mockDB,
				opts: pgx.TxOptions{IsoLevel: pgx.ReadCommitted},
			},
			assertion: assert.NoError,
		},
		{
			name:      "ok default options",
			args:      args{db: mockDB},
			want:      Runner{db: mockDB},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRunner(tt.args.db, tt.args.opts)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRun_Commit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMocktxStarter(ctrl)
	mockTx := NewMockTx(ctrl)

	mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).AnyTimes()

	runner, err := NewRunner(mockDB, pgx.TxOptions{})
	assert.NoError(t, err)

	err = runner.Run(context.Background(), func(ctx context.Context) error {
		tx := TxFromContext(ctx)
		assert.Equal(t, mockTx, tx)
		return nil
	})
	assert.NoError(t, err)
}

func TestRun_Rollback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMocktxStarter(ctrl)
	mockTx := NewMockTx(ctrl)

	mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).AnyTimes()

	runner, err := NewRunner(mockDB, pgx.TxOptions{})
	assert.NoError(t, err)

	err = runner.Run(context.Background(), func(ctx context.Context) error {
		return errTest
	})
	assert.Error(t, err)
}

func TestRun_BeginError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMocktxStarter(ctrl)

	mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(nil, errTest)

	runner, err := NewRunner(mockDB, pgx.TxOptions{})
	assert.NoError(t, err)

	err = runner.Run(context.Background(), func(ctx context.Context) error {
		return nil
	})
	assert.Error(t, err)
}
