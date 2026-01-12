package balance

import (
	"context"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BalanceRepositoryInterface interface {
	Withdraw(ctx context.Context, userID uint, orderNumber string, sum float64) error
	GetByUser(ctx context.Context, userID uint) (models.UserBalance, error)
	GetWithdrawalsByUser(ctx context.Context, userID uint) ([]models.UserBalanceFlow, error)
}

type BalanceRepository struct {
	DBPool *pgxpool.Pool
}

func NewBalanceRepository(dbPool *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{
		DBPool: dbPool,
	}
}
