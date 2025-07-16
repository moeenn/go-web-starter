package main

import (
	"context"
	"log/slog"
	"os"
	dbmodels "sandbox/db/models"
	"sandbox/internal/config"
	"sandbox/internal/controller"
	customMiddleware "sandbox/internal/lib/middleware"
	"sandbox/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func run(ctx context.Context, logger *slog.Logger) error {
	globalConfig, err := config.NewConfig()
	if err != nil {
		return err
	}

	dbConn, err := pgx.Connect(ctx, globalConfig.Database.Uri)
	if err != nil {
		return err
	}
	defer dbConn.Close(ctx)
	if err := dbConn.Ping(ctx); err != nil {
		return err
	}
	db := dbmodels.New(dbConn)

	authService := &service.AuthService{
		TokenCookieName: globalConfig.Auth.TokenCookieName,
		Logger:          logger,
		DB:              db,
		Config:          globalConfig,
	}

	authMiddleware := customMiddleware.NewAuthMiddleware(
		globalConfig.Auth.TokenCookieName,
		globalConfig.Auth.JwtSecret,
	)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Static("/public", "public")

	// register all controllers.
	controller.NewPublicController(logger).RegisterRoutes(e)
	authController := &controller.AuthController{
		Logger:         logger,
		AuthService:    authService,
		AuthMiddleware: authMiddleware,
	}
	authController.RegisterRoutes(e)
	controller.NewDashboardController(logger).RegisterRoutes(e)

	// start the web server.
	address := globalConfig.Server.Address()
	logger.Info("starting server", "address", address)
	return e.Start(address)
}

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(ctx, logger); err != nil {
		logger.Error("error", "details", err.Error())
		os.Exit(1)
	}
}
