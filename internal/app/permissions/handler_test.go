package permissions

import (
	"context"
	"testing"

	"ssl-custom-api/internal/app/auth"

	"github.com/google/uuid"
)

type stubPermissionService struct {
	PermissionService
	gotUserID  uuid.UUID
	gotIsAdmin bool
	modules    []ModulePermitted
}

func (s *stubPermissionService) GetModulePermittedForUser(
	_ context.Context,
	userID uuid.UUID,
	isAdmin bool,
) ([]ModulePermitted, error) {
	s.gotUserID = userID
	s.gotIsAdmin = isAdmin
	return s.modules, nil
}

func TestGetModulePermittedUsesPrincipal(t *testing.T) {
	userID := uuid.New()
	svc := &stubPermissionService{}
	h := NewHandler(svc)

	ctx := auth.WithPrincipal(context.Background(), auth.Principal{UserID: userID, IsAdmin: true})
	out, err := h.GetModulePermittedForUser(ctx, &ModulePermittedRequestCurrentUser{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if svc.gotUserID != userID {
		t.Errorf("service got user %v, want %v", svc.gotUserID, userID)
	}
	if !svc.gotIsAdmin {
		t.Errorf("expected isAdmin true to reach service")
	}
	if out == nil {
		t.Errorf("expected non-nil output")
	}
}

func TestGetModulePermittedPassesNonAdmin(t *testing.T) {
	svc := &stubPermissionService{}
	h := NewHandler(svc)

	userID := uuid.New()

	ctx := auth.WithPrincipal(context.Background(), auth.Principal{UserID: userID, IsAdmin: false})
	if _, err := h.GetModulePermittedForUser(ctx, &ModulePermittedRequestCurrentUser{}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if svc.gotIsAdmin {
		t.Errorf("expected isAdmin false to reach service")
	}
}

func TestGetModulePermittedPanicsWithoutPrincipal(t *testing.T) {
	h := NewHandler(&stubPermissionService{})

	defer func() {
		if recover() == nil {
			t.Errorf("expected panic when principal is missing")
		}
	}()
	_, _ = h.GetModulePermittedForUser(context.Background(), &ModulePermittedRequestCurrentUser{})
}
