package sap

import (
	"net/http"
	"ssl-custom-api/internal/sap/sales"

	"github.com/danielgtaylor/huma/v2"
)

func SetupSalesRoutes(
	base huma.API,
	salesHandler *sales.Handler,
) {
	tags := []string{"sap:sales"}
	salesGroup := huma.NewGroup(base, "/sales")
	huma.Register(salesGroup, huma.Operation{
		OperationID: "getTopOutstandingCustomers",
		Method:      http.MethodGet,
		Path:        "/top",
		Summary:     "Get top outstanding customers",
		Description: "Returns the top customers by outstanding balance from " +
			"SAP Business One (SAP HANA), including sales order vs billing " +
			"variance details.",
		Tags:   tags,
		Errors: []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, salesHandler.GetTopOutStandingCustomers)
}
