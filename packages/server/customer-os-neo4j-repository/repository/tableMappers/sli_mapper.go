package tableMappers

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository/types"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func SliToTable(table *godog.Table) []types.SLI {
	header := make([]string, len(table.Rows[0].Cells))
	for i, cell := range table.Rows[0].Cells {
		header[i] = capitalizeFirstLetter(cell.Value)
	}

	var sliArray []types.SLI
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i].Cells
		item := types.SLI{}

		for j := 0; j < len(header); j++ {
			fieldName := header[j]
			fieldValue := row[j].Value

			// Use reflection to set the struct field dynamically
			setStructField(&item, fieldName, fieldValue)
		}

		sliArray = append(sliArray, item)
	}
	return sliArray
}
func setStructField(item *types.SLI, fieldName string, fieldValue string) error {
	// Use reflection to set the struct field dynamically
	structValue := reflect.ValueOf(item).Elem()
	field := structValue.FieldByName(fieldName)

	if !field.IsValid() {
		return nil // Skip if the field doesn't exist in the struct
	}

	// Convert the table cell value to the appropriate type of the struct field
	switch field.Kind() {
	case reflect.String:
		field.SetString(fieldValue)
	case reflect.Int, reflect.Int64:
		intValue, err := strconv.Atoi(fieldValue)
		if err == nil {
			field.SetInt(int64(intValue))
		}
	case reflect.Float64:
		intValue, err := strconv.Atoi(fieldValue)
		if err == nil {
			field.SetFloat(float64(intValue))
		}
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			// Assuming the date format is "2006-01-02"
			dateValue, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", fieldValue)

			if err == nil {
				field.Set(reflect.ValueOf(dateValue))
			}
		} else {
			return fmt.Errorf("Unsupported struct field type for %s", fieldName)
		}

	}
	return nil
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
