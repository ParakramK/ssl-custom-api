package router

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func getSpec(t *testing.T) map[string]any {
	t.Helper()

	r := Setup(testHandlers())

	req, _ := http.NewRequest(http.MethodGet, "/openapi.json", nil)
	resp, err := r.Test(req)
	if err != nil {
		t.Fatalf("GET /openapi.json failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for /openapi.json, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var spec map[string]any
	if err := json.Unmarshal(body, &spec); err != nil {
		t.Fatalf("openapi.json is not valid JSON: %v", err)
	}

	return spec
}

func operation(t *testing.T, spec map[string]any, path string) map[string]any {
	t.Helper()

	paths, _ := spec["paths"].(map[string]any)
	item, ok := paths[path].(map[string]any)
	if !ok {
		t.Fatalf("spec is missing path %s", path)
	}

	op, ok := item["get"].(map[string]any)
	if !ok {
		t.Fatalf("spec path %s is missing GET operation", path)
	}

	return op
}

func responseSchemaRef(t *testing.T, op map[string]any, code string) string {
	t.Helper()

	responses, _ := op["responses"].(map[string]any)
	resp, ok := responses[code].(map[string]any)
	if !ok {
		t.Fatalf("operation is missing %s response: %v", code, responses)
	}

	content, _ := resp["content"].(map[string]any)
	for _, media := range content {
		mediaObj, _ := media.(map[string]any)
		schema, _ := mediaObj["schema"].(map[string]any)
		if ref, ok := schema["$ref"].(string); ok {
			return ref
		}
	}

	t.Fatalf("operation %s response has no $ref schema: %v", code, resp)

	return ""
}

func queryParamNames(t *testing.T, op map[string]any) map[string]bool {
	t.Helper()

	names := map[string]bool{}

	params, _ := op["parameters"].([]any)
	for _, p := range params {
		param, _ := p.(map[string]any)
		if param["in"] != "query" {
			t.Errorf("expected query parameter, got %v", param)
		}

		name, _ := param["name"].(string)
		names[name] = true
	}

	return names
}

func TestOpenAPIInfo(t *testing.T) {
	spec := getSpec(t)

	info, _ := spec["info"].(map[string]any)

	if info["title"] != "SSL Custom API" {
		t.Errorf("expected title SSL Custom API, got %v", info["title"])
	}

	if info["version"] != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %v", info["version"])
	}
}

func TestOpenAPIContainsAllEndpoints(t *testing.T) {
	spec := getSpec(t)

	for _, path := range []string{
		"/api/v1/health",
		"/api/v1/sap/customers/aging",
		"/api/v1/sap/sales/top",
		"/api/v1/gatepass/quality-report",
	} {
		operation(t, spec, path)
	}
}

func TestOpenAPISalesTop(t *testing.T) {
	spec := getSpec(t)
	op := operation(t, spec, "/api/v1/sap/sales/top")

	if op["operationId"] != "getTopOutstandingCustomers" {
		t.Errorf("unexpected operationId: %v", op["operationId"])
	}

	names := queryParamNames(t, op)
	if !names["CompanyDB"] || !names["Limit"] {
		t.Errorf("expected CompanyDB and Limit query parameters, got %v", names)
	}

	if ref := responseSchemaRef(t, op, "200"); ref != "#/components/schemas/TopCustomerResponse" {
		t.Errorf("unexpected 200 schema: %s", ref)
	}

	schemas, _ := spec["components"].(map[string]any)["schemas"].(map[string]any)
	row, _ := schemas["SalesOrderVsBillingRow"].(map[string]any)
	props, _ := row["properties"].(map[string]any)

	for field, want := range map[string]string{
		"customerCode": "string", "customerName": "string", "salesEmployee": "string",
		"orderedQuantity": "number", "billedQuantity": "number", "quantityVariance": "number",
		"orderedAmount": "number", "billedAmount": "number", "amountVariance": "number",
	} {
		prop, ok := props[field].(map[string]any)
		if !ok {
			t.Errorf("SalesOrderVsBillingRow is missing property %s", field)
			continue
		}

		if prop["type"] != want {
			t.Errorf("property %s: expected type %s, got %v", field, want, prop["type"])
		}
	}
}

func TestOpenAPICustomerAging(t *testing.T) {
	spec := getSpec(t)
	op := operation(t, spec, "/api/v1/sap/customers/aging")

	if op["operationId"] != "getCustomerAging" {
		t.Errorf("unexpected operationId: %v", op["operationId"])
	}

	names := queryParamNames(t, op)
	if !names["CardCode"] || !names["CompanyDB"] {
		t.Errorf("expected CardCode and CompanyDB query parameters, got %v", names)
	}

	if ref := responseSchemaRef(t, op, "200"); ref != "#/components/schemas/CustomerAgingDetail" {
		t.Errorf("unexpected 200 schema: %s", ref)
	}
}

func TestDocsUI(t *testing.T) {
	r := Setup(testHandlers())

	req, _ := http.NewRequest(http.MethodGet, "/docs", nil)
	resp, err := r.Test(req)
	if err != nil {
		t.Fatalf("GET /docs failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for /docs, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	if !strings.Contains(string(body), "openapi") {
		t.Error("/docs does not reference the OpenAPI document")
	}
}

func TestOpenAPIQualityReport(t *testing.T) {
	spec := getSpec(t)
	op := operation(t, spec, "/api/v1/gatepass/quality-report")

	if op["operationId"] != "getQualityReport" {
		t.Errorf("unexpected operationId: %v", op["operationId"])
	}

	names := queryParamNames(t, op)
	if !names["ssl_no"] {
		t.Errorf("expected ssl_no query parameter, got %v", names)
	}

	if ref := responseSchemaRef(t, op, "200"); ref != "#/components/schemas/QualityReportResponse" {
		t.Errorf("unexpected 200 schema: %s", ref)
	}
}

func TestOpenAPIScrapData(t *testing.T) {
	spec := getSpec(t)
	op := operation(t, spec, "/api/v1/gatepass/scrap-data")

	if op["operationId"] != "getScrapData" {
		t.Errorf("unexpected operationId: %v", op["operationId"])
	}

	names := queryParamNames(t, op)
	if !names["DbName"] || !names["DocumentNo"] {
		t.Errorf("expected DbName and DocumentNo query parameters, got %v", names)
	}

	if ref := responseSchemaRef(t, op, "200"); ref != "#/components/schemas/ScrapDataResponse" {
		t.Errorf("unexpected 200 schema: %s", ref)
	}
}
