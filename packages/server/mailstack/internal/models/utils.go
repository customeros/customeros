package models

import (
	"database/sql/driver"
	"encoding/json"
)

// JSONMap represents a JSON object that can be stored in PostgreSQL
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface for JSONMap
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONMap
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}
