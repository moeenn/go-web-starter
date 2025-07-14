package main

import (
	"log/slog"
	"net/http"
	"os"
	"sandbox/internal/config"
	"sandbox/internal/lib"
	"sandbox/views/pages"
)

func run(logger *slog.Logger) error {
	serverConfig, err := config.NewServerConfig()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("GET /auth/login", LoginPageHandler)
	mux.HandleFunc("GET /auth/forgot-password", ForgotPasswordPageHandler)

	// serve static files
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public", fs))

	// start the web server
	address := serverConfig.Address()
	logger.Info("starting server", "address", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		return err
	}

	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("error", "details", err.Error())
		os.Exit(1)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	h := pages.HomePage(pages.HomePageProps{
		IsLoggedIn:     false,
		FlashedMessage: lib.Ref("Welcome to our simple home-page"),
	})

	h.Render(r.Context(), w)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	h := pages.LoginPage()
	h.Render(r.Context(), w)
}

func ForgotPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	h := pages.ForgotPasswordPage()
	h.Render(r.Context(), w)
}
