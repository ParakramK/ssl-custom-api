package main

import (
	"gorm.io/gen"

	appmodels "ssl-custom-api/internal/app/models"
)

// go run ./cmd/gentool
func main() {

	app := gen.NewGenerator(gen.Config{
		OutPath:      "internal/app/query",
		ModelPkgPath: "internal/app/models",
	})

	app.ApplyBasic(
		appmodels.All()...,
	)
	app.Execute()
}
