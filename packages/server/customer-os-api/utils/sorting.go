package utils

import (
	"fmt"
	"reflect"
	"strings"
)

type Sortings struct {
	properties []*Sorting
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
	sorting.nodeProperty = props["property"]
	sorting.descending = "DESC" == direction
	sorting.supportCaseSensitive = props["supportCaseSensitive"] == "true"
	sorting.caseSensitive = caseSensitive
	s.properties = append(s.properties, sorting)
	return nil
}

func (s *Sortings) CypherFragment(nodeAlias string) string {
	if len(s.properties) == 0 {
		return ""
	}
	query := " ORDER BY "
	for i := 0; i < len(s.properties); i++ {
		sortingProperty := s.properties[i]
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
		tag, ok := structField.Tag.Lookup("neo4jDb")
		if ok {
			tags := strings.Split(tag, ";")
			if len(tags) > 0 {
				m := make(map[string]string)
				for _, v := range tags {
					kvs := strings.Split(v, ":")
					m[kvs[0]] = kvs[1]
				}
				if val, ok := m["exposedName"]; ok {
					if val == exposedName {
						return m, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("Given field %s not identified", exposedName)
}
