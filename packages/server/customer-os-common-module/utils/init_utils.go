package utils

import (
	"log"
	"reflect"
)

func IsInitialized(data any) bool {
	v := reflect.ValueOf(data)

	// Check if input is a pointer
	if v.Kind() != reflect.Ptr {
		log.Printf("Input must be a pointer to struct")
		return false
	}

	// Check if pointer is nil
	if v.IsNil() {
		log.Printf("Input is nil")
		return false
	}

	val := v.Elem()
	// Check if pointing to a struct
	if val.Kind() != reflect.Struct {
		log.Printf("Input must be a pointer to struct")
		return false
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		if field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
			if field.IsNil() {
				log.Printf("Field %s is not initialized", fieldName)
				return false
			}
		}
	}
	return true
}
