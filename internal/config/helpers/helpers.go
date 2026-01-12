package helpers

import "os"

func CheckAndGetEnvValue(value string, key string, defaultValue string) string {
	if len(value) == 0 {
		envValue, ok := os.LookupEnv(key)
		if ok {
			value = envValue
		}
	}

	if len(value) == 0 && len(defaultValue) > 0 {
		value = defaultValue
	}

	return value
}
