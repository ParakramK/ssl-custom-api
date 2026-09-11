package permissions

import (
	"context"
	"fmt"

	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"
	"ssl-custom-api/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ssl-custom-api/internal/app/query"
)

type PermissionRepository interface {
	GetModulePermittedForUser(ctx context.Context, userID uuid.UUID) ([]ModulePermitted, error)
	AddModulePermissions(ctx context.Context, userID uuid.UUID, grants []ModulePermissionGrant) error
}

func NewPermissionRepository(db *gorm.DB, jwtSvc *auth.JWT) PermissionRepository {
	return &permissionRepository{db: db, jwtSvc: jwtSvc}
}

type permissionRepository struct {
	db     *gorm.DB
	jwtSvc *auth.JWT
}

func (r *permissionRepository) GetModulePermittedForUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]ModulePermitted, error) {
	q := query.Use(r.db)

	var role struct {
		IsAdmin bool
	}

	err := q.User.WithContext(ctx).
		Join(
			q.Role,
			q.Role.ID.EqCol(q.User.RoleID),
		).
		Select(
			q.Role.IsAdmin.As("IsAdmin"),
		).
		Where(q.User.ID.Eq(userID)).
		Scan(&role)

	if err != nil {
		return nil, err
	}

	moduleQuery := q.Module.WithContext(ctx).
		Select(
			q.Module.Name.As("ModuleName"),
			q.Module.Code.As("ModuleCode"),
			q.Module.ID.As("ModuleID"),
			q.UserModulePermission.Permission.As("Permission"),
		).
		Distinct()

	if !role.IsAdmin {
		moduleQuery = moduleQuery.
			Join(
				q.UserModulePermission,
				q.UserModulePermission.ModuleID.EqCol(q.Module.ID),
			).
			Where(q.UserModulePermission.UserID.Eq(userID))
	}

	var modules []ModulePermitted
	if err := moduleQuery.Scan(&modules); err != nil {
		return nil, err
	}

	return modules, nil
}
func (r *permissionRepository) AddModulePermissions(
	ctx context.Context,
	userID uuid.UUID,
	grants []ModulePermissionGrant,
) error {
	q := query.Use(r.db)

	// Validate module IDs upfront so unknown IDs yield a 400
	// instead of a raw foreign-key violation (23503) from Postgres.
	ids := make([]uuid.UUID, 0, len(grants))
	requested := make(map[uuid.UUID]struct{}, len(grants))
	for _, grant := range grants {
		ids = append(ids, grant.ModuleID)
		requested[grant.ModuleID] = struct{}{}
	}
	if len(requested) > 0 {
		var found []uuid.UUID
		if err := r.db.WithContext(ctx).
			Model(&models.Module{}).
			Where("id IN ?", ids).
			Pluck("id", &found).Error; err != nil {
			return err
		}
		if len(found) != len(requested) {
			known := make(map[uuid.UUID]struct{}, len(found))
			for _, id := range found {
				known[id] = struct{}{}
			}
			for id := range requested {
				if _, ok := known[id]; !ok {
					return fmt.Errorf("%w: %q", constants.ErrUnknownModule, id.String())
				}
			}
		}
	}

	existing, err := q.UserModulePermission.
		WithContext(ctx).
		Where(q.UserModulePermission.UserID.Eq(userID)).
		Find()
	if err != nil {
		return err
	}

	existingSet := make(map[uuid.UUID]map[models.Permission]struct{}, len(existing))
	for _, grant := range existing {
		perms := existingSet[grant.ModuleID]
		if perms == nil {
			perms = make(map[models.Permission]struct{})
			existingSet[grant.ModuleID] = perms
		}
		perms[grant.Permission] = struct{}{}
	}

	for _, grant := range grants {
		if perms := existingSet[grant.ModuleID]; len(perms) == 1 {
			if _, ok := perms[grant.Permission]; ok {
				continue
			}
		}
		if err := r.db.WithContext(ctx).
			Where("user_id = ? AND module_id = ?", userID, grant.ModuleID).
			Delete(&models.UserModulePermission{}).Error; err != nil {
			return err
		}

		if err := q.UserModulePermission.
			WithContext(ctx).
			Create(&models.UserModulePermission{
				UserID:     userID,
				ModuleID:   grant.ModuleID,
				Permission: grant.Permission,
			}); err != nil {
			return err
		}
		existingSet[grant.ModuleID] = map[models.Permission]struct{}{
			grant.Permission: {},
		}
	}

	return nil
}
