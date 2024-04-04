package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func buildSort(sortBy []*model.SortBy, T reflect.Type) (*utils.CypherSort, error) {
	transformedSorting := new(utils.CypherSort)
	if sortBy != nil {
		for _, v := range sortBy {
			orderBy := transformedSorting.NewSortRule(v.By, v.Direction.String(), utils.IfNotNilBool(v.CaseSensitive), T)
			if !orderBy.IsValid() {
				return nil, errors.New("Invalid sorting rule")
			}
		}
	}
	return transformedSorting, nil
}

type SortMultipleEntitiesDefinition struct {
	EntityPrefix   string
	EntityMapping  reflect.Type
	EntityAlias    string
	EntityDefaults []SortMultipleEntitiesDefinitionDefault
}

type SortMultipleEntitiesDefinitionDefault struct {
	PropertyName string
	AscDefault   string
	DescDefault  string
}

func buildSortMultipleEntities(sortBy []*model.SortBy, mapping []SortMultipleEntitiesDefinition) (*utils.Cypher, error) {
	transformedSorting := new(utils.CypherSort)

	if sortBy != nil {
		for _, v := range sortBy {

			entity, sort, found := strings.Cut(v.By, "_")
			if !found {
				continue
			}

			var mappingFound *SortMultipleEntitiesDefinition
			for _, mapping := range mapping {
				if entity == mapping.EntityPrefix {
					mappingFound = &mapping
					break
				}
			}

			if mappingFound == nil {
				return nil, errors.New("Entity not found in mapping")
			}

			orderBy := transformedSorting.NewSortRule(sort, v.Direction.String(), utils.IfNotNilBool(v.CaseSensitive), mappingFound.EntityMapping)
			if !orderBy.IsValid() {
				return nil, errors.New("Invalid sorting rule")
			}

			var defaultIfNil string
			for _, value := range mappingFound.EntityDefaults {
				if value.PropertyName == sort {
					if v.Direction.String() == "ASC" {
						defaultIfNil = value.AscDefault
					} else {
						defaultIfNil = value.DescDefault
					}
					break
				}
			}

			var aliases []string
			for _, v := range mapping {
				aliases = append(aliases, v.EntityAlias)
			}

			//generating cypher frangment like below
			//WITH c, i ORDER BY i.status DESC
			//WITH c, i WITH c, i , CASE WHEN c.endedAt IS NULL THEN date('2100-01-01') ELSE c.endedAt END AS endedAt_FOR_SORTING ORDER BY endedAt_FOR_SORTING
			fragment := transformedSorting.SortingCypherFragmentWithDefaultIfNil(strings.Join(aliases, ", "), mappingFound.EntityAlias, defaultIfNil)
			return &fragment, nil
		}
	}
	return nil, nil
}
