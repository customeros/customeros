package api_country

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type countryService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCountryService(log logger.Logger, repository *repository.Repositories) cosapi_interfaces.CountryService {
	return &countryService{
		log:          log,
		repositories: repository,
	}
}

func (s *countryService) GetCountriesForPhoneNumbers(ctx context.Context, ids []string) (*neo4j_entity.CountryEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryService.GetCountriesForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	countryDbNodes, err := s.repositories.Neo4jRepositories.CountryReadRepository.GetAllForPhoneNumbers(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	CountryEntities := make(neo4j_entity.CountryEntities, 0, len(countryDbNodes))
	for _, v := range countryDbNodes {
		countryEntity := s.mapDbNodeToCountryEntity(*v.Node)
		countryEntity.DataloaderKey = v.LinkedNodeId
		CountryEntities = append(CountryEntities, *countryEntity)
	}
	return &CountryEntities, nil
}

func (s *countryService) mapDbNodeToCountryEntity(node dbtype.Node) *neo4j_entity.CountryEntity {
	props := utils.GetPropsFromNode(node)
	return &neo4j_entity.CountryEntity{
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CodeA2:    utils.GetStringPropOrEmpty(props, "codeA2"),
		CodeA3:    utils.GetStringPropOrEmpty(props, "codeA3"),
		PhoneCode: utils.GetStringPropOrEmpty(props, "phoneCode"),
	}
}
