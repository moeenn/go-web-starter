package controller

import (
	"log/slog"
	"net/http"
	"sandbox/internal/form"
	"sandbox/internal/lib"
	"sandbox/internal/lib/htmx"
	customMiddleware "sandbox/internal/lib/middleware"
	"sandbox/internal/service"
	"sandbox/views/components"
	"sandbox/views/pages"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	Logger         *slog.Logger
	AuthService    *service.AuthService
	AuthMiddleware *customMiddleware.AuthMiddleware
}

func (c *AuthController) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/auth")
	g.GET("/login", c.AuthMiddleware.NotLoggedIn(c.LoginPage))
	g.POST("/login", c.AuthMiddleware.NotLoggedIn(c.ProcessLoginRequest))
	g.GET("/register", c.AuthMiddleware.NotLoggedIn(c.RegisterPage))
	g.POST("/register", c.AuthMiddleware.NotLoggedIn(c.ProcessRegisterRequest))
	g.GET("/forgot-password", c.AuthMiddleware.NotLoggedIn(c.ForgotPasswordPage))
	g.POST("/forgot-password", c.AuthMiddleware.NotLoggedIn(c.ProcessForgotPasswordRequest))
	g.GET("/logout", c.AuthMiddleware.LoggedIn(c.ProcessLogoutRequest))
}

func (c AuthController) LoginPage(ctx echo.Context) error {
	html := pages.LoginPage(&pages.LoginPageProps{})
	return render(ctx, html)
}

func (c AuthController) ProcessLoginRequest(ctx echo.Context) error {
	form := form.LoginFormFromContext(ctx)
	if err := form.Validate(); err != nil {
		html := components.LoginForm(components.LoginFormProps{
			Errors: err,
			Values: components.LoginFormFields{
				Email:    form.Email,
				Password: form.Password,
			},
		})
		return render(ctx, html)
	}

	loginResult, err := c.AuthService.Login(ctx.Request().Context(), &form)
	if err != nil {
		c.Logger.Error("failed to log in", "error", err.Error())
		html := components.LoginForm(components.LoginFormProps{
			Message: lib.Ref("Invalid email or password"),
			Values: components.LoginFormFields{
				Email: form.Email,
			},
		})
		return render(ctx, html)
	}

	if err := c.AuthService.SetAuthCookies(ctx, loginResult); err != nil {
		c.Logger.Error("failed to set auth cookies", "error", err.Error())
		html := components.LoginForm(components.LoginFormProps{
			Message: lib.Ref("Something went wrong. Please try again."),
		})
		return render(ctx, html)
	}

	return htmx.Redirect(ctx, "/dashboard")
}

func (c AuthController) RegisterPage(ctx echo.Context) error {
	html := pages.RegisterPage()
	return render(ctx, html)
}

func (c AuthController) ProcessRegisterRequest(ctx echo.Context) error {
	form := form.RegisterFormFromContext(ctx)
	if err := form.Validate(); err != nil {
		html := components.RegisterForm(components.RegisterFormProps{
			Errors: err,
			Values: components.RegisterFormFields{
				Email: form.Email,
			},
		})
		return render(ctx, html)
	}

	if err := c.AuthService.CreateAccount(ctx.Request().Context(), &form); err != nil {
		html := components.RegisterForm(components.RegisterFormProps{
			Message: lib.Ref(err.Error()),
		})
		return render(ctx, html)
	}

	ctx.Response().Header().Set("HX-Location", "/auth/login")
	return ctx.NoContent(http.StatusNoContent)
}

func (c AuthController) ForgotPasswordPage(ctx echo.Context) error {
	html := pages.ForgotPasswordPage(pages.ForgotPasswordPageProps{})
	return render(ctx, html)
}

func (c AuthController) ProcessForgotPasswordRequest(ctx echo.Context) error {
	form := form.ForgotPasswordFormFromContext(ctx)
	if err := form.Validate(); err != nil {
		html := components.ForgotPasswordForm(components.ForgotPasswordFormProps{
			Errors: err,
		})
		return render(ctx, html)
	}

	message := "You request has been submitted. You will receive an email shortly with instructions to reset your password"
	html := components.ForgotPasswordForm(components.ForgotPasswordFormProps{
		Message: &message,
	})
	return render(ctx, html)
}

func (c AuthController) ProcessLogoutRequest(ctx echo.Context) error {
	c.AuthService.RemoveAuthCookies(ctx)
	return htmx.Redirect(ctx, "/")
}
