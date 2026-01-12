package handlers

import configHelpers "github.com/acya-skulskaya/go-musthave-diploma/internal/config/helpers"

type Config struct {
	ServerAddress  string
	AccrualAddress string
}

func Get(config Config) Config {
	config.ServerAddress = configHelpers.CheckAndGetEnvValue(config.ServerAddress, "RUN_ADDRESS", "localhost:8080")
	config.AccrualAddress = configHelpers.CheckAndGetEnvValue(config.AccrualAddress, "ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8081")

	return config
}
