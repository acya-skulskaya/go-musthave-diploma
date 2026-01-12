package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (repo *UserRepository) Create(ctx context.Context, login string, password string) (models.User, error) {
	sql := "INSERT INTO " + models.UserTable + " (login, password) VALUES ($1, $2)"
	_, err := repo.DBPool.Exec(ctx, sql, login, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return models.User{}, ErrLoginAlreadyExists
		} else {
			return models.User{}, fmt.Errorf(ErrMsgCouldNotCreate+": %w", err)
		}
	}

	user, err := repo.Get(ctx, login)
	if err != nil {
		return models.User{}, fmt.Errorf("could not get created user: %w", err)
	}

	return user, nil
}
