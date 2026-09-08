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
			SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui"
			})
		}
	</script>
</body>
</html>
`)
	})
}
