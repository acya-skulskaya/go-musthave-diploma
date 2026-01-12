package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	"github.com/jackc/pgx/v5"
	"reflect"
)

func (repo *UserRepository) Get(ctx context.Context, id any) (models.User, error) {
	field := "id"
	if reflect.TypeOf(id).Kind() == reflect.String {
		field = "login"
	} else if reflect.TypeOf(id).Kind() != reflect.Uint {
		return models.User{}, fmt.Errorf("id should by of type int or string, got %s", reflect.TypeOf(id).Kind().String())
	}

	sql := "SELECT id, login, password, created_at FROM " + models.UserTable + " WHERE " + field + "= $1"
	var user models.User
	row := repo.DBPool.QueryRow(ctx, sql, id)
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	} else if err != nil {
		return models.User{}, fmt.Errorf("could not get user: %w", err)
	}

	return user, nil
}
