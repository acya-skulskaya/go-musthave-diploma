package router

import (
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config"
	handlerAuth "github.com/acya-skulskaya/go-musthave-diploma/internal/handlers/api/user/auth"
	handlerBalance "github.com/acya-skulskaya/go-musthave-diploma/internal/handlers/api/user/balance"
	handlerOrders "github.com/acya-skulskaya/go-musthave-diploma/internal/handlers/api/user/orders"
	handlerWithdrawals "github.com/acya-skulskaya/go-musthave-diploma/internal/handlers/api/user/withdrawals"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	serviceAuth "github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
	serviceBalance "github.com/acya-skulskaya/go-musthave-diploma/internal/service/balance"
	serviceOrder "github.com/acya-skulskaya/go-musthave-diploma/internal/service/order"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func New(cfg *config.Config, authService *serviceAuth.Service, orderService *serviceOrder.Service, balanceService *serviceBalance.Service) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestEncoder)
	router.Use(middleware.RequestLogger)
	router.Use(chiMiddleware.Compress(5))

	router.Route("/api/user", func(router chi.Router) {
		router.Post("/register", handlerAuth.PostAPIUserRegister(authService)) // — регистрация пользователя;
		router.Post("/login", handlerAuth.PostAPIUserLogin(authService))       //— аутентификация пользователя;

		router.With(middleware.HandleCookieAuth(cfg, authService)).Route("/", func(router chi.Router) {
			router.Route("/orders", func(router chi.Router) {
				router.Get("/", handlerOrders.GetAPIUserOrders(orderService))   // — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
				router.Post("/", handlerOrders.PostAPIUserOrders(orderService)) //— загрузка пользователем номера заказа для расчёта;
			})
			router.Route("/balance", func(router chi.Router) {
				router.Get("/", handlerBalance.GetAPIUserBalance(balanceService))                   //  — получение текущего баланса счёта баллов лояльности пользователя;
				router.Post("/withdraw", handlerBalance.PostAPIUserBalanceWithdraw(balanceService)) //— запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
			})
			router.Get("/withdrawals", handlerWithdrawals.GetAPIUserWithdrawals(balanceService)) // — получение информации о выводе средств с накопительного счёта пользователем.
		})
	})

	return router
}
