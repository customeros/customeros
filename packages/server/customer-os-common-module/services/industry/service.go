package industry

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
)

type industryService struct {
	log   logger.Logger
	neo4j *neo4j_repository.Repositories
}

func NewIndustryService(log logger.Logger, neo4j *neo4j_repository.Repositories) interfaces.IndustryService {
	return &industryService{
		log:   log,
		neo4j: neo4j,
	}
}

func (s *industryService) GetAllForOrganizationIds(ctx context.Context, organizationIds []string) (*neo4j_entity.IndustryEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IndustryService.GetAllForOrganizationIds")
	defer spans.Finish()

	spans.LogKV("organizationIds", fmt.Sprintf("%v", organizationIds))

	industryDbNodes, err := s.neo4j.IndustryReadRepository.GetAllForOrganizationIds(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	industryEntities := make(neo4j_entity.IndustryEntities, 0)
	for _, v := range industryDbNodes {
		industryEntity := neo4jmapper.MapDbNodeToIndustryEntity(v.Node)
		industryEntity.DataloaderKey = v.LinkedNodeId
		industryEntities = append(industryEntities, *industryEntity)
	}
	return &industryEntities, nil
}

// Returns the industry entity by code, nil if not found
func (s *industryService) GetByCode(ctx context.Context, code string) (*neo4j_entity.IndustryEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IndustryService.GetByCode")
	defer spans.Finish()

	spans.LogKV("code", code)

	industryDbNode, err := s.neo4j.IndustryReadRepository.GetByCode(ctx, code)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if industryDbNode == nil {
		return nil, nil
	}
	industryEntity := neo4jmapper.MapDbNodeToIndustryEntity(industryDbNode)
	return industryEntity, nil
}

func (s *industryService) GetClosestByCode(ctx context.Context, code string) (*neo4j_entity.IndustryEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IndustryService.GetClosestByCode")
	defer spans.Finish()

	spans.LogKV("code", code)

	queryCode := code
	// If the code is not found, strip the last character and try again until found or empty
	for {
		if len(queryCode) == 0 {
			break
		}
		industryEntity, err := s.GetByCode(ctx, queryCode)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}
		if industryEntity != nil {
			return industryEntity, nil
		}
		queryCode = queryCode[:len(queryCode)-1]
	}

	return nil, nil
}

func (s *industryService) GetInUseIndustries(ctx context.Context) (*neo4j_entity.IndustryEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IndustryService.GetInUseIndustries")
	defer spans.Finish()

	industryDbNodes, err := s.neo4j.IndustryReadRepository.GetInUseIndustries(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return nil, err
	}
	industryEntities := make(neo4j_entity.IndustryEntities, 0)
	for _, v := range industryDbNodes {
		industryEntity := neo4jmapper.MapDbNodeToIndustryEntity(v)
		industryEntities = append(industryEntities, *industryEntity)
	}
	return &industryEntities, nil
}
