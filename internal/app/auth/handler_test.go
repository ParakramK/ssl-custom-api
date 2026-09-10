package auth

import (
	"context"
)

type stubUserRepository struct {
	createCalls int
	createErr   error
}

func (s *stubUserRepository) LoginUser(_ context.Context, _, _ string) (*LoginResponse, error) {
	return &LoginResponse{}, nil
}

func (s *stubUserRepository) ChangeUserPassword(_, _, _ string) error { return nil }
