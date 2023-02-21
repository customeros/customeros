package utils

import (
	"reflect"
	"strings"
)

type CypherSort struct {
	sorts []*orderBy
}

type orderBy struct {
	nodeProperty         string
	supportCaseSensitive bool
	caseSensitive        bool
	descending           bool
	dbNodePropertyProps  map[string]string
}

func (s *CypherSort) NewSortRule(lookupName, direction string, caseSensitive bool, T reflect.Type) error {
	orderBy := new(orderBy)

	props, err := GetPropertyDetailsByLookupName(T, lookupName)
	if err != nil {
		return err
	}

	orderBy.nodeProperty = props[TagProperty]
	orderBy.descending = "DESC" == direction
	orderBy.supportCaseSensitive = props[TagSupportCaseSensitive] == "true"
	orderBy.caseSensitive = caseSensitive

	s.sorts = append(s.sorts, orderBy)

	return nil
}

func (s *CypherSort) SortingCypherFragment(nodeAlias string) Cypher {
	if len(s.sorts) == 0 {
		return ""
	}
	var cypherStr strings.Builder
	cypherStr.WriteString(" ORDER BY ")
	for i := 0; i < len(s.sorts); i++ {
		sortingProperty := s.sorts[i]
		if i > 0 {
			cypherStr.WriteString(" , ")
		}
		toLower := sortingProperty.supportCaseSensitive && !sortingProperty.caseSensitive
		if toLower {
			cypherStr.WriteString("toLower(")
		}
		cypherStr.WriteString(nodeAlias)
		cypherStr.WriteString(".")
		cypherStr.WriteString(sortingProperty.nodeProperty)
		if toLower {
			cypherStr.WriteString(")")
		}
		if sortingProperty.descending {
			cypherStr.WriteString(" DESC ")
		}
	}

	return Cypher(cypherStr.String())
}
