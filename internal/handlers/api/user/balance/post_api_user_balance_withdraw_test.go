package balance_test

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"
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

func TestPostAPIUserBalanceWithdraw(t *testing.T) {
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

	orderNumberNoErr := "2377225624"
	orderNumberErrOrderAlreadyExists := "12345678903"
	balanceRepo := mocksBalance.NewMockBalanceRepositoryInterface(t)
	setupMockBalance := func(m *mocksBalance.MockBalanceRepositoryInterface) {
		m.EXPECT().
			GetByUser(mock.Anything, successUserID).
			Return(models.UserBalance{
				UserID:    123,
				Current:   1000,
				Withdrawn: 0,
			}, nil)

		m.EXPECT().Withdraw(mock.Anything, mock.Anything, orderNumberNoErr, mock.Anything).Return(nil)
		m.EXPECT().Withdraw(mock.Anything, mock.Anything, orderNumberErrOrderAlreadyExists, mock.Anything).Return(balance.ErrWithdrawalForThisOrderAlreadyExists)
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
		data     string
		wantCode int
	}{
		{
			name:     "200 — успешная обработка запроса",
			data:     `{"order": "` + orderNumberNoErr + `","sum": 751}`,
			wantCode: http.StatusOK,
		}, {
			name:     "402 — на счету недостаточно средств",
			data:     `{"order": "` + orderNumberNoErr + `","sum": 1751}`,
			wantCode: http.StatusPaymentRequired,
		}, {
			name:     "422 — неверный номер заказа",
			data:     `{"order": "` + orderNumberErrOrderAlreadyExists + `","sum": 751}`,
			wantCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader = strings.NewReader(tt.data)
			request, err = http.NewRequest(http.MethodPost, testServer.URL+"/api/user/balance/withdraw", bodyReader)
			require.NoError(t, err)
			for _, cookie := range cookies {
				request.AddCookie(cookie)
			}
			//nolint:bodyclose // не понятно почему тут ругается, все закрывается
			res, err = testServer.Client().Do(request)
			require.NoError(t, err)
			res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}
