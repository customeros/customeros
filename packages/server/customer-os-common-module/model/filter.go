package model

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type Filter struct {
	Not    *Filter     `json:"NOT,omitempty"`
	And    []*Filter   `json:"AND,omitempty"`
	Or     []*Filter   `json:"OR,omitempty"`
	Filter *FilterItem `json:"filter,omitempty"`
}

type FilterItem struct {
	Property      string             `json:"property"`
	Operation     ComparisonOperator `json:"operation"`
	CaseSensitive *bool              `json:"caseSensitive,omitempty"`
	IncludeEmpty  *bool              `json:"includeEmpty,omitempty"`
	Value         AnyTypeValue       `json:"-"`
	JsonValue     interface{}        `json:"value"`
}

type ComparisonOperator string

const (
	ComparisonOperatorEq         ComparisonOperator = "EQ"
	ComparisonOperatorContains   ComparisonOperator = "CONTAINS"
	ComparisonOperatorStartsWith ComparisonOperator = "STARTS_WITH"
	ComparisonOperatorLte        ComparisonOperator = "LTE"
	ComparisonOperatorGte        ComparisonOperator = "GTE"
	ComparisonOperatorIn         ComparisonOperator = "IN"
	ComparisonOperatorBetween    ComparisonOperator = "BETWEEN"
	ComparisonOperatorIsNull     ComparisonOperator = "IS_NULL"
	ComparisonOperatorIsEmpty    ComparisonOperator = "IS_EMPTY"
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
	ArrayTime *[]time.Time
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
		case float64:
			// By default, the encoding/json package in Go unmarshals JSON numbers into float64 when the destination type is interface{}.
			// Forcing converting to array of int if the float64 value is actually an integer.
			var arrayInt []int64
			for _, v := range input {
				if v.(float64) == float64(int64(v.(float64))) {
					arrayInt = append(arrayInt, int64(v.(float64)))
				} else {
					return AnyTypeValue{}, fmt.Errorf("unknown type for input: %s", input)
				}
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
			var arrayTime []time.Time
			for _, v := range input {
				dateTime, err := utils.UnmarshalDateTime(v.(string))
				if err == nil {
					arrayTime = append(arrayTime, *dateTime)
					continue
				} else {
					arrayStr = append(arrayStr, v.(string))
				}
			}
			if len(arrayTime) > 0 {
				return AnyTypeValue{ArrayTime: &arrayTime}, nil
			} else {
				return AnyTypeValue{ArrayStr: &arrayStr}, nil
			}
		default:
			return AnyTypeValue{}, fmt.Errorf("unknown type for input: %s", input)
		}
	default:
		return AnyTypeValue{}, fmt.Errorf("unknown type for input: %s", input)
	}
}

func UnmarshalFilter(input string) (*Filter, error) {
	var filter Filter

	if input == "" {
		return &Filter{}, nil
	}

	err := json.Unmarshal([]byte(input), &filter)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Filter: %w", err)
	}

	// Recursively process the filter structure
	err = processFilter(&filter)
	if err != nil {
		return nil, fmt.Errorf("failed to process Filter: %w", err)
	}

	return &filter, nil
}

func processFilter(filter *Filter) error {
	if filter.Not != nil {
		err := processFilter(filter.Not)
		if err != nil {
			return err
		}
	}

	for _, andFilter := range filter.And {
		err := processFilter(andFilter)
		if err != nil {
			return err
		}
	}

	for _, orFilter := range filter.Or {
		err := processFilter(orFilter)
		if err != nil {
			return err
		}
	}

	if filter.Filter != nil {
		value, err := UnmarshalAnyTypeValue(filter.Filter.JsonValue)
		if err != nil {
			return fmt.Errorf("failed to unmarshal AnyTypeValue: %w", err)
		}
		filter.Filter.Value = value
	}

	return nil
}
