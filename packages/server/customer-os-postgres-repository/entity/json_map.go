package postgres_entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONMap is a helper type for JSONB columns that automatically unmarshals
// into a map[string]interface{}.
type JSONMap map[string]interface{}

// Scan implements the sql.Scanner interface.
func (j *JSONMap) Scan(src interface{}) error {
	if src == nil {
		*j = nil
		return nil
	}

	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
	return json.Unmarshal(data, j)
}

// Value implements the driver.Valuer interface.
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
