package gatepass

import (
	"ssl-custom-api/internal/app"

	"github.com/danielgtaylor/huma/v2"
)

func SetupGatepassRoutes(base huma.API, handlers *app.Handlers) {
	gatepassRoutes := huma.NewGroup(base, "/gatepass")
	SetupSalesRoutes(gatepassRoutes, handlers.Gatepass.Sales)
	SetupScrapRoutes(gatepassRoutes, handlers.Gatepass.Scrap)
}
