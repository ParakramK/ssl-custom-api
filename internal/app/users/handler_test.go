package users

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/constants"
	"ssl-custom-api/internal/app/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type stubUserRepository struct {
	createCalls int
	createErr   error
}

func (s *stubUserRepository) GetUserByUsername(_ string) (*models.User, error) {
	return nil, nil
}

func (s *stubUserRepository) CreateUser(_ context.Context, input *CreateUserRequest) (*CreateUserResponse, error) {
	s.createCalls++
	if s.createErr != nil {
		return nil, s.createErr
	}
	return &CreateUserResponse{ID: uuid.New(), Username: input.Username}, nil
}

func (s *stubUserRepository) UpdateUser(_ *models.User) error { return nil }

func (s *stubUserRepository) DeleteUser(_ *models.User) error { return nil }

func (s *stubUserRepository) ChangeUserPassword(_, _, _ string) error { return nil }

func (s *stubUserRepository) ListUsers(_ context.Context, _ *uuid.UUID) (*ListUsersResponse, error) {
	return &ListUsersResponse{
		Users: []UserRow{
			{ID: uuid.New(), Username: "user1", RoleName: "role1"},
			{ID: uuid.New(), Username: "user2", RoleName: "role2"},
		},
	}, nil
}

func TestCreateUserAdmin(t *testing.T) {
	repo := &stubUserRepository{}
	h := NewHandler(NewService(repo))

	ctx := auth.WithPrincipal(context.Background(), auth.Principal{UserID: uuid.New(), IsAdmin: true})
	out, err := h.CreateUser(ctx, &CreateUserInput{Body: CreateUserRequest{Username: "newuser"}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.createCalls != 1 {
		t.Errorf("expected repository to be called once, got %d", repo.createCalls)
	}
	if out.Body.Username != "newuser" {
		t.Errorf("unexpected output: %+v", out.Body)
	}
}

func TestCreateUserNonAdminForbidden(t *testing.T) {
	repo := &stubUserRepository{}
	h := NewHandler(NewService(repo))

	ctx := auth.WithPrincipal(context.Background(), auth.Principal{UserID: uuid.New(), IsAdmin: false})
	_, err := h.CreateUser(ctx, &CreateUserInput{Body: CreateUserRequest{Username: "newuser"}})
	if err == nil {
		t.Fatalf("expected forbidden error")
	}
	se, ok := err.(huma.StatusError)
	if !ok || se.GetStatus() != http.StatusForbidden {
		t.Errorf("expected 403 error, got %v", err)
	}
	if repo.createCalls != 0 {
		t.Errorf("repository must not be called for non-admin, got %d calls", repo.createCalls)
	}
}

func TestCreateUserDuplicateIsConflict(t *testing.T) {
	repo := &stubUserRepository{
		// Mirror the repository: sentinel wrapped with %w.
		createErr: fmt.Errorf("%w: %q", constants.ErrUserExists, "newuser"),
	}
	h := NewHandler(NewService(repo))

	ctx := auth.WithPrincipal(context.Background(), auth.Principal{UserID: uuid.New(), IsAdmin: true})
	_, err := h.CreateUser(ctx, &CreateUserInput{Body: CreateUserRequest{Username: "newuser"}})
	if err == nil {
		t.Fatalf("expected conflict error")
	}
	se, ok := err.(huma.StatusError)
	if !ok || se.GetStatus() != http.StatusConflict {
		t.Errorf("expected 409 error, got %v", err)
	}
}

func TestCreateUserPanicsWithoutPrincipal(t *testing.T) {
	h := NewHandler(NewService(&stubUserRepository{}))

	defer func() {
		if recover() == nil {
			t.Errorf("expected panic when principal is missing")
		}
	}()
	_, _ = h.CreateUser(context.Background(), &CreateUserInput{})
}
