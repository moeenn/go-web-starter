package config

import (
	"fmt"
	"sandbox/internal/lib"
	"time"
)

const (
	defaultServerHost       string        = "0.0.0.0"
	defaultServerPort       string        = "3000"
	ddefaultSeverTimeout    time.Duration = time.Second * 10
	defaultJwtExpiryMinutes int           = 60
	defaultTokenCookieName  string        = "auth.token"
)

type serverConfig struct {
	Host    string
	Port    string
	Timeout time.Duration
}

func newServerConfig() (*serverConfig, error) {
	host, err := readString("SERVER_HOST", lib.Ref(defaultServerHost))
	if err != nil {
		return nil, err
	}

	port, err := readString("SERVER_PORT", lib.Ref(defaultServerPort))
	if err != nil {
		return nil, err
	}

	config := &serverConfig{host, port, ddefaultSeverTimeout}
	return config, nil
}

func (c serverConfig) Address() string {
	return c.Host + ":" + c.Port
}

type databaseConfig struct {
	Uri string
}

func newDatabaseConfig() (*databaseConfig, error) {
	uri, err := readString("GOOSE_DBSTRING", nil)
	if err != nil {
		return nil, err
	}

	config := &databaseConfig{
		Uri: uri,
	}

	return config, nil
}

type authConfig struct {
	JwtSecret        []byte
	JwtExpiryMinutes time.Duration
	TokenCookieName  string
}

func newAuthConfig() (*authConfig, error) {
	secret, err := readString("JWT_SECRET", nil)
	if err != nil {
		return nil, err
	}

	expiryMinutes, err := readInt("JWT_EXPIRY_MINUTES", lib.Ref(defaultJwtExpiryMinutes))
	if err != nil {
		return nil, err
	}

	tokenCookieName, err := readString("AUTH_TOKEN_COOKIE_NAME", lib.Ref(defaultTokenCookieName))
	if err != nil {
		return nil, err
	}

	config := &authConfig{
		JwtSecret:        []byte(secret),
		JwtExpiryMinutes: time.Minute * time.Duration(expiryMinutes),
		TokenCookieName:  tokenCookieName,
	}

	return config, nil
}

type Config struct {
	Server   *serverConfig
	Database *databaseConfig
	Auth     *authConfig
}

func NewConfig() (*Config, error) {
	server, err := newServerConfig()
	if err != nil {
		return nil, fmt.Errorf("server config: %w", err)
	}

	db, err := newDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("database config: %w", err)
	}

	auth, err := newAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("auth config: %w", err)
	}

	config := &Config{
		Server:   server,
		Database: db,
		Auth:     auth,
	}

	return config, nil
}
