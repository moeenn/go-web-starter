package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sandbox/db"
	"sandbox/internal/config"
	"sandbox/internal/controller"
	"sandbox/internal/lib/middleware"
	"sandbox/internal/service"
)

func run(ctx context.Context, logger *slog.Logger) error {
	globalConfig, err := config.NewConfig()
	if err != nil {
		return err
	}

	// connect to database.
	dbConn, models, err := db.Connect(ctx, globalConfig.Database.Uri)
	if err != nil {
		return err
	}
	defer dbConn.Close(ctx)

	authService := &service.AuthService{
		TokenCookieName: globalConfig.Auth.TokenCookieName,
		Logger:          logger,
		DB:              models,
		Config:          globalConfig,
	}

	authMiddleware := middleware.NewAuthMiddleware(
		globalConfig.Auth.TokenCookieName,
		globalConfig.Auth.JwtSecret,
	)

	mux := http.NewServeMux()

	// serve public assets.
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public", fs))

	// register all controllers.
	controller.NewPublicController(logger).RegisterRoutes(mux)
	controller.NewAuthController(logger, authService, authMiddleware).RegisterRoutes(mux)
	controller.NewDashboardController(logger, authMiddleware).RegisterRoutes(mux)

	// start the web server.
	address := globalConfig.Server.Address()
	logger.Info("starting server", "address", address)
	handler := middleware.Logger(logger, authMiddleware.SetClaimsContext(mux))

	//nolint: exhaustruct
	server := &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadTimeout:       globalConfig.Server.Timeout,
		WriteTimeout:      globalConfig.Server.Timeout,
		IdleTimeout:       globalConfig.Server.Timeout,
		ReadHeaderTimeout: globalConfig.Server.Timeout,
	}

	return server.ListenAndServe()
}

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(ctx, logger); err != nil {
		logger.Error("error", "details", err.Error())
		os.Exit(1)
	}
}
