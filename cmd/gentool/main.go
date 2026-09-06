package main

import (
	"gorm.io/gen"

	appmodels "ssl-custom-api/internal/app/models"
	"ssl-custom-api/internal/gatepass/models"
)

// go run ./cmd/gentool
func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "internal/gatepass/query",
		ModelPkgPath: "internal/gatepass/models",
	})

	g.ApplyBasic(
		models.GateEntry{},
		models.LoadingEntry{},
		models.LoadingEntryDetails{},
		models.GateLoadingSlip{},
		models.Thulokata{},
		models.ThulokataLoading{},
		models.Sanokata{},
		models.SHook{},
		models.QualityReport{},
		models.QualityReportDetails{},
		models.Scrap{},
		models.Vendors{},
		models.User{},
	)

	g.Execute()

	app := gen.NewGenerator(gen.Config{
		OutPath:      "internal/app/query",
		ModelPkgPath: "internal/app/models",
	})

	app.ApplyBasic(
		appmodels.All()...,
	)
	app.Execute()
}
