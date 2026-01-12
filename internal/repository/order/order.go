package order

import (
	"context"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepositoryInterface interface {
	Create(ctx context.Context, userID uint, orderNumber string) (models.Order, error)
	AllByUser(ctx context.Context, userID uint) ([]models.Order, error)
	Show(ctx context.Context, orderNumber string) (models.Order, error)
	GetUnprocessed(ctx context.Context, limit uint) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderNumber string, status string) error
	Finish(ctx context.Context, orderNumber string, status string, accrual float64, userID uint) error
}

type OrderRepository struct {
	DBPool *pgxpool.Pool
}

func NewOrderRepository(dbPool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		DBPool: dbPool,
	}
}
