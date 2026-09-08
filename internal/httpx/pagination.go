package httpx

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3/middleware/paginate"
)

func Pagination(ctx context.Context) (*paginate.PageInfo, error) {
	pi, ok := paginate.FromContext(ctx)
	if !ok || pi == nil {
		return nil, fmt.Errorf("pagination middleware is not mounted")
	}
	return pi, nil
}
