package gatepass

import (
	"net/http"

	"ssl-custom-api/internal/gatepass/sales"
	"ssl-custom-api/internal/gatepass/scrap"

	"github.com/danielgtaylor/huma/v2"
)

func SetupScrapRoutes(
	base huma.API,
	scrapHandler *scrap.Handler,
) {

	huma.Register(base, huma.Operation{
		OperationID: "getQualityReport",
		Method:      http.MethodGet,
		Path:        "/quality-report",
		Summary:     "Get quality report",
		Description: "Returns the quality report for an SSL slip number, " +
			"including supplier and billing details plus grading details " +
			"with payable lines.",
		Tags:   []string{"gatepass"},
		Errors: []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, scrapHandler.GetQualityReport)

}
func SetupSalesRoutes(
	base huma.API,
	salesHandler *sales.Handler,
) {

	huma.Register(base, huma.Operation{
		OperationID: "getTagData",
		Method:      http.MethodGet,
		Path:        "/tag",
		Summary:     "Get tag data",
		Description: "Returns the tag data for a Tag Bundle No",

		Tags:   []string{"gatepass"},
		Errors: []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, salesHandler.GetTagData)
	huma.Register(base, huma.Operation{
		OperationID: "getPackingList",
		Method:      http.MethodGet,
		Path:        "/packing-list",
		Summary:     "Get packing list",
		Description: "Returns the packing list for an SSL slip number",

		Tags:   []string{"gatepass"},
		Errors: []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, salesHandler.GetPackingList)

}
