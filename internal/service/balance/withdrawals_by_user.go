package balance

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/repository/balance"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
)

func (s *Service) WithdrawalsByUser(ctx context.Context) ([]models.UserBalanceFlow, error) {
	user, ok := ctx.Value(middleware.ContextKeyUser).(models.User)
	if !ok {
		return []models.UserBalanceFlow{}, auth.ErrCouldNotGetUserFromContext
	}

	userBalanceFlow, err := s.Repo.GetWithdrawalsByUser(ctx, user.ID)
	if err != nil {
		return []models.UserBalanceFlow{}, fmt.Errorf("could not get user balance flow: %w", err)
	}

	if len(userBalanceFlow) == 0 {
		return []models.UserBalanceFlow{}, balance.ErrNoWithdrawals
	}

	return userBalanceFlow, nil
}
