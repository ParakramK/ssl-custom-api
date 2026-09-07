package apikey

import (
	"context"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/app/query"

	"gorm.io/gorm"
)

type ApiKeyRepository interface {
	CreateApiKey(apiKey *models.ApiKey) error
	GetApiKeyByKey(ctx context.Context, key string) (*models.ApiKey, error)
	DeleteApiKey(apiKey *models.ApiKey) error
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
