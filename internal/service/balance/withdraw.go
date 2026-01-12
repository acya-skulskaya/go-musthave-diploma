package balance

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
)

func (s *Service) Withdraw(ctx context.Context, orderNumber string, sum float64) error {
	user, ok := ctx.Value(middleware.ContextKeyUser).(models.User)
	if !ok {
		return auth.ErrCouldNotGetUserFromContext
	}

	userBalance, err := s.Repo.GetByUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("could not get user balance: %w", err)
	}

	if userBalance.Current < sum {
		return balance.ErrNotEnoughBalanceToWithdraw
	}

	err = s.Repo.Withdraw(ctx, user.ID, orderNumber, sum)
	if err != nil {
		if errors.Is(err, balance.ErrNotEnoughBalanceToWithdraw) {
			return balance.ErrNotEnoughBalanceToWithdraw
		} else if errors.Is(err, balance.ErrWithdrawalForThisOrderAlreadyExists) {
			return balance.ErrWithdrawalForThisOrderAlreadyExists
		}

		return fmt.Errorf("could not withdraw user balance: %w", err)
	}

	return nil
}
