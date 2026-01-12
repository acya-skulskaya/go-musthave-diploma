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
	var successUserID uint
	successUserID = 123
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
	orderId200 := "66886516672865"
	orderId202 := "65273138"
	orderId409 := "22230379"
	setupMockOrder := func(m *mocksOrder.MockOrderRepositoryInterface) {
		m.EXPECT().
			Create(mock.Anything, mock.Anything, orderId200).
			Return(models.Order{
				OrderNumber: orderId200,
				Status:      models.OrderStatusNew,
				UserID:      successUserID,
				Accrual:     0,
			}, order.ErrOrderAlreadyExistsByCurrentUser)

		m.EXPECT().
			Create(mock.Anything, mock.Anything, orderId202).
			Return(models.Order{
				OrderNumber: orderId202,
				Status:      models.OrderStatusNew,
				UserID:      successUserID,
				Accrual:     0,
			}, nil)

		m.EXPECT().
			Create(mock.Anything, mock.Anything, orderId409).
			Return(models.Order{
				OrderNumber: orderId409,
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
		orderId  string
		wantCode int
	}{
		{
			name:     "200 — номер заказа уже был загружен этим пользователем",
			orderId:  orderId200,
			wantCode: http.StatusOK,
		}, {
			name:     "202 — новый номер заказа принят в обработку",
			orderId:  orderId202,
			wantCode: http.StatusAccepted,
		}, {
			name:     "409 — номер заказа уже был загружен другим пользователем",
			orderId:  orderId409,
			wantCode: http.StatusConflict,
		}, {
			name:     "422 — неверный формат номера заказа",
			orderId:  "123",
			wantCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader = strings.NewReader(tt.orderId)
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
