package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (repo *OrderRepository) Finish(ctx context.Context, orderNumber string, status string, accrual float64, userID uint) error {
	traceID, ok := ctx.Value(middleware.ContextKeyTraceID).(string)
	if !ok {
		traceID = ""
	}

	tx, err := repo.DBPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin transaction to finish order: %w", err)
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(ctx); rErr != nil {
				if !errors.Is(rErr, pgx.ErrTxClosed) {
					logger.Log.Error("could not roll back transaction to finish order",
						zap.String(middleware.TraceID, traceID),
						zap.Error(err))
				}
			}
		} else {
			if rErr := tx.Commit(ctx); rErr != nil {
				logger.Log.Error("could not commit transaction to finish order",
					zap.String(middleware.TraceID, traceID),
					zap.Error(err))
			}
		}
	}()

	sql := "UPDATE " + models.OrderTable + " SET status=$1, accrual=$2, processed_at=NOW() WHERE order_number = $3"
	_, err = tx.Exec(ctx, sql, status, accrual, orderNumber)
	if err != nil {
		return fmt.Errorf("could not update order: %w", err)
	}

	userBalanceExists := true
	sql = "SELECT user_id, current FROM " + models.UserBalanceTable + " WHERE user_id = $1 FOR UPDATE"
	var userBalance models.UserBalance
	row := tx.QueryRow(ctx, sql, userID)
	err = row.Scan(&userBalance.UserID, &userBalance.Current)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			userBalanceExists = false
		} else {
			return fmt.Errorf("could not get user balance: %w", err)
		}
	}

	sql = "INSERT INTO " + models.UserBalanceFlowTable + " (user_id, order_number, amount) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, sql, userID, orderNumber, accrual)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return balance.ErrAccrualForThisOrderAlreadyExists
		} else {
			return fmt.Errorf("could not insert accrual into user_balance_flow: %w", err)
		}
	}

	if !userBalanceExists {
		sql = "INSERT INTO " + models.UserBalanceTable + " (user_id, current) VALUES ($1, $2)"
		_, err = tx.Exec(ctx, sql, userID, accrual)
		if err != nil {
			return fmt.Errorf("could not insert user balance: %w", err)
		}
	} else {
		userBalance.Current += accrual

		sql = "UPDATE " + models.UserBalanceTable + " SET current=$1, updated_at=NOW() WHERE user_id = $2"
		_, err = tx.Exec(ctx, sql, userBalance.Current, userID)
		if err != nil {
			return fmt.Errorf("could not update user balance: %w", err)
		}
	}

	return nil
}
