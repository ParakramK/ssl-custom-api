package customer

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

func (h *Handler) GetCustomerAging(ctx context.Context, input *AgingInput) (*AgingOutput, error) {
	if input.CardCode == "" {
		return nil, huma.Error400BadRequest("CardCode is required")
	}

	if input.CompanyDB == "" {
		return nil, huma.Error400BadRequest("CompanyDB is required")
	}

	response, err := h.service.GetCustomerAging(
		ctx,
		input.CardCode,
		input.CompanyDB,
	)
	if err != nil {
		log.Printf(
			"customer aging error: cardCode=%s companyDB=%s error=%v",
			input.CardCode,
			input.CompanyDB,
			err,
		)

		return nil, huma.Error500InternalServerError("failed to fetch customer aging")
	}

	return &AgingOutput{Body: response}, nil
}
