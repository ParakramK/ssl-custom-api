package auth

import (
	"context"

	"github.com/google/uuid"
)

type userIDKey struct{}

// WithUserID stores the authenticated user's ID in the context.
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserIDFromContext returns the authenticated user's ID.
// The second return value is false when no validated JWT provided one,
// in which case handlers must reject the request.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return userID, ok
}
