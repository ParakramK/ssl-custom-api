package sales

import (
	"context"
	"errors"
	"fmt"
	"math"

	"ssl-custom-api/internal/gatepass/query"
	"ssl-custom-api/internal/utils"

	"gorm.io/gorm"
)

var ErrTagNotFound = errors.New("tag not found")

var ErrPackingListNotFound = errors.New("packing list not found")

type Repository interface {
	GetTagData(ctx context.Context, sslNo string) (*BundleDetails, error)
	GetPackingList(ctx context.Context, sslNo string) (*PackingList, error)
}

type salesRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &salesRepository{db: db}
}

func (r *salesRepository) GetTagData(
	ctx context.Context,
	tagNo string,
) (*BundleDetails, error) {
	if tagNo == "" {
		return nil, errors.New("tag number is required")
	}

	q := query.Use(r.db)

	var row TagDataRow

	err := q.Sanokata.
		WithContext(ctx).
		Select(
			q.Sanokata.Code,
			q.Sanokata.Lot,
			q.Sanokata.SizeMM,
			q.Sanokata.Bundles,
			q.Sanokata.Pieces,
			q.Sanokata.NetWeight,
			q.Sanokata.CreatedAt,
			q.Sanokata.LoadingDate.As("LoadingDate"),

			q.GateEntry.DocumentNo.As("DocumentNo"),
			q.GateEntry.LoadingSlipNo,
			q.GateEntry.VehicleNo.As("VehicleNo"),

			q.LoadingEntry.LoadingSlipNo,
			q.LoadingEntry.SalesOrder,
			q.LoadingEntry.PartyName,
			q.LoadingEntry.ShippingAddress,
		).
		Join(
			q.GateEntry,
			q.Sanokata.DocumentNo.EqCol(
				q.GateEntry.DocumentNo,
			),
		).
		Join(
			q.LoadingEntry,
			q.GateEntry.LoadingSlipNo.EqCol(
				q.LoadingEntry.LoadingSlipNo,
			),
		).
		Where(
			q.Sanokata.Code.Eq(tagNo),
		).
		Order(
			q.Sanokata.CreatedAt.Desc(),
		).
		Scan(&row)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}

		return nil, fmt.Errorf("get tag data: %w", err)
	}

	// Scan reports no rows as a nil error with a zero row, so
	// detect the miss via the filtered column.
	if row.Code == "" {
		return nil, ErrTagNotFound
	}

	bundle := TagDetails(row)

	return &bundle, nil
}

