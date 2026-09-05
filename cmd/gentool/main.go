package main

import (
	"gorm.io/gen"

	"ssl-custom-api/internal/gatepass/models"
)

// Generates type-safe query code for the quality-report models into
// internal/gatepass/query. Run from the repo root:
//
//	go run ./cmd/gentool
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
}
