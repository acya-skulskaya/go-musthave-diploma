package order

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
)

func (s *Service) All(ctx context.Context) ([]models.Order, error) {
	user, ok := ctx.Value(middleware.ContextKeyUser).(models.User)
	if !ok {
		return []models.Order{}, auth.ErrCouldNotGetUserFromContext
	}

	orders, err := s.Repo.AllByUser(ctx, user.ID)
	if err != nil {
		return []models.Order{}, fmt.Errorf("could not get all orders by user: %w", err)
	}

	return orders, nil
}
