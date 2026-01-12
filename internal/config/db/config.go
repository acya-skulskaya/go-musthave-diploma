package db

import configHelpers "github.com/acya-skulskaya/go-musthave-diploma/internal/config/helpers"

type Config struct {
	DSN            string
	MigrationsPath string
}

func Get(config Config) Config {
	config.DSN = configHelpers.CheckAndGetEnvValue(config.DSN, "DATABASE_URI", "postgres://user:pass@localhost:5432/test-db?sslmode=disable")
	config.MigrationsPath = configHelpers.CheckAndGetEnvValue(config.MigrationsPath, "MIGRATIONS_PATH", "file://./migrations")
	//config.MigrationsPath = configHelpers.CheckAndGetEnvValue(config.MigrationsPath, "MIGRATIONS_PATH", "file://./go-musthave-diploma/migrations")

	return config
}
