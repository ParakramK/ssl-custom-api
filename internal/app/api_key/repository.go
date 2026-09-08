package apikey

import (
	"context"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/app/paging"
	"ssl-custom-api/internal/app/query"

	"github.com/google/uuid"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type ApiKeyRepository interface {
	CreateApiKey(apiKey *models.ApiKey) error
	GetApiKeyByKey(ctx context.Context, key string) (*models.ApiKey, error)
	DeleteApiKey(apiKey *models.ApiKey) error
	ListApiKeys(ctx context.Context, userID uuid.UUID, pg paging.Query) ([]ApiKeyRow, bool, error)
}

func NewApiKeyRepository(db *gorm.DB) ApiKeyRepository {
	return &apiKeyRepository{db: db}
}

type apiKeyRepository struct {
	db *gorm.DB
}

func (r *apiKeyRepository) CreateApiKey(apiKey *models.ApiKey) error {
	result := r.db.Create(apiKey)
	return result.Error
}

func (r *apiKeyRepository) GetApiKeyByKey(ctx context.Context, key string) (*models.ApiKey, error) {
	q := query.Use(r.db)
	apiKey, err := q.ApiKey.WithContext(ctx).
		Select(
			q.ApiKey.ID,
			q.ApiKey.Key,
			q.ApiKey.UserID,
		).
		Where(q.ApiKey.Key.Eq(key)).
		First()
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

func (r *apiKeyRepository) DeleteApiKey(apiKey *models.ApiKey) error {
	result := r.db.Delete(apiKey)
	return result.Error
}

func (r *apiKeyRepository) ListApiKeys(
	ctx context.Context,
	userID uuid.UUID,
	pg paging.Query,
) ([]ApiKeyRow, bool, error) {
	q := query.Use(r.db)
	var apiKeyRows []ApiKeyRow
	stmt := q.ApiKey.WithContext(ctx).
		Select(
			q.ApiKey.ID,
			q.ApiKey.Key,
			q.User.Username.As("user_name"),
		).
		LeftJoin(q.User, q.User.ID.EqCol(q.ApiKey.UserID))
	if userID != uuid.Nil {
		stmt = stmt.Where(q.ApiKey.UserID.Eq(userID))
	}
	for _, s := range pg.Sorts {
		stmt = stmt.Order(sortColumn(q, s.Field, s.Desc))
	}
	if pg.Cursor != nil {
		stmt = stmt.Where(keysetCondition(q, pg.LeadingSort(), pg.Cursor))
	}
	limit := pg.Limit
	if limit < 1 {
		limit = 10
	}
	// Fetch one extra row to know whether another page exists.
	if err := stmt.Limit(limit + 1).Scan(&apiKeyRows); err != nil {
		return nil, false, err
	}

	if len(apiKeyRows) > limit {
		return apiKeyRows[:limit], true, nil
	}
	return apiKeyRows, false, nil
}

// sortColumn maps a middleware sort field to a real column.
// Unknown fields fall back to id so a bad value can never inject SQL.
func sortColumn(q *query.Query, f string, desc bool) field.Expr {
	var col field.Expr = q.ApiKey.ID
	if f == "key" {
		col = q.ApiKey.Key
	}
	if desc {
		return col.Desc()
	}
	return col.Asc()
}

// keysetCondition positions the scan after the cursor row on the leading
// sort column, with id as the tiebreak.
func keysetCondition(q *query.Query, lead paging.Sort, cur *paging.Cursor) field.Expr {
	if lead.Field == "key" {
		tie := field.And(q.ApiKey.Key.Eq(cur.Key), q.ApiKey.ID.Gt(cur.ID))
		if lead.Desc {
			return field.Or(q.ApiKey.Key.Lt(cur.Key), tie)
		}
		return field.Or(q.ApiKey.Key.Gt(cur.Key), tie)
	}
	if lead.Desc {
		return q.ApiKey.ID.Lt(cur.ID)
	}
	return q.ApiKey.ID.Gt(cur.ID)
}
