package utils

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func AnyToJSONB(data any) (datatypes.JSON, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	jsonData := datatypes.JSON(jsonBytes)
	return jsonData, err
}

func JSONBToAny(jsonData datatypes.JSON) (any, error) {
	var data any
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
