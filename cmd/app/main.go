package main

import (
	"log/slog"
	"os"
	"sandbox/internal/config"
	"sandbox/internal/controller"

	"github.com/labstack/echo/v4"
)

func run(logger *slog.Logger) error {
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		return err
	}

	e := echo.New()
	e.Static("/public", "public")
	controller.NewPublicController(logger).RegisterRoutes(e)
	controller.NewAuthController(logger).RegisterRoutes(e)

	// start the web server.
	address := serverConfig.Address()
	logger.Info("starting server", "address", address)
	return e.Start(address)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("error", "details", err.Error())
		os.Exit(1)
	}
}
