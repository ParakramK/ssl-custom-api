package httpx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/paginate"
)

type pingInput struct{}

type pingOutput struct {
	Body pingBody
}

type pingBody struct {
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

func TestPaginationReachesHumaHandler(t *testing.T) {
	app := fiber.New(fiber.Config{PassLocalsToContext: true})
	app.Use(paginate.New(paginate.Config{
		SortKey:      "sort",
		DefaultSort:  "id",
		AllowedSorts: []string{"id", "key"},
	}))
	api := humafiber.New(app, huma.DefaultConfig("test", "1.0.0"))
	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/ping",
	}, func(ctx context.Context, _ *pingInput) (*pingOutput, error) {
		pi, err := Pagination(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		sort := ""
		if len(pi.Sort) > 0 {
			sort = pi.Sort[0].Field
		}
		return &pingOutput{Body: pingBody{Limit: pi.Limit, Sort: sort}}, nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/ping?limit=5&sort=-key", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var decoded pingBody
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Limit != 5 || decoded.Sort != "key" {
		t.Errorf("PageInfo did not propagate: %+v", decoded)
	}
}
