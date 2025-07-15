package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sandbox/db/models"
	"sandbox/internal/config"
	"sandbox/internal/form"
	"sandbox/internal/lib/jwt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	TokenCookieName string
	Logger          *slog.Logger
	DB              *models.Queries
	Config          *config.Config
}

type LoginResult struct {
	User   *models.User
	Token  string
	Expiry int64
}

func (s *AuthService) Login(ctx context.Context, form *form.LoginForm) (*LoginResult, error) {
	user, err := s.DB.GetUserByEmail(ctx, form.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, expiry, err := jwt.NewExpiringToken(&jwt.ExpiringTokenArgs{
		UserId:        user.ID.String(),
		Email:         user.Email,
		Role:          string(user.Role),
		JwtSecret:     s.Config.Auth.JwtSecret,
		ExpiryMinutes: s.Config.Auth.JwtExpiryMinutes,
	})

	result := &LoginResult{
		User:   &user,
		Token:  token,
		Expiry: expiry,
	}

	return result, nil
}

func (s AuthService) SetAuthCookies(c echo.Context, result *LoginResult) error {
	expiry := time.Now().Add(s.Config.Auth.JwtExpiryMinutes)

	tokenCookie := &http.Cookie{
		Name:     s.TokenCookieName,
		Value:    result.Token,
		Expires:  expiry,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(tokenCookie)

	// TODO: enable.
	// userJson, err := json.Marshal(result.User)
	// if err != nil {
	// 	s.Logger.Error("failed to json encode login result for auth user cookie", "error", err.Error())
	// 	return errors.New("failed to set auth cookies")
	// }

	// encodedUser := base64.StdEncoding.EncodeToString(userJson)
	// userCookie := &http.Cookie{
	// 	Name:    authUserCookieName,
	// 	Value:   encodedUser,
	// 	Expires: expiry,
	// }
	// c.SetCookie(userCookie)

	return nil
}

func (s AuthService) RemoveAuthCookies(c echo.Context) {
	expiry := time.Unix(0, 0)
	c.SetCookie(&http.Cookie{
		Name:     s.TokenCookieName,
		Expires:  expiry,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// c.SetCookie(&http.Cookie{Name: authUserCookieName, Expires: expiry}) // TODO: enable.
}

func (s AuthService) CreateAccount(ctx context.Context, form *form.RegisterForm) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("failed to hash user password", "error", err.Error())
		return errors.New("failed to register user")
	}

	err = s.DB.CreateUser(ctx, models.CreateUserParams{
		ID:       uuid.New(),
		Email:    form.Email,
		Role:     models.UserRoleADMIN,
		Password: string(passwordHash),
		Name:     pgtype.Text{Valid: false},
	})

	if err != nil {
		// TODO: handle email_unique constraint.
		s.Logger.Error("failed to save user to database", "error", err.Error())
		return errors.New("failed to register user")
	}

	return nil
}
