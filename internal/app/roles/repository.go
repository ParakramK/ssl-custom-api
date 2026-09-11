package roles

import (
	"context"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/app/query"
	"ssl-custom-api/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	ListRoles(ctx context.Context) (*RolesResponse, error)
	GetRoleById(ctx context.Context, id uuid.UUID) (*RolesResponse, error)
	GetRoleByName(ctx context.Context, name string) (*RolesResponse, error)
	CreateRole(ctx context.Context, input *CreateRoleRequest) (*CreateRoleResponse, error)
}

type roleRepository struct {
	db  *gorm.DB
	jwt *auth.JWT
}

func NewRoleRepository(db *gorm.DB, jwt *auth.JWT) RoleRepository {
	return &roleRepository{db: db, jwt: jwt}
}
func (r *roleRepository) ListRoles(ctx context.Context) (*RolesResponse, error) {
	q := query.Use(r.db)
	var roles []RoleRow
	err := q.Role.WithContext(ctx).
		Select(q.Role.ID.As("ID"), q.Role.Name.As("Name"), q.Role.IsAdmin.As("IsAdmin")).
		Scan(&roles)
	if err != nil {
		return nil, err
	}
	return &RolesResponse{
		Roles: roles,
	}, nil

}
func (r *roleRepository) GetRoleById(ctx context.Context, id uuid.UUID) (*RolesResponse, error) {
	q := query.Use(r.db)
	var role RoleRow
	err := q.Role.WithContext(ctx).
		Where(q.Role.ID.Eq(id)).
		Scan(&role)
	if err != nil {
		return nil, err
	}
	return &RolesResponse{
		Roles: []RoleRow{role},
	}, nil
}
func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*RolesResponse, error) {
	q := query.Use(r.db)
	var role RoleRow
	err := q.Role.WithContext(ctx).
		Where(q.Role.Name.Eq(name)).
		Scan(&role)
	if err != nil {
		return nil, err
	}
	return &RolesResponse{
		Roles: []RoleRow{role},
	}, nil
}
func (r *roleRepository) CheckRoleExists(ctx context.Context, name string) (bool, error) {
	q := query.Use(r.db)
	count, err := q.Role.WithContext(ctx).
		Where(q.Role.Name.Eq(name)).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *roleRepository) CreateRole(ctx context.Context, input *CreateRoleRequest) (*CreateRoleResponse, error) {

	q := query.Use(r.db)
	if exists, err := r.CheckRoleExists(ctx, input.Name); err != nil {
		return nil, err
	} else if exists {
		return nil, constants.ErrRoleExists
	}

	role := &models.Role{
		ID:      utils.NewV7ID(),
		Name:    input.Name,
		IsAdmin: input.IsAdmin,
	}

	if err := q.Role.WithContext(ctx).Create(role); err != nil {
		return nil, err
	}

	return &CreateRoleResponse{
		ID:   role.ID,
		Name: role.Name,
	}, nil
}
