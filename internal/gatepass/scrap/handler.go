package scrap

import (
	"context"
	"errors"
	"log"
	"ssl-custom-api/internal/app/constants"

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

func (h *Handler) GetQualityReport(ctx context.Context, input *QualityReportInput) (*QualityReportOutput, error) {
	if input.DocumentNo == "" {
		return nil, huma.Error400BadRequest("ssl_sno is required")
	}

	response, err := h.service.GetQualityReport(ctx, input.DocumentNo)
	if err != nil {
		if errors.Is(err, constants.ErrQualityReportNotFound) {
			return nil, huma.Error404NotFound("No quality report found for the given vendor code and bill number")
		}

		log.Printf(
			"quality report error: documentNo=%s error=%v",
			input.DocumentNo,
			err,
		)

		return nil, huma.Error500InternalServerError("failed to fetch quality report")
	}

	return &QualityReportOutput{Body: response}, nil
}
