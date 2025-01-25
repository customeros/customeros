package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// StructToMap converts any struct to a map[string]interface{}
// based on JSON tags and fields.
func StructToMap(input interface{}) (map[string]interface{}, error) {
	// First, marshal the struct to JSON
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Then unmarshal into a map
	var output map[string]interface{}
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// MapToStruct takes a map[string]interface{} and fills 'outStruct'.
// 'outStruct' must be a pointer to a struct. JSON tags determine which
// map keys populate which struct fields.
func MapToStruct(input map[string]any, outStruct any) error {
	// Ensure outStruct is a pointer
	v := reflect.ValueOf(outStruct)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("outStruct must be a non-nil pointer to a struct")
	}

	// Marshal the map to JSON
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}

	// Unmarshal JSON back into 'outStruct'
	err = json.Unmarshal(data, outStruct)
	if err != nil {
		return err
	}

	return nil
}
