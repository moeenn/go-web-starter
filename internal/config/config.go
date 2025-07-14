package config

import (
	"fmt"
	"os"
	"sandbox/internal/lib"
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

const (
	defaultServerHost string = "0.0.0.0"
	defaultServerPort string = "3000"
)

type ServerConfig struct {
	Host string
	Port string
}

func NewServerConfig() (*ServerConfig, error) {
	host, err := readString("SERVER_HOST", lib.Ref(defaultServerHost))
	if err != nil {
		return nil, fmt.Errorf("server config error: %w", err)
	}

	port, err := readString("SERVER_PORT", lib.Ref(defaultServerPort))
	if err != nil {
		return nil, fmt.Errorf("server config error: %w", err)
	}

	config := &ServerConfig{host, port}
	return config, nil
}

func (c ServerConfig) Address() string {
	return c.Host + ":" + c.Port
}
