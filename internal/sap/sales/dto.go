package sales

type TopCustomersInput struct {
	CompanyDB string `query:"CompanyDB" doc:"Company database (HANA schema) name"`
	Limit     int    `query:"Limit" doc:"Number of top customers to return"`
}

type TopCustomersOutput struct {
	Body TopCustomerResponse
}

type TopCustomerResponse struct {
	Details []SalesOrderVsBillingRow `json:"details"`
	Cached  bool                     `json:"cached"`
}
