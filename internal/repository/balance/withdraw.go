package balance

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (repo *BalanceRepository) Withdraw(ctx context.Context, userID uint, orderNumber string, sum float64) error {
	traceID, ok := ctx.Value(middleware.ContextKeyTraceID).(string)
	if !ok {
		traceID = ""
	}

	tx, err := repo.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin transaction to withdraw user balance: %w", err)
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(ctx); rErr != nil {
				if !errors.Is(rErr, pgx.ErrTxClosed) {
					logger.Log.Error("could not roll back transaction to withdraw user balance",
						zap.String(middleware.TraceID, traceID),
						zap.Error(err))
				}
			}
		} else {
			if rErr := tx.Commit(ctx); rErr != nil {
				logger.Log.Error("could not commit transaction to withdraw",
					zap.String(middleware.TraceID, traceID),
					zap.Error(err))
			}
		}
	}()

	sql := "SELECT user_id, current, withdrawn, updated_at FROM " + models.UserBalanceTable + " WHERE user_id = $1 FOR UPDATE"
	var userBalance models.UserBalance
	var withdrawn pgtype.Numeric
	row := tx.QueryRow(ctx, sql, userID)
	err = row.Scan(&userBalance.UserID, &userBalance.Current, &withdrawn, &userBalance.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotEnoughBalanceToWithdraw
		}

		return fmt.Errorf("could not get user balance: %w", err)
	}

	if !withdrawn.NaN {
		withdrawnFloat8, err := withdrawn.Float64Value()
		if err != nil {
			return fmt.Errorf("user %d balance: could not convert withdrawn value: %w", userID, err)
		}
		userBalance.Withdrawn = withdrawnFloat8.Float64
	}

	if userBalance.Current < sum {
		return ErrNotEnoughBalanceToWithdraw
	}

	sql = "INSERT INTO " + models.UserBalanceFlowTable + " (user_id, order_number, amount) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, sql, userID, orderNumber, sum*-1)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrWithdrawalForThisOrderAlreadyExists
		} else {
			return fmt.Errorf("could not insert into user_balance_flow: %w", err)
		}
	}

	userBalance.Current -= sum
	userBalance.Withdrawn += sum

	sql = "UPDATE " + models.UserBalanceTable + " SET current = $1, withdrawn = $2, updated_at = NOW()  WHERE user_id= $3"
	_, err = tx.Exec(ctx, sql, userBalance.Current, userBalance.Withdrawn, userID)
	if err != nil {
		return fmt.Errorf("could not update user_balance: %w", err)
	}

	return nil
}
