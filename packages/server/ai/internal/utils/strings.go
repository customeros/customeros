package utils

import "reflect"

func IfNotNilString(check any, valueExtractor ...func() string) string {
	if reflect.ValueOf(check).Kind() == reflect.String {
		return check.(string)
	}
	if reflect.ValueOf(check).Kind() == reflect.Pointer && reflect.ValueOf(check).IsNil() {
		return ""
	}
	if len(valueExtractor) > 0 {
		return valueExtractor[0]()
	}
	out := check.(*string)
	return *out
}
