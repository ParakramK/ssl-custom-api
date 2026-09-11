package constants

import "errors"

var ErrUnauthorized = errors.New("user is not authorized to perform this action")

// Sentinel errors for API key creation
var (
	ErrScopesRequired      = errors.New("at least one scope is required")
	ErrUnknownResource     = errors.New("unknown resource")
	ErrInvalidPermission   = errors.New("invalid permission")
	ErrCannotCreateAPIKey  = errors.New("user may not create API keys")
	ErrScopeNotDelegatable = errors.New("scope exceeds user's delegatable permissions")
)

// Sentinel errors for modules and resources
var ErrModuleExists = errors.New("module with the same name or code already exists")
var ErrResourceExists = errors.New("resource with the same name or code already exists")
var ErrUnknownModule = errors.New("unknown module")

var ErrRoleExists = errors.New("role with the same name already exists")
var ErrUserExists = errors.New("user with the same username or email already exists")

var ErrQualityReportNotFound = errors.New("quality report not found")
