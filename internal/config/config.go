package config

import (
	"flag"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config/auth"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config/db"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config/handlers"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config/logging"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Handlers handlers.Config
	DB       db.Config
	Logging  logging.Config
	Auth     auth.Config
}

func GetConfig() (Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.Handlers.ServerAddress, "a", "", "address of HTTP server")
	flag.StringVar(&cfg.DB.DSN, "d", "", "connection settings for pgsql")
	flag.StringVar(&cfg.Handlers.AccrualAddress, "r", "", "accrual system address")
	dotEnvFilePath := ".env"
	flag.StringVar(&dotEnvFilePath, "e", ".env", "path to .env file")
	flag.Parse()

	if _, err := os.Stat(dotEnvFilePath); err == nil {
		err := godotenv.Load()
		if err != nil {
			return cfg, fmt.Errorf("error loading dot env file: %w", err)
		}
	}

	cfg.Handlers = handlers.Get(cfg.Handlers)
	cfg.DB = db.Get(cfg.DB)
	cfg.Logging = logging.Get(cfg.Logging)
	cfg.Auth = auth.Get(cfg.Auth)

	return cfg, nil
}
