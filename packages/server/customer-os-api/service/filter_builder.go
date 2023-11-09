package service

import (
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
)

func buildFilter(modelFilter *model.Filter, T reflect.Type) (*utils.CypherFilter, error) {
	if modelFilter == nil {
		return nil, nil
	}
	cypherFilter := new(utils.CypherFilter)
	cypherFilter.Negate = false

	foundAtCurrentLevel := false

	if modelFilter.Not != nil {
		foundAtCurrentLevel = true
		cypherFilter.Negate = true
		innerFilter, err := buildFilter(modelFilter.Not, T)
		if err != nil {
			return nil, err
		}
		cypherFilter.Filters = append(cypherFilter.Filters, innerFilter)
	}
	if modelFilter.And != nil {
		if foundAtCurrentLevel {
			return nil, newFilterError()
		}
		foundAtCurrentLevel = true
		if len(modelFilter.And) < 2 {
			return nil, newFilterErrorf("at least 2 filters expected in AND group")
		}
		cypherFilter.LogicalOperator = utils.AND
		for _, v := range modelFilter.And {
			innerFilter, err := buildFilter(v, T)
			if err != nil {
				return nil, err
			}
			cypherFilter.Filters = append(cypherFilter.Filters, innerFilter)
		}
	}
	if modelFilter.Or != nil {
		if foundAtCurrentLevel {
			return nil, newFilterError()
		}
		foundAtCurrentLevel = true
		if len(modelFilter.Or) < 2 {
			return nil, newFilterErrorf("at least 2 filters expected in OR group")
		}
		cypherFilter.LogicalOperator = utils.OR
		for _, v := range modelFilter.Or {
			innerFilter, err := buildFilter(v, T)
			if err != nil {
				return nil, err
			}
			cypherFilter.Filters = append(cypherFilter.Filters, innerFilter)
		}
	}
	if modelFilter.Filter != nil {
		if foundAtCurrentLevel {
			return nil, newFilterError()
		}
		props, err := utils.GetPropertyDetailsByLookupName(T, modelFilter.Filter.Property)
		if err != nil {
			return nil, err
		}
		cypherFilterItem := utils.CypherFilterItem{
			NodeProperty:         props[utils.TagProperty],
			SupportCaseSensitive: props[utils.TagSupportCaseSensitive] == "true",
			CaseSensitive:        *modelFilter.Filter.CaseSensitive == true,
			Value:                modelFilter.Filter.Value.RealValue(),
			ComparisonOperator:   comparisonOperatorToEnum(modelFilter.Filter.Operation),
			DbNodePropertyProps:  props,
		}
		cypherFilter.Details = &cypherFilterItem
	}
	return cypherFilter, nil
}

func comparisonOperatorToEnum(co model.ComparisonOperator) utils.ComparisonOperator {
	switch co {
	case model.ComparisonOperatorContains:
		return utils.CONTAINS
	case model.ComparisonOperatorStartsWith:
		return utils.STARTS_WITH
	case model.ComparisonOperatorIn:
		return utils.IN
	case model.ComparisonOperatorLte:
		return utils.LTE
	case model.ComparisonOperatorGte:
		return utils.GTE
	case model.ComparisonOperatorBetween:
		return utils.BETWEEN
	default:
		return utils.EQUALS
	}
}

func newFilterError() error {
	return errors.New("incorrect filter formatting")
}

func newFilterErrorf(msg string) error {
	return errors.New(fmt.Sprintf("incorrect filter formatting: %s", msg))
}
