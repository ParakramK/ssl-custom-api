package auth

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) LoginUser(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	if input.Body.Username == "" {
		return nil, huma.Error400BadRequest("username is required")
	}
	if input.Body.Password == "" {
		return nil, huma.Error400BadRequest("password is required")
	}

	response, err := h.service.LoginUser(ctx, LoginInput{
		Body: LoginRequest{
			Username: input.Body.Username,
			Password: input.Body.Password,
		},
	})
	if err != nil {
		return nil, huma.Error401Unauthorized("invalid credentials")
	}

	return &LoginOutput{Body: *response}, nil
}
