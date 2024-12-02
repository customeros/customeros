package utils

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"strconv"
	"strings"
)

const (
	paramPrefix = "param_"
)

type LogicalOperator int

const (
	L_NONE LogicalOperator = iota
	AND
	OR
)

func (l LogicalOperator) String() string {
	switch l {
	case L_NONE:
		return "NONE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	default:
		return fmt.Sprintf("%d", int(l))
	}
}

type CypherFilterItem struct {
	NodeProperty         string
	SupportCaseSensitive bool
	CaseSensitive        bool
	Value                any
	DbNodePropertyProps  map[string]string
	ComparisonOperator   model.ComparisonOperator
}

type CypherFilter struct {
	Negate          bool
	LogicalOperator LogicalOperator
	Filters         []*CypherFilter
	Details         *CypherFilterItem
	paramCount      int
}

func (f CypherFilter) String() string {
	var res strings.Builder
	res.WriteString(fmt.Sprintf("Negate: %v ", f.Negate))
	res.WriteString(fmt.Sprintf("LogicalOperator: %v ", f.LogicalOperator.String()))
	if f.Details != nil {
		res.WriteString(fmt.Sprintf("Details: {%v} ", f.Details.String()))
	}
	var filtersRes strings.Builder
	for _, v := range f.Filters {
		filtersRes.WriteString("{")
		filtersRes.WriteString(v.String())
		filtersRes.WriteString("}")
	}
	res.WriteString(fmt.Sprintf("Filters: [%v] ", filtersRes.String()))
	return res.String()
}

func (f CypherFilterItem) String() string {
	var res strings.Builder
	res.WriteString(fmt.Sprintf("NodeProperty: %v ", f.NodeProperty))
	res.WriteString(fmt.Sprintf("SupportCaseSensitive: %v ", f.SupportCaseSensitive))
	res.WriteString(fmt.Sprintf("CaseSensitive: %v ", f.CaseSensitive))
	res.WriteString(fmt.Sprintf("Value: %v ", f.Value))
	res.WriteString(fmt.Sprintf("DbNodePropertyProps: %v ", f.DbNodePropertyProps))
	res.WriteString(fmt.Sprintf("ComparisonOperator: %s ", f.ComparisonOperator))
	return res.String()
}

func CreateStringCypherFilter(propertyName string, searchTerm any, comparator model.ComparisonOperator) *CypherFilter {
	filter := CypherFilter{}
	filter.Details = new(CypherFilterItem)
	filter.Details.NodeProperty = propertyName
	filter.Details.Value = &searchTerm
	filter.Details.ComparisonOperator = comparator
	filter.Details.SupportCaseSensitive = true
	return &filter
}

func CreateCypherFilter(propertyName string, searchTerm any, comparator model.ComparisonOperator) *CypherFilter {
	filter := CypherFilter{}
	filter.Details = new(CypherFilterItem)
	filter.Details.NodeProperty = propertyName
	filter.Details.Value = &searchTerm
	filter.Details.ComparisonOperator = comparator
	filter.Details.SupportCaseSensitive = false
	return &filter
}

func CreateCypherFilterIsNull(propertyName string) *CypherFilter {
	return CreateCypherFilter(propertyName, "", model.ComparisonOperatorIsNull)
}

func CreateCypherFilterIsNotNull(propertyName string) *CypherFilter {
	return CreateCypherFilter(propertyName, "", model.ComparisonOperatorIsNotNull)
}

func CreateCypherFilterIn(propertyName string, arrayValues any) *CypherFilter {
	return CreateCypherFilter(propertyName, arrayValues, model.ComparisonOperatorIn)
}

func CreateCypherFilterEq(propertyName string, value any) *CypherFilter {
	return CreateCypherFilter(propertyName, value, model.ComparisonOperatorEq)
}

func CreateCypherFilterNotEq(propertyName string, value any) *CypherFilter {
	return CreateCypherFilter(propertyName, value, model.ComparisonOperatorNotEquals)
}

func (f *CypherFilter) CypherFilterFragment(nodeAlias string) (Cypher, map[string]any) {
	if f == nil || (f.Details == nil && (f.Filters == nil || len(f.Filters) == 0)) {
		return "", map[string]any{}
	}

	f.paramCount = 0

	var cypherStr strings.Builder
	cypherStr.WriteString(" WHERE ")
	innerCypherStr, params := f.BuildCypherFilterFragment(nodeAlias)
	cypherStr.WriteString(innerCypherStr)

	return Cypher(cypherStr.String()), params
}

