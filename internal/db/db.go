package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config/db"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(ctx context.Context, cfg db.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("could not create pgx connection pool: %w", err)
	}

	if err = applyMigrations(&cfg); err != nil {
		return nil, fmt.Errorf("could not apply migrations: %w", err)
	}

	return pool, nil
}

func applyMigrations(cfg *db.Config) error {
	migrations, err := migrate.New(cfg.MigrationsPath, cfg.DSN)
	if err != nil {
		return fmt.Errorf("could not init migrations: %w", err)
	}

	err = migrations.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		logger.Log.Info("nothing to migrate")
	} else if err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	return nil
}
