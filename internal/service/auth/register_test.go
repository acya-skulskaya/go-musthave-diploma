package auth

import (
	"context"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/user"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/user/mocks"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

func TestService_Register(t *testing.T) {
	ctx := context.Background()
	secretKey := "test"
	login := "test"
	password := "test"

	tests := []struct {
		name       string
		userCreds  request.UserCredits
		setupMock  func(m *mocks.MockUserRepositoryInterface)
		want       models.User
		wantErr    error
		wantHasErr bool
	}{
		{
			name: "success",
			userCreds: request.UserCredits{
				Login:    login,
				Password: password,
			},
			setupMock: func(m *mocks.MockUserRepositoryInterface) {
				m.EXPECT().
					Create(ctx, login, mock.Anything).
					Return(models.User{Login: login}, nil).
					Once()
			},
			want: models.User{
				Login: login,
			},
			wantErr:    nil,
			wantHasErr: false,
		}, {
			name: "could not hash password",
			userCreds: request.UserCredits{
				Login:    login,
				Password: strings.Repeat(password, 1000),
			},
			setupMock: func(m *mocks.MockUserRepositoryInterface) {
				//
			},
			want:       models.User{},
			wantErr:    bcrypt.ErrPasswordTooLong,
			wantHasErr: true,
		}, {
			name: "could not create",
			userCreds: request.UserCredits{
				Login:    login,
				Password: password,
			},
			setupMock: func(m *mocks.MockUserRepositoryInterface) {
				m.EXPECT().
					Create(ctx, login, mock.Anything).
					Return(models.User{}, user.ErrLoginAlreadyExists).
					Once()
			},
			want:       models.User{},
			wantErr:    user.ErrLoginAlreadyExists,
			wantHasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepositoryInterface(t)
			tt.setupMock(mockRepo)

			as := New(secretKey, mockRepo)
			model, err := as.Register(ctx, &tt.userCreds)

			if tt.wantHasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !tt.wantHasErr {
				assert.Equal(t, tt.want.Login, model.Login)
			}

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
