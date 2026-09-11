package router

import "github.com/danielgtaylor/huma/v2"

var tagDisplayNames = map[string]string{
	"app:auth":       "Authentication",
	"app:modules":    "Modules",
	"app:keys":       "API Keys",
	"app:roles":      "Roles",
	"app:users":      "Users",
	"gatepass:sales": "Sales",
	"gatepass:scrap": "Scrap",
	"sap:customer":   "Customer",
	"sap:sales":      "Sales",
	"health":         "Health",
}

func configureOpenAPI(config *huma.Config) {
	config.OpenAPI.Info.Description =
		"Custom API For Integrating Gatepass and SAP Data and other Internal Systems."

	config.CreateHooks = nil
}

func configureOpenAPITags(api huma.API) {
	openapi := api.OpenAPI()

	openapi.Extensions = map[string]any{
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

	for _, tag := range openapi.Tags {
		if displayName, ok := tagDisplayNames[tag.Name]; ok {
			if tag.Extensions == nil {
				tag.Extensions = map[string]any{}
			}

			tag.Extensions["x-displayName"] = displayName
		}
	}
}
