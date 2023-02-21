package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"reflect"
)

func buildSort(sortBy []*model.SortBy, T reflect.Type) (*utils.CypherSort, error) {
	transformedSorting := new(utils.CypherSort)
	if sortBy != nil {
		for _, v := range sortBy {
			err := transformedSorting.NewSortRule(v.By, v.Direction.String(), *v.CaseSensitive, T)
			if err != nil {
				return nil, err
			}
		}
	}
	return transformedSorting, nil
}
