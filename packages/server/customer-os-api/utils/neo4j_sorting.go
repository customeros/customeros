package utils

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tagKey                  = "neo4jDb"
	tagLookupName           = "lookupName"
	tagProperty             = "property"
	tagSupportCaseSensitive = "supportCaseSensitive"
)

type Sortings struct {
	sorting []*Sorting
}

type Sorting struct {
	nodeProperty         string
	supportCaseSensitive bool
	caseSensitive        bool
	descending           bool
	dbNodePropertyProps  map[string]string
}

func (s *Sortings) NewSortingProperty(exposedName, direction string, caseSensitive bool, T reflect.Type) error {
	sorting := new(Sorting)
	props, err := getFieldByExposedName(T, exposedName)
	if err != nil {
		return err
	}
	sorting.nodeProperty = props[tagProperty]
	sorting.descending = "DESC" == direction
	sorting.supportCaseSensitive = props[tagSupportCaseSensitive] == "true"
	sorting.caseSensitive = caseSensitive
	s.sorting = append(s.sorting, sorting)
	return nil
}

func (s *Sortings) CypherFragment(nodeAlias string) string {
	if len(s.sorting) == 0 {
		return ""
	}
	query := " ORDER BY "
	for i := 0; i < len(s.sorting); i++ {
		sortingProperty := s.sorting[i]
		if i > 0 {
			query += " , "
		}
		toLower := sortingProperty.supportCaseSensitive && !sortingProperty.caseSensitive
		if toLower {
			query += "toLower("
		}
		query += nodeAlias
		query += "."
		query += sortingProperty.nodeProperty
		if toLower {
			query += ")"
		}
		if sortingProperty.descending {
			query += " DESC "
		}
	}
	return query
}

func getFieldByExposedName(T reflect.Type, exposedName string) (map[string]string, error) {
	for i := 0; i < T.NumField(); i++ {
		structField := T.Field(i)
		tag, ok := structField.Tag.Lookup(tagKey)
		if ok {
			tags := strings.Split(tag, ";")
			if len(tags) > 0 {
				m := make(map[string]string)
				for _, v := range tags {
					kvs := strings.Split(v, ":")
					if len(kvs) == 2 {
						m[kvs[0]] = kvs[1]
					}
				}
				if val, ok := m[tagLookupName]; ok {
					if val == exposedName {
						return m, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("Given field %s not identified", exposedName)
}
