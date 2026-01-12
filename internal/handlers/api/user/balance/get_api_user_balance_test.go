package balance_test

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	mocksBalance "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance/mocks"
	mocksOrder "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order/mocks"
	mocksUser "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/user/mocks"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/router"
	serviceAuth "github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
	serviceBalance "github.com/acya-skulskaya/go-musthave-diploma/internal/service/balance"
	serviceOrder "github.com/acya-skulskaya/go-musthave-diploma/internal/service/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAPIUserBalance(t *testing.T) {
	var successUserID uint = 123
	successLogin := "success_login"
	successPassword := "success_password"
	successPasswordHashed, _ := serviceAuth.HashPassword(successPassword)
	userRepo := mocksUser.NewMockUserRepositoryInterface(t)
	setupMockUser := func(m *mocksUser.MockUserRepositoryInterface) {
		m.EXPECT().
			Get(mock.Anything, mock.Anything).
			Return(models.User{Login: successLogin, Password: successPasswordHashed, ID: successUserID}, nil)
	}
	setupMockUser(userRepo)
	balanceRepo := mocksBalance.NewMockBalanceRepositoryInterface(t)
	setupMockBalance := func(m *mocksBalance.MockBalanceRepositoryInterface) {
		m.EXPECT().
			GetByUser(mock.Anything, successUserID).
			Return(models.UserBalance{
				UserID:    123,
				Current:   123,
				Withdrawn: 0,
			}, nil)
	}
	setupMockBalance(balanceRepo)
	orderRepo := mocksOrder.NewMockOrderRepositoryInterface(t)

	authService := serviceAuth.New("test", userRepo)
	orderService := serviceOrder.New(orderRepo)
	balanceService := serviceBalance.New(balanceRepo)

	cfg := config.Config{}
	cfg.Auth.SecretKey = "test"

	r := router.New(&cfg, authService, orderService, balanceService)
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	bodyReader := strings.NewReader(`{"login": "` + successLogin + `","password": "` + successPassword + `"}`)
	request, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/user/login", bodyReader)
	require.NoError(t, err)
	res, err := testServer.Client().Do(request)
	require.NoError(t, err)
	res.Body.Close()
	cookies := res.Cookies()

	tests := []struct {
		name     string
		wantCode int
	}{
		{
			name:     "200 — успешная обработка запроса",
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err = http.NewRequest(http.MethodGet, testServer.URL+"/api/user/balance", nil)
			require.NoError(t, err)
			for _, cookie := range cookies {
				request.AddCookie(cookie)
			}
			//nolint:bodyclose // не понятно почему тут ругается, все закрывается
			res, err = testServer.Client().Do(request)
			defer res.Body.Close()
			require.NoError(t, err)

			// проверяем код ответа
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}
