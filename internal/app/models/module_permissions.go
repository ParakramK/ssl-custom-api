package models

import "github.com/google/uuid"

type ModulePermission struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ModuleId    uuid.UUID `gorm:"type:uuid;not null"`
	ApiKeyId    uuid.UUID `gorm:"type:uuid;not null"`
	AllowAccess bool      `gorm:"type:boolean;not null"`
}
