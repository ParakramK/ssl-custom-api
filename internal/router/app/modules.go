package app

import (
	"net/http"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/modules"
	"ssl-custom-api/internal/router/middleware"

	"github.com/danielgtaylor/huma/v2"
)

func SetupModulesRoutes(base huma.API, handlers *modules.Handler, jwtSvc *auth.JWT, resolver middleware.PrincipalResolver) {
	modulesBase := huma.NewGroup(base, "/modules")
	modulesBase.UseMiddleware(middleware.JWTAuth(base, jwtSvc, resolver))
	tags := []string{"app:modules"}
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

	huma.Register(modulesBase, huma.Operation{
		OperationID: "createModule",
		Method:      http.MethodPost,
		Path:        "",
		Summary:     "Create a new module",
		Description: "Creates a new module with the provided details",
		Tags:        tags,
		Security:    security,
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusConflict, http.StatusInternalServerError},
	}, handlers.CreateModule)
	huma.Register(modulesBase, huma.Operation{
		OperationID: "getModule",
		Method:      http.MethodGet,
		Path:        "",
		Summary:     "Get module details",
		Description: "Retrieves details of a module by its ID or code",
		Tags:        tags,
		Security:    security,
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError},
	}, handlers.GetModuleInfo)
}
