package sales

import (
	"ssl-custom-api/internal/utils"
)

func TagDetails(tag TagDataRow) BundleDetails {
	return BundleDetails{
		BundleTag:       tag.Code,
		SizeMM:          tag.SizeMM,
		Lot:             tag.Lot,
		Weight:          tag.NetWeight,
		Pieces:          tag.Pieces,
		Bundles:         tag.Bundles,
		LoadingDate:     tag.LoadingDate,
		SslSno:          tag.DocumentNo,
		LoadingSlipNo:   tag.LoadingSlipNo,
		VehicleNumber:   tag.VehicleNo,
		SalesOrder:      utils.StrVal(tag.SalesOrder),
		PartyName:       utils.StrVal(tag.PartyName),
		ShippingAddress: utils.StrVal(tag.ShippingAddress),
	}
}

func packThulokataRow(row ThulokataRow) ThulokataPackingRow {
	return ThulokataPackingRow{

		Materials:     row.Materials,
		InitialWeight: row.InitialWeight,
		FinalWeight:   row.FinalWeight,
		NetWeight:     row.NetWeight,
		Bundles:       row.Bundles,
	}
}

func packSanokataRow(row SanokataRow) SanokataPackingRow {
	return SanokataPackingRow{
		Code:        row.Code,
		SizeMM:      row.SizeMM,
		Bundles:     row.Bundles,
		Pieces:      row.Pieces,
		NetWeight:   row.NetWeight,
		LoadingDate: row.LoadingDate,
		KataNo:      row.KataNo,
	}
}
func packShookRow(row SHookLine) SHookRow {
	return SHookRow{
		DocumentNo: row.DocumentNo,
		Size:       row.Size,
		Quantity:   row.Quantity,
	}
}
