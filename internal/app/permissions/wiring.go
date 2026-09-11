package permissions

import (
	"ssl-custom-api/internal/app/auth"

	"gorm.io/gorm"
)

func New(db *gorm.DB, jwtSvc *auth.JWT) *Handler {
	repo := NewPermissionRepository(db, jwtSvc)
	service := NewPermissionService(repo)

	return NewHandler(service)
}
