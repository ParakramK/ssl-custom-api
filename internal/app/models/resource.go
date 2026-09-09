package models

import "strings"

// ResourceCode identifies a single API resource surface.
// Always fully qualified (<module>.<resource>): "sales" alone is
// ambiguous because both Gatepass and SAP contain one.
//
// The database UUID is database identity; this code is the stable
// application identifier used in authorization.
type ResourceCode string

// AllResources lists every known resource code.
func AllResources() []ResourceCode {
	return []ResourceCode{
		ResourceGatepassSales,
		ResourceGatepassScrap,
		ResourceSAPCustomer,
		ResourceSAPSales,
	}
}

// Valid reports whether the code is a known resource.
func (r ResourceCode) Valid() bool {
	for _, known := range AllResources() {
		if r == known {
			return true
		}
	}
	return false
}

// ModuleCode returns the parent module code ("gatepass.sales" -> "gatepass").
func (r ResourceCode) ModuleCode() string {
	if i := strings.Index(string(r), "."); i >= 0 {
		return string(r)[:i]
	}
	return string(r)
}
