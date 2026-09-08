package apikey

import (
	"github.com/google/uuid"
)

type CreateApiKeyRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
type CreateApiKeyInput struct {
	Body CreateApiKeyRequest `json:"body"`
}

type CreateApiKeyResponse struct {
	ID     uuid.UUID `json:"id"`
	Key    string    `json:"key"`
	UserID uuid.UUID `json:"user_id"`
}
type CreateApiKeyOutput struct {
	Body CreateApiKeyResponse `json:"body"`
}

type ApiKeyListResponse struct {
	ApiKeys []ApiKeyRow `json:"api_keys"`
}

type ApiKeyListOutput struct {
	Body ApiKeyListResponse `json:"body"`
}

type ListAllApiKeysByUserIDInput struct {
	UserID uuid.UUID `query:"user_id" doc:"User ID"`
}
type ListAllApiKeysByUserIDOutput struct {
	Body ApiKeyListResponse `json:"body"`
}

type ListAllApiKeysInput struct{}
