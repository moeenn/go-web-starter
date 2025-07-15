package controller

import (
	"log/slog"
	"net/http"
	"sandbox/views/components"
	"sandbox/views/pages"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	Logger *slog.Logger
}

func NewAuthController(logger *slog.Logger) *AuthController {
	return &AuthController{
		Logger: logger,
	}
}

func (c *AuthController) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/auth")
	g.GET("/login", c.LoginPage)
	g.POST("/login", c.ProcessLoginRequest)
	g.GET("/forgot-password", c.ForgotPasswordPage)
	g.POST("/forgot-password", c.ProcessForgotPasswordRequest)
}

func (c AuthController) LoginPage(ctx echo.Context) error {
	html := pages.LoginPage()
	return html.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func (c AuthController) ProcessLoginRequest(ctx echo.Context) error {
	form := LoginFormFromContext(ctx)
	if err := form.Validate(); err != nil {
		html := components.LoginForm(components.LoginFormProps{
			Errors: err,
			Values: components.LoginFormFields{
				Email:    form.Email,
				Password: form.Password,
			},
		})
		return html.Render(ctx.Request().Context(), ctx.Response().Writer)
	}

	ctx.Response().Header().Set("HX-Redirect", "/") // TODO: redirect to dashboard.
	return ctx.NoContent(http.StatusNoContent)
}

func (c AuthController) ForgotPasswordPage(ctx echo.Context) error {
	html := pages.ForgotPasswordPage(pages.ForgotPasswordPageProps{})
	return html.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func (c AuthController) ProcessForgotPasswordRequest(ctx echo.Context) error {
	form := ForgotPasswordFormFromContext(ctx)
	if err := form.Validate(); err != nil {
		html := components.ForgotPasswordForm(components.ForgotPasswordFormProps{
			Errors: err,
		})
		return html.Render(ctx.Request().Context(), ctx.Response().Writer)
	}

	message := "You request has been submitted. You will receive an email shortly with instructions to reset your password"
	html := components.ForgotPasswordForm(components.ForgotPasswordFormProps{
		Message: &message,
	})
	return html.Render(ctx.Request().Context(), ctx.Response().Writer)
}
