package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

var entityType = reflect.TypeOf(neo4jentity.ContactEntity{})
var defaultErrorMessage = "incorrect filter formatting"
var defaultStringFilterItem = model.FilterItem{
	Property:  "NAME",
	Operation: model.ComparisonOperatorContains,
	Value: model.AnyTypeValue{
		Str: utils.StringPtr("testValue"),
	},
	CaseSensitive: utils.BoolPtr(true),
}
var defaultIntFilterItem = model.FilterItem{
	Property:  "NAME",
	Operation: model.ComparisonOperatorEq,
	Value: model.AnyTypeValue{
		Int: utils.Int64Ptr(100),
	},
	CaseSensitive: utils.BoolPtr(false),
}

func TestBuilderFilter_InputIsNil(t *testing.T) {
	cypherFilter, err := buildFilter(nil, entityType)
	require.Nil(t, err)
	require.Nil(t, cypherFilter)
}

func TestBuilderFilter_NotWithAnd_NotAllowed(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		Not: &model.Filter{},
		And: []*model.Filter{{}},
	}, entityType)
	require.Nil(t, cypherFilter)
	require.NotNil(t, err)
	require.Equal(t, defaultErrorMessage, err.Error())
}

func TestBuilderFilter_NotWithOr_NotAllowed(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		Not: &model.Filter{},
		Or:  []*model.Filter{{}},
	}, entityType)
	require.Nil(t, cypherFilter)
	require.NotNil(t, err)
	require.Equal(t, defaultErrorMessage, err.Error())
}

func TestBuilderFilter_NotWithFilterItem_NotAllowed(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		Not:    &model.Filter{},
		Filter: &model.FilterItem{},
	}, entityType)
	require.Nil(t, cypherFilter)
	require.NotNil(t, err)
	require.Equal(t, defaultErrorMessage, err.Error())
}

func TestBuilderFilter_AndWithOr_NotAllowed(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		And: []*model.Filter{{}, {}},
		Or:  []*model.Filter{{}, {}},
	}, entityType)
	require.Nil(t, cypherFilter)
	require.NotNil(t, err)
	require.Equal(t, defaultErrorMessage, err.Error())
}

func TestBuilderFilter_AndWithFilterItem_NotAllowed(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		And:    []*model.Filter{{}, {}},
		Filter: &model.FilterItem{},
	}, entityType)
	require.Nil(t, cypherFilter)
	require.NotNil(t, err)
	require.Equal(t, defaultErrorMessage, err.Error())
}

func TestBuilderFilter_OrWithFilterItem_NotAllowed(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		Or:     []*model.Filter{{}, {}},
		Filter: &model.FilterItem{},
	}, entityType)
	require.Nil(t, cypherFilter)
	require.NotNil(t, err)
	require.Equal(t, defaultErrorMessage, err.Error())
}

func TestBuilderFilter_AndWithOneFilter_NotAllowed(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		And: []*model.Filter{{}},
	}, entityType)
	require.Nil(t, cypherFilter)
	require.NotNil(t, err)
	require.Equal(t, "incorrect filter formatting: at least 2 filters expected in AND group", err.Error())
}

func TestBuilderFilter_OrWithOneFilter_NotAllowed(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		Or: []*model.Filter{{}},
	}, entityType)
	require.Nil(t, cypherFilter)
	require.NotNil(t, err)
	require.Equal(t, "incorrect filter formatting: at least 2 filters expected in OR group", err.Error())
}

func TestBuilderFilter_NegationOfNegation(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		Not: &model.Filter{
			Not: &model.Filter{
				Filter: &defaultStringFilterItem,
			},
		},
	}, entityType)
	require.Nil(t, err)

	require.Equal(t, true, cypherFilter.Negate)
	require.Nil(t, cypherFilter.Details)
	require.Equal(t, 1, len(cypherFilter.Filters))
	require.Equal(t, utils.L_NONE, cypherFilter.LogicalOperator)
	innerFilter1 := cypherFilter.Filters[0]
	require.Equal(t, true, innerFilter1.Negate)
	require.Equal(t, utils.L_NONE, innerFilter1.LogicalOperator)
	require.Nil(t, innerFilter1.Details)
	require.Equal(t, 1, len(innerFilter1.Filters))
	innerFilter2 := innerFilter1.Filters[0]
	require.Equal(t, false, innerFilter2.Negate)
	require.Equal(t, utils.L_NONE, innerFilter2.LogicalOperator)
	require.Empty(t, innerFilter2.Filters)
	require.NotNil(t, innerFilter2.Details)
	details := innerFilter2.Details
	require.Equal(t, "name", details.NodeProperty)
	require.Equal(t, true, details.SupportCaseSensitive)
	require.Equal(t, true, details.CaseSensitive)
	require.Equal(t, utils.CONTAINS, details.ComparisonOperator)
	require.Equal(t, "testValue", details.Value)
	require.NotEmpty(t, details.DbNodePropertyProps)
}

