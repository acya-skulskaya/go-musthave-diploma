package order

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
)

func (repo *OrderRepository) UpdateOrderStatus(ctx context.Context, orderNumber string, status string) error {
	sql := "UPDATE " + models.OrderTable + " SET status=$1 WHERE order_number = $2"

	_, err := repo.DBPool.Exec(ctx, sql, status, orderNumber)
	if err != nil {
		return fmt.Errorf("could not update order: %w", err)
	}

	return nil
}
