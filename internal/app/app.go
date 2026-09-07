package app

import (
	"ssl-custom-api/internal/app/auth"
	gatepassSales "ssl-custom-api/internal/gatepass/sales"
	"ssl-custom-api/internal/gatepass/scrap"
	"ssl-custom-api/internal/providers/hana"
	"ssl-custom-api/internal/providers/mysql"
	"ssl-custom-api/internal/providers/postgres"
	"ssl-custom-api/internal/sap/customer"
	"ssl-custom-api/internal/sap/sales"
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
	SAP      *SAPHandlers
	Gatepass *GatepassHandlers
	App      *AppHandlers
}

type AppHandlers struct {
	Auth *auth.Handler
}

func New(hanaProvider *hana.Provider,
	mysqlProvider *mysql.Provider,
	postgresProvider *postgres.Provider,
	jwt *auth.JWT,
) *Handlers {
	customerRepo := customer.NewRepository(hanaProvider)
	customerService := customer.NewService(customerRepo)
	customerHandler := customer.NewHandler(customerService)

	salesRepo := sales.NewRepository(hanaProvider)
	salesService := sales.NewService(salesRepo)
	salesHandler := sales.NewHandler(salesService)

	scrapRepo := scrap.NewRepository(mysqlProvider.DB())
	scrapService := scrap.NewService(scrapRepo)
	scrapHandler := scrap.NewHandler(scrapService)

	gatepassSalesRepo := gatepassSales.NewRepository(mysqlProvider.DB())
	gatepassSalesService := gatepassSales.NewService(gatepassSalesRepo)
	gatepassSalesHandler := gatepassSales.NewHandler(gatepassSalesService)

	appUserRepo := auth.NewUserRepository(postgresProvider.DB(), jwt)
	appUserService := auth.NewService(appUserRepo)
	appUserHandler := auth.NewHandler(appUserService)

	return &Handlers{
		SAP: &SAPHandlers{
			Customer: customerHandler,
			Sales:    salesHandler,
		},
		Gatepass: &GatepassHandlers{
			Scrap: scrapHandler,
			Sales: gatepassSalesHandler,
		},
		App: &AppHandlers{
			Auth: appUserHandler,
		},
	}
}
