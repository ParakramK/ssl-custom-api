package sales

import (
	"context"
	"database/sql"
	"fmt"

	"ssl-custom-api/internal/providers/hana"
	"ssl-custom-api/internal/utils"

	"github.com/SAP/go-hdb/driver"
)

type Repository interface {
	GetTopOutStandingCustomers(ctx context.Context, companyDB string, limit int) ([]SalesOrderVsBillingRow, error)
}

type salesRepository struct {
	hana *hana.Provider
}

func NewRepository(hanaProvider *hana.Provider) Repository {
	return &salesRepository{
		hana: hanaProvider,
	}
}

func (r *salesRepository) GetTopOutStandingCustomers(
	ctx context.Context,
	companyDB string,
	limit int,
) ([]SalesOrderVsBillingRow, error) {

	var result []SalesOrderVsBillingRow

	// HANA does not accept bind placeholders in LIMIT, so the already
	// server-clamped limit is interpolated as an integer literal.
	query := fmt.Sprintf(SALES_ORDER_VS_BILLING_QUERY, limit)

	err := r.hana.WithSchema(ctx, companyDB, func(conn *sql.Conn) error {

		rows, err := conn.QueryContext(
			ctx,
			query,
		)
		if err != nil {
			return fmt.Errorf("get sales order vs billing: %w", err)
		}
		defer rows.Close()

		for rows.Next() {

			var (
				orderQuantity    driver.Decimal
				billedQuantity   driver.Decimal
				orderAmount      driver.Decimal
				billedAmount     driver.Decimal
				quantityVariance driver.Decimal
				amountVariance   driver.Decimal
			)

			var row SalesOrderVsBillingRow

			err := rows.Scan(
				&row.CustomerCode,
				&row.CustomerName,
				&row.SalesEmployee,
				&orderQuantity,
				&billedQuantity,
				&quantityVariance,
				&orderAmount,
				&billedAmount,
				&amountVariance,
			)
			if err != nil {
				return fmt.Errorf("scan customer order vs billing: %w", err)
			}
			row.OrderedQuantity = utils.DecimalToFloat64(orderQuantity)
			row.BilledQuantity = utils.DecimalToFloat64(billedQuantity)
			row.OrderedAmount = utils.DecimalToFloat64(orderAmount)
			row.BilledAmount = utils.DecimalToFloat64(billedAmount)
			row.QuantityVariance = utils.DecimalToFloat64(quantityVariance)
			row.AmountVariance = utils.DecimalToFloat64(amountVariance)

			result = append(result, row)

		}

		if err := rows.Err(); err != nil {
			return fmt.Errorf("iterate customer order vs billing: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
