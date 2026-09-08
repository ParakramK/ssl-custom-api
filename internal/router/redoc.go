package router

import "github.com/gofiber/fiber/v3"

func setupRedoc(r *fiber.App) {
	r.Get("/redoc", func(c fiber.Ctx) error {
		return c.Type("html").SendString(`
<!doctype html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>SSL Custom API - ReDoc</title>
</head>

<body>
	<redoc spec-url="/openapi.json"></redoc>

	<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>
`)
	})
}
