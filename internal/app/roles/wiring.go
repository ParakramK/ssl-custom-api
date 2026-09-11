package roles

import (
	"ssl-custom-api/internal/app/auth"

	"gorm.io/gorm"
)

func New(db *gorm.DB, jwt *auth.JWT) *Handler {
	repo := NewRoleRepository(db, jwt)
	service := NewService(repo)

	return NewHandler(service)
}
