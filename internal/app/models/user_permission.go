package models

import (
	"github.com/google/uuid"
)

// UserModulePermission grants a user a permission level on a whole
// Control Plane module (e.g. gatepass means all of Gatepass).
type UserModulePermission struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_user_module_perm,priority:1"`
	ModuleID   uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:uq_user_module_perm,priority:2"`
	Permission Permission `gorm:"type:varchar(20);not null;uniqueIndex:uq_user_module_perm,priority:3"`

	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Module Module `gorm:"foreignKey:ModuleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
