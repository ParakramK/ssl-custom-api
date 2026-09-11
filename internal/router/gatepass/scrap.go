package gatepass

import (
	"net/http"

	"ssl-custom-api/internal/gatepass/scrap"

	"github.com/danielgtaylor/huma/v2"
)

func SetupScrapRoutes(
	base huma.API,
	scrapHandler *scrap.Handler,
) {
	tags := []string{"gatepass:scrap"}
	huma.Register(base, huma.Operation{
		OperationID: "getQualityReport",
		Method:      http.MethodGet,
		Path:        "/quality-report",
		Summary:     "Get quality report",
		Description: "Returns the quality report for an SSL entry number, " +
			"including supplier and billing details plus grading details " +
			"with payable lines.",
		Tags:   tags,
		Errors: []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, scrapHandler.GetQualityReport)

}
