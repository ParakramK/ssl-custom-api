package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestErrorHandlerFiberError(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Get("/missing", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "Not Found")
	})

	req, _ := http.NewRequest(http.MethodGet, "/missing", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}

	if decoded["error"] != "Not Found" {
		t.Errorf("expected error body, got %v", decoded)
	}
}

func TestErrorHandlerGenericError(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Get("/panic", func(c fiber.Ctx) error {
		return errors.New("raw failure")
	})

	req, _ := http.NewRequest(http.MethodGet, "/panic", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}

	if decoded["error"] == "" || decoded["error"] == "raw failure" {
		t.Errorf("expected generic client-safe message, got %v", decoded)
	}
}
