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
	ListAllApiKeysByUserID(ctx context.Context, id uuid.UUID) (ApiKeyListResponse, error)
	ListAllApiKeys(ctx context.Context) (ApiKeyListResponse, error)
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

func (r *apiKeyRepository) ListAllApiKeysByUserID(ctx context.Context, id uuid.UUID) (ApiKeyListResponse, error) {
	var apiKeyRows []ApiKeyRow
	q := query.Use(r.db)
	err := q.ApiKey.WithContext(ctx).
		Select(
			q.ApiKey.ID,
			q.ApiKey.Key,
		).
		Where(q.ApiKey.UserID.Eq(id)).
		Scan(&apiKeyRows)
	if err != nil {
		return ApiKeyListResponse{}, err
	}
	return ApiKeyListResponse{ApiKeys: apiKeyRows}, nil
}

func (r *apiKeyRepository) ListAllApiKeys(ctx context.Context) (ApiKeyListResponse, error) {
	q := query.Use(r.db)
	var apiKeyRows []ApiKeyRow
	err := q.ApiKey.WithContext(ctx).
		Select(
			q.ApiKey.ID,
			q.ApiKey.Key,
			q.User.Username.As("user_name"),
		).
		LeftJoin(q.User, q.User.ID.EqCol(q.ApiKey.UserID)).
		Scan(&apiKeyRows)
	if err != nil {
		return ApiKeyListResponse{}, err
	}

	return ApiKeyListResponse{ApiKeys: apiKeyRows}, nil
}
