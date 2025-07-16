package controller

import (
	"log/slog"
	"net/http"
	"sandbox/internal/lib"
	"sandbox/views/pages"
)

type PublicController struct {
	Logger *slog.Logger
}

func NewPublicController(logger *slog.Logger) *PublicController {
	return &PublicController{
		Logger: logger,
	}
}

func (c *PublicController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", c.HomePage)
}

func (c PublicController) HomePage(w http.ResponseWriter, r *http.Request) {
	html := pages.HomePage(pages.HomePageProps{
		IsLoggedIn:     false,
		FlashedMessage: lib.Ref("Welcome to our simple home-page"),
	})

	html.Render(r.Context(), w)
}
