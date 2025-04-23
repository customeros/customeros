package utils

import (
	"reflect"
)

func IsInitialized(data any, skipTypes ...reflect.Type) bool {
	if data == nil {
		return false
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := field.Type()

		// Skip specified types
		shouldSkip := false
		for _, skipType := range skipTypes {
			if fieldType == skipType {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			continue
		}

		if (field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface) && field.IsNil() {
			return false
		}
	}

	return true
}
