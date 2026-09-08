package models

import (
	"time"

	"github.com/google/uuid"
)

type Module struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Code        string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Resource struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ModuleID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Code        string    `gorm:"type:varchar(100);not null"`
	Description *string   `gorm:"type:text"`

	Module Module `gorm:"foreignKey:ModuleID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}
