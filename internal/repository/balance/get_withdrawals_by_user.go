package balance

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
)

func (repo *BalanceRepository) GetWithdrawalsByUser(ctx context.Context, userID uint) ([]models.UserBalanceFlow, error) {
	sql := "SELECT user_id, order_number, amount, processed_at FROM " + models.UserBalanceFlowTable + " WHERE user_id = $1 AND amount < 0 ORDER BY processed_at DESC"

	rows, err := repo.DBPool.Query(ctx, sql, userID)
	if err != nil {
		return []models.UserBalanceFlow{}, fmt.Errorf("could not query user_balance_flow: %w", err)
	}
	defer rows.Close()

	var flow []models.UserBalanceFlow

	for rows.Next() {
		var userBalanceFlow models.UserBalanceFlow
		err = rows.Scan(&userBalanceFlow.UserID, &userBalanceFlow.OrderNumber, &userBalanceFlow.Amount, &userBalanceFlow.ProcessedAt)
		if err != nil {
			return []models.UserBalanceFlow{}, fmt.Errorf("could not scan user_balance_flow row: %w", err)
		}

		flow = append(flow, userBalanceFlow)
	}

	return flow, nil
}
