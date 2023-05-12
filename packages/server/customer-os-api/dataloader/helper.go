package dataloader

import (
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
	"reflect"
)

func sortKeys(keys dataloader.Keys) ([]string, map[string]int) {
	var ids []string
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	for ix, key := range keys {
		ids = append(ids, key.String())
		keyOrder[key.String()] = ix
	}
	return ids, keyOrder
}

func assertEntitiesType(results []*dataloader.Result, expectedType reflect.Type) error {
	for _, res := range results {
		if reflect.TypeOf(res.Data) != expectedType {
			return errors.New(fmt.Sprintf("Not expected type: %v", reflect.TypeOf(res.Data)))
		}
	}
	return nil
}

func assertEntitiesPtrType(results []*dataloader.Result, expectedType reflect.Type, allowNils bool) error {
	for _, res := range results {
		if res.Data == nil && allowNils {
			continue
		} else if res.Data == nil && !allowNils {
			return errors.New("Not expected nil")
		}

		value := reflect.ValueOf(res.Data)
		if value.Kind() != reflect.Ptr {
			return errors.New(fmt.Sprintf("Not expected type: %v", value.Kind()))
		}

		elemType := value.Elem().Type()
		if elemType != expectedType {
			return errors.New(fmt.Sprintf("Not expected type: %v", elemType))
		}
	}
	return nil
}
