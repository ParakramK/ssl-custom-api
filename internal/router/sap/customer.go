package sap

import (
	"net/http"

	"ssl-custom-api/internal/sap/customer"

	"github.com/danielgtaylor/huma/v2"
)

func SetupCustomerRoutes(
	base huma.API,
	customerHandler *customer.Handler,
) {
	tags := []string{"sap:customer"}
	customerGroup := huma.NewGroup(base, "/customers")
	huma.Register(customerGroup, huma.Operation{
		OperationID: "getCustomerAging",
		Method:      http.MethodGet,
		Path:        "/aging",
		Summary:     "Get customer aging report",
		Description: "Returns the accounts receivable aging report for a customer " +
			"from SAP Business One (SAP HANA), including invoice details, aging " +
			"buckets, payment-term summary and totals.",
		Tags:   tags,
		Errors: []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, customerHandler.GetCustomerAging)
}
