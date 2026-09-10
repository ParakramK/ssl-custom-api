package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"ssl-custom-api/internal/app/auth"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/uuid"
)

type pingOutput struct {
	Body struct {
		UserID  string `json:"user_id"`
		RoleID  string `json:"role_id"`
		IsAdmin bool   `json:"is_admin"`
	}
}

// stubResolver is a fake PrincipalResolver keyed by user ID.
type stubResolver struct {
	principals map[uuid.UUID]auth.Principal
	err        error
}

func (s *stubResolver) GetPrincipal(_ context.Context, userID uuid.UUID) (auth.Principal, error) {
	if s.err != nil {
		return auth.Principal{}, s.err
	}
	if p, ok := s.principals[userID]; ok {
		return p, nil
	}
	return auth.Principal{}, auth.ErrPrincipalNotFound
}

func setupProtectedAPI(t *testing.T, jwtSvc *auth.JWT, resolver PrincipalResolver) humatest.TestAPI {
	t.Helper()

	_, api := humatest.New(t)
	protected := huma.NewGroup(api, "/protected")
	protected.UseMiddleware(JWTAuth(api, jwtSvc, resolver))

	huma.Register(protected, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/ping",
	}, func(ctx context.Context, _ *struct{}) (*pingOutput, error) {
		principal, ok := auth.PrincipalFromContext(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("missing principal")
		}
		out := &pingOutput{}
		out.Body.UserID = principal.UserID.String()
		out.Body.RoleID = principal.RoleID.String()
		out.Body.IsAdmin = principal.IsAdmin
		return out, nil
	})

	return api
}

func testJWT(t *testing.T) *auth.JWT {
	t.Helper()
	jwtSvc, err := auth.NewJWT("test-secret")
	if err != nil {
		t.Fatal(err)
	}
	return jwtSvc
}

func TestJWTAuthMissingHeader(t *testing.T) {
	api := setupProtectedAPI(t, testJWT(t), &stubResolver{})

	resp := api.Get("/protected/ping")
	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestJWTAuthMalformedHeader(t *testing.T) {
	api := setupProtectedAPI(t, testJWT(t), &stubResolver{})

	for _, header := range []string{
		"Authorization: Token abc",
		"Authorization: Bearer ",
		"Authorization: Bearer",
	} {
		resp := api.Get("/protected/ping", header)
		if resp.Code != http.StatusUnauthorized {
			t.Errorf("header %q: expected 401, got %d", header, resp.Code)
		}
	}
}

func TestJWTAuthInvalidToken(t *testing.T) {
	api := setupProtectedAPI(t, testJWT(t), &stubResolver{})

	resp := api.Get("/protected/ping", "Authorization: Bearer invalid.token.here")
	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.Code)
	}
}

func TestJWTAuthValidTokenBuildsPrincipal(t *testing.T) {
	jwtSvc := testJWT(t)

	userID := uuid.New()
	roleID := uuid.New()
	resolver := &stubResolver{principals: map[uuid.UUID]auth.Principal{
		userID: {UserID: userID, RoleID: roleID, IsAdmin: true},
	}}
	api := setupProtectedAPI(t, jwtSvc, resolver)

	token, err := jwtSvc.GenerateToken(userID, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	resp := api.Get("/protected/ping", "Authorization: Bearer "+token)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	body := resp.Body.String()
	for _, want := range []string{userID.String(), roleID.String()} {
		if !strings.Contains(body, want) {
			t.Errorf("expected body to contain %s, got %s", want, body)
		}
	}
	if !strings.Contains(body, `"is_admin":true`) {
		t.Errorf("expected body to carry is_admin true, got %s", body)
	}
}

func TestJWTAuthUnknownUserIsUnauthorized(t *testing.T) {
	jwtSvc := testJWT(t)
	api := setupProtectedAPI(t, jwtSvc, &stubResolver{})

	token, err := jwtSvc.GenerateToken(uuid.New(), time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	resp := api.Get("/protected/ping", "Authorization: Bearer "+token)
	if resp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unknown user, got %d", resp.Code)
	}
}

func TestJWTAuthResolverErrorIsInternal(t *testing.T) {
	jwtSvc := testJWT(t)
	api := setupProtectedAPI(t, jwtSvc, &stubResolver{err: errors.New("db down")})

	token, err := jwtSvc.GenerateToken(uuid.New(), time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	resp := api.Get("/protected/ping", "Authorization: Bearer "+token)
	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for resolver failure, got %d", resp.Code)
	}
}
