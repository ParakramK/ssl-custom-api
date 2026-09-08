package apikey

import (
	"context"

	"ssl-custom-api/internal/app/paging"
	"ssl-custom-api/internal/httpx"

	"github.com/google/uuid"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v3/middleware/paginate"
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
	pi, err := httpx.Pagination(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("pagination is not configured")
	}

	pg := paging.Query{Limit: pi.Limit}
	for _, s := range pi.Sort {
		pg.Sorts = append(pg.Sorts, paging.Sort{
			Field: s.Field,
			Desc:  s.Order == paginate.DESC,
		})
	}
	if vals := pi.CursorValues(); vals != nil {
		idStr, _ := vals["id"].(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid cursor")
		}
		key, _ := vals["key"].(string)
		pg.Cursor = &paging.Cursor{ID: id, Key: key}
	}

	rows, hasMore, err := h.service.ListApiKeys(ctx, input.UserID, pg)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list API keys")
	}

	response := ApiKeyListResponse{ApiKeys: rows, HasMore: hasMore}
	if hasMore {
		last := rows[len(rows)-1]
		if err := pi.SetNextCursor(map[string]any{
			"id":  last.ID.String(),
			"key": last.Key,
		}); err != nil {
			return nil, huma.Error500InternalServerError("failed to encode cursor")
		}
		response.NextCursor = pi.NextCursor
	}

	return &ApiKeyListOutput{Body: response}, nil
}
