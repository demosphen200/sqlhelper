package utils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func ConvertAnyToFloat64(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		return math.NaN(), errors.New(fmt.Sprintf("can not convert value %+v to float", unk))
	}
}

func ConvertAnyToFloat64OrDefault(value any, def float64) float64 {
	floatValue, err := ConvertAnyToFloat64(value)
	if err != nil {
		return def
	}
	return floatValue
}
