package app

import (
	"net/http"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/permissions"
	"ssl-custom-api/internal/router/middleware"

	"github.com/danielgtaylor/huma/v2"
)

func SetupPermissionRoutes(base huma.API, handlers *permissions.Handler, jwtSvc *auth.JWT, resolver middleware.PrincipalResolver) {
	permissionBase := huma.NewGroup(base, "/permissions")
	permissionBase.UseMiddleware(middleware.JWTAuth(base, jwtSvc, resolver))
	tags := []string{"app:permissions"}
	openapi := base.OpenAPI()
	if openapi.Components == nil {
		openapi.Components = &huma.Components{}
	}
	if openapi.Components.SecuritySchemes == nil {
		openapi.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}
	}
	openapi.Components.SecuritySchemes["BearerAuth"] = &huma.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "JWT access token from POST /api/v1/app/auth/login",
	}

	security := []map[string][]string{{"BearerAuth": {}}}

	huma.Register(permissionBase, huma.Operation{
		OperationID: "createModulePermission",
		Method:      http.MethodPost,
		Path:        "",
		Summary:     "Create a new module permission for a user",
		Description: "Creates a new module with the provided details",
		Tags:        tags,
		Security:    security,
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusConflict, http.StatusInternalServerError},
	}, handlers.AddModulePermissions)

	huma.Register(permissionBase, huma.Operation{
		OperationID: "getModulePermission",
		Method:      http.MethodGet,
		Path:        "",
		Summary:     "List all module permissions for a user",
		Description: "Retrieves details of a module permission by its ID",
		Tags:        tags,
		Security:    security,
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError},
	}, handlers.GetModulePermittedForUser)
	huma.Register(permissionBase, huma.Operation{
		OperationID: "getModulePermissionById",
		Method:      http.MethodGet,
		Path:        "/user",
		Summary:     "Get a module permission by user ID",
		Description: "Retrieves details of a module permission by the user's ID",
		Tags:        tags,
		Security:    security,
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError},
	}, handlers.GetModulePermittedForUserById)

}
