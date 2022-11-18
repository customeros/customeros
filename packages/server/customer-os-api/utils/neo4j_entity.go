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

func getPropertyDetailsByLookupName(T reflect.Type, lookupName string) (map[string]string, error) {
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
					if val == lookupName {
						return m, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("Given field %s not found", lookupName)
}
