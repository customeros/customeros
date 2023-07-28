package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j/entity"
)

type CacheService interface {
	InitCache()
	GetStates() []*model.GCliItem
}

type cacheService struct {
	services *Services

	States                    []*model.GCliItem
	OrganizationRelationships []*model.GCliItem
}

func NewCacheService(services *Services) CacheService {
	return &cacheService{
		services: services,
	}
}

func (s *cacheService) InitCache() {

	//cache US states for the gCliCache
	gCliStatesCache := make([]*model.GCliItem, 0)
	gCliOrganizationRelationshipsCache := make([]*model.GCliItem, 0)

	countries := []*commonEntity.CountryEntity{}
	countries = append(countries, &commonEntity.CountryEntity{Id: "1", CodeA3: "USA"})

	for _, country := range countries {
		states, err := s.services.CommonServices.StateService.GetStatesByCountryId(context.Background(), country.Id)
		if err != nil {
			//todo: log error
		}

		for _, v := range states {
			item := mapper.MapStateToGCliItem(*v)
			gCliStatesCache = append(gCliStatesCache, &item)
		}
	}

	for _, organizationRelationship := range entity.AllOrganizationRelationship {
		item := mapper.MapOrganizationRelationshipToGCliItem(organizationRelationship)
		gCliOrganizationRelationshipsCache = append(gCliOrganizationRelationshipsCache, &item)
	}

	s.States = gCliStatesCache
	s.OrganizationRelationships = gCliOrganizationRelationshipsCache
}

func (s *cacheService) GetStates() []*model.GCliItem {
	return s.States
}
