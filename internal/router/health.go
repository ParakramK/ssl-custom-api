package router

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// HealthResponse is the health check response body.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthOutput is the typed Huma output for the health endpoint.
type HealthOutput struct {
	Body HealthResponse
}

func healthCheck(ctx context.Context, _ *struct{}) (*HealthOutput, error) {
	return &HealthOutput{Body: HealthResponse{Status: "ok"}}, nil
}

func setupHealthRoutes(base huma.API) {
	huma.Register(base, huma.Operation{
		OperationID: "getHealth",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check",
		Description: "Returns the service health status.",
		Tags:        []string{"health"},
	}, healthCheck)
}
