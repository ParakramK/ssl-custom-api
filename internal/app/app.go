package app

import (
	gatepassSales "ssl-custom-api/internal/gatepass/sales"
	"ssl-custom-api/internal/gatepass/scrap"
	"ssl-custom-api/internal/providers/hana"
	"ssl-custom-api/internal/providers/mysql"
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
}

func New(hanaProvider *hana.Provider, mysqlProvider *mysql.Provider) *Handlers {
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

	return &Handlers{
		SAP: &SAPHandlers{
			Customer: customerHandler,
			Sales:    salesHandler,
		},
		Gatepass: &GatepassHandlers{
			Scrap: scrapHandler,
			Sales: gatepassSalesHandler,
		},
	}
}
