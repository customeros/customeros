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
	lookupName = strings.ToUpper(lookupName)
	for i := 0; i < T.NumField(); i++ {
		structField := T.Field(i)

		// if field is struct, but not time.Time, then recursively call this function
		if structField.Type.Kind() == reflect.Struct &&
			!(structField.Type.Name() == "Time" && structField.Type.PkgPath() == "time") {
			m, _ := GetPropertyDetailsByLookupName(structField.Type, lookupName)
			if m != nil {
				return m, nil
			}
		} else {
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
	return nil, fmt.Errorf("given field %s not found", lookupName)
}
