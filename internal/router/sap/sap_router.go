package sap

import (
	"ssl-custom-api/internal/app"

	"github.com/danielgtaylor/huma/v2"
)

func SetupSAPRoutes(base huma.API, handlers *app.Handlers) {
	gatepassRoutes := huma.NewGroup(base, "/sap")
	SetupSalesRoutes(gatepassRoutes, handlers.SAP.Sales)
	SetupCustomerRoutes(gatepassRoutes, handlers.SAP.Customer)
}
