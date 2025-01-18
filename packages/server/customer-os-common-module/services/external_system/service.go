package externalsystem

import (
	"context"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type externalSystemService struct {
	log    logger.Logger
	neo4j  *neo4j_repository.Repositories
	events *events.EventsService
}

func NewExternalSystemService(log logger.Logger, neo4j *neo4j_repository.Repositories, events *events.EventsService) interfaces.ExternalSystemService {
	return &externalSystemService{
		log:    log,
		neo4j:  neo4j,
		events: events,
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

	err := s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenant, externalSystem, externalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (s *externalSystemService) SetPrimaryExternalId(ctx context.Context, externalSystem, externalId string, linkWith common_srv.LinkWith) error {
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
	err = s.neo4j.ExternalSystemWriteRepository.SetPrimaryExternalId(ctx, tenant, externalSystem, externalId, linkWith.Type.Neo4jLabel(), linkWith.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Send completion event if link with is an organization
	if linkWith.Type == model.ORGANIZATION {
		s.events.Publisher.PublishEventCompleted(ctx, tenant, linkWith.Id, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
	}
	return nil
}

func (s *externalSystemService) GetExternalSystemsForEntities(ctx context.Context, ids []string, entityType model.EntityType) (*neo4jentity.ExternalSystemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.GetExternalSystemsForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	dbRecords, err := s.neo4j.ExternalSystemReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), ids, entityType.Neo4jLabel())
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
