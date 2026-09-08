package apikey

import (
	"context"

	"github.com/google/uuid"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service ApiKeyService
}

func NewHandler(service ApiKeyService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateApiKey(ctx context.Context, input *CreateApiKeyInput) (*CreateApiKeyOutput, error) {
	if input.Body.UserID == uuid.Nil {
		return nil, huma.Error400BadRequest("user_id is required")
	}

	response, err := h.service.CreateApiKey(ctx, &CreateApiKeyRequest{
		UserID: input.Body.UserID,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create API key")
	}

	return &CreateApiKeyOutput{Body: *response}, nil
}
func (h *Handler) ListApiKeys(ctx context.Context, input *ListApiKeysInput) (*ApiKeyListOutput, error) {
	response, err := h.service.ListApiKeys(ctx, input.UserID, input.Limit, input.After)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list API keys")
	}

	return &ApiKeyListOutput{Body: response}, nil
}
