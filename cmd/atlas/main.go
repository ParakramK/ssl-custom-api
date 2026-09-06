package main

// Atlas external-schema loader for the app-owned PostgreSQL tables.
// Models come from models.All(), the same source the gorm/gen
// generator uses, so the two can never drift apart. The legacy
// MySQL gatepass schema belongs to another app and is never
// managed here.
//
// Run from the repo root:
//
//	go run ./cmd/atlas

import (
	"fmt"
	"os"

	appmodels "ssl-custom-api/internal/app/models"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		appmodels.All()...,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(stmts)
}
