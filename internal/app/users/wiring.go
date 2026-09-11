package users

import (
	"gorm.io/gorm"
)

func New(db *gorm.DB) *Handler {
	repo := NewUserRepository(db, nil)
	service := NewService(repo)

	return NewHandler(service)
}
