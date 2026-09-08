package router

import (
	"ssl-custom-api/internal/app"
	"ssl-custom-api/internal/httpx"

	appRouter "ssl-custom-api/internal/router/app"
	"ssl-custom-api/internal/router/gatepass"
	"ssl-custom-api/internal/router/sap"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/paginate"
)

func Setup(handlers *app.Handlers) *fiber.App {
	r := fiber.New(fiber.Config{
		ErrorHandler:        httpx.ErrorHandler,
		PassLocalsToContext: true,
	})

	setupSwagger(r)
	setupScalar(r)
	setupRedoc(r)

	config := huma.DefaultConfig("SSL Custom API", "1.0.0")
	config.OpenAPI.Info.Description = "Custom reporting API over SAP Business One data on SAP HANA: " +
		"customer accounts receivable aging and sales order vs billing dashboards."
	// Disable the $schema link transformer so response bodies keep their
	// exact existing JSON structures.
	config.CreateHooks = nil

	api := humafiber.New(r, config)
	v1 := huma.NewGroup(api, "/api/v1")
	sapRoutes := huma.NewGroup(v1, "/sap")
	gatepassRoutes := huma.NewGroup(v1, "/gatepass")
	appRoutes := huma.NewGroup(v1, "/app")

	r.Use("/api/v1/sap", apiKeyAuth(handlers.APIKey))
	r.Use("/api/v1/gatepass", apiKeyAuth(handlers.APIKey))

	r.Use("/api/v1/app", paginate.New(paginate.Config{
		SortKey:      "sort",
		DefaultSort:  "id",
		AllowedSorts: []string{"id", "key"},
	}))

	sap.SetupCustomerRoutes(sapRoutes, handlers.SAP.Customer)
	sap.SetupSalesRoutes(sapRoutes, handlers.SAP.Sales)
	gatepass.SetupScrapRoutes(gatepassRoutes, handlers.Gatepass.Scrap)
	gatepass.SetupSalesRoutes(gatepassRoutes, handlers.Gatepass.Sales)
	appRouter.SetupAuthRoutes(appRoutes, handlers.App.Auth)
	appRouter.SetupApiKeysRoutes(appRoutes, handlers.App.ApiKey, handlers.JWT)
	setupHealthRoutes(v1)

	return r
}
