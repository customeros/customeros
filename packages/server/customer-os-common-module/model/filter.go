package model

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
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
	//deprecated
	ComparisonOperatorEq          ComparisonOperator = "EQ"
	ComparisonOperatorEquals      ComparisonOperator = "EQUALS"
	ComparisonOperatorNotEquals   ComparisonOperator = "NOT_EQUALS"
	ComparisonOperatorContains    ComparisonOperator = "CONTAINS"
	ComparisonOperatorNotContains ComparisonOperator = "NOT_CONTAINS"
	ComparisonOperatorStartsWith  ComparisonOperator = "STARTS_WITH"
	ComparisonOperatorLte         ComparisonOperator = "LTE"
	ComparisonOperatorGte         ComparisonOperator = "GTE"
	ComparisonOperatorIn          ComparisonOperator = "IN"
	ComparisonOperatorNotIn       ComparisonOperator = "NOT_IN"
	ComparisonOperatorBetween     ComparisonOperator = "BETWEEN"
	ComparisonOperatorIsNull      ComparisonOperator = "IS_NULL"
	ComparisonOperatorIsNotNull   ComparisonOperator = "IS_NOT_NULL"
	ComparisonOperatorIsEmpty     ComparisonOperator = "IS_EMPTY"
	ComparisonOperatorIsNotEmpty  ComparisonOperator = "IS_NOT_EMPTY"
	ComparisonOperatorLt          ComparisonOperator = "LT"
	ComparisonOperatorGt          ComparisonOperator = "GT"

	ComparisonOperatorCountRelation ComparisonOperator = "COUNT_RELATION"
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
				dateTime, err := UnmarshalDateTime(v.(string))
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

const customLayout1 = "2006-01-02 15:04:05"
const customLayout2 = "2006-01-02T15:04:05.000-0700"
const customLayout3 = "2006-01-02T15:04:05-07:00"
const customLayout4 = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
const customLayout5 = "Mon, 2 Jan 2006 15:04:05 MST"
const customLayout6 = "Mon, 2 Jan 2006 15:04:05 -0700"
const customLayout7 = "Mon, 2 Jan 2006 15:04:05 +0000 (GMT)"
const customLayout8 = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
const customLayout9 = "2 Jan 2006 15:04:05 -0700"

func UnmarshalDateTime(input string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, input)
	if err == nil {
		// Parsed as RFC3339
		return &t, nil
	}

	// Try custom layouts
	customLayouts := []string{customLayout1, customLayout2, customLayout4, customLayout5, customLayout6, customLayout7, customLayout8, customLayout9}

	for _, layout := range customLayouts {
		t, err = time.Parse(layout, input)
		if err == nil {
			return &t, nil
		}
	}
	inputForLayout3 := input
	if !strings.Contains(input, "[UTC]") {
		index := strings.Index(input, "[")
		// If found, strip off the timezone information
		if index != -1 {
			inputForLayout3 = input[:index]
		}
	}
	t, err = time.Parse(customLayout3, inputForLayout3)
	if err == nil {
		return &t, nil
	}

	return nil, errors.New(fmt.Sprintf("cannot parse input as date time %s", input))
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
