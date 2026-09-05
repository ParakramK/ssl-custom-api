package utils

import (
	"math"
	"math/big"

	"github.com/SAP/go-hdb/driver"
)

func Round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func DecimalToFloat64(d driver.Decimal) float64 {
	rat := (*big.Rat)(&d)

	value, _ := rat.Float64()

	return value
}

func StrVal(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func ValueOrZero(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
