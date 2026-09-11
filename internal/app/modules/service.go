package modules

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ModuleService interface {
	Create(ctx context.Context, input ModuleCreateRequest) (*ModuleCreateResponse, error)
	CreateResource(ctx context.Context, input ResourceCreateRequest) (*ResourceCreateResponse, error)
	GetModule(ctx context.Context, input *ModuleInfoRequest) (*ModuleInfoResponse, error)
	ListModules(ctx context.Context) ([]ModuleInfoResponse, error)
}

type moduleService struct {
	repo ModuleRepository
}

func NewModuleService(repo ModuleRepository) ModuleService {
	return &moduleService{
		repo: repo,
	}
}

func (s *moduleService) Create(
	ctx context.Context,
	input ModuleCreateRequest,
) (*ModuleCreateResponse, error) {
	module, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return &ModuleCreateResponse{
		Name:    module.Name,
		Code:    module.Code,
		Message: "Module created successfully",
	}, nil
}

func (s *moduleService) CreateResource(
	ctx context.Context,
	input ResourceCreateRequest,
) (*ResourceCreateResponse, error) {
	resource, err := s.repo.CreateResource(ctx, input)
	if err != nil {
		return nil, err
	}

	return &ResourceCreateResponse{
		Name:    resource.Name,
		Code:    resource.Code,
		Message: "Resource created successfully",
	}, nil
}

func (s *moduleService) GetModule(
	ctx context.Context,
	input *ModuleInfoRequest,
) (*ModuleInfoResponse, error) {
	var (
		module ModuleRow
		err    error
	)

	if input.ID != uuid.Nil {
		module, err = s.repo.GetByID(ctx, input.ID)
		if err != nil {
			return nil, err
		}
	} else if input.Code != "" {
		module, err = s.repo.GetByCode(ctx, input.Code)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("either ID or Code must be provided")
	}

	if module == (ModuleRow{}) {
		return nil, errors.New("module not found")
	}

	return &ModuleInfoResponse{
		ID:          module.Id,
		Name:        module.Name,
		Code:        module.Code,
		Description: module.Description,
	}, nil
}

func (s *moduleService) ListModules(ctx context.Context) ([]ModuleInfoResponse, error) {
	modules, err := s.repo.ListModules(ctx)
	if err != nil {
		return nil, err
	}

	var responses []ModuleInfoResponse
	for _, module := range modules {
		responses = append(responses, ModuleInfoResponse{
			ID:          module.Id,
			Name:        module.Name,
			Code:        module.Code,
			Description: module.Description,
		})
	}

	return responses, nil
}
