package permissions

import (
	"context"

	"ssl-custom-api/internal/app/constants"
	"ssl-custom-api/internal/app/models"

	"github.com/google/uuid"
)

type PermissionService interface {
	GetModulePermittedForUser(ctx context.Context, userID uuid.UUID, isAdmin bool) ([]ModulePermitted, error)
	AddModulePermissions(ctx context.Context, userID uuid.UUID, grants []ModulePermissionGrant) error
}

func NewPermissionService(repo PermissionRepository) PermissionService {
	return &permissionService{repo: repo}
}

type permissionService struct {
	repo PermissionRepository
}

func (s *permissionService) GetModulePermittedForUser(
	ctx context.Context,
	userID uuid.UUID,
	_ bool,
) ([]ModulePermitted, error) {
	return s.repo.GetModulePermittedForUser(ctx, userID)
}
func (s *permissionService) AddModulePermissions(
	ctx context.Context,
	userID uuid.UUID,
	grants []ModulePermissionGrant,
) error {
	for _, grant := range grants {
		if !models.ValidPermission(grant.Permission) {
			return constants.ErrInvalidPermission
		}
	}
	return s.repo.AddModulePermissions(ctx, userID, grants)
}
