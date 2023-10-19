package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
	"reflect"
)

func buildSort(sortBy []*model.SortBy, T reflect.Type) (*utils.CypherSort, error) {
	transformedSorting := new(utils.CypherSort)
	if sortBy != nil {
		for _, v := range sortBy {
			orderBy := transformedSorting.NewSortRule(v.By, v.Direction.String(), *v.CaseSensitive, T)
			if !orderBy.IsValid() {
				return nil, errors.New("Invalid sorting rule")
			}
		}
	}
	return transformedSorting, nil
}
