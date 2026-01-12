package user

import (
	"context"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, login string, password string) (models.User, error)
	Get(ctx context.Context, id any) (models.User, error)
}

type UserRepository struct {
	DBPool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DBPool: dbPool,
	}
}
