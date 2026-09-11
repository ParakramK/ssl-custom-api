package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestPrincipalFromContextPresent(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()
	ctx := WithPrincipal(context.Background(), Principal{UserID: userID, RoleID: roleID, IsAdmin: true})

	got, ok := PrincipalFromContext(ctx)
	if !ok {
		t.Fatalf("expected principal present")
	}
	if got.UserID != userID || got.RoleID != roleID || !got.IsAdmin {
		t.Errorf("unexpected principal: %+v", got)
	}
}

func TestPrincipalFromContextAbsent(t *testing.T) {
	_, ok := PrincipalFromContext(context.Background())
	if ok {
		t.Errorf("expected no principal in empty context")
	}
}

func TestMustPrincipal(t *testing.T) {
	userID := uuid.New()
	ctx := WithPrincipal(context.Background(), Principal{UserID: userID})

	if got := MustPrincipal(ctx); got.UserID != userID {
		t.Errorf("unexpected user ID: %v", got.UserID)
	}

	defer func() {
		if recover() == nil {
			t.Errorf("expected panic for missing principal")
		}
	}()
	MustPrincipal(context.Background())
}

func TestUserIDCompat(t *testing.T) {
	userID := uuid.New()

	// WithUserID stores a principal readable via both helpers.
	ctx := WithUserID(context.Background(), userID)
	if p, ok := PrincipalFromContext(ctx); !ok || p.UserID != userID {
		t.Errorf("WithUserID should populate principal, got %+v %v", p, ok)
	}
	if got, ok := UserIDFromContext(ctx); !ok || got != userID {
		t.Errorf("UserIDFromContext should read principal, got %v %v", got, ok)
	}

	// WithUserID preserves an existing RoleID for forward compatibility.
	roleID := uuid.New()
	ctx = WithPrincipal(context.Background(), Principal{UserID: uuid.New(), RoleID: roleID})
	ctx = WithUserID(ctx, userID)
	if p, _ := PrincipalFromContext(ctx); p.UserID != userID || p.RoleID != roleID {
		t.Errorf("WithUserID should preserve RoleID, got %+v", p)
	}

	// Absent everywhere.
	if _, ok := UserIDFromContext(context.Background()); ok {
		t.Errorf("expected no user ID in empty context")
	}
}
