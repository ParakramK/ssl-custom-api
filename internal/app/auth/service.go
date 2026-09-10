package auth

import (
	"context"
)

type Service struct {
	repository UserRepository
}

func NewService(repository UserRepository) *Service {
	return &Service{
		repository: repository,
	}

}
func (S *Service) LoginUser(ctx context.Context, req LoginInput) (*LoginResponse, error) {
	return S.repository.LoginUser(ctx, req.Body.Username, req.Body.Password)
}
