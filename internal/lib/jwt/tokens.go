package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ExpiringTokenArgs struct {
	Claims        JwtClaims
	JwtSecret     []byte
	ExpiryMinutes time.Duration
}

func NewExpiringToken(args *ExpiringTokenArgs) (string, int64, error) {
	expiry := time.Now().Add(args.ExpiryMinutes).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":    expiry,
		"iat":    time.Now().Unix(),
		"userId": args.Claims.UserId,
		"email":  args.Claims.Email,
		"role":   args.Claims.Role,
	})

	tokenString, err := token.SignedString(args.JwtSecret)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign jwt: %w", err)
	}

	return tokenString, expiry, nil
}
