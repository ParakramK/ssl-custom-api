package users

import (
	"context"
	"errors"
	"fmt"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/app/query"
	"ssl-custom-api/internal/utils"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(ctx context.Context, input *CreateUserRequest) (*CreateUserResponse, error)
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
}

func NewUserRepository(db *gorm.DB, jwt *auth.JWT) UserRepository {
	return &userRepository{db: db, jwt: jwt}
}

type userRepository struct {
	db  *gorm.DB
	jwt *auth.JWT
}

func (r *userRepository) UpdateUser(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

func (r *userRepository) DeleteUser(user *models.User) error {
	result := r.db.Delete(user)
	return result.Error
}
func (r *userRepository) CreateUser(ctx context.Context,
	input *CreateUserRequest) (*CreateUserResponse, error) {
	q := query.Use(r.db)
	if existing, err := r.CheckUserExists(ctx, &models.User{
		Username: input.Username,
		Email:    input.Email,
	}); err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	} else if existing != nil {
		return nil, fmt.Errorf("user already exists: %w", constants.ErrUserExists)
	}
	user := &models.User{
		ID:       utils.NewV7ID(),
		Username: input.Username,
		Password: input.Password,
		RoleID:   input.RoleId,
		Email:    input.Email,
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return &CreateUserResponse{}, err
	}
	user.Password = hashedPassword
	if err := q.User.WithContext(ctx).Create(user); err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (r *userRepository) CheckUserExists(
	ctx context.Context,
	user *models.User,
) (*models.User, error) {
	q := query.Use(r.db)

	existing, err := q.User.WithContext(ctx).
		Where(
			q.User.Username.Eq(user.Username),
		).
		Or(
			q.User.Email.Eq(user.Email),
		).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return existing, fmt.Errorf("%w: %q", constants.ErrUserExists, user.Username)
}
