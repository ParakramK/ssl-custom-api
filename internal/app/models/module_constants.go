package models

// Stable application identifiers for Control Plane modules and API
// resources. Codes are immutable: the UI may rename a module without
// affecting authorization.

// Module codes (Control Plane product level).
const (
	ModuleCodeGatepass = "gatepass"
	ModuleCodeSAP      = "sap"
)

// Resource codes (API resource level). Always fully qualified:
// "sales" alone is ambiguous because both Gatepass and SAP contain one.
const (
	ResourceGatepassSales ResourceCode = "gatepass.sales"
	ResourceGatepassScrap ResourceCode = "gatepass.scrap"

	ResourceSAPCustomer ResourceCode = "sap.customer"
	ResourceSAPSales    ResourceCode = "sap.sales"
)
