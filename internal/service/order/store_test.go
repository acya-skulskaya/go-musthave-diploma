package order

import (
	"context"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order/mocks"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_Store(t *testing.T) {
	user := models.User{
		Login:    "test",
		Password: "test",
		ID:       123,
	}
	ctx := context.WithValue(context.Background(), middleware.ContextKeyUser, user)

	tests := []struct {
		name       string
		ctx        context.Context
		setupMock  func(m *mocks.MockOrderRepositoryInterface)
		wantErr    error
		wantHasErr bool
	}{
		{
			name: "success",
			ctx:  ctx,
			setupMock: func(m *mocks.MockOrderRepositoryInterface) {
				m.EXPECT().Create(ctx, mock.Anything, mock.Anything).Return(models.Order{
					OrderNumber: "123",
					Status:      models.OrderStatusNew,
					UserID:      123,
					Accrual:     0,
				}, nil).Once()
			},
			wantErr:    nil,
			wantHasErr: false,
		}, {
			name: "could not create order and get it",
			ctx:  ctx,
			setupMock: func(m *mocks.MockOrderRepositoryInterface) {
				m.EXPECT().Create(ctx, mock.Anything, mock.Anything).Return(models.Order{}, order.ErrOrderNumberAlreadyExists).Once()
				m.EXPECT().Show(ctx, mock.Anything).Return(models.Order{}, assert.AnError)
			},
			wantErr:    assert.AnError,
			wantHasErr: true,
		}, {
			name: "no user in context",
			ctx:  context.Background(),
			setupMock: func(m *mocks.MockOrderRepositoryInterface) {

			},
			wantErr:    auth.ErrCouldNotGetUserFromContext,
			wantHasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockOrderRepositoryInterface(t)
			tt.setupMock(mockRepo)

			os := New(mockRepo)

			_, err := os.Store(tt.ctx, tt.name)
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
