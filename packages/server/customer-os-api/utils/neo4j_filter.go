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
	EQUALS
	CONTAINS
)

func (c ComparisonOperator) String() string {
	switch c {
	case C_NONE:
		return "NONE"
	case EQUALS:
		return "EQUALS"
	case CONTAINS:
		return "CONTAINS"
	default:
		return fmt.Sprintf("%d", int(c))
	}
}

func (c ComparisonOperator) CypherString() string {
	switch c {
	case C_NONE:
		return ""
	case EQUALS:
		return "="
	case CONTAINS:
		return "CONTAINS"
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
	innerCypherStr, params := f.buildCypherFilterFragment(nodeAlias)
	cypherStr.WriteString(innerCypherStr)

	return Cypher(cypherStr.String()), params
}

func (f *CypherFilter) buildCypherFilterFragment(nodeAlias string) (string, map[string]any) {
	var cypherStr strings.Builder
	var params = map[string]any{}

	if f.Negate {
		cypherStr.WriteString(" NOT ")
		f.Filters[0].paramCount = f.paramCount
		innerCypherStr, innerParams := f.Filters[0].buildCypherFilterFragment(nodeAlias)
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
			innerCypherStr, innerParams := v.buildCypherFilterFragment(nodeAlias)
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
		f.paramCount++
		paramSuffix := strconv.Itoa(f.paramCount)
		cypherStr.WriteString("$" + paramPrefix + paramSuffix)
		if params == nil {
			params = map[string]any{paramPrefix + paramSuffix: f.Details.Value}
		} else {
			params[paramPrefix+paramSuffix] = f.Details.Value
		}

		if toLower {
			cypherStr.WriteString(")")
		}
	}

	return cypherStr.String(), params
}
