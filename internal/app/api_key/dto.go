package apikey

import (
	"github.com/google/uuid"
)

type CreateApiKeyRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type CreateApiKeyResponse struct {
	ID     uuid.UUID `json:"id"`
	Key    string    `json:"key"`
	UserID uuid.UUID `json:"user_id"`
}
