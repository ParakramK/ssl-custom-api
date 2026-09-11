package roles

import (
	"context"
	"errors"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *RoleService
}

func NewHandler(service *RoleService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateRole(ctx context.Context, input *CreateRoleInput) (*CreateRoleOutput, error) {
	principal := auth.MustPrincipal(ctx)
	if !principal.IsAdmin {
		return nil, huma.Error403Forbidden("admin access required")
	}

	role, err := h.service.CreateRole(ctx, *input)
	if err != nil {
		if errors.Is(err, constants.ErrRoleExists) {
			return nil, huma.Error409Conflict(err.Error())
		}
		return nil, err
	}

	return &CreateRoleOutput{
		Body: *role,
	}, nil
}
func (h *Handler) ListRoles(ctx context.Context, input *ListRolesRequest) (*RolesOutput, error) {
	principal := auth.MustPrincipal(ctx)
	if !principal.IsAdmin {
		return nil, huma.Error403Forbidden("admin access required")
	}

	roles, err := h.service.ListRoles(ctx, input)
	if err != nil {
		return nil, err
	}

	return &RolesOutput{
		Body: *roles,
	}, nil
}
