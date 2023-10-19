package utils

import (
	"reflect"
	"strings"
)

type CypherSort struct {
	sorts []*OrderBy
}

type OrderBy struct {
	nodeProperty         string
	supportCaseSensitive bool
	caseSensitive        bool
	descending           bool
	nodeAlias            string
	coalesce             bool
	valid                bool
}

func (ob *OrderBy) WithAlias(nodeAlias string) *OrderBy {
	ob.nodeAlias = nodeAlias
	return ob
}

func (ob *OrderBy) WithCoalesce() *OrderBy {
	ob.coalesce = true
	return ob
}

func (ob *OrderBy) WithDescending() *OrderBy {
	ob.descending = true
	return ob
}

func (ob *OrderBy) IsValid() bool {
	return ob.valid
}

func (s *CypherSort) NewSortRule(lookupName, direction string, caseSensitive bool, T reflect.Type) *OrderBy {
	orderBy := new(OrderBy)
	orderBy.valid = true

	props, err := GetPropertyDetailsByLookupName(T, lookupName)
	if err != nil {
		orderBy.valid = false
		return orderBy
	}

	orderBy.nodeProperty = props[TagProperty]
	orderBy.descending = "DESC" == direction
	orderBy.supportCaseSensitive = props[TagSupportCaseSensitive] == "true"
	orderBy.caseSensitive = caseSensitive

	s.sorts = append(s.sorts, orderBy)

	return orderBy
}

func (s *CypherSort) normalize() {
	sorts := make([]*OrderBy, 0)
	for _, sorting := range s.sorts {
		if sorting.valid {
			sorts = append(sorts, sorting)
		}
	}
	s.sorts = sorts
}

func (s *CypherSort) SortingCypherFragment(nodeAlias string) Cypher {
	s.normalize()
	if len(s.sorts) == 0 {
		return ""
	}
	var cypherStr strings.Builder
	cypherStr.WriteString(" ORDER BY ")
	inCoalesce := false
	isLast := false
	isNextCoalesce := false
	for i := 0; i < len(s.sorts); i++ {
		sortingProperty := s.sorts[i]
		isLast = i == len(s.sorts)-1
		isNextCoalesce = !isLast && s.sorts[i+1].coalesce
		if i > 0 {
			cypherStr.WriteString(" , ")
		}
		if !inCoalesce && sortingProperty.coalesce {
			cypherStr.WriteString("COALESCE(")
			inCoalesce = true
		}
		toLower := sortingProperty.supportCaseSensitive && !sortingProperty.caseSensitive
		if toLower {
			cypherStr.WriteString("toLower(")
		}
		cypherStr.WriteString(StringFirstNonEmpty(sortingProperty.nodeAlias, nodeAlias))
		cypherStr.WriteString(".")
		cypherStr.WriteString(sortingProperty.nodeProperty)
		if toLower {
			cypherStr.WriteString(")")
		}
		if inCoalesce && !isNextCoalesce {
			cypherStr.WriteString(")")
			inCoalesce = false
		}
		if !inCoalesce && sortingProperty.descending {
			cypherStr.WriteString(" DESC ")
		}
	}

	return Cypher(cypherStr.String())
}
