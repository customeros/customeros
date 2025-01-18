package interaction_session

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/exp/slices"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type interactionSessionService struct {
	neo4j *neoRepo.Repositories
}

func NewInteractionSessionService(neo4j *neoRepo.Repositories) interfaces.InteractionSessionService {
	return &interactionSessionService{
		neo4j: neo4j,
	}
}

func (s *interactionSessionService) GetById(ctx context.Context, id string) (*neo4jentity.InteractionSessionEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionService.GetById")
	defer span.Finish()

	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	byId, err := s.neo4j.CommonReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), id, commonModel.NodeLabelInteractionSession)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if byId == nil {
		return nil, nil
	}

	return neo4jmapper.MapDbNodeToInteractionSessionEntity(byId), nil
}

func (s *interactionSessionService) GetAttendedByParticipantsForInteractionSessions(ctx context.Context, ids []string) (*neo4jentity.InteractionSessionParticipants, error) {
	records, err := s.neo4j.InteractionSessionReadRepository.GetAttendedByParticipantsForInteractionSessions(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToInteractionSessionParticipants(records)

	return &interactionEventParticipants, nil
}

func (s *interactionSessionService) GetInteractionSessionsForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.InteractionSessionEntities, error) {
	interactionSessions, err := s.neo4j.InteractionSessionReadRepository.GetAllForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	interactionSessionEntities := neo4jentity.InteractionSessionEntities{}
	for _, v := range interactionSessions {
		interactionSessionEntity := neo4jmapper.MapDbNodeToInteractionSessionEntity(v.Node)
		interactionSessionEntity.DataloaderKey = v.LinkedNodeId
		interactionSessionEntities = append(interactionSessionEntities, *interactionSessionEntity)
	}
	return &interactionSessionEntities, nil
}

func (s *interactionSessionService) Create(ctx context.Context, newInteractionSession *neo4jentity.InteractionSessionEntity) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		id, err := s.CreateInTx(ctx, tx, newInteractionSession)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return id, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//TODO EDI
	//for _, v := range newInteractionSession.AttendedBy {
	//	if v.EntityId != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, *v.EntityId)
	//	}
	//	if v.Email != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByEmail(ctx, *v.Email)
	//	}
	//	if v.PhoneNumber != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByPhoneNumber(ctx, *v.PhoneNumber)
	//	}
	//}

	return queryResult.(*string), nil
}

func (s *interactionSessionService) CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, data *neo4jentity.InteractionSessionEntity) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionService.CreateInTx")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	interactionSessionId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelInteractionSession)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	err = s.neo4j.InteractionSessionWriteRepository.CreateInTx(ctx, tx, common.GetTenantFromContext(ctx), interactionSessionId, *data)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &interactionSessionId, nil
}

func (s *interactionSessionService) mapDbRelationshipToParticipantDetails(relationship dbtype.Relationship) neo4jentity.InteractionSessionParticipantDetails {
	props := utils.GetPropsFromRelationship(relationship)
	details := neo4jentity.InteractionSessionParticipantDetails{
		Type: utils.GetStringPropOrEmpty(props, "type"),
	}
	return details
}

func (s *interactionSessionService) convertDbNodesToInteractionSessionParticipants(records []*utils.DbNodeWithRelationAndId) neo4jentity.InteractionSessionParticipants {
	interactionSessionParticipants := neo4jentity.InteractionSessionParticipants{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, commonModel.NodeLabelEmail) {
			participant := neo4jmapper.MapDbNodeToEmailEntity(v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionSessionParticipants = append(interactionSessionParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelPhoneNumber) {
			participant := neo4jmapper.MapDbNodeToPhoneNumberEntity(v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionSessionParticipants = append(interactionSessionParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelUser) {
			participant := neo4jmapper.MapDbNodeToUserEntity(v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionSessionParticipants = append(interactionSessionParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelContact) {
			participant := neo4jmapper.MapDbNodeToContactEntity(v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionSessionParticipants = append(interactionSessionParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelJobRole) {
			participant := neo4jmapper.MapDbNodeToJobRoleEntity(v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionSessionParticipants = append(interactionSessionParticipants, participant)
		}
	}
	return interactionSessionParticipants
}

func (s *interactionSessionService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.neo4j.Neo4jDriver
}
