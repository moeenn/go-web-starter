package main

import (
	"app/internal/config"
	"app/internal/controller"
	"app/internal/lib/middleware"
	"app/internal/repo"
	"app/internal/service"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func run(ctx context.Context, logger *slog.Logger) error {
	globalConfig, err := config.NewConfig()
	if err != nil {
		return err
	}

	// connect to database.
	db, err := sqlx.Open("postgres", globalConfig.Database.Uri)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// create server multiplexer.
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public", fs))

	// -------------------------------------------------------------------------
	//
	// initialize repos.
	//
	// -------------------------------------------------------------------------
	userRepo := repo.NewUserRepo(db, logger)

	// -------------------------------------------------------------------------
	//
	// initialize services.
	//
	// -------------------------------------------------------------------------
	authService := &service.AuthService{
		TokenCookieName: globalConfig.Auth.TokenCookieName,
		Logger:          logger,
		UserRepo:        userRepo,
		Config:          globalConfig,
	}

	// -------------------------------------------------------------------------
	//
	// initialize middleware.
	//
	// -------------------------------------------------------------------------
	authMiddleware := middleware.NewAuthMiddleware(
		globalConfig.Auth.TokenCookieName,
		globalConfig.Auth.JwtSecret,
	)

	// -------------------------------------------------------------------------
	//
	// register controllers.
	//
	// -------------------------------------------------------------------------
	controller.NewPublicController(logger).RegisterRoutes(mux)
	controller.NewAuthController(logger, authService, authMiddleware).RegisterRoutes(mux)
	controller.NewDashboardController(logger, authMiddleware, userRepo).RegisterRoutes(mux)

	// -------------------------------------------------------------------------
	//
	// start web server.
	//
	// -------------------------------------------------------------------------
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
