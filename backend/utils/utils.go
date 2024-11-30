package utils

import (
	"backend/loggers"
	"os"
)

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		loggers.Error.Fatalf("Environment variable %s not set", key)
	}
	loggers.Info.Printf("Environment variable %s set to %s", key, value)
	return value
}
