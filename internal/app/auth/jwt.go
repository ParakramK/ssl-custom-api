package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT struct {
	secret []byte
}

func NewJWT(secret string) (*JWT, error) {
	if secret == "" {
		return nil, fmt.Errorf("JWT secret is empty")
	}

	return &JWT{
		secret: []byte(secret),
	}, nil
}

func (j *JWT) GenerateToken(
	userID uuid.UUID,
	validFor time.Duration,
) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(validFor).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(j.secret)
}
