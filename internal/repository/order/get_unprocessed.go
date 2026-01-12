package order

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
)

func (repo *OrderRepository) GetUnprocessed(ctx context.Context, limit uint) ([]models.Order, error) {
	sql := "SELECT order_number, user_id, status, uploaded_at FROM " + models.OrderTable + " WHERE status=$1 OR status=$2 order by uploaded_at asc LIMIT $3"

	var orders []models.Order

	rows, err := repo.DBPool.Query(ctx, sql, models.OrderStatusNew, models.OrderStatusProcessing, limit)
	if err != nil {
		return []models.Order{}, fmt.Errorf("could not query orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.OrderNumber, &order.UserID, &order.Status, &order.UploadedAt)
		if err != nil {
			return []models.Order{}, fmt.Errorf("could not scan order row: %w", err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}
