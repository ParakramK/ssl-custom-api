package modules

import (
	"context"
	"errors"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service ModuleService
}

func NewHandler(service ModuleService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateModule(
	ctx context.Context,
	input *ModuleCreateInput,
) (*ModuleCreateOutput, error) {
	principal := auth.MustPrincipal(ctx)
	if !principal.IsAdmin {
		return nil, huma.Error403Forbidden("admin access required")
	}

	module, err := h.service.Create(ctx, input.Body)
	if err != nil {
		if errors.Is(err, constants.ErrModuleExists) {
			return nil, huma.Error409Conflict(err.Error())
		}
		return nil, err
	}

	return &ModuleCreateOutput{
		Body: ModuleCreateResponse{
			Name:    module.Name,
			Code:    module.Code,
			Message: module.Message,
		},
	}, nil
}

func (h *Handler) CreateResource(
	ctx context.Context,
	input *ResourceCreateInput,
) (*ResourceCreateOutput, error) {
	principal := auth.MustPrincipal(ctx)
	if !principal.IsAdmin {
		return nil, huma.Error403Forbidden("admin access required")
	}

	resource, err := h.service.CreateResource(ctx, input.Body)
	if err != nil {
		if errors.Is(err, constants.ErrResourceExists) {
			return nil, huma.Error409Conflict(err.Error())
		}
		return nil, err
	}

	return &ResourceCreateOutput{
		Body: ResourceCreateResponse{
			Name:    resource.Name,
			Code:    resource.Code,
			Message: resource.Message,
		},
	}, nil
}
func (h *Handler) GetModuleInfo(
	ctx context.Context,
	input *ModuleInfoRequest,
) (*ModuleInfoOutput, error) {
	module, err := h.service.GetModule(ctx, input)
	if err != nil {
		return nil, err
	}

	return &ModuleInfoOutput{
		Body: ModuleInfoResponse{
			ID:          module.ID,
			Name:        module.Name,
			Code:        module.Code,
			Description: module.Description,
		},
	}, nil
}
