package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	paramPrefix = "param_"
)

type ComparisonOperator int
type LogicalOperator int

const (
	C_NONE ComparisonOperator = iota
	IS_NULL
	EQUALS
	CONTAINS
	STARTS_WITH
	LTE
	GTE
	IN
	BETWEEN
)

func (c ComparisonOperator) String() string {
	switch c {
	case C_NONE:
		return "NONE"
	case IS_NULL:
		return "IS_NULL"
	case EQUALS:
		return "EQUALS"
	case CONTAINS:
		return "CONTAINS"
	case STARTS_WITH:
		return "STARTS WITH"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case IN:
		return "IN"
	case BETWEEN:
		return "BETWEEN"
	default:
		return fmt.Sprintf("%d", int(c))
	}
}

func (c ComparisonOperator) CypherString() string {
	switch c {
	case C_NONE:
		return ""
	case IS_NULL:
		return "is null"
	case EQUALS:
		return "="
	case CONTAINS:
		return "CONTAINS"
	case STARTS_WITH:
		return "STARTS WITH"
	case LTE:
		return "<="
	case GTE:
		return ">="
	case IN:
		return "IN"
	case BETWEEN:
		return "BETWEEN"
	default:
		return "="
	}
}

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
	ComparisonOperator   ComparisonOperator
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
	res.WriteString(fmt.Sprintf("ComparisonOperator: %v ", f.ComparisonOperator.String()))
	return res.String()
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
		cypherStr.WriteString(SurroundWithSpaces(f.Details.ComparisonOperator.CypherString()))
		if toLower {
			cypherStr.WriteString("toLower(")
		}

		if f.Details.ComparisonOperator != IS_NULL {
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

	return cypherStr.String(), params
}
