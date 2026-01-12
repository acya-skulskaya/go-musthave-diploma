package auth_test

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	mocksBalance "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance/mocks"
	mocksOrder "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order/mocks"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/user"
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

func TestPostAPIUserRegister(t *testing.T) {
	userSuccess := "user_success"
	userAlreadyExists := "user_already_exists"
	userInternalServerError := "user_internalServerError"
	userRepo := mocksUser.NewMockUserRepositoryInterface(t)
	setupMock := func(m *mocksUser.MockUserRepositoryInterface) {
		m.EXPECT().
			Create(mock.Anything, userSuccess, mock.Anything).
			Return(models.User{Login: userSuccess}, nil)
		m.EXPECT().
			Create(mock.Anything, userAlreadyExists, mock.Anything).
			Return(models.User{}, user.ErrLoginAlreadyExists)
		m.EXPECT().
			Create(mock.Anything, userInternalServerError, mock.Anything).
			Return(models.User{}, assert.AnError)
	}
	setupMock(userRepo)

	balanceRepo := mocksBalance.NewMockBalanceRepositoryInterface(t)
	orderRepo := mocksOrder.NewMockOrderRepositoryInterface(t)

	authService := serviceAuth.New("test", userRepo)
	orderService := serviceOrder.New(orderRepo)
	balanceService := serviceBalance.New(balanceRepo)

	cfg := config.Config{}

	r := router.New(&cfg, authService, orderService, balanceService)
	testServer := httptest.NewServer(r)
	defer testServer.Close()

	tests := []struct {
		name     string
		data     string
		wantCode int
	}{
		{
			name:     "200 — пользователь успешно зарегистрирован и аутентифицирован",
			data:     `{"login": "` + userSuccess + `","password": "password"}`,
			wantCode: http.StatusOK,
		}, {
			name:     "400 — неверный формат запроса",
			data:     `{"login1": "","password1": "password"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "409 — логин уже занят",
			data:     `{"login": "` + userAlreadyExists + `","password": "password"}`,
			wantCode: http.StatusConflict,
		}, {
			name:     "500 — внутренняя ошибка сервера",
			data:     `{"login": "` + userInternalServerError + `","password": "password"}`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.data)
			request, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/user/register", bodyReader)
			require.NoError(t, err)
			res, err := testServer.Client().Do(request)
			require.NoError(t, err)
			res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}
