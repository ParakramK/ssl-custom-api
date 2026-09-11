package apikey

import (
	"ssl-custom-api/internal/app/permissions"

	"gorm.io/gorm"
)

func New(db *gorm.DB, checker permissions.Checker) *Handler {
	repo := NewApiKeyRepository(db)
	service := NewApiKeyService(repo, checker)

	return NewHandler(service)
}
