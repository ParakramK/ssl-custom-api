package apikey

import (
	"context"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/app/query"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApiKeyRepository interface {
	CreateApiKey(apiKey *models.ApiKey) error
	GetApiKeyByKey(ctx context.Context, key string) (*models.ApiKey, error)
	DeleteApiKey(apiKey *models.ApiKey) error
	ListApiKeys(ctx context.Context, userID uuid.UUID, limit int, after uuid.UUID) (ApiKeyListResponse, error)
}

func NewApiKeyRepository(db *gorm.DB) ApiKeyRepository {
	return &apiKeyRepository{db: db}
}

type apiKeyRepository struct {
	db *gorm.DB
}

func (r *apiKeyRepository) CreateApiKey(apiKey *models.ApiKey) error {
	result := r.db.Create(apiKey)
	return result.Error
}

func (r *apiKeyRepository) GetApiKeyByKey(ctx context.Context, key string) (*models.ApiKey, error) {
	q := query.Use(r.db)
	apiKey, err := q.ApiKey.WithContext(ctx).
		Select(
			q.ApiKey.ID,
			q.ApiKey.Key,
			q.ApiKey.UserID,
		).
		Where(q.ApiKey.Key.Eq(key)).
		First()
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

func (r *apiKeyRepository) DeleteApiKey(apiKey *models.ApiKey) error {
	result := r.db.Delete(apiKey)
	return result.Error
}

func (r *apiKeyRepository) ListApiKeys(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	after uuid.UUID,
) (ApiKeyListResponse, error) {
	q := query.Use(r.db)
	var apiKeyRows []ApiKeyRow
	stmt := q.ApiKey.WithContext(ctx).
		Select(
			q.ApiKey.ID,
			q.ApiKey.Key,
			q.User.Username.As("user_name"),
		).
		LeftJoin(q.User, q.User.ID.EqCol(q.ApiKey.UserID))
	if userID != uuid.Nil {
		stmt = stmt.Where(q.ApiKey.UserID.Eq(userID))
	}
	if after != uuid.Nil {
		stmt = stmt.Where(q.ApiKey.ID.Gt(after))
	}
	// Fetch one extra row to know whether another page exists.
	if err := stmt.Order(q.ApiKey.ID.Asc()).Limit(limit + 1).Scan(&apiKeyRows); err != nil {
		return ApiKeyListResponse{}, err
	}

	response := ApiKeyListResponse{ApiKeys: apiKeyRows}
	if len(apiKeyRows) > limit {
		response.ApiKeys = apiKeyRows[:limit]
		last := apiKeyRows[limit-1].ID
		response.NextCursor = &last
	}

	return response, nil
}
