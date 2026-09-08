package models

import "github.com/google/uuid"

type Role struct {
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name    string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	IsAdmin bool      `gorm:"type:boolean;not null;default:false"`
}
