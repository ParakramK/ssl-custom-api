package auth

import (
	"context"
	"errors"
	"fmt"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/app/query"
	"time"

	"gorm.io/gorm"
)

var ErrUserExists = errors.New("user with the same username or email already exists")

type UserRepository interface {
	LoginUser(ctx context.Context, username, password string) (*LoginResponse, error)
	ChangeUserPassword(username string, currentPassword string, newPassword string) error
}

func NewUserRepository(db *gorm.DB, jwt *JWT) UserRepository {
	return &userRepository{db: db, jwt: jwt}
}

type userRepository struct {
	db  *gorm.DB
	jwt *JWT
}

func (r *userRepository) LoginUser(ctx context.Context,
	username, password string) (*LoginResponse, error) {
	q := query.Use(r.db)
	user, err := q.User.WithContext(ctx).
		Where(q.User.Username.Eq(username)).
		First()
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !VerifyPassword(password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}
	if r.jwt == nil {
		return nil, fmt.Errorf("jwt signer is not configured")
	}
	accessToken, err := r.jwt.GenerateToken(user.ID, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := r.jwt.GenerateToken(user.ID, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (r *userRepository) ChangeUserPassword(username string,
	currentPassword string, newPassword string) error {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return result.Error
	}

	if !VerifyPassword(currentPassword, user.Password) {
		return fmt.Errorf("invalid current password")
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	result = r.db.Save(&user)
	return result.Error
}
