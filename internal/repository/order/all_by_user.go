package order

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func (repo *OrderRepository) AllByUser(ctx context.Context, userID uint) ([]models.Order, error) {
	sql := "SELECT order_number, user_id, status, accrual, uploaded_at, processed_at FROM " + models.OrderTable + " WHERE user_id = $1 order by uploaded_at desc"

	var orders []models.Order

	rows, err := repo.DBPool.Query(ctx, sql, userID)
	if err != nil {
		return []models.Order{}, fmt.Errorf("could not query orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		var accrual pgtype.Numeric
		err = rows.Scan(&order.OrderNumber, &order.UserID, &order.Status, &accrual, &order.UploadedAt, &order.ProcessedAt)
		if err != nil {
			return []models.Order{}, fmt.Errorf("could not scan order row: %w", err)
		}

		if !accrual.NaN {
			accrualFloat8, err := accrual.Float64Value()
			if err != nil {
				return []models.Order{}, fmt.Errorf("order %s: could not convert accrual value: %w", order.OrderNumber, err)
			}
			order.Accrual = accrualFloat8.Float64
		}

		orders = append(orders, order)
	}

	return orders, nil
}
