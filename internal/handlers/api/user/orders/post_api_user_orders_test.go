package orders_test

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	mocksBalance "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance/mocks"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order"
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

func TestPostAPIUserOrders(t *testing.T) {
	var successUserID uint = 123
	successLogin := "success_login"
	password := "success_password"
	passwordHashed, _ := serviceAuth.HashPassword(password)
	userRepo := mocksUser.NewMockUserRepositoryInterface(t)
	setupMockUser := func(m *mocksUser.MockUserRepositoryInterface) {
		m.EXPECT().
			Get(mock.Anything, successLogin).
			Return(models.User{Login: successLogin, Password: passwordHashed, ID: successUserID}, nil)

	}
	setupMockUser(userRepo)

	balanceRepo := mocksBalance.NewMockBalanceRepositoryInterface(t)

	orderRepo := mocksOrder.NewMockOrderRepositoryInterface(t)
	orderID200 := "66886516672865"
	orderID202 := "65273138"
	orderID409 := "22230379"
	setupMockOrder := func(m *mocksOrder.MockOrderRepositoryInterface) {
		m.EXPECT().
			Create(mock.Anything, mock.Anything, orderID200).
			Return(models.Order{
				OrderNumber: orderID200,
				Status:      models.OrderStatusNew,
				UserID:      successUserID,
				Accrual:     0,
			}, order.ErrOrderAlreadyExistsByCurrentUser)

		m.EXPECT().
			Create(mock.Anything, mock.Anything, orderID202).
			Return(models.Order{
				OrderNumber: orderID202,
				Status:      models.OrderStatusNew,
				UserID:      successUserID,
				Accrual:     0,
			}, nil)

		m.EXPECT().
			Create(mock.Anything, mock.Anything, orderID409).
			Return(models.Order{
				OrderNumber: orderID409,
				Status:      models.OrderStatusNew,
				UserID:      successUserID,
				Accrual:     0,
			}, order.ErrOrderAlreadyExistsFromAnotherUser)

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

	bodyReader := strings.NewReader(`{"login": "` + successLogin + `","password": "` + password + `"}`)
	request, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/user/login", bodyReader)
	require.NoError(t, err)
	res, err := testServer.Client().Do(request)
	require.NoError(t, err)
	res.Body.Close()
	cookies := res.Cookies()

	tests := []struct {
		name     string
		orderID  string
		wantCode int
	}{
		{
			name:     "200 — номер заказа уже был загружен этим пользователем",
			orderID:  orderID200,
			wantCode: http.StatusOK,
		}, {
			name:     "202 — новый номер заказа принят в обработку",
			orderID:  orderID202,
			wantCode: http.StatusAccepted,
		}, {
			name:     "409 — номер заказа уже был загружен другим пользователем",
			orderID:  orderID409,
			wantCode: http.StatusConflict,
		}, {
			name:     "422 — неверный формат номера заказа",
			orderID:  "123",
			wantCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader = strings.NewReader(tt.orderID)
			request, err = http.NewRequest(http.MethodPost, testServer.URL+"/api/user/orders", bodyReader)
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
