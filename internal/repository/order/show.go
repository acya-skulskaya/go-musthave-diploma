package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (repo *OrderRepository) Show(ctx context.Context, orderNumber string) (models.Order, error) {
	sql := "SELECT order_number, user_id, status, accrual, uploaded_at, processed_at FROM " + models.OrderTable + " WHERE order_number = $1"
	var order models.Order
	var accrual pgtype.Numeric
	row := repo.DBPool.QueryRow(ctx, sql, orderNumber)
	err := row.Scan(&order.OrderNumber, &order.UserID, &order.Status, &accrual, &order.UploadedAt, &order.ProcessedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Order{}, ErrNotFound
	} else if err != nil {
		return models.Order{}, fmt.Errorf("could not get order: %w", err)
	}

	if !accrual.NaN {
		accrualFloat8, err := accrual.Float64Value()
		if err != nil {
			return models.Order{}, fmt.Errorf("order %s: could not convert accrual value: %w", order.OrderNumber, err)
		}
		order.Accrual = accrualFloat8.Float64
	}

	return order, nil
}
