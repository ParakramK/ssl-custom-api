package authorization

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/models"

	"github.com/google/uuid"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envOr("PG_HOST", "localhost"),
		envOr("PG_PORT", "5432"),
		envOr("PG_USER", "postgres"),
		envOr("PG_PASSWORD", "postgres"),
		envOr("PG_DATABASE", "ssl_db"),
		envOr("PG_SSL_MODE", "disable"),
	)
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("postgres not available: %v", err)
	}

	tx := db.Begin()
	t.Cleanup(func() { tx.Rollback() })

	if err := tx.AutoMigrate(&models.Role{}, &models.User{}); err != nil {
		t.Fatalf("migrate test schema: %v", err)
	}
	return tx
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func seedUser(t *testing.T, tx *gorm.DB, isAdmin bool) (uuid.UUID, uuid.UUID) {
	t.Helper()

	suffix := uuid.NewString()[:8]
	role := models.Role{Name: "role-" + suffix, IsAdmin: isAdmin}
	if err := tx.Create(&role).Error; err != nil {
		t.Fatal(err)
	}
	user := models.User{
		Username: "user-" + suffix,
		Email:    suffix + "@example.com",
		Password: "x",
		RoleID:   role.ID,
	}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	return user.ID, role.ID
}

func TestGetPrincipalAdmin(t *testing.T) {
	tx := testDB(t)
	userID, roleID := seedUser(t, tx, true)

	svc := NewService(NewAuthorizationRepository(tx))
	got, err := svc.GetPrincipal(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.UserID != userID {
		t.Errorf("user ID: got %v, want %v", got.UserID, userID)
	}
	if got.RoleID != roleID {
		t.Errorf("role ID: got %v, want %v", got.RoleID, roleID)
	}
	if !got.IsAdmin {
		t.Errorf("expected IsAdmin true")
	}
}

func TestGetPrincipalNonAdmin(t *testing.T) {
	tx := testDB(t)
	userID, roleID := seedUser(t, tx, false)

	svc := NewService(NewAuthorizationRepository(tx))
	got, err := svc.GetPrincipal(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.UserID != userID || got.RoleID != roleID {
		t.Errorf("unexpected principal: %+v", got)
	}
	if got.IsAdmin {
		t.Errorf("expected IsAdmin false")
	}
}

func TestGetPrincipalUnknownUser(t *testing.T) {
	tx := testDB(t)

	svc := NewService(NewAuthorizationRepository(tx))
	_, err := svc.GetPrincipal(context.Background(), uuid.New())
	if !errors.Is(err, auth.ErrPrincipalNotFound) {
		t.Errorf("expected ErrPrincipalNotFound, got %v", err)
	}
}
