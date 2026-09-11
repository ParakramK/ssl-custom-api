package permissions

import (
	"ssl-custom-api/internal/app/models"

	"github.com/google/uuid"
)

type ModulePermittedOutput struct {
	Body ModulePermittedResponse `json:"body"`
}
type ModulePermittedResponse struct {
	Modules []ModulePermitted `json:"modules"`
}
type ModulePermittedRequestCurrentUser struct{}
type ModulePermittedRequest struct {
	UserID uuid.UUID `query:"user_id" required:"true" doc:"User ID" minLength:"1" example:"00000000-0000-0000-0000-000000000000"`
}

type CreateModulePermissionInput struct {
	Body CreateModulePermissionRequest `json:"body"`
}

// ModulePermissionGrant assigns a permission level on a single module.
type ModulePermissionGrant struct {
	ModuleID   uuid.UUID         `json:"module_id" doc:"Module ID" example:"00000000-0000-0000-0000-000000000000"`
	Permission models.Permission `json:"permission" doc:"Permission level on this module (read, write, delete, admin)" example:"read"`
}
type CreateModulePermissionRequest struct {
	UserID  uuid.UUID               `json:"user_id" doc:"User ID" minLength:"1" example:"00000000-0000-0000-0000-000000000000"`
	Modules []ModulePermissionGrant `json:"modules" doc:"Per-module permission grants" minLength:"1"`
}
type CreateModulePermissionOutput struct {
	Body CreateModulePermissionResponse `json:"body"`
}
type CreateModulePermissionResponse struct {
	Message string `json:"message" doc:"Response message" example:"module permission added successfully"`
}
