package utils

import (
	"fmt"
	"strings"
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
	return "", nil
	//if len(f.filters) == 0 {
	//	return ""
	//}
	//query := " WHERE "
	//for i := 0; i < len(f.filters); i++ {
	//	//filterProperty := f.filters[i]
	//	if i > 0 {
	//		query += " , "
	//	}
	//	//toLower := filterProperty.supportCaseSensitive && !filterProperty.caseSensitive
	//	//if toLower {
	//	//	query += "toLower("
	//	//}
	//	//query += nodeAlias
	//	//query += "."
	//	//query += filterProperty.nodeProperty
	//	//if toLower {
	//	//	query += ")"
	//	//}
	//}
	//return query
}
