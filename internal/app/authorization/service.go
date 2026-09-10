package authorization

import (
	"context"
	"ssl-custom-api/internal/app/auth"

	"github.com/google/uuid"
)

type Service interface {
	GetPrincipal(ctx context.Context, userID uuid.UUID) (auth.Principal, error)
}

type service struct {
	repository AuthorizationRepository
}

func NewService(repository AuthorizationRepository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetPrincipal(
	ctx context.Context,
	userID uuid.UUID,
) (auth.Principal, error) {
	return s.repository.GetPrincipal(ctx, userID)
}
