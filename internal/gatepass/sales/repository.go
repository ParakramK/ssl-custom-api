package sales

import (
	"context"
	"errors"
	"fmt"

	"ssl-custom-api/internal/gatepass/models"
	"ssl-custom-api/internal/gatepass/query"

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

// func (r *salesRepository) GetTagData(
// 	ctx context.Context,
// 	tagNo string,
// ) (*BundleDetails, error) {

// 	if tagNo == "" {
// 		return nil, errors.New("tag number is required")
// 	}

// 	qc := query.Use(r.db)

// 	tag, err := qc.Sanokata.WithContext(ctx).
// 		Where(qc.Sanokata.Code.Eq(tagNo)).
// 		Order(qc.Sanokata.CreatedAt.Desc()).
// 		First()
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrTagNotFound
// 		}
// 		return nil, fmt.Errorf("get tag data: %w", err)
// 	}

// 	ge, err := qc.GateEntry.WithContext(ctx).
// 		Where(qc.GateEntry.DocumentNo.Eq(tag.DocumentNo)).
// 		First()
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrTagNotFound
// 		}
// 		return nil, fmt.Errorf("get gate entry: %w", err)
// 	}

// 	le, err := qc.LoadingEntry.WithContext(ctx).
// 		Where(qc.LoadingEntry.LoadingSlipNo.Eq(ge.LoadingSlipNo)).
// 		First()
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, ErrTagNotFound
// 		}
// 		return nil, fmt.Errorf("get loading entry: %w", err)
// 	}

//		bundle := TagDetails(*tag, ge, le)
//		return &bundle, nil
//	}
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

	var loadings []models.ThulokataLoading
	if err := r.db.WithContext(ctx).
		Model(&models.ThulokataLoading{}).
		Joins("Thulokata").
		Joins("Thulokata.GateEntry").
		Joins("Thulokata.GateEntry.LoadingEntry").
		Where("`Thulokata__GateEntry`.`document_number` = ?", sslNo).
		Find(&loadings).Error; err != nil {
		return nil, fmt.Errorf("get thulokata packing lines: %w", err)
	}

	// Sanokata loadings: bundles loaded one by one, ordered by
	// size then bundles.
	var orders []models.Sanokata
	if err := r.db.WithContext(ctx).
		Model(&models.Sanokata{}).
		Joins("User").
		Joins("GateEntry").
		Joins("GateEntry.LoadingEntry").
		Where("sales_orders.internal_numbering = ?", sslNo).
		Order("sales_orders.size_mm, sales_orders.bundles").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("get sanokata packing lines: %w", err)
	}

	var shookRows []ShookRow
	err_shook := r.db.WithContext(ctx).
		Select(
			q.SHook.DocumentNo,
			q.SHook.Size,
			q.SHook.Quantity,
		).Where(
		q.SHook.DocumentNo.Eq(sslNo),
	).Scan(&shookRows)

	if err_shook.Error != nil {
		return nil, fmt.Errorf("get shook lines: %w", err_shook.Error)
	}
	list := &PackingList{
		ThulokataLines: []ThulokataPackingRow{},
		SanokataLines:  make([]SanokataPackingRow, 0, len(orders)),
		ShookLines:     make([]ShookRow, 0, len(shookRows)),
	}
	for _, thl := range loadings {
		list.ThulokataLines = append(list.ThulokataLines, packThulokataRow(thl))
	}
	for _, so := range orders {
		list.SanokataLines = append(list.SanokataLines, packSanokataRow(so))
	}
	for _, sh := range shookRows {
		list.ShookLines = append(list.ShookLines, packShookRow(sh))
	}

	// All three sets are nullable; shook is purely additional. At
	// least thulokata or sanokata lines must exist.
	if len(list.ThulokataLines) == 0 && len(list.SanokataLines) == 0 {
		return nil, ErrPackingListNotFound
	}

	return list, nil
}
