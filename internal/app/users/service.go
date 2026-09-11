package users

import (
	"context"
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
