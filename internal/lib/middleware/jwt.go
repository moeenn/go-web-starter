package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	jwtClaimsContextKey string = "jwtClaims"
)

type JwtClaims struct {
	UserId string
	Email  string
	Role   string
}

func jwtClaimsFromMap(claims map[string]any) (*JwtClaims, error) {
	userId, userIdOk := claims["userId"].(string)
	if !userIdOk {
		return nil, errors.New("invalid or missing userId")
	}

	email, emailOk := claims["email"].(string)
	if !emailOk {
		return nil, errors.New("invalid or missing email")
	}

	role, roleOk := claims["role"].(string)
	if !roleOk {
		return nil, errors.New("invalid or missing role")
	}

	return &JwtClaims{
		UserId: userId,
		Email:  email,
		Role:   role,
	}, nil
}

func validateAndParseJwtClaims(jwtSecret []byte, bearerToken string) (*JwtClaims, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt claims")
	}

	// Check token expiration
	exp, err := claims.GetExpirationTime()
	if err != nil || exp.Time.Before(time.Now()) {
		return nil, errors.New("expired jwt")
	}

	// Access claims
	parsedClaims, err := jwtClaimsFromMap(claims)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt claims: %w", err)
	}

	return parsedClaims, nil
}

func GetJwtClaims(c echo.Context) *JwtClaims {
	claims := c.Get(jwtClaimsContextKey)
	claimsCast, ok := claims.(*JwtClaims)
	if !ok {
		return nil
	}
	return claimsCast
}
