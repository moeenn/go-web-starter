package controller

import (
	"log/slog"
	"net/http"
	"sandbox/internal/lib"
	"sandbox/internal/lib/jwt"
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
	_, ok := jwt.GetJwtClaims(r)
	html := pages.HomePage(pages.HomePageProps{
		IsLoggedIn:     ok,
		FlashedMessage: lib.Ref("Welcome to our simple home-page"),
	})

	html.Render(r.Context(), w)
}
