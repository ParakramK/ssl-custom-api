package app

import (
	apikey "ssl-custom-api/internal/app/api_key"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/authorization"
	"ssl-custom-api/internal/app/modules"
	"ssl-custom-api/internal/app/permissions"
	"ssl-custom-api/internal/app/roles"
	"ssl-custom-api/internal/app/users"
	gatepassSales "ssl-custom-api/internal/gatepass/sales"
	"ssl-custom-api/internal/gatepass/scrap"
	"ssl-custom-api/internal/providers/hana"
	"ssl-custom-api/internal/providers/mysql"
	"ssl-custom-api/internal/providers/postgres"
	"ssl-custom-api/internal/sap/customer"
	"ssl-custom-api/internal/sap/sales"

	"gorm.io/gorm"
)

type SAPHandlers struct {
	Customer *customer.Handler
	Sales    *sales.Handler
}

type GatepassHandlers struct {
	Scrap *scrap.Handler
	Sales *gatepassSales.Handler
}

type Handlers struct {
	SAP           *SAPHandlers
	Gatepass      *GatepassHandlers
	App           *AppHandlers
	JWT           *auth.JWT
	APIKey        string
	Authorization authorization.Service
}

type AppHandlers struct {
	Auth   *auth.Handler
	ApiKey *apikey.Handler
	Module *modules.Handler
	Role   *roles.Handler
	Users  *users.Handler
}

func New(
	hanaProvider *hana.Provider,
	mysqlProvider *mysql.Provider,
	postgresProvider *postgres.Provider,
	jwt *auth.JWT,
	apiKey string,
) *Handlers {
	postgresDB := postgresProvider.DB()

	return &Handlers{
		SAP:      newSAPHandlers(hanaProvider),
		Gatepass: newGatepassHandlers(mysqlProvider.DB()),
		App:      newAppHandlers(postgresDB, jwt),
		JWT:      jwt,
		APIKey:   apiKey,
		Authorization: authorization.NewService(
			authorization.NewAuthorizationRepository(postgresDB),
		),
	}
}

func newSAPHandlers(provider *hana.Provider) *SAPHandlers {
	return &SAPHandlers{
		Customer: customer.New(provider),
		Sales:    sales.New(provider),
	}
}

func newGatepassHandlers(db *gorm.DB) *GatepassHandlers {
	return &GatepassHandlers{
		Scrap: scrap.New(db),
		Sales: gatepassSales.New(db),
	}
}

func newAppHandlers(db *gorm.DB, jwt *auth.JWT) *AppHandlers {
	permissionChecker := permissions.NewDBChecker(db)

	return &AppHandlers{
		Auth:   auth.New(db, jwt),
		ApiKey: apikey.New(db, permissionChecker),
		Module: modules.New(db),
		Role:   roles.New(db, jwt),
		Users:  users.New(db),
	}
}
