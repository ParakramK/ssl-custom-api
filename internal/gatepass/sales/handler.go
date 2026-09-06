package sales

import (
	"context"
	"errors"
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

func (h *Handler) GetTagData(ctx context.Context, input *TagDataInput) (*TagDataOutput, error) {
	if input.BundleTag == "" {
		return nil, huma.Error400BadRequest("bundle_tag is required")
	}

	response, err := h.service.GetTagData(ctx, input.BundleTag)
	if err != nil {
		if errors.Is(err, ErrTagNotFound) {
			return nil, huma.Error404NotFound("No tag data found for the given bundle tag")
		}

		log.Printf(
			"tag data error: bundleTag=%s error=%v",
			input.BundleTag,
			err,
		)

		return nil, huma.Error500InternalServerError("failed to fetch tag data")
	}

	return &TagDataOutput{Body: response}, nil
}

func (h *Handler) GetPackingList(ctx context.Context, input *PackingListInput) (*PackingListOutput, error) {
	if input.SslSno == "" {
		return nil, huma.Error400BadRequest("ssl_no is required")
	}

	response, err := h.service.GetPackingList(ctx, input.SslSno)
	if err != nil {
		if errors.Is(err, ErrPackingListNotFound) {
			return nil, huma.Error404NotFound("No packing list found for the given ssl number")
		}

		log.Printf(
			"packing list error: sslSno=%s error=%v",
			input.SslSno,
			err,
		)

		return nil, huma.Error500InternalServerError("failed to fetch packing list")
	}
	if response == nil {
		response = &PackingList{}
	}

	return &PackingListOutput{Body: *response}, nil
}
