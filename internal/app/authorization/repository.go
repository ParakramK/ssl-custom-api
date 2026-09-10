package authorization

import (
	"context"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/query"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthorizationRepository interface {
	GetPrincipal(ctx context.Context, userID uuid.UUID) (auth.Principal, error)
}

type authorizationRepository struct {
	db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) AuthorizationRepository {
	return &authorizationRepository{db: db}
}

// principalRow holds the authorization state for one user. role_id is
// NOT NULL, so a zero RoleID after the lookup means no row matched.
type principalRow struct {
	RoleID  uuid.UUID `gorm:"column:role_id"`
	IsAdmin bool      `gorm:"column:is_admin"`
}

func (r *authorizationRepository) GetPrincipal(
	ctx context.Context,
	userID uuid.UUID,
) (auth.Principal, error) {
	q := query.Use(r.db)

	var row principalRow

	err := q.User.WithContext(ctx).
		Select(q.User.RoleID, q.Role.IsAdmin).
		Join(q.Role, q.Role.ID.EqCol(q.User.RoleID)).
		Where(q.User.ID.Eq(userID)).
		Scan(&row)

	if err != nil {
		return auth.Principal{}, err
	}

	if row.RoleID == uuid.Nil {
		return auth.Principal{}, auth.ErrPrincipalNotFound
	}

	return auth.Principal{
		UserID:  userID,
		RoleID:  row.RoleID,
		IsAdmin: row.IsAdmin,
	}, nil
}
