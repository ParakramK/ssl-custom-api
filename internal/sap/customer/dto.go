package customer

import "time"

type AgingInput struct {
	CardCode  string `query:"CardCode" doc:"SAP Business One customer code"`
	CompanyDB string `query:"CompanyDB" doc:"Company database (HANA schema) name"`
}

type AgingOutput struct {
	Body CustomerAgingDetail
}

type AgingSummary struct {
	Z0To30Days  float64 `json:"z0to30Days"`
	Z31To60Days float64 `json:"z31to60Days"`
	Z61To90Days float64 `json:"z61to90Days"`
	Z91PlusDays float64 `json:"z91PlusDays"`
}

type PaymentSummary struct {
	BG    float64 `json:"BG"`
	LC    float64 `json:"LC"`
	Other float64 `json:"Other"`
}

type TotalSummary struct {
	AccountBalance      float64 `json:"AccountBalance"`
	UnreconciledBalance float64 `json:"UnreconciledBalance"`
}

type AgingDetail struct {
	InvoiceNumber     int        `json:"InvoiceNumber"`
	InvoiceDate       *time.Time `json:"InvoiceDate"`
	InvoiceAmount     float64    `json:"InvoiceAmount"`
	DueDate           *time.Time `json:"DueDate"`
	PaymentTerms      string     `json:"PaymentTerms"`
	PaymentTermGroup  string     `json:"PaymentTermGroup"`
	OutstandingDays   int        `json:"OutstandingDays"`
	PaidAmount        float64    `json:"PaidAmount"`
	OutstandingAmount float64    `json:"OutstandingAmount"`
	Z0To30Days        float64    `json:"z0to30Days"`
	Z31To60Days       float64    `json:"z31to60Days"`
	Z61To90Days       float64    `json:"z61to90Days"`
	Z91PlusDays       float64    `json:"z91PlusDays"`
}

type CustomerAgingDetail struct {
	Status         string         `json:"Status"`
	CompanyDB      string         `json:"CompanyDB"`
	CardCode       string         `json:"CardCode"`
	CardName       string         `json:"CardName"`
	Balance        float64        `json:"Balance"`
	TaxNumber      string         `json:"TaxNumber"`
	SalesEmployee  string         `json:"SalesEmployee"`
	LastFetched    time.Time      `json:"LastFetched"`
	Details        []AgingDetail  `json:"details"`
	AgingSummary   AgingSummary   `json:"aging_summary"`
	PaymentSummary PaymentSummary `json:"payment_summary"`
	Total          TotalSummary   `json:"total"`
	Cached         bool           `json:"cached"`
}
