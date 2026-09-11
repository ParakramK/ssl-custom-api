package app

import (
	"net/http"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/roles"
	"ssl-custom-api/internal/router/middleware"

	"github.com/danielgtaylor/huma/v2"
)

func SetupRolesRoutes(base huma.API, handlers *roles.Handler, jwtSvc *auth.JWT, resolver middleware.PrincipalResolver) {
	rolesBase := huma.NewGroup(base, "/roles")
	rolesBase.UseMiddleware(middleware.JWTAuth(base, jwtSvc, resolver))
	tags := []string{"app:roles"}
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

	huma.Register(rolesBase, huma.Operation{
		OperationID: "createRole",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "Create a new role",
		Description: "Creates a new role with the specified name and admin status.",
		Tags:        tags,
		Security:    security,
	}, handlers.CreateRole)
	huma.Register(rolesBase, huma.Operation{
		OperationID: "listRoles",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "List all roles",
		Description: "Retrieves a list of all available roles.",
		Tags:        tags,
		Security:    security,
	}, handlers.ListRoles)
}
