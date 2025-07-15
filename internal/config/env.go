package config

import (
	"fmt"
	"os"
	"strconv"
)

func readString(key string, fallback *string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		if fallback == nil {
			return "", fmt.Errorf("environment variable not found: %q", key)
		}
		return *fallback, nil
	}

	return value, nil
}

func readInt(key string, fallback *int) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		if fallback == nil {
			return 0, fmt.Errorf("environment variable not found: %q", key)
		}
		return *fallback, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q value is not a valid int", key)
	}

	return parsedValue, nil
}
