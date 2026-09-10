package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"ssl-custom-api/internal/app/auth"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

// PrincipalResolver builds the request Principal for an authenticated
// user ID. It is satisfied by the authorization service.
type PrincipalResolver interface {
	GetPrincipal(ctx context.Context, userID uuid.UUID) (auth.Principal, error)
}

// JWTAuth is huma middleware that requires a valid Bearer token.
// Requests without one are rejected with 401 and never reach the handler.
// After validating the token it resolves the caller's authorization state
// into a Principal and stores it in the request context. It performs no
// business authorization itself: handlers decide whether the Principal
// may perform the requested operation (403 on denial).
func JWTAuth(api huma.API, jwtSvc *auth.JWT, resolver PrincipalResolver) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		header := ctx.Header("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") || parts[1] == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized",
				fmt.Errorf("missing or malformed Authorization header"))
			return
		}

		userID, err := jwtSvc.ValidateToken(parts[1])
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		principal, err := resolver.GetPrincipal(ctx.Context(), userID)
		if err != nil {
			if errors.Is(err, auth.ErrPrincipalNotFound) {
				_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", err)
			} else {
				_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "Internal Server Error", err)
			}
			return
		}

		next(huma.WithContext(ctx, auth.WithPrincipal(ctx.Context(), principal)))
	}
}
