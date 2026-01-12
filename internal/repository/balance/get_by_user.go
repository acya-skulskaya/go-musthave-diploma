package balance

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (repo *BalanceRepository) GetByUser(ctx context.Context, userID uint) (models.UserBalance, error) {
	sql := "SELECT user_id, current, withdrawn, updated_at FROM " + models.UserBalanceTable + " WHERE user_id = $1"
	var userBalance models.UserBalance
	var withdrawn pgtype.Numeric
	row := repo.DBPool.QueryRow(ctx, sql, userID)
	err := row.Scan(&userBalance.UserID, &userBalance.Current, &withdrawn, &userBalance.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.UserBalance{
			UserID:    userID,
			Current:   0,
			Withdrawn: 0,
		}, nil
	} else if err != nil {
		return models.UserBalance{}, fmt.Errorf("could not get user balance: %w", err)
	}

	if !withdrawn.NaN {
		withdrawnFloat8, err := withdrawn.Float64Value()
		if err != nil {
			return models.UserBalance{}, fmt.Errorf("user %d balance: could not convert withdrawn value: %w", userID, err)
		}
		userBalance.Withdrawn = withdrawnFloat8.Float64
	}

	return userBalance, nil
}
