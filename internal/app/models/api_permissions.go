package models

import (
	"github.com/google/uuid"
)

// APIPermission grants an API key a permission on a single API resource
// (e.g. gatepass.sales + write). Grants are always resource-level: a
// key must never receive a broad module-level permission such as
// "gatepass write", which would authorize every Gatepass endpoint.
type APIPermission struct {
	ID           uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ApiKeyID     uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:uq_api_key_perm,priority:1"`
	ResourceCode ResourceCode `gorm:"type:varchar(100);not null;uniqueIndex:uq_api_key_perm,priority:2"`
	Permission   Permission   `gorm:"type:varchar(20);not null;uniqueIndex:uq_api_key_perm,priority:3"`

	ApiKey ApiKey `gorm:"foreignKey:ApiKeyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
