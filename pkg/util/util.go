package util

import (
	"os"
)

func ReadEnvOrDie(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		panic("env:" + key + " not set")
	}

	return value
}