func (r *salesRepository) GetPackingList(
	ctx context.Context,
	sslNo string,
) (*PackingList, error) {
	if sslNo == "" {
		return nil, errors.New("ssl number is required")
	}
	q := query.Use(r.db)

	var thulokataRows []ThulokataRow
	if err := q.ThulokataLoading.
		WithContext(ctx).
		Select(
			q.ThulokataLoading.Materials, q.ThulokataLoading.InitialWeight, q.ThulokataLoading.FinalWeight,
			q.ThulokataLoading.NetWeight, q.ThulokataLoading.Bundles,
			q.Thulokata.PartyName.As("PartyName"), q.Thulokata.DocumentNo.As("EntryNo"),
			q.LoadingEntry.ShippingAddress.As("ShippingAddress"), q.LoadingEntry.SalesOrder.As("SalesOrder"),
			q.LoadingEntry.VehicleNo.As("VehicleNumber"),
		).
		Join(
			q.Thulokata,
			q.ThulokataLoading.ThuloKataEntryNo.EqCol(
				q.Thulokata.Sn,
			),
		).
		Join(
			q.GateEntry,
			q.Thulokata.DocumentNo.EqCol(
				q.GateEntry.DocumentNo,
			),
		).
		LeftJoin(
			q.LoadingEntry,
			q.GateEntry.LoadingSlipNo.EqCol(
				q.LoadingEntry.LoadingSlipNo,
			),
		).
		Where(q.GateEntry.DocumentNo.Eq(sslNo)).
		Scan(&thulokataRows); err != nil {
		return nil, fmt.Errorf("get thulokata packing lines: %w", err)
	}

	// Sanokata loadings: bundles loaded one by one, ordered by
	// size then bundles. User, gate entry and loading entry are
	// supplementary, so they stay left joins.
	var sanokataRows []SanokataRow
	if err := q.Sanokata.
		WithContext(ctx).
		Select(
			q.Sanokata.Code, q.Sanokata.SizeMM, q.Sanokata.Bundles,
			q.Sanokata.Pieces, q.Sanokata.NetWeight, q.Sanokata.Lot,
			q.Sanokata.LoadingDate.As("LoadingDate"), q.Sanokata.DocumentNo.As("DocumentNo"),
			q.User.KataNo.As("KataNo"), q.GateEntry.VehicleNo.As("VehicleNumber"),
			q.LoadingEntry.SalesOrder.As("SalesOrder"),
			q.LoadingEntry.PartyName.As("PartyName"),
			q.LoadingEntry.ShippingAddress.As("ShippingAddress"),
		).
		LeftJoin(
			q.User, q.Sanokata.UserId.EqCol(q.User.Id),
		).
		LeftJoin(
			q.GateEntry,
			q.Sanokata.DocumentNo.EqCol(
				q.GateEntry.DocumentNo,
			),
		).
		LeftJoin(
			q.LoadingEntry,
			q.GateEntry.LoadingSlipNo.EqCol(
				q.LoadingEntry.LoadingSlipNo,
			),
		).
		Where(q.Sanokata.DocumentNo.Eq(sslNo)).
		Order(
			q.Sanokata.SizeMM.Asc(), q.Sanokata.Bundles.Asc(),
		).
		Scan(&sanokataRows); err != nil {
		return nil, fmt.Errorf("get sanokata packing lines: %w", err)
	}

	var shookRows []SHookLine
	if err := q.SHook.
		WithContext(ctx).
		Select(
			q.SHook.DocumentNo.As("DocumentNo"), q.SHook.Size, q.SHook.Quantity,
		).
		Where(
			q.SHook.DocumentNo.Eq(sslNo),
		).
		Scan(&shookRows); err != nil {
		return nil, fmt.Errorf("get shook lines: %w", err)
	}
	sizeBundles := make(map[string]int)
	totalWeightForSize := make(map[string]float64)
	list := &PackingList{
		Header:          PackingListHeader{},
		ThulokataLines:  []ThulokataPackingRow{},
		SanokataLines:   make([]SanokataPackingRow, 0, len(sanokataRows)),
		ShookLines:      make([]SHookRow, 0, len(shookRows)),
		SanokataSummary: make([]SanokataSummaryRow, 0, len(sanokataRows)),
	}
	for _, thl := range thulokataRows {
		list.ThulokataLines = append(list.ThulokataLines, packThulokataRow(thl))
	}
	for _, so := range sanokataRows {
		list.SanokataLines = append(list.SanokataLines, packSanokataRow(so))
		size := so.SizeMM

		sizeBundles[size] += so.Bundles
		totalWeightForSize[size] += so.NetWeight
	}
	for _, shook := range shookRows {
		list.ShookLines = append(list.ShookLines, packShookRow(shook))
	}
	for _, thl := range thulokataRows {
		list.Header.EntryNo = thl.EntryNo
		list.Header.PartyName = thl.PartyName
		list.Header.SalesOrder = utils.StrVal(thl.SalesOrder)
		list.Header.VehicleNumber = utils.StrVal(thl.VehicleNumber)
		break
	}
	for _, so := range sanokataRows {
		list.Header.EntryNo = so.DocumentNo
		list.Header.PartyName = utils.StrVal(so.PartyName)
		list.Header.SalesOrder = utils.StrVal(so.SalesOrder)
		list.Header.VehicleNumber = utils.StrVal(so.VehicleNumber)

		break
	}
	for size, bundlesCount := range sizeBundles {
		roundedWeight := math.Round(totalWeightForSize[size])

		list.SanokataSummary = append(
			list.SanokataSummary,
			SanokataSummaryRow{
				SizeMM:    size,
				Bundles:   bundlesCount,
				NetWeight: roundedWeight,
			},
		)
	}

	// All three sets are nullable; shook is purely additional. At
	// least thulokata or sanokata lines must exist.
	if len(list.ThulokataLines) == 0 && len(list.SanokataLines) == 0 {
		return nil, ErrPackingListNotFound
	}

	return list, nil
}
