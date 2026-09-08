package app

import (
	"net/http"
	apikey "ssl-custom-api/internal/app/api_key"
	"ssl-custom-api/internal/app/auth"

	"ssl-custom-api/internal/router/middleware"

	"github.com/danielgtaylor/huma/v2"
)

func SetupApiKeysRoutes(base huma.API, apikeyHandler *apikey.Handler, jwtSvc *auth.JWT) {
	keysBase := huma.NewGroup(base, "/keys")
	keysBase.UseMiddleware(middleware.JWTAuth(base, jwtSvc))
	tags := []string{"app:keys"}

	openAPI := base.OpenAPI()
	if openAPI.Components == nil {
		openAPI.Components = &huma.Components{}
	}
	if openAPI.Components.SecuritySchemes == nil {
		openAPI.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}
	}
	openAPI.Components.SecuritySchemes["BearerAuth"] = &huma.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "JWT access token from POST /api/v1/app/auth/login",
	}

	security := []map[string][]string{{"BearerAuth": {}}}

	huma.Register(keysBase, huma.Operation{
		OperationID: "listApiKeys",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "List API keys",
		Description: "Retrieve API keys with keyset pagination. Query params: limit (default 10, max 100), cursor (opaque token from next_cursor), sort (id,key; prefix with - for desc, e.g. ?sort=-key), user_id (filter by user).",
		Tags:        tags,
		Security:    security,
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError},
	}, apikeyHandler.ListApiKeys)
}
