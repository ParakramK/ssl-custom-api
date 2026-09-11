package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrPrincipalNotFound is returned when resolving a Principal for a user
// ID that no longer exists (e.g. deleted after the JWT was issued).
// Callers should treat this as failed authentication (401), unlike other
// resolution errors which are internal failures.
var ErrPrincipalNotFound = errors.New("auth: principal not found")

type Principal struct {
	UserID  uuid.UUID
	RoleID  uuid.UUID
	IsAdmin bool
}

type principalKey struct{}

type userIDKey struct{}

// WithPrincipal stores the authenticated principal in the context.
func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, principal)
}

// PrincipalFromContext returns the authenticated principal.
func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalKey{}).(Principal)
	return principal, ok
}

// MustPrincipal returns the authenticated principal.
func MustPrincipal(ctx context.Context) Principal {
	principal, ok := PrincipalFromContext(ctx)
	if !ok {
		panic("auth: no principal in context; JWT middleware misconfigured or missing")
	}
	return principal
}

// Kept for backward compatibility; new code should use WithPrincipal.
// It preserves any existing authorization state on the principal.
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	if existing, ok := PrincipalFromContext(ctx); ok {
		existing.UserID = userID
		return WithPrincipal(ctx, existing)
	}
	return WithPrincipal(ctx, Principal{UserID: userID})
}

// UserIDFromContext returns the authenticated user's ID.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	if principal, ok := PrincipalFromContext(ctx); ok {
		return principal.UserID, true
	}
	userID, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return userID, ok
}
