package customer

import (
	"time"
)

type AgingRow struct {
	DocEntry          int
	InvoiceNumber     int
	InvoiceDate       *time.Time
	CustomerCode      string
	CustomerName      string
	Balance           float64
	TaxNumber         string
	DueDate           *time.Time
	PaymentTerms      string
	SalesEmployee     string
	PaymentTermGroup  string
	OutstandingDays   int
	InvoiceAmount     float64
	PaidAmount        float64
	OutstandingAmount float64
	Z0To30Days        float64
	Z31To60Days       float64
	Z61To90Days       float64
	Z91PlusDays       float64
}
