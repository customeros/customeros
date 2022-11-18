package utils

import (
	"reflect"
)

type Filters struct {
	filters []*Filter
}

type ComparisonOperator int

const (
	EQUALS ComparisonOperator = iota
	CONTAINS
)

type Filter struct {
	nodeProperty         string
	supportCaseSensitive bool
	caseSensitive        bool
	dbNodePropertyProps  map[string]string
}

func (f *Filters) NewFilterRule(lookupName string, caseSensitive bool, T reflect.Type) error {
	filter := new(Filter)
	props, err := getPropertyDetailsByLookupName(T, lookupName)
	if err != nil {
		return err
	}
	filter.nodeProperty = props[tagProperty]
	filter.supportCaseSensitive = props[tagSupportCaseSensitive] == "true"
	filter.caseSensitive = caseSensitive
	f.filters = append(f.filters, filter)
	return nil
}

func (f *Filters) FilterCypherFragment(nodeAlias string) string {
	if len(f.filters) == 0 {
		return ""
	}
	query := " WHERE "
	for i := 0; i < len(f.filters); i++ {
		filterProperty := f.filters[i]
		if i > 0 {
			query += " , "
		}
		toLower := filterProperty.supportCaseSensitive && !filterProperty.caseSensitive
		if toLower {
			query += "toLower("
		}
		query += nodeAlias
		query += "."
		query += filterProperty.nodeProperty
		if toLower {
			query += ")"
		}
	}
	return query
}
