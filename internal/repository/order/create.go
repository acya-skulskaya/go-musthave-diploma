package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (repo *OrderRepository) Create(ctx context.Context, userID uint, orderNumber string) (models.Order, error) {
	sql := "INSERT INTO " + models.OrderTable + " (order_number, user_id, status) VALUES ($1, $2, $3)"
	_, err := repo.DBPool.Exec(ctx, sql, orderNumber, userID, string(models.OrderStatusNew))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return models.Order{}, ErrOrderNumberAlreadyExists
		} else {
			return models.Order{}, fmt.Errorf(ErrMsgCouldNotCreate+": %w", err)
		}
	}

	order, err := repo.Show(ctx, orderNumber)
	if err != nil {
		return models.Order{}, fmt.Errorf("could not get created user: %w", err)
	}

	return order, nil
}
