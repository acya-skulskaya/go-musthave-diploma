package balance

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
)

func (s *Service) Get(ctx context.Context) (models.UserBalance, error) {
	user, ok := ctx.Value(middleware.ContextKeyUser).(models.User)
	if !ok {
		return models.UserBalance{}, auth.ErrCouldNotGetUserFromContext
	}

	userBalance, err := s.Repo.GetByUser(ctx, user.ID)
	if err != nil {
		return models.UserBalance{}, fmt.Errorf("could not get user balance: %w", err)
	}

	return userBalance, nil
}
