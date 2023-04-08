package zmath

import "github.com/shopspring/decimal"

func FloatAdd(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Add(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}
func FloatSub(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Sub(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}
func FloatMul(d1, d2 float64) float64 {
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Mul(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}
func FloatDiv(d1, d2 float64) float64 {
	if d2 == 0 {
		d2 = 1
	}
	decimalD1 := decimal.NewFromFloat(d1)
	decimalD2 := decimal.NewFromFloat(d2)
	decimalResult := decimalD1.Div(decimalD2)
	float64Result, _ := decimalResult.Float64()
	return float64Result
}

func Round(num float64, round int32) float64 {
	d, _ := decimal.NewFromFloat(num).Round(round).Float64()
	return d
}
