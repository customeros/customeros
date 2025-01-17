package api_cache

import (
	"context"

	mapper2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
	cosapi_interfaces "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
)

type cacheService struct {
	States []*model.GCliItem
	neo4j  *repository.Repositories
}

func NewCacheService(neo4j *repository.Repositories) cosapi_interfaces.CacheService {
	return &cacheService{
		neo4j: neo4j,
	}
}

func (s *cacheService) InitCache() {
	// cache US states for the gCliCache
	gCliStatesCache := make([]*model.GCliItem, 0)

	// read states from db for USA
	states, err := s.neo4j.StateReadRepository.GetStatesByCountryId(context.Background(), "1")
	if err != nil {
		// todo: log error
	}

	for _, v := range states {
		stateEntity := mapper2.MapDbNodeToStateEntity(*v)
		item := mapper.MapStateToGCliItem(*stateEntity)
		gCliStatesCache = append(gCliStatesCache, &item)
	}

	s.States = gCliStatesCache
}

func (s *cacheService) GetStates() []*model.GCliItem {
	return s.States
}
