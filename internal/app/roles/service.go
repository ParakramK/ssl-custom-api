package roles

import (
	"context"

	"github.com/google/uuid"
)

type RoleService struct {
	repository RoleRepository
}

func NewService(repository RoleRepository) *RoleService {
	return &RoleService{
		repository: repository,
	}
}

func (s *RoleService) ListRoles(
	ctx context.Context,
	input *ListRolesRequest,
) (*RolesResponse, error) {

	if input != nil && input.Id != uuid.Nil {
		role, err := s.repository.GetRoleById(ctx, input.Id)
		if err != nil {
			return nil, err
		}
		return role, nil
	}

	if input != nil && input.Name != "" {
		role, err := s.repository.GetRoleByName(ctx, input.Name)
		if err != nil {
			return nil, err
		}

		return role, nil
	}

	return s.repository.ListRoles(ctx)
}

func (S *RoleService) CreateRole(ctx context.Context, req CreateRoleInput) (*CreateRoleResponse, error) {
	return S.repository.CreateRole(ctx, &req.Body)
}
