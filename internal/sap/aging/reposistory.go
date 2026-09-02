package aging

import (
	"context"
	"database/sql"
	"fmt"
	"ssl-custom-api/internal/providers/hana"

	"github.com/SAP/go-hdb/driver"
)

type Repository struct {
	hana *hana.Provider
}

func NewRepository(hanaProvider *hana.Provider) *Repository {
	return &Repository{
		hana: hanaProvider,
	}
}

func (r *Repository) GetCustomerAging(
	ctx context.Context,
	cardCode string,
	companyDB string,
) ([]AgingRow, error) {

	var result []AgingRow

	err := r.hana.WithSchema(ctx, companyDB, func(conn *sql.Conn) error {

		rows, err := conn.QueryContext(
			ctx,
			CustomerAgingQuery,
			cardCode,
		)
		if err != nil {
			return fmt.Errorf("get customer aging: %w", err)
		}
		defer rows.Close()

		for rows.Next() {

			var (
				balance           driver.Decimal
				invoiceAmount     driver.Decimal
				paidAmount        driver.Decimal
				outstandingAmount driver.Decimal
				z0To30Days        driver.Decimal
				z31To60Days       driver.Decimal
				z61To90Days       driver.Decimal
				z91PlusDays       driver.Decimal
			)

			var row AgingRow

			err := rows.Scan(
				&row.DocEntry,
				&row.InvoiceNumber,
				&row.InvoiceDate,
				&row.CustomerCode,
				&row.CustomerName,
				&balance,
				&row.TaxNumber,
				&row.DueDate,
				&row.PaymentTerms,
				&row.SalesEmployee,
				&row.PaymentTermGroup,
				&row.OutstandingDays,
				&invoiceAmount,
				&paidAmount,
				&outstandingAmount,
				&z0To30Days,
				&z31To60Days,
				&z61To90Days,
				&z91PlusDays,
			)
			if err != nil {
				return fmt.Errorf("scan customer aging: %w", err)
			}

			row.Balance = decimalToFloat64(balance)
			row.InvoiceAmount = decimalToFloat64(invoiceAmount)
			row.PaidAmount = decimalToFloat64(paidAmount)
			row.OutstandingAmount = decimalToFloat64(outstandingAmount)
			row.Z0To30Days = decimalToFloat64(z0To30Days)
			row.Z31To60Days = decimalToFloat64(z31To60Days)
			row.Z61To90Days = decimalToFloat64(z61To90Days)
			row.Z91PlusDays = decimalToFloat64(z91PlusDays)
			result = append(result, row)

		}

		if err := rows.Err(); err != nil {
			return fmt.Errorf("iterate customer aging: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
