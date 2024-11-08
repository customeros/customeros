package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemService interface {
	MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error
	SetPrimaryExternalId(ctx context.Context, externalSystem, externalId string, linkWith LinkWith) error
	GetPrimaryExternalId(ctx context.Context, externalSystem, linkedWithId string, linkedWithEntityType model.EntityType) (string, error)
	GetExternalSystemsForEntities(ctx context.Context, ids []string, entityType model.EntityType) (*neo4jentity.ExternalSystemEntities, error)
}

type externalSystemService struct {
	log      logger.Logger
	services *Services
}

func NewExternalSystemService(log logger.Logger, services *Services) ExternalSystemService {
	return &externalSystemService{
		log:      log,
		services: services,
	}
}

func (s *externalSystemService) MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.MergeExternalSystem")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("externalSystem", externalSystem))

	if externalSystem == "" {
		return nil
	}

	err := s.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenant, externalSystem, externalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (s *externalSystemService) SetPrimaryExternalId(ctx context.Context, externalSystem, externalId string, linkWith LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.SetPrimaryExternalId")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	// Create external system if it doesn't exist
	err := s.MergeExternalSystem(ctx, tenant, externalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Set primary external id
	err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.SetPrimaryExternalId(ctx, tenant, externalSystem, externalId, linkWith.Type.Neo4jLabel(), linkWith.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Send completion event if link with is an organization
	if linkWith.Type == model.ORGANIZATION {
		utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), linkWith.Id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	}
	return nil
}

func (s *externalSystemService) GetExternalSystemsForEntities(ctx context.Context, ids []string, entityType model.EntityType) (*neo4jentity.ExternalSystemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.GetExternalSystemsForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	dbRecords, err := s.services.Neo4jRepositories.ExternalSystemReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), ids, entityType.Neo4jLabel())
	if err != nil {
		return nil, err
	}
	externalSystemEntities := make(neo4jentity.ExternalSystemEntities, 0, len(dbRecords))
	for _, v := range dbRecords {
		externalSystemEntity := neo4jmapper.MapDbNodeToExternalSystem(v.Node)
		neo4jmapper.AddDbRelationshipToExternalSystemEntity(*v.Relationship, externalSystemEntity)
		externalSystemEntity.DataloaderKey = v.LinkedNodeId
		externalSystemEntities = append(externalSystemEntities, *externalSystemEntity)
	}
	return &externalSystemEntities, nil
}

func (s *externalSystemService) GetPrimaryExternalId(ctx context.Context, externalSystem, linkedWithId string, linkedWithEntityType model.EntityType) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.GetPrimaryExternalId")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem), log.String("linkedWithId", linkedWithId), log.String("linkedWithEntityType", linkedWithEntityType.String()))

	externalSystemEntities, err := s.GetExternalSystemsForEntities(ctx, []string{linkedWithId}, linkedWithEntityType)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	primaryExternalId := ""
	if externalSystemEntities != nil {
		for _, externalSystemEntity := range *externalSystemEntities {
			if externalSystemEntity.ExternalSystemId.String() == externalSystem && externalSystemEntity.Relationship.Primary {
				primaryExternalId = externalSystemEntity.Relationship.ExternalId
			}
		}
	}
	return primaryExternalId, nil
}
