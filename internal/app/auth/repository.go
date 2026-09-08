package auth

import (
	"context"
	"fmt"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/app/query"
	"ssl-custom-api/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
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

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = utils.NewV7ID()
	}
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	result := r.db.Create(user)
	return result.Error
}

func (r *userRepository) UpdateUser(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

func (r *userRepository) DeleteUser(user *models.User) error {
	result := r.db.Delete(user)
	return result.Error
}
func (r *userRepository) LoginUser(ctx context.Context, username, password string) (*LoginResponse, error) {
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

func (r *userRepository) ChangeUserPassword(username string, currentPassword string, newPassword string) error {
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
