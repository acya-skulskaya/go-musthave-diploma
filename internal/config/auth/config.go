package auth

import configHelpers "github.com/acya-skulskaya/go-musthave-diploma/internal/config/helpers"

type Config struct {
	SecretKey string
}

func Get(config Config) Config {
	config.SecretKey = configHelpers.CheckAndGetEnvValue(config.SecretKey, "SECRET_KEY", "secret")

	return config
}
