package modules

import (
	"context"
	"errors"
	"fmt"

	"ssl-custom-api/internal/app/constants"
	"ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/app/query"

	"ssl-custom-api/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	Create(ctx context.Context, input ModuleCreateRequest) (*models.Module, error)
	CreateResource(ctx context.Context, input ResourceCreateRequest) (*models.Resource, error)
	GetByID(ctx context.Context, id uuid.UUID) (ModuleRow, error)
	GetByCode(ctx context.Context, code string) (ModuleRow, error)
	ListModules(ctx context.Context) ([]ModuleRow, error)
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}

type moduleRepository struct {
	db *gorm.DB
}

func (r *moduleRepository) CheckModuleExists(
	ctx context.Context,
	module *models.Module,
) (*models.Module, error) {
	q := query.Use(r.db)

	existing, err := q.Module.WithContext(ctx).
		Where(
			q.Module.Code.Eq(module.Code),
		).
		Or(
			q.Module.Name.Eq(module.Name),
		).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return existing, fmt.Errorf("%w: %q", constants.ErrModuleExists, module.Code)
}
func (r *moduleRepository) CheckResourceExists(
	ctx context.Context,
	resource *models.Resource,
) (*models.Resource, error) {
	q := query.Use(r.db)

	existing, err := q.Resource.WithContext(ctx).
		Where(
			q.Resource.Code.Eq(resource.Code),
		).
		Or(
			q.Resource.Name.Eq(resource.Name),
		).
		Where(
			q.Resource.ModuleID.Eq(resource.ModuleID),
		).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return existing, fmt.Errorf("%w: %q", constants.ErrResourceExists, resource.Code)
}

func (r *moduleRepository) Create(
	ctx context.Context,
	input ModuleCreateRequest,
) (*models.Module, error) {
	q := query.Use(r.db)

	existing, err := r.CheckModuleExists(ctx, &models.Module{
		Code: input.Code,
		Name: input.Name,
	})
	if err != nil {
		return nil, err
	}
	_ = existing

	module := &models.Module{
		ID:          utils.NewV7ID(),
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
	}

	if err := q.Module.WithContext(ctx).Create(module); err != nil {
		return nil, err
	}

	return module, nil
}

func (r *moduleRepository) GetByCode(ctx context.Context, code string) (ModuleRow, error) {
	q := query.Use(r.db)
	module := ModuleRow{}
	err := q.Module.WithContext(ctx).
		Where(q.Module.Code.Eq(code)).
		Scan(&module)
	if err != nil {
		return ModuleRow{}, err
	}

	return module, nil
}
func (r *moduleRepository) GetByID(ctx context.Context, id uuid.UUID) (ModuleRow, error) {
	module := ModuleRow{}
	q := query.Use(r.db)
	err := q.Module.WithContext(ctx).
		Where(q.Module.ID.Eq(id)).
		Scan(&module)

	if err != nil {
		return ModuleRow{}, err
	}
	return module, nil
}

func (r *moduleRepository) CreateResource(ctx context.Context, input ResourceCreateRequest) (*models.Resource, error) {
	q := query.Use(r.db)
	existing, err := r.CheckResourceExists(ctx, &models.Resource{
		Code:     input.Code,
		Name:     input.Name,
		ModuleID: input.ModuleID,
	})
	if err != nil {
		return nil, err
	}
	_ = existing

	resource := &models.Resource{
		ID:          utils.NewV7ID(),
		Code:        input.Code,
		Name:        input.Name,
		ModuleID:    input.ModuleID,
		Description: input.Description,
	}
	if err := q.Resource.WithContext(ctx).Create(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func (r *moduleRepository) ListModules(ctx context.Context) ([]ModuleRow, error) {
	q := query.Use(r.db)
	var modules []ModuleRow
	err := q.Module.WithContext(ctx).
		Scan(&modules)
	if err != nil {
		return nil, err
	}
	return modules, nil
}
