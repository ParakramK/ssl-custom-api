package router

import "github.com/gofiber/fiber/v3"

func setupScalar(r *fiber.App) {
	r.Get("/scalar", func(c fiber.Ctx) error {
		return c.Type("html").SendString(`
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>SSL Custom API Documentation</title>
   <style>
        /* Hide the "Powered by Scalar" link container */
      
        a[href*="scalar.com"] {
            display: none !important;
        }
    </style>
</head>

<body>
    <script
        id="api-reference"
        data-url="/openapi.json"
        data-configuration='{
            "hideModels": false,
            "hideDownloadButton": true,
            "showSidebar": true,
            "darkMode": true
        }'
    ></script>

    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>
`)
	})
}
