package main

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/db"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	repoBalance "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"
	repoOrder "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order"
	repoUser "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/user"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/router"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/accrual"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
	serviceBalance "github.com/acya-skulskaya/go-musthave-diploma/internal/service/balance"
	serviceOrder "github.com/acya-skulskaya/go-musthave-diploma/internal/service/order"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/workers/orderprocesser"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	if err := logger.Init(cfg.Logging); err != nil {
		return fmt.Errorf("could not init logging: %w", err)
	}
	logger.Log.Info("config loaded", zap.Any("cfg", cfg))

	dbPool, err := db.New(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	userRepo := repoUser.NewUserRepository(dbPool)
	balanceRepo := repoBalance.NewBalanceRepository(dbPool)
	orderRepo := repoOrder.NewOrderRepository(dbPool)

	authService := auth.New(cfg.Auth.SecretKey, userRepo)
	orderService := serviceOrder.New(orderRepo)
	balanceService := serviceBalance.New(balanceRepo)

	accrualService := accrual.New(cfg.Handlers)
	orderProcessor := orderprocesser.New(orderRepo, accrualService)
	go orderProcessor.Run(ctx)

	r := router.New(&cfg, authService, orderService, balanceService)
	srv := &http.Server{
		Addr:    cfg.Handlers.ServerAddress,
		Handler: r,
	}

	//nolint:wrapcheck //no need
	return srv.ListenAndServe()
}
