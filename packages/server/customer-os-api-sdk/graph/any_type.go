package graph

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
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
		logrus.Errorf("failed unmarshal time field: %s", *a.Str)
	}
	a.Time = &timeValue
	a.Str = nil
}

func (a *AnyTypeValue) StrToBool() {
	boolValue, err := strconv.ParseBool(*a.Str)
	if err != nil {
		logrus.Error(err)
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
		logrus.Errorf("%s is not an int", *a.Str)
	} else {
		a.Int = &intVal
		a.Str = nil
	}
}

func (a *AnyTypeValue) StrToFloat() {
	if floatVal, err := strconv.ParseFloat(*a.Str, 64); err != nil {
		logrus.Errorf("%s is not a float", *a.Str)
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
		return AnyTypeValue{Str: &input}, nil
	case int64:
		return AnyTypeValue{Int: &input}, nil
	case time.Time:
		return AnyTypeValue{Time: &input}, nil
	case bool:
		return AnyTypeValue{Bool: &input}, nil
	case float64:
		return AnyTypeValue{Float: &input}, nil
	default:
		return AnyTypeValue{}, fmt.Errorf("unknown type for input: %s", input)
	}
}
