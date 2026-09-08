package app

import (
	"net/http"
	"ssl-custom-api/internal/app/auth"

	"github.com/danielgtaylor/huma/v2"
)

func SetupAuthRoutes(base huma.API, authHandler *auth.Handler) {
	authBase := huma.NewGroup(base, "/auth")
	tags := []string{"app:auth"}
	huma.Register(authBase, huma.Operation{
		OperationID: "loginUser",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "Login user",
		Description: "Authenticate user and return a JWT token",
		Tags:        tags,
		Errors:      []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError},
	}, authHandler.LoginUser)

}
