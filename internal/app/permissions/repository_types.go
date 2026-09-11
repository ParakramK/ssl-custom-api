package permissions

import "github.com/google/uuid"

type ModulePermitted struct {
	ModuleName string
	ModuleCode string
	ModuleID   uuid.UUID
	Permission string
}
