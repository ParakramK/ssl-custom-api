package sales

import (
	"context"
	"log"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetTopOutStandingCustomers(ctx context.Context, input *TopCustomersInput) (*TopCustomersOutput, error) {
	if input.CompanyDB == "" {
		return nil, huma.Error400BadRequest("CompanyDB is required")
	}

	limit := input.Limit

	if limit == 0 {
		limit = 10
	}

	limit = min(limit, 100)

	response, err := h.service.GetTopOutStandingCustomers(
		ctx,
		input.CompanyDB,
		limit,
	)

	if err != nil {
		log.Printf(
			"top outstanding customers error: companyDB=%s error=%v",
			input.CompanyDB,
			err,
		)

		return nil, huma.Error500InternalServerError("failed to fetch top outstanding customers")
	}

	return &TopCustomersOutput{Body: response}, nil
}
