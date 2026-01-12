package logging

import configHelpers "github.com/acya-skulskaya/go-musthave-diploma/internal/config/helpers"

type Config struct {
	Level string
}

func Get(config Config) Config {
	//config.Level = configHelpers.CheckAndGetEnvValue(config.Level, "LOG_LEVEL", "INFO")
	config.Level = configHelpers.CheckAndGetEnvValue(config.Level, "LOG_LEVEL", "DEBUG")

	return config
}
