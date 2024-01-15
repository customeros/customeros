package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
)

type CacheService interface {
	InitCache()
	GetStates() []*model.GCliItem
}

type cacheService struct {
	services *Services

	States []*model.GCliItem
}

func NewCacheService(services *Services) CacheService {
	return &cacheService{
		services: services,
	}
}

func (s *cacheService) InitCache() {

	//cache US states for the gCliCache
	gCliStatesCache := make([]*model.GCliItem, 0)

	//read states from db for USA
	states, err := s.services.CommonServices.StateService.GetStatesByCountryId(context.Background(), "1")
	if err != nil {
		//todo: log error
	}

	for _, v := range states {
		item := mapper.MapStateToGCliItem(*v)
		gCliStatesCache = append(gCliStatesCache, &item)
	}

	s.States = gCliStatesCache
}

func (s *cacheService) GetStates() []*model.GCliItem {
	return s.States
}
