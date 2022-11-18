package utils

import (
	"reflect"
	"strings"
)

type Sorts struct {
	sorts []*orderBy
}

type orderBy struct {
	nodeProperty         string
	supportCaseSensitive bool
	caseSensitive        bool
	descending           bool
	dbNodePropertyProps  map[string]string
}

func (s *Sorts) NewSortRule(lookupName, direction string, caseSensitive bool, T reflect.Type) error {
	orderBy := new(orderBy)

	props, err := getPropertyDetailsByLookupName(T, lookupName)
	if err != nil {
		return err
	}

	orderBy.nodeProperty = props[tagProperty]
	orderBy.descending = "DESC" == direction
	orderBy.supportCaseSensitive = props[tagSupportCaseSensitive] == "true"
	orderBy.caseSensitive = caseSensitive

	s.sorts = append(s.sorts, orderBy)

	return nil
}

func (s *Sorts) SortingCypherFragment(nodeAlias string) string {
	if len(s.sorts) == 0 {
		return ""
	}
	var query strings.Builder
	query.WriteString(" ORDER BY ")
	for i := 0; i < len(s.sorts); i++ {
		sortingProperty := s.sorts[i]
		if i > 0 {
			query.WriteString(" , ")
		}
		toLower := sortingProperty.supportCaseSensitive && !sortingProperty.caseSensitive
		if toLower {
			query.WriteString("toLower(")
		}
		query.WriteString(nodeAlias)
		query.WriteString(".")
		query.WriteString(sortingProperty.nodeProperty)
		if toLower {
			query.WriteString(")")
		}
		if sortingProperty.descending {
			query.WriteString(" DESC ")
		}
	}
	return query.String()
}
