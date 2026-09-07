package apikey

import (
	"context"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/utils"
)

type ApiKeyService interface {
	CreateApiKey(ctx context.Context, req *CreateApiKeyRequest) (*CreateApiKeyResponse, error)
	GetApiKeyByKey(ctx context.Context, key string) (*models.ApiKey, error)
	DeleteApiKey(ctx context.Context, apiKey *models.ApiKey) error
}

type apiKeyService struct {
	repo ApiKeyRepository
}

func NewApiKeyService(repo ApiKeyRepository) ApiKeyService {
	return &apiKeyService{repo: repo}
}

func (s *apiKeyService) CreateApiKey(ctx context.Context, req *CreateApiKeyRequest) (*CreateApiKeyResponse, error) {
	key, err := utils.GenerateKey()
	if err != nil {
		return nil, err
	}

	apiKey := &models.ApiKey{
		Key:    key,
		UserID: req.UserID,
	}

	err = s.repo.CreateApiKey(apiKey)
	if err != nil {
		return nil, err
	}

	return &CreateApiKeyResponse{
		ID:     apiKey.ID,
		Key:    apiKey.Key,
		UserID: apiKey.UserID,
	}, nil
}

func (s *apiKeyService) GetApiKeyByKey(ctx context.Context, key string) (*models.ApiKey, error) {
	return s.repo.GetApiKeyByKey(ctx, key)
}

func (s *apiKeyService) DeleteApiKey(ctx context.Context, apiKey *models.ApiKey) error {
	return s.repo.DeleteApiKey(apiKey)
}
