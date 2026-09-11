package models

import (
	"github.com/google/uuid"
)

// UserModulePermission grants a user a permission level on a whole
// Control Plane module (e.g. gatepass means all of Gatepass).
// One row per (user, module): re-granting replaces the level.
type UserModulePermission struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_user_module,priority:1"`
	ModuleID   uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_user_module,priority:2"`
	Permission Permission `gorm:"type:varchar(20);not null"`

	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Module Module `gorm:"foreignKey:ModuleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
