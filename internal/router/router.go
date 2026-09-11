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

	config.OpenAPI.Info.Description = "Custom API For Integrating Gatepass and SAP Data and other Internal Systems."

	config.CreateHooks = nil

	config.Extensions = map[string]any{
		"x-tagGroups": []map[string]any{
			{
				"name": "Application",
				"tags": []string{
					"app:auth",
					"app:modules",
					"app:keys",
					"app:roles",
					"app:users",
				},
			},
			{
				"name": "Gatepass",
				"tags": []string{
					"gatepass:sales",
					"gatepass:scrap",
				},
			},
			{
				"name": "SAP",
				"tags": []string{
					"sap:customer",
					"sap:sales",
				},
			},
			{
				"name": "System",
				"tags": []string{
					"health",
				},
			},
		},
	}

	api := humafiber.New(r, config)
	v1 := huma.NewGroup(api, "/api/v1")

	// r.Use("/api/v1/sap", apiKeyAuth(handlers.APIKey))
	// r.Use("/api/v1/gatepass", apiKeyAuth(handlers.APIKey))

	r.Use("/api/v1/app", paginate.New(paginate.Config{
		SortKey:      "sort",
		DefaultSort:  "id",
		AllowedSorts: []string{"id", "key"},
	}))

	sap.SetupSAPRoutes(v1, handlers)
	gatepass.SetupGatepassRoutes(v1, handlers)
	appRouter.SetupAppRoutes(v1, handlers)

	setupHealthRoutes(v1)

	return r
}
