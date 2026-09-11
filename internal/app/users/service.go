package users

import (
	"context"

	"github.com/google/uuid"
)

type UserService struct {
	repository UserRepository
}

func NewService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}

}

// CreateUser creates a user. The caller (handler) enforces admin-only
// access from the request Principal; this service performs no
// authorization lookup.
func (S *UserService) CreateUser(ctx context.Context, req CreateUserInput) (*CreateUserResponse, error) {
	return S.repository.CreateUser(ctx, &req.Body)
}

func (S *UserService) ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	var roleID *uuid.UUID
	if req.RoleID != uuid.Nil {
		roleID = &req.RoleID
	}
	return S.repository.ListUsers(ctx, roleID)
}