func TestBuilderFilter_NegationOfGroup(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		Not: &model.Filter{
			And: []*model.Filter{
				{Filter: &defaultStringFilterItem},
				{Filter: &defaultStringFilterItem},
			},
		},
	}, entityType)
	require.Nil(t, err)

	require.Equal(t, true, cypherFilter.Negate)
	require.Nil(t, cypherFilter.Details)
	require.Equal(t, 1, len(cypherFilter.Filters))
	require.Equal(t, utils.L_NONE, cypherFilter.LogicalOperator)

	andGroupFilter := cypherFilter.Filters[0]
	require.Equal(t, false, andGroupFilter.Negate)
	require.Equal(t, utils.AND, andGroupFilter.LogicalOperator)
	require.Nil(t, andGroupFilter.Details)
	require.Equal(t, 2, len(andGroupFilter.Filters))

	andFirstCondition := andGroupFilter.Filters[0]
	require.Equal(t, false, andFirstCondition.Negate)
	require.Equal(t, utils.L_NONE, andFirstCondition.LogicalOperator)
	require.Empty(t, andFirstCondition.Filters)
	require.NotNil(t, andFirstCondition.Details)
	require.Equal(t, "name", andFirstCondition.Details.NodeProperty)

	andSecondCondition := andGroupFilter.Filters[0]
	require.Equal(t, false, andSecondCondition.Negate)
	require.Equal(t, utils.L_NONE, andSecondCondition.LogicalOperator)
	require.Empty(t, andSecondCondition.Filters)
	require.NotNil(t, andSecondCondition.Details)
	require.Equal(t, "name", andSecondCondition.Details.NodeProperty)
}

func TestBuilderFilter_GroupOfGroupAndItem(t *testing.T) {
	cypherFilter, err := buildFilter(&model.Filter{
		And: []*model.Filter{
			{Or: []*model.Filter{
				{Filter: &defaultStringFilterItem},
				{Filter: &defaultStringFilterItem},
			}},
			{Filter: &defaultIntFilterItem},
		},
	}, entityType)
	require.Nil(t, err)

	require.Equal(t, false, cypherFilter.Negate)
	require.Nil(t, cypherFilter.Details)
	require.Equal(t, 2, len(cypherFilter.Filters))
	require.Equal(t, utils.AND, cypherFilter.LogicalOperator)

	andFirstFilter := cypherFilter.Filters[0]
	require.Equal(t, false, andFirstFilter.Negate)
	require.Equal(t, utils.OR, andFirstFilter.LogicalOperator)
	require.Nil(t, andFirstFilter.Details)
	require.Equal(t, 2, len(andFirstFilter.Filters))

	orFirstFilter := andFirstFilter.Filters[0]
	require.Equal(t, false, orFirstFilter.Negate)
	require.Equal(t, utils.L_NONE, orFirstFilter.LogicalOperator)
	require.Nil(t, orFirstFilter.Filters)
	require.NotNil(t, orFirstFilter.Details)
	require.Equal(t, "testValue", orFirstFilter.Details.Value)

	orSecondFilter := andFirstFilter.Filters[1]
	require.Equal(t, false, orSecondFilter.Negate)
	require.Equal(t, utils.L_NONE, orSecondFilter.LogicalOperator)
	require.Nil(t, orSecondFilter.Filters)
	require.NotNil(t, orSecondFilter.Details)
	require.Equal(t, "testValue", orSecondFilter.Details.Value)

	andSecondFilter := cypherFilter.Filters[1]
	require.Equal(t, false, andSecondFilter.Negate)
	require.Equal(t, utils.L_NONE, andSecondFilter.LogicalOperator)
	require.Empty(t, andSecondFilter.Filters)
	require.NotNil(t, andSecondFilter.Details)
	require.Equal(t, int64(100), andSecondFilter.Details.Value)
	require.Equal(t, false, andSecondFilter.Details.CaseSensitive)
}
