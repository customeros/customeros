package externalsystem

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

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
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemService.MergeExternalSystem")
	defer spans.Finish()

	spans.LogKV("externalSystem", externalSystem)

	if externalSystem == "" {
		return nil
	}

	err := s.neo4j.ExternalSystemWriteRepository.CreateIfNotExists(ctx, nil, tenant, externalSystem, externalSystem)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (s *externalSystemService) SetPrimaryExternalId(ctx context.Context, externalSystem, externalId string, linkWith common_srv.LinkWith) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemService.SetPrimaryExternalId")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	// Create external system if it doesn't exist
	err := s.MergeExternalSystem(ctx, tenant, externalSystem)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// Set primary external id
	err = s.neo4j.ExternalSystemWriteRepository.SetPrimaryExternalId(ctx, tenant, externalSystem, externalId, linkWith.Type.Neo4jLabel(), linkWith.Id)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// Send completion event if link with is an organization
	if linkWith.Type == model.ORGANIZATION {
		s.events.Publisher.PublishNotification(ctx, tenant, linkWith.Id, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
	}
	return nil
}

func (s *externalSystemService) GetExternalSystemsForEntities(ctx context.Context, ids []string, entityType model.EntityType) (*neo4jentity.ExternalSystemEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemService.GetExternalSystemsForEntities")
	defer spans.Finish()

	spans.LogFields(log.Object("ids", ids))

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
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemService.GetPrimaryExternalId")
	defer spans.Finish()

	spans.LogKV("externalSystem", externalSystem, "linkedWithId", linkedWithId, "linkedWithEntityType", linkedWithEntityType.String())

	externalSystemEntities, err := s.GetExternalSystemsForEntities(ctx, []string{linkedWithId}, linkedWithEntityType)
	if err != nil {
		spans.TraceError(err)
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
