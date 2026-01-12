package auth_test

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

func TestPostAPIUserLogin(t *testing.T) {
	successLogin := "success_login"
	successPassword := "success_password"
	successPasswordHashed, _ := serviceAuth.HashPassword(successPassword)
	internalServerErrorLogin := "internalServerError_login"
	userRepo := mocksUser.NewMockUserRepositoryInterface(t)
	setupMock := func(m *mocksUser.MockUserRepositoryInterface) {
		m.EXPECT().
			Get(mock.Anything, successLogin).
			Return(models.User{Login: successLogin, Password: successPasswordHashed}, nil)
		m.EXPECT().
			Get(mock.Anything, internalServerErrorLogin).
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
			name:     "200 — пользователь успешно аутентифицирован",
			data:     `{"login": "` + successLogin + `","password": "` + successPassword + `"}`,
			wantCode: http.StatusOK,
		}, {
			name:     "400 — неверный формат запроса",
			data:     `{"login1": "","password2": ""}`,
			wantCode: http.StatusBadRequest,
		}, {
			name:     "401 — неверная пара логин/пароль",
			data:     `{"login": "` + successLogin + `","password": "` + successPassword + `_1"}`,
			wantCode: http.StatusUnauthorized,
		}, {
			name:     "500 — внутренняя ошибка сервера",
			data:     `{"login": "` + internalServerErrorLogin + `","password": "pass123"}`,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := strings.NewReader(tt.data)
			request, err := http.NewRequest(http.MethodPost, testServer.URL+"/api/user/login", bodyReader)
			require.NoError(t, err)
			res, err := testServer.Client().Do(request)
			require.NoError(t, err)
			res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, tt.wantCode, res.StatusCode)
		})
	}
}
