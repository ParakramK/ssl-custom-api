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

func (j *JWT) ValidateToken(tokenString string) (uuid.UUID, error) {
	if j == nil {
		return uuid.Nil, fmt.Errorf("jwt signer is not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return uuid.Nil, fmt.Errorf("token is missing subject")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token subject: %w", err)
	}

	return userID, nil
}
