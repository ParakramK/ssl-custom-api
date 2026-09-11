package users

import (
	"context"
	"errors"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *UserService
}

func NewHandler(service *UserService) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	principal := auth.MustPrincipal(ctx)
	if !principal.IsAdmin {
		return nil, huma.Error403Forbidden("admin access required")
	}

	response, err := h.service.CreateUser(ctx, *input)
	if err != nil {
		if errors.Is(err, constants.ErrUserExists) {
			return nil, huma.Error409Conflict(err.Error())
		}
		return nil, err
	}

	return &CreateUserOutput{Body: *response}, nil
}

func (h *Handler) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersOutput, error) {
	principal := auth.MustPrincipal(ctx)
	if !principal.IsAdmin {
		return nil, huma.Error403Forbidden("admin access required")
	}

	response, err := h.service.ListUsers(ctx, *req)
	if err != nil {
		return nil, err
	}

	return &ListUsersOutput{Body: *response}, nil
}
