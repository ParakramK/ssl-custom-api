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
	// Cursor for the next page (ID of the last key returned).
	// Omitted when there are no more results.
	NextCursor *uuid.UUID `json:"next_cursor,omitempty"`
}

type ApiKeyListOutput struct {
	Body ApiKeyListResponse `json:"body"`
}

type ListApiKeysInput struct {
	// Optional filter; omit (or Nil) to list all keys.
	UserID uuid.UUID `query:"user_id" doc:"Filter by user ID (omit to list all)"`
	// Keyset cursor: ID of the last key from the previous page.
	// Omit to start from the beginning.
	After uuid.UUID `query:"after" doc:"Return keys with ID greater than this cursor"`
	Limit int       `query:"limit" default:"10" minimum:"1" maximum:"100" doc:"Max keys per page"`
}
