package modules

import (
	"gorm.io/gorm"
)

func New(db *gorm.DB) *Handler {
	repo := NewModuleRepository(db)
	service := NewModuleService(repo)

	return NewHandler(service)
}
