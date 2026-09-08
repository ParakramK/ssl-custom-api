package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"ssl-custom-api/internal/app/auth"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// UserIDFromContext returns the authenticated user's ID stored by JWTAuth.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}

// JWTAuth is huma middleware that requires a valid Bearer token.
// Requests without one are rejected with 401 and never reach the handler.
func JWTAuth(api huma.API, jwtSvc *auth.JWT) func(huma.Context, func(huma.Context)) {
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

		next(huma.WithValue(ctx, userIDKey, userID))
	}
}
