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
