package app

import (
	"ssl-custom-api/internal/app"

	"github.com/danielgtaylor/huma/v2"
)

func SetupAppRoutes(base huma.API, handlers *app.Handlers) {
	appRoutes := huma.NewGroup(base, "/app")
	SetupAuthRoutes(appRoutes, handlers.App.Auth)
	SetupApiKeysRoutes(appRoutes, handlers.App.ApiKey, handlers.JWT, handlers.Authorization)
	SetupModulesRoutes(appRoutes, handlers.App.Module, handlers.JWT, handlers.Authorization)
	SetupRolesRoutes(appRoutes, handlers.App.Role, handlers.JWT, handlers.Authorization)
	SetupUsersRoutes(appRoutes, handlers.App.Users, handlers.JWT, handlers.Authorization)
	SetupPermissionRoutes(appRoutes, handlers.App.Permissions, handlers.JWT, handlers.Authorization)
}
