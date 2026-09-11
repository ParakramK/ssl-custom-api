package router

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"ssl-custom-api/internal/app"
	apikey "ssl-custom-api/internal/app/api_key"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/app/permissions"
	"ssl-custom-api/internal/gatepass/scrap"
	"ssl-custom-api/internal/sap/customer"
	"ssl-custom-api/internal/sap/sales"

	"github.com/google/uuid"
)

func strPtr(v string) *string {
	return &v
}

type stubCustomerRepo struct{}

func (stubCustomerRepo) GetCustomerAging(
	ctx context.Context,
	cardCode string,
	companyDB string,
) ([]customer.AgingRow, error) {
	return nil, nil
}

type stubSalesRepo struct{}

func (stubSalesRepo) GetTopOutStandingCustomers(
	ctx context.Context,
	companyDB string,
	limit int,
) ([]sales.SalesOrderVsBillingRow, error) {
	return nil, nil
}

type stubScrapRepo struct {
	data     *scrap.QualityReportData
	err      error
	scrapErr error
}

func (s stubScrapRepo) GetQualityReport(
	ctx context.Context,
	sslNo string,
) (*scrap.QualityReportData, error) {
	return s.data, s.err
}

func testHandlers() *app.Handlers {
	testJWT, _ := auth.NewJWT("test-secret")
	return &app.Handlers{
		SAP: &app.SAPHandlers{
			Customer: customer.NewHandler(customer.NewService(stubCustomerRepo{})),
			Sales:    sales.NewHandler(sales.NewService(stubSalesRepo{})),
		},
		Gatepass: &app.GatepassHandlers{
			Scrap: scrap.NewHandler(scrap.NewService(stubScrapRepo{
				data: &scrap.QualityReportData{
					Header: scrap.QualityReportHeader{
						SslSno:       "SSL001",
						SupplierName: "Acme",
						VendorCode:   "C001",
					},
					Details: []scrap.QualityReportDetailResult{
						{ScrapType: strPtr("A"), Qty: float64Ptr(100.5), Rate: float64Ptr(10)},
					},
				},
				scrapErr: scrap.ErrQualityReportNotFound,
			})),
		},
		App: &app.AppHandlers{
			Auth: auth.NewHandler(auth.NewService(auth.NewUserRepository(nil, testJWT))),
			ApiKey: apikey.NewHandler(apikey.NewApiKeyService(
				apikey.NewApiKeyRepository(nil),
				permissions.NewAllowAllChecker(),
			)),
		},
		JWT:           testJWT,
		APIKey:        "test-key",
		Authorization: permissions.NewAllowAllChecker(),
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}

func doRequest(t *testing.T, path string) (int, map[string]any) {
	t.Helper()

	r := Setup(testHandlers())

	req, _ := http.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("X-API-Key", "test-key")
	resp, err := r.Test(req)
	if err != nil {
		t.Fatalf("GET %s failed: %v", path, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var decoded map[string]any
	_ = json.Unmarshal(body, &decoded)

	return resp.StatusCode, decoded
}

func TestApiKeysRequiresAuth(t *testing.T) {
	status, _ := doRequest(t, "/api/v1/app/keys/")

	if status != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", status)
	}
}

func TestApiKeysInvalidCursor(t *testing.T) {
	jwtSvc, err := auth.NewJWT("test-secret")
	if err != nil {
		t.Fatal(err)
	}
	token, err := jwtSvc.GenerateToken(uuid.New(), time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	r := Setup(testHandlers())

	// Malformed cursor must be rejected by the paginate middleware
	// with 400 before the handler runs (no DB needed).
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/app/keys/?cursor=not-valid!!!", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := r.Test(req)
	if err != nil {
		t.Fatalf("GET keys failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestMachineRoutesRequireAPIKey(t *testing.T) {
	r := Setup(testHandlers())

	for _, path := range []string{
		"/api/v1/sap/customers/aging?CardCode=C001&CompanyDB=DB1",
		"/api/v1/gatepass/quality-report?ssl_no=SSL001",
	} {
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		resp, err := r.Test(req)
		if err != nil {
			t.Fatalf("GET %s failed: %v", path, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("GET %s without key: expected 401, got %d", path, resp.StatusCode)
		}

		req, _ = http.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("X-API-Key", "wrong-key")
		resp, err = r.Test(req)
		if err != nil {
			t.Fatalf("GET %s failed: %v", path, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("GET %s with wrong key: expected 401, got %d", path, resp.StatusCode)
		}
	}
}

func TestHealth(t *testing.T) {
	status, body := doRequest(t, "/api/v1/health")

	if status != http.StatusOK {
		t.Errorf("expected 200, got %d", status)
	}

	if len(body) != 1 || body["status"] != "ok" {
		t.Errorf("expected exactly {status: ok}, got %v", body)
	}
}

func TestCustomerAgingValidation(t *testing.T) {
	for path, detail := range map[string]string{
		"/api/v1/sap/customers/aging":                          "CardCode is required",
		"/api/v1/sap/customers/aging?CardCode=C001":            "CompanyDB is required",
		"/api/v1/sap/customers/aging?CompanyDB=DB1":            "CardCode is required",
		"/api/v1/sap/customers/aging?CardCode=C001&CompanyDB=": "CompanyDB is required",
	} {
		status, body := doRequest(t, path)

		if status != http.StatusBadRequest {
			t.Errorf("GET %s: expected 400, got %d", path, status)
		}

		if body["detail"] != detail {
			t.Errorf("GET %s: expected detail %q, got %v", path, detail, body)
		}
	}
}

func TestSalesTopValidation(t *testing.T) {
	status, body := doRequest(t, "/api/v1/sap/sales/top")

	if status != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", status)
	}

	if body["detail"] != "CompanyDB is required" {
		t.Errorf("expected CompanyDB detail, got %v", body)
	}
}

func TestSalesTopSuccessShape(t *testing.T) {
	r := Setup(testHandlers())

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/sap/sales/top?CompanyDB=DB1", nil)
	req.Header.Set("X-API-Key", "test-key")
	resp, err := r.Test(req)
	if err != nil {
		t.Fatalf("GET sales/top failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}

	if _, ok := decoded["details"]; !ok {
		t.Errorf("expected details key, got %v", decoded)
	}
}

func TestQualityReportValidation(t *testing.T) {
	for path, detail := range map[string]string{
		"/api/v1/gatepass/quality-report":        "ssl_sno is required",
		"/api/v1/gatepass/quality-report?ssl_no": "ssl_sno is required",
	} {
		status, body := doRequest(t, path)

		if status != http.StatusBadRequest {
			t.Errorf("GET %s: expected 400, got %d", path, status)
		}

		if body["detail"] != detail {
			t.Errorf("GET %s: expected detail %q, got %v", path, detail, body)
		}
	}
}

func TestQualityReportSuccessShape(t *testing.T) {
	status, body := doRequest(t, "/api/v1/gatepass/quality-report?ssl_no=SSL001")

	if status != http.StatusOK {
		t.Fatalf("expected 200, got %d", status)
	}

	if body["success"] != true {
		t.Errorf("expected success true, got %v", body)
	}

	if body["ssl_sno"] != "SSL001" {
		t.Errorf("expected ssl_sno SSL001, got %v", body["ssl_sno"])
	}

	if body["payable_weight"] != "10.00" {
		t.Errorf("expected payable_weight 10.00, got %v", body["payable_weight"])
	}

	details, ok := body["details"].([]any)
	if !ok || len(details) != 1 {
		t.Errorf("expected 1 detail, got %v", body["details"])
	}

	if _, ok := body["payable_details"]; !ok {
		t.Errorf("expected payable_details key, got %v", body)
	}
}

func TestQualityReportNotFound(t *testing.T) {
	handlers := testHandlers()
	handlers.Gatepass.Scrap = scrap.NewHandler(scrap.NewService(stubScrapRepo{
		err: scrap.ErrQualityReportNotFound,
	}))
	r := Setup(handlers)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/gatepass/quality-report?ssl_no=SSL999", nil)
	req.Header.Set("X-API-Key", "test-key")
	resp, err := r.Test(req)
	if err != nil {
		t.Fatalf("GET quality-report failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}
