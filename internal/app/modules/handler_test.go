package modules

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type stubModuleService struct {
	createCalls         int
	createResourceCalls int
	createErr           error
	createResourceErr   error
}

func (s *stubModuleService) Create(_ context.Context, input ModuleCreateRequest) (*ModuleCreateResponse, error) {
	s.createCalls++
	if s.createErr != nil {
		return nil, s.createErr
	}
	return &ModuleCreateResponse{Name: input.Name, Code: input.Code, Message: "ok"}, nil
}

func (s *stubModuleService) CreateResource(_ context.Context, input ResourceCreateRequest) (*ResourceCreateResponse, error) {
	s.createResourceCalls++
	if s.createResourceErr != nil {
		return nil, s.createResourceErr
	}
	return &ResourceCreateResponse{Name: input.Name, Code: input.Code, Message: "ok"}, nil
}

func (s *stubModuleService) GetModule(_ context.Context, _ *ModuleInfoRequest) (*ModuleInfoResponse, error) {
	return &ModuleInfoResponse{}, nil
}
func (s *stubModuleService) ListModules(_ context.Context) ([]ModuleInfoResponse, error) {
	return []ModuleInfoResponse{
		{Name: "Module1", Code: "M1"},
		{Name: "Module2", Code: "M2"},
	}, nil
}
func adminCtx() context.Context {
	return auth.WithPrincipal(context.Background(), auth.Principal{
		UserID:  uuid.New(),
		RoleID:  uuid.New(),
		IsAdmin: true,
	})
}

func memberCtx() context.Context {
	return auth.WithPrincipal(context.Background(), auth.Principal{
		UserID:  uuid.New(),
		RoleID:  uuid.New(),
		IsAdmin: false,
	})
}

func TestCreateModuleAdmin(t *testing.T) {
	svc := &stubModuleService{}
	h := NewHandler(svc)

	out, err := h.CreateModule(adminCtx(), &ModuleCreateInput{Body: ModuleCreateRequest{Name: "M", Code: "m"}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if svc.createCalls != 1 {
		t.Errorf("expected service to be called once, got %d", svc.createCalls)
	}
	if out.Body.Code != "m" {
		t.Errorf("unexpected output: %+v", out.Body)
	}
}

func TestCreateModuleNonAdminForbidden(t *testing.T) {
	svc := &stubModuleService{}
	h := NewHandler(svc)

	_, err := h.CreateModule(memberCtx(), &ModuleCreateInput{Body: ModuleCreateRequest{Name: "M", Code: "m"}})
	assertForbidden(t, err)
	if svc.createCalls != 0 {
		t.Errorf("service must not be called for non-admin, got %d calls", svc.createCalls)
	}
}

func TestCreateResourceNonAdminForbidden(t *testing.T) {
	svc := &stubModuleService{}
	h := NewHandler(svc)

	_, err := h.CreateResource(memberCtx(), &ResourceCreateInput{Body: ResourceCreateRequest{Name: "R", Code: "r"}})
	assertForbidden(t, err)
	if svc.createResourceCalls != 0 {
		t.Errorf("service must not be called for non-admin, got %d calls", svc.createResourceCalls)
	}
}

func TestCreateModuleDuplicateIsConflict(t *testing.T) {
	svc := &stubModuleService{
		// Mirror the repository: sentinel wrapped with %w.
		createErr: fmt.Errorf("%w: %q", constants.ErrModuleExists, "gatepass"),
	}
	h := NewHandler(svc)

	_, err := h.CreateModule(adminCtx(), &ModuleCreateInput{Body: ModuleCreateRequest{Name: "G", Code: "gatepass"}})
	assertConflict(t, err)
}

func TestCreateResourceDuplicateIsConflict(t *testing.T) {
	svc := &stubModuleService{
		createResourceErr: fmt.Errorf("%w: %q", constants.ErrResourceExists, "sales"),
	}
	h := NewHandler(svc)

	_, err := h.CreateResource(adminCtx(), &ResourceCreateInput{Body: ResourceCreateRequest{Name: "S", Code: "sales"}})
	assertConflict(t, err)
}

func TestCreateModulePanicsWithoutPrincipal(t *testing.T) {
	h := NewHandler(&stubModuleService{})

	defer func() {
		if recover() == nil {
			t.Errorf("expected panic when principal is missing")
		}
	}()
	_, _ = h.CreateModule(context.Background(), &ModuleCreateInput{})
}

func assertConflict(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected conflict error")
	}
	se, ok := err.(huma.StatusError)
	if !ok || se.GetStatus() != http.StatusConflict {
		t.Errorf("expected 409 error, got %v", err)
	}
}

func assertForbidden(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected forbidden error")
	}
	se, ok := err.(huma.StatusError)
	if !ok || se.GetStatus() != http.StatusForbidden {
		t.Errorf("expected 403 error, got %v", err)
	}
	if !strings.Contains(err.Error(), "admin access required") {
		t.Errorf("expected admin message, got %v", err)
	}
}
