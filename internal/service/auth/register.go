package auth

import (
	"context"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/request"
)

func (s *Service) Register(ctx context.Context, req *request.UserCredits) (models.User, error) {
	password, err := HashPassword(req.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("hashing password failed: %w", err)
	}

	user, err := s.Repo.Create(ctx, req.Login, password)
	if err != nil {
		return models.User{}, fmt.Errorf("could not craete user: %w", err)
	}

	return user, nil
}
