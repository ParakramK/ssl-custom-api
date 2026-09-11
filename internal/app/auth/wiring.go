package auth

import (
	"gorm.io/gorm"
)

func New(db *gorm.DB, jwt *JWT) *Handler {
	repo := NewUserRepository(db, jwt)
	service := NewService(repo)

	return NewHandler(service)
}
