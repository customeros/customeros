package utils

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	TagKey                  = "neo4jDb"
	TagLookupName           = "lookupName"
	TagProperty             = "property"
	TagSupportCaseSensitive = "supportCaseSensitive"
)

func GetPropertyDetailsByLookupName(T reflect.Type, lookupName string) (map[string]string, error) {
	for i := 0; i < T.NumField(); i++ {
		structField := T.Field(i)

		switch structField.Type.Kind() {
		case reflect.Struct:
			m, _ := GetPropertyDetailsByLookupName(structField.Type, lookupName)
			if m != nil {
				return m, nil
			}
		default:
			tag, ok := structField.Tag.Lookup(TagKey)
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
					if val, ok := m[TagLookupName]; ok {
						if val == lookupName {
							return m, nil
						}
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("Given field %s not found", lookupName)
}
