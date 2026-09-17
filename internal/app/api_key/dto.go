package apikey

import (
	"ssl-custom-api/internal/app/models"

	"github.com/google/uuid"
)

type ApiKeyScope struct {
	Resource   models.ResourceCode `json:"resource"`
	Permission models.Permission   `json:"permission"`
}

type CreateApiKeyRequest struct {
	UserID uuid.UUID     `json:"user_id" binding:"required"`
	Scopes []ApiKeyScope `json:"scopes"`
}
type CreateApiKeyInput struct {
	Body CreateApiKeyRequest `json:"body"`
}

type CreateApiKeyResponse struct {
	ID     uuid.UUID     `json:"id"`
	Key    string        `json:"key"`
	UserID uuid.UUID     `json:"user_id"`
	Scopes []ApiKeyScope `json:"scopes"`
}
type CreateApiKeyOutput struct {
	Body CreateApiKeyResponse `json:"body"`
}

type ApiKeyListResponse struct {
	ApiKeys    []ApiKeyRow `json:"api_keys"`
	HasMore    bool        `json:"has_more"`
	NextCursor string      `json:"next_cursor,omitempty"`
}

type ApiKeyListOutput struct {
	Body ApiKeyListResponse `json:"body"`
}

type ListApiKeysInput struct {
	UserID uuid.UUID `query:"user_id" doc:"Filter by user ID (omit to list all)"`
}

type DeleteApiKeyInput struct {
	ID string `path:"id" doc:"API key ID"`
}

type DeleteApiKeyResponse struct {
	Message string `json:"message" doc:"Response message" example:"api key revoked"`
}

type DeleteApiKeyOutput struct {
	Body DeleteApiKeyResponse `json:"body"`
}
