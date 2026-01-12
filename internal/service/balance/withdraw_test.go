package balance

import (
	"context"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance/mocks"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_Withdraw(t *testing.T) {
	ctx := context.WithValue(context.Background(), middleware.ContextKeyUser, models.User{
		Login:    "test",
		Password: "test",
		ID:       123,
	})

	tests := []struct {
		name        string
		ctx         context.Context
		orderNumber string
		sum         float64
		setupMock   func(m *mocks.MockBalanceRepositoryInterface)
		wantErr     error
		wantHasErr  bool
	}{
		{
			name:        "success",
			ctx:         ctx,
			orderNumber: "12345678903",
			sum:         1,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().
					Withdraw(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).
					Once()
				m.EXPECT().GetByUser(ctx, mock.Anything).Return(models.UserBalance{
					UserID:    123,
					Current:   123,
					Withdrawn: 0,
				}, nil).Once()
			},
			wantErr:    nil,
			wantHasErr: false,
		}, {
			name:        "could not get balance",
			ctx:         ctx,
			orderNumber: "12345678903",
			sum:         1,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().GetByUser(ctx, mock.Anything).Return(models.UserBalance{}, assert.AnError).Once()
			},
			wantErr:    assert.AnError,
			wantHasErr: true,
		}, {
			name:        "not enough balance",
			ctx:         ctx,
			orderNumber: "12345678903",
			sum:         1000,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().GetByUser(ctx, mock.Anything).Return(models.UserBalance{
					UserID:    123,
					Current:   1,
					Withdrawn: 0,
				}, nil).Once()
			},
			wantErr:    balance.ErrNotEnoughBalanceToWithdraw,
			wantHasErr: true,
		}, {
			name:        "not enough balance from repo",
			ctx:         ctx,
			orderNumber: "12345678903",
			sum:         1,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().GetByUser(ctx, mock.Anything).Return(models.UserBalance{
					UserID:    123,
					Current:   1000,
					Withdrawn: 0,
				}, nil).Once()
				m.EXPECT().
					Withdraw(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(balance.ErrNotEnoughBalanceToWithdraw).
					Once()
			},
			wantErr:    balance.ErrNotEnoughBalanceToWithdraw,
			wantHasErr: true,
		}, {
			name:        "withdrawal for order already exists",
			ctx:         ctx,
			orderNumber: "12345678903",
			sum:         1,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().GetByUser(ctx, mock.Anything).Return(models.UserBalance{
					UserID:    123,
					Current:   1000,
					Withdrawn: 0,
				}, nil).Once()
				m.EXPECT().
					Withdraw(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(balance.ErrWithdrawalForThisOrderAlreadyExists).
					Once()
			},
			wantErr:    balance.ErrWithdrawalForThisOrderAlreadyExists,
			wantHasErr: true,
		}, {
			name:        "err to withdraw",
			ctx:         ctx,
			orderNumber: "12345678903",
			sum:         1,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().GetByUser(ctx, mock.Anything).Return(models.UserBalance{
					UserID:    123,
					Current:   1000,
					Withdrawn: 0,
				}, nil).Once()
				m.EXPECT().
					Withdraw(ctx, mock.Anything, mock.Anything, mock.Anything).
					Return(assert.AnError).
					Once()
			},
			wantErr:    assert.AnError,
			wantHasErr: true,
		}, {
			name:        "no user in context",
			ctx:         context.Background(),
			orderNumber: "12345678903",
			sum:         1,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {

			},
			wantErr:    auth.ErrCouldNotGetUserFromContext,
			wantHasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockBalanceRepositoryInterface(t)
			tt.setupMock(mockRepo)

			bs := New(mockRepo)

			err := bs.Withdraw(tt.ctx, tt.orderNumber, tt.sum)
			if tt.wantHasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
