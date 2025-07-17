package controller

import (
	"app/internal/lib/jwt"
	"app/views/components"
	"app/views/pages"
	"log/slog"
	"net/http"
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
		IsLoggedIn: ok,
		FlashedMessage: &components.Message{
			Message: "Welcome to our simple home-page",
			Type:    components.MessageTypeInfo,
		},
	})

	html.Render(r.Context(), w)
}
