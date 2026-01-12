package orders_test

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

func TestGetAPIUserOrders(t *testing.T) {
	var successUserID uint
	successUserID = 123
	successLogin := "success_login"
	var noContentUserID uint
	noContentUserID = 12345
	noContentLogin := "no_content_login"
	var internalServerErrorUserID uint
	internalServerErrorUserID = 123456
	internalServerErrorLogin := "internal_server_error_login"
	password := "success_password"
	passwordHashed, _ := serviceAuth.HashPassword(password)
	userRepo := mocksUser.NewMockUserRepositoryInterface(t)
	setupMockUser := func(m *mocksUser.MockUserRepositoryInterface) {
		m.EXPECT().
			Get(mock.Anything, successLogin).
			Return(models.User{Login: successLogin, Password: passwordHashed, ID: successUserID}, nil)
		m.EXPECT().
			Get(mock.Anything, successUserID).
			Return(models.User{Login: successLogin, Password: passwordHashed, ID: successUserID}, nil)

		m.EXPECT().
			Get(mock.Anything, noContentLogin).
			Return(models.User{Login: noContentLogin, Password: passwordHashed, ID: noContentUserID}, nil)
		m.EXPECT().
			Get(mock.Anything, noContentUserID).
			Return(models.User{Login: noContentLogin, Password: passwordHashed, ID: noContentUserID}, nil)

		m.EXPECT().
			Get(mock.Anything, internalServerErrorLogin).
			Return(models.User{Login: internalServerErrorLogin, Password: passwordHashed, ID: internalServerErrorUserID}, nil)
		m.EXPECT().
			Get(mock.Anything, internalServerErrorUserID).
			Return(models.User{Login: internalServerErrorLogin, Password: passwordHashed, ID: internalServerErrorUserID}, nil)

	}
	setupMockUser(userRepo)

	balanceRepo := mocksBalance.NewMockBalanceRepositoryInterface(t)

	orderRepo := mocksOrder.NewMockOrderRepositoryInterface(t)
	setupMockOrder := func(m *mocksOrder.MockOrderRepositoryInterface) {
		m.EXPECT().
			AllByUser(mock.Anything, successUserID).
			Return([]models.Order{{
				OrderNumber: "123",
				Status:      models.OrderStatusNew,
				UserID:      successUserID,
				Accrual:     0,
			}, {
				OrderNumber: "1234",
				Status:      models.OrderStatusNew,
				UserID:      successUserID,
				Accrual:     0,
			}}, nil)

		m.EXPECT().
			AllByUser(mock.Anything, noContentUserID).
			Return([]models.Order{}, nil)

		m.EXPECT().
			AllByUser(mock.Anything, internalServerErrorUserID).
			Return([]models.Order{}, assert.AnError)

	}
	setupMockOrder(orderRepo)

	authService := serviceAuth.New("test", userRepo)
	orderService := serviceOrder.New(orderRepo)
	balanceService := serviceBalance.New(balanceRepo)

	cfg := config.Config{}
	cfg.Auth.SecretKey = "test"

	r := router.New(&cfg, authService, orderService, balanceService)
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	tests := []struct {
		name     string
		login    string
		wantCode int
	}{
		{
			name:     "200 — успешная обработка запроса",
			login:    successLogin,
			wantCode: http.StatusOK,
		}, {
			name:     "204 — нет данных для ответа",
			login:    noContentLogin,
			wantCode: http.StatusNoContent,
		}, {
			name:     "500 — внутренняя ошибка сервера",
			login:    internalServerErrorLogin,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(`{"login": "` + tt.login + `","password": "` + password + `"}`)
			request, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/user/login", bodyReader)
			require.NoError(t, err)
			res, err := testServer.Client().Do(request)
			require.NoError(t, err)
			res.Body.Close()
			cookies := res.Cookies()

			request, err = http.NewRequest(http.MethodGet, testServer.URL+"/api/user/orders", nil)
			require.NoError(t, err)
			for _, cookie := range cookies {
				request.AddCookie(cookie)
			}
			res, err = testServer.Client().Do(request)
			require.NoError(t, err)
			res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}
