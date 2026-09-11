package gatepass

import (
	"net/http"

	"ssl-custom-api/internal/gatepass/sales"

	"github.com/danielgtaylor/huma/v2"
)

func SetupSalesRoutes(
	base huma.API,
	salesHandler *sales.Handler,
) {
	tags := []string{"gatepass:sales"}
	huma.Register(base, huma.Operation{
		OperationID: "getTagData",
		Method:      http.MethodGet,
		Path:        "/tag",
		Summary:     "Get tag data",
		Description: "Returns the tag data for a Tag Bundle No",

		Tags:   tags,
		Errors: []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, salesHandler.GetTagData)
	huma.Register(base, huma.Operation{
		OperationID: "getPackingList",
		Method:      http.MethodGet,
		Path:        "/packing-list",
		Summary:     "Get packing list",
		Description: "Returns the packing list for an SSL entry number",

		Tags:   tags,
		Errors: []int{http.StatusBadRequest, http.StatusNotFound, http.StatusInternalServerError},
	}, salesHandler.GetPackingList)

}
