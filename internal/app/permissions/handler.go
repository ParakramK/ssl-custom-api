package permissions

import (
	"context"
	"errors"

	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service PermissionService
}

func NewHandler(service PermissionService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetModulePermittedForUser(
	ctx context.Context,
	_ *ModulePermittedRequestCurrentUser,
) (*ModulePermittedOutput, error) {
	principal := auth.MustPrincipal(ctx)

	modules, err := h.service.GetModulePermittedForUser(ctx, principal.UserID, principal.IsAdmin)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list permitted modules")
	}

	return &ModulePermittedOutput{
		Body: ModulePermittedResponse{
			Modules: modules,
		},
	}, nil
}
func (h *Handler) GetModulePermittedForUserById(
	ctx context.Context,
	input *ModulePermittedRequest,
) (*ModulePermittedOutput, error) {
	principal := auth.MustPrincipal(ctx)

	modules, err := h.service.GetModulePermittedForUser(ctx, input.UserID, principal.IsAdmin)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list permitted modules")
	}

	return &ModulePermittedOutput{
		Body: ModulePermittedResponse{
			Modules: modules,
		},
	}, nil
}

func (h *Handler) AddModulePermissions(
	ctx context.Context,
	input *CreateModulePermissionInput,
) (*CreateModulePermissionOutput, error) {
	principal := auth.MustPrincipal(ctx)
	if !principal.IsAdmin {
		return nil, huma.Error403Forbidden("admin access required")
	}

	err := h.service.AddModulePermissions(ctx, input.Body.UserID, input.Body.Modules)
	if err != nil {
		if errors.Is(err, constants.ErrUnknownModule) || errors.Is(err, constants.ErrInvalidPermission) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, huma.Error500InternalServerError("failed to add module permission for user")
	}

	return &CreateModulePermissionOutput{
		Body: CreateModulePermissionResponse{
			Message: "module permission added successfully",
		},
	}, nil
}
