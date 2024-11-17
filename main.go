package main

import (
	"log/slog"
	"net/http"
	"os"
	"sandbox/views/pages"
)

const (
	ADDRESS = "0.0.0.0:3000"
)

func run(logger *slog.Logger) error {
	mux := http.NewServeMux()

	// register all handler functions here
	mux.HandleFunc("/", HomeHandler)

	// serve static files
	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public", fs))

	// start the web server
	logger.Info("starting server", "address", ADDRESS)
	if err := http.ListenAndServe(ADDRESS, mux); err != nil {
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
	h := pages.HomePage("Admin")
	h.Render(r.Context(), w)
}
