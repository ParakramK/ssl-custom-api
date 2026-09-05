package sales

type SalesOrderVsBillingRow struct {
	CustomerCode     string  `json:"customerCode"`
	CustomerName     string  `json:"customerName"`
	SalesEmployee    string  `json:"salesEmployee"`
	OrderedQuantity  float64 `json:"orderedQuantity"`
	BilledQuantity   float64 `json:"billedQuantity"`
	QuantityVariance float64 `json:"quantityVariance"`
	OrderedAmount    float64 `json:"orderedAmount"`
	BilledAmount     float64 `json:"billedAmount"`
	AmountVariance   float64 `json:"amountVariance"`
}
