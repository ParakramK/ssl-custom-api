package app

import (
	"net/http"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/users"
	"ssl-custom-api/internal/router/middleware"

	"github.com/danielgtaylor/huma/v2"
)

func SetupUsersRoutes(base huma.API, handlers *users.Handler, jwtSvc *auth.JWT, resolver middleware.PrincipalResolver) {
	usersBase := huma.NewGroup(base, "/users")
	usersBase.UseMiddleware(middleware.JWTAuth(base, jwtSvc, resolver))
	tags := []string{"app:users"}
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

	huma.Register(usersBase, huma.Operation{
		OperationID: "createUser",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "Create a new user",
		Description: "Creates a new user with the specified details.",
		Tags:        tags,
		Security:    security,
	}, handlers.CreateUser)
	huma.Register(usersBase, huma.Operation{
		OperationID: "listUsers",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "List all users",
		Description: "Retrieves a list of all available users.",
		Tags:        tags,
		Security:    security,
	}, handlers.ListUsers)

}
