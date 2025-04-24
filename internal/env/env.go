package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return num
}
