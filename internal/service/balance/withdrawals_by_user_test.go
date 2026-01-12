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

func TestService_WithdrawalsByUser(t *testing.T) {
	user := models.User{
		Login:    "test",
		Password: "test",
		ID:       123,
	}
	ctx := context.WithValue(context.Background(), middleware.ContextKeyUser, user)

	tests := []struct {
		name       string
		ctx        context.Context
		setupMock  func(m *mocks.MockBalanceRepositoryInterface)
		wantErr    error
		wantHasErr bool
	}{
		{
			name: "success",
			ctx:  ctx,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().GetWithdrawalsByUser(ctx, mock.Anything).Return([]models.UserBalanceFlow{
					{
						OrderNumber: "123",
						UserID:      123,
						Amount:      1,
					}, {
						OrderNumber: "1234",
						UserID:      123,
						Amount:      13,
					},
				}, nil).Once()
			},
			wantErr:    nil,
			wantHasErr: false,
		}, {
			name: "no withdrawals",
			ctx:  ctx,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().GetWithdrawalsByUser(ctx, mock.Anything).Return([]models.UserBalanceFlow{}, nil).Once()
			},
			wantErr:    balance.ErrNoWithdrawals,
			wantHasErr: true,
		}, {
			name: "could not get withdrawals",
			ctx:  ctx,
			setupMock: func(m *mocks.MockBalanceRepositoryInterface) {
				m.EXPECT().GetWithdrawalsByUser(ctx, mock.Anything).Return([]models.UserBalanceFlow{}, assert.AnError).Once()
			},
			wantErr:    assert.AnError,
			wantHasErr: true,
		}, {
			name: "no user in context",
			ctx:  context.Background(),
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

			_, err := bs.WithdrawalsByUser(tt.ctx)
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
