package model

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"go.uber.org/zap"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type AnyTypeValue struct {
	Str   *string
	Int   *int64
	Time  *time.Time
	Bool  *bool
	Float *float64

	ArrayStr  *[]string
	ArrayInt  *[]int64
	ArrayBool *[]bool
}

func (a *AnyTypeValue) TimeToStr() {
	a.Str = utils.StringPtr((*a.Time).String())
	a.Time = nil
}

func (a *AnyTypeValue) BoolToStr() {
	a.Str = utils.StringPtr(strconv.FormatBool(*a.Bool))
	a.Bool = nil
}

func (a *AnyTypeValue) IntToStr() {
	a.Str = utils.StringPtr(strconv.FormatInt(*a.Int, 10))
	a.Int = nil
}

func (a *AnyTypeValue) FloatToStr() {
	a.Str = utils.StringPtr(fmt.Sprintf("%f", *a.Float))
	a.Float = nil
}

func (a *AnyTypeValue) StrToTime() {
	timeValue, err := graphql.UnmarshalTime(*a.Str)
	if err != nil {
		zap.L().Sugar().Errorf("failed unmarshal time field: %s", *a.Str)
	}
	a.Time = &timeValue
	a.Str = nil
}

func (a *AnyTypeValue) StrToBool() {
	boolValue, err := strconv.ParseBool(*a.Str)
	if err != nil {
		zap.L().Sugar().Error(err)
	}
	a.Bool = &boolValue
	a.Str = nil
}

func (a *AnyTypeValue) IntToBool() {
	boolValue := *a.Int != 0
	a.Bool = &boolValue
	a.Int = nil
}

func (a *AnyTypeValue) StrToInt() {
	if intVal, err := strconv.ParseInt(*a.Str, 10, 64); err != nil {
		zap.L().Sugar().Errorf("%s is not an int", *a.Str)
	} else {
		a.Int = &intVal
		a.Str = nil
	}
}

func (a *AnyTypeValue) StrToFloat() {
	if floatVal, err := strconv.ParseFloat(*a.Str, 64); err != nil {
		zap.L().Sugar().Errorf("%s is not a float", *a.Str)
	} else {
		a.Float = &floatVal
		a.Str = nil
	}
}

func (a *AnyTypeValue) IntToFloat() {
	floatVal := float64(*a.Int)
	a.Float = &floatVal
	a.Int = nil
}

func (a *AnyTypeValue) FloatToInt() {
	intVal := int64(*a.Float)
	a.Int = &intVal
	a.Float = nil
}

func (a *AnyTypeValue) RealValue() any {
	if a.Int != nil {
		return *a.Int
	} else if a.Float != nil {
		return *a.Float
	} else if a.Time != nil {
		return *a.Time
	} else if a.Bool != nil {
		return *a.Bool
	} else if a.ArrayStr != nil {
		return *a.ArrayStr
	} else if a.ArrayInt != nil {
		return *a.ArrayInt
	} else if a.ArrayBool != nil {
		return *a.ArrayBool
	} else {
		return *a.Str
	}
}

func MarshalAnyTypeValue(atv AnyTypeValue) graphql.Marshaler {
	if atv.Time != nil {
		return graphql.MarshalTime(*atv.Time)
	} else if atv.Int != nil {
		return graphql.MarshalInt64(*atv.Int)
	} else if atv.Float != nil {
		return graphql.MarshalFloat(*atv.Float)
	} else if atv.Bool != nil {
		return graphql.MarshalBoolean(*atv.Bool)
	}
	return graphql.MarshalString(*atv.Str)
}

func UnmarshalAnyTypeValue(input any) (AnyTypeValue, error) {
	switch input := input.(type) {
	case string:
		pt, err := time.Parse(time.RFC3339, input)
		if err != nil {
			return AnyTypeValue{Str: &input}, nil
		} else {
			return AnyTypeValue{Time: &pt}, nil
		}
	case int64:
		return AnyTypeValue{Int: &input}, nil
	case json.Number:
		intVal, err := input.Int64()
		if err != nil {
			return AnyTypeValue{}, err
		}
		return AnyTypeValue{Int: &intVal}, nil
	case time.Time:
		return AnyTypeValue{Time: &input}, nil
	case bool:
		return AnyTypeValue{Bool: &input}, nil
	case float64:
		return AnyTypeValue{Float: &input}, nil
	case []interface{}:
		if len(input) == 0 {
			return AnyTypeValue{}, nil
		}
		switch input[0].(type) {
		case int64:
			var arrayInt []int64
			for _, v := range input {
				arrayInt = append(arrayInt, v.(int64))
			}
			return AnyTypeValue{ArrayInt: &arrayInt}, nil
		case json.Number:
			var arrayInt []int64
			for _, v := range input {
				intVal, err := v.(json.Number).Int64()
				if err != nil {
					return AnyTypeValue{}, err
				}
				arrayInt = append(arrayInt, intVal)
			}
			return AnyTypeValue{ArrayInt: &arrayInt}, nil
		case bool:
			var arrayBool []bool
			for _, v := range input {
				arrayBool = append(arrayBool, v.(bool))
			}
			return AnyTypeValue{ArrayBool: &arrayBool}, nil
		case string:
			var arrayStr []string
			for _, v := range input {
				arrayStr = append(arrayStr, v.(string))
			}
			return AnyTypeValue{ArrayStr: &arrayStr}, nil
		default:
			return AnyTypeValue{}, fmt.Errorf("unknown type for input: %s", input)
		}
	default:
		return AnyTypeValue{}, fmt.Errorf("unknown type for input: %s", input)
	}
}