func (f *CypherFilter) BuildCypherFilterFragment(nodeAlias string) (string, map[string]any) {
	return f.BuildCypherFilterFragmentWithParamName(nodeAlias, paramPrefix)
}
func (f *CypherFilter) BuildCypherFilterFragmentWithParamName(nodeAlias string, customParamPrefix string) (string, map[string]any) {
	var cypherStr strings.Builder
	var params = map[string]any{}

	// convert IS_EMPTY to IS_NULL + EQUALS empty string
	if f.Details != nil && f.Details.ComparisonOperator == model.ComparisonOperatorIsEmpty {
		nodeProperty := f.Details.NodeProperty
		dbNodePropertyProps := f.Details.DbNodePropertyProps
		f.LogicalOperator = OR
		f.Filters = append(f.Filters,
			&CypherFilter{
				Details: &CypherFilterItem{
					NodeProperty:        nodeProperty,
					ComparisonOperator:  model.ComparisonOperatorIsNull,
					DbNodePropertyProps: dbNodePropertyProps,
				},
				LogicalOperator: L_NONE,
			},
			&CypherFilter{
				Details: &CypherFilterItem{
					NodeProperty:        nodeProperty,
					ComparisonOperator:  model.ComparisonOperatorEq,
					Value:               "",
					DbNodePropertyProps: dbNodePropertyProps,
				},
				LogicalOperator: L_NONE,
			})
		f.Details = nil
	} else if f.Details != nil && f.Details.ComparisonOperator == model.ComparisonOperatorIsNotEmpty {
		nodeProperty := f.Details.NodeProperty
		dbNodePropertyProps := f.Details.DbNodePropertyProps
		f.LogicalOperator = AND
		f.Filters = append(f.Filters,
			&CypherFilter{
				Details: &CypherFilterItem{
					NodeProperty:        nodeProperty,
					ComparisonOperator:  model.ComparisonOperatorIsNotNull,
					DbNodePropertyProps: dbNodePropertyProps,
				},
				LogicalOperator: L_NONE,
			},
			&CypherFilter{
				Details: &CypherFilterItem{
					NodeProperty:        nodeProperty,
					ComparisonOperator:  model.ComparisonOperatorNotEquals,
					Value:               "",
					DbNodePropertyProps: dbNodePropertyProps,
				},
				LogicalOperator: L_NONE,
			})
		f.Details = nil
	} else if f.Details != nil && f.Details.ComparisonOperator == model.ComparisonOperatorNotContains {
		nodeProperty := f.Details.NodeProperty
		dbNodePropertyProps := f.Details.DbNodePropertyProps
		f.LogicalOperator = AND
		f.Negate = true
		f.Filters = append(f.Filters,
			&CypherFilter{
				Details: &CypherFilterItem{
					NodeProperty:         nodeProperty,
					ComparisonOperator:   model.ComparisonOperatorContains,
					DbNodePropertyProps:  dbNodePropertyProps,
					Value:                f.Details.Value,
					SupportCaseSensitive: f.Details.SupportCaseSensitive,
					CaseSensitive:        f.Details.CaseSensitive,
				},
				LogicalOperator: L_NONE,
			},
		)
		f.Details = nil
	}
	if f.Negate {
		cypherStr.WriteString(" NOT ")
		f.Filters[0].paramCount = f.paramCount
		innerCypherStr, innerParams := f.Filters[0].BuildCypherFilterFragmentWithParamName(nodeAlias, customParamPrefix)
		f.paramCount = f.Filters[0].paramCount
		MergeMapToMap(innerParams, params)
		cypherStr.WriteString(SurroundWithRoundParentheses(innerCypherStr))
	} else if f.LogicalOperator != L_NONE {
		cypherStr.WriteString("(")
		i := 0
		for _, v := range f.Filters {
			if i > 0 {
				cypherStr.WriteString(SurroundWithSpaces(f.LogicalOperator.String()))
			}
			v.paramCount = f.paramCount
			innerCypherStr, innerParams := v.BuildCypherFilterFragmentWithParamName(nodeAlias, customParamPrefix)
			f.paramCount = v.paramCount
			MergeMapToMap(innerParams, params)
			cypherStr.WriteString(SurroundWithRoundParentheses(innerCypherStr))
			i++
		}
		cypherStr.WriteString(")")
	} else {

		if f.Details.ComparisonOperator == model.ComparisonOperatorCountRelation {
			cypherStr.WriteString(f.Details.NodeProperty) //hack. you need the full condition here
		} else {
			toLower := f.Details.SupportCaseSensitive && !f.Details.CaseSensitive
			if toLower {
				cypherStr.WriteString("toLower(")
			}
			cypherStr.WriteString(nodeAlias)
			cypherStr.WriteString(".")
			cypherStr.WriteString(f.Details.NodeProperty)
			if toLower {
				cypherStr.WriteString(")")
			}
			cypherStr.WriteString(SurroundWithSpaces(CypherString(f.Details.ComparisonOperator)))
			if toLower {
				cypherStr.WriteString("toLower(")
			}

			if f.Details.ComparisonOperator != model.ComparisonOperatorIsNull && f.Details.ComparisonOperator != model.ComparisonOperatorIsNotNull {
				f.paramCount++
				paramSuffix := strconv.Itoa(f.paramCount)
				cypherStr.WriteString("$" + customParamPrefix + paramSuffix)
				if params == nil {
					params = map[string]any{customParamPrefix + paramSuffix: f.Details.Value}
				} else {
					params[customParamPrefix+paramSuffix] = f.Details.Value
				}
			}

			if toLower {
				cypherStr.WriteString(")")
			}
		}
	}

	return cypherStr.String(), params
}

func CypherString(c model.ComparisonOperator) string {
	switch c {
	case model.ComparisonOperatorIsNull:
		return "is null"
	case model.ComparisonOperatorIsNotNull:
		return "is not null"
	case model.ComparisonOperatorEq:
		return "="
	case model.ComparisonOperatorEquals:
		return "="
	case model.ComparisonOperatorNotEquals:
		return "<>"
	case model.ComparisonOperatorContains:
		return "CONTAINS"
	case model.ComparisonOperatorStartsWith:
		return "STARTS WITH"
	case model.ComparisonOperatorGte:
		return ">="
	case model.ComparisonOperatorGt:
		return ">"
	case model.ComparisonOperatorIn:
		return "IN"
	case model.ComparisonOperatorBetween:
		return "BETWEEN"
	case model.ComparisonOperatorLt:
		return "<"
	case model.ComparisonOperatorLte:
		return "<="
	default:
		return "="
	}
}
