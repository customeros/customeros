package utils

import (
	"fmt"
	"strings"
)

const (
	paramPrefix = "param"
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
	nodeAlias       string
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
	var cypherStr strings.Builder

	if f.Details == nil && len(f.Filters) == 0 {
		return Cypher(""), nil
	}
	f.nodeAlias = nodeAlias
	f.paramCount = 0

	cypherStr.WriteString(" WHERE ")
	innerCypherStr, params := f.buildCypherFilterFragment()
	cypherStr.WriteString(innerCypherStr)

	return Cypher(cypherStr.String()), params
}

func (f *CypherFilter) buildCypherFilterFragment() (string, map[string]any) {
	var cypherStr strings.Builder
	var params map[string]any

	if f.Negate {
		cypherStr.WriteString(" NOT ")
		innerCypherStr, innerParams := f.buildCypherFilterFragment()
		AddMapToMap(innerParams, params)
		cypherStr.WriteString(SurroundWithRoundParentheses(innerCypherStr))
	} else if f.LogicalOperator != L_NONE {
		i := 0
		cypherStr.WriteString("(")
		for _, v := range f.Filters {
			if i > 0 {
				cypherStr.WriteString(SurroundWithSpaces(f.LogicalOperator.String()))
			}
			innerCypherStr, innerParams := v.buildCypherFilterFragment()
			AddMapToMap(innerParams, params)
			cypherStr.WriteString(SurroundWithRoundParentheses(innerCypherStr))
			i++
		}
		cypherStr.WriteString(")")
	} else {
		// ITEM
	}

	return cypherStr.String(), params
}
