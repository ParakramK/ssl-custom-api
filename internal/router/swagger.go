package router

import "github.com/gofiber/fiber/v3"

func setupSwagger(r *fiber.App) {
	r.Get("/swagger", func(c fiber.Ctx) error {
		return c.Type("html").SendString(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>SSL Custom API - Swagger UI</title>

	<link
		rel="stylesheet"
		href="https://unpkg.com/swagger-ui-dist/swagger-ui.css"
	>
</head>

<body>
	<div id="swagger-ui"></div>

	<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>

	<script>
		window.onload = () => {
			const tagOrder = [
				"health",
				"app:auth",
				"app:modules",
				"app:keys",
				"gatepass:sales",
				"gatepass:scrap",
				"sap:customer",
				"sap:sales",
			]

			SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui",

				tagsSorter: (a, b) => {
					const ai = tagOrder.indexOf(a)
					const bi = tagOrder.indexOf(b)

					// Known tags follow custom order.
					// Unknown tags go after them and remain alphabetical.
					if (ai !== -1 && bi !== -1) {
						return ai - bi
					}

					if (ai !== -1) return -1
					if (bi !== -1) return 1

					return a.localeCompare(b)
				}
			})
		}
	</script>
</body>
</html>
`)
	})
}
