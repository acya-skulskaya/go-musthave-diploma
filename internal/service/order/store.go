package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	orderRepo "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/auth"
)

func (s *Service) Store(ctx context.Context, orderNumber string) (models.Order, error) {
	user, ok := ctx.Value(middleware.ContextKeyUser).(models.User)
	if !ok {
		return models.Order{}, auth.ErrCouldNotGetUserFromContext
	}

	order, err := s.Repo.Create(ctx, user.ID, orderNumber)
	if err != nil {
		if errors.Is(err, orderRepo.ErrOrderNumberAlreadyExists) {
			order, err = s.Repo.Show(ctx, orderNumber)
			if err != nil {
				return models.Order{}, fmt.Errorf("could not create order and get it: %w", err)
			}

			if order.UserID == user.ID {
				return order, orderRepo.ErrOrderAlreadyExistsByCurrentUser
			} else {
				return order, orderRepo.ErrOrderAlreadyExistsFromAnotherUser
			}
		} else {
			return models.Order{}, fmt.Errorf("could not create order: %w", err)
		}
	}

	return order, nil
}
