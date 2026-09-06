package sap

import (
	"net/http"

	"ssl-custom-api/internal/sap/customer"
	"ssl-custom-api/internal/sap/sales"

	"github.com/danielgtaylor/huma/v2"
)

func SetupCustomerRoutes(
	base huma.API,
	customerHandler *customer.Handler,
) {
	customerGroup := huma.NewGroup(base, "/customers")
	huma.Register(customerGroup, huma.Operation{
		OperationID: "getCustomerAging",
		Method:      http.MethodGet,
		Path:        "/aging",
		Summary:     "Get customer aging report",
		Description: "Returns the accounts receivable aging report for a customer " +
			"from SAP Business One (SAP HANA), including invoice details, aging " +
			"buckets, payment-term summary and totals.",
		Tags:   []string{"sap"},
		Errors: []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, customerHandler.GetCustomerAging)
}

func SetupSalesRoutes(
	base huma.API,
	salesHandler *sales.Handler,
) {
	salesGroup := huma.NewGroup(base, "/sales")
	huma.Register(salesGroup, huma.Operation{
		OperationID: "getTopOutstandingCustomers",
		Method:      http.MethodGet,
		Path:        "/top",
		Summary:     "Get top outstanding customers",
		Description: "Returns the top customers by outstanding balance from " +
			"SAP Business One (SAP HANA), including sales order vs billing " +
			"variance details.",
		Tags:   []string{"sap"},
		Errors: []int{http.StatusBadRequest, http.StatusInternalServerError},
	}, salesHandler.GetTopOutStandingCustomers)
}
