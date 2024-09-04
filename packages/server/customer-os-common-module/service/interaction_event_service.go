package service

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/exp/slices"
)

type InteractionEventService interface {
	GetById(ctx context.Context, id string) (*neo4jentity.InteractionEventEntity, error)
	GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error)
	GetInteractionEventsForMeetings(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error)
	GetInteractionEventsForIssues(ctx context.Context, issueIds []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error)
	GetSentByParticipantsForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.InteractionEventParticipants, error)
	GetSentToParticipantsForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.InteractionEventParticipants, error)
	GetReplyToInteractionsEventForInteractionEvents(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error)

	Create(ctx context.Context, data *InteractionEventCreateData) (*string, error)
	CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, data *InteractionEventCreateData) (*string, error)
}

type InteractionEventCreateData struct {
	InteractionEventEntity *neo4jentity.InteractionEventEntity
	SessionIdentifier      *string
	MeetingIdentifier      *string
	RepliesTo              *string
	SentBy                 []InteractionEventParticipantData
	SentTo                 []InteractionEventParticipantData
	SentCc                 []InteractionEventParticipantData
	SentBcc                []InteractionEventParticipantData
	ExternalSystem         *neo4jentity.ExternalSystemEntity
	Source                 neo4jentity.DataSource
	SourceOfTruth          neo4jentity.DataSource
}

type InteractionEventParticipantData struct {
	Email       *string
	PhoneNumber *string
	ContactId   *string
	UserId      *string
}

type interactionEventService struct {
	services *Services
}

func NewInteractionEventService(services *Services) InteractionEventService {
	return &interactionEventService{
		services: services,
	}
}

func (s *interactionEventService) GetById(ctx context.Context, id string) (*neo4jentity.InteractionEventEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetById")
	defer span.Finish()

	byId, err := s.services.Neo4jRepositories.CommonReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), id, commonModel.NodeLabelInteractionEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if byId == nil {
		return nil, nil
	}

	return neo4jmapper.MapDbNodeToInteractionEventEntity(byId), nil
}

func (s *interactionEventService) GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetInteractionEventsForInteractionSessions")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	interactionEvents, err := s.services.Neo4jRepositories.InteractionEventReadRepository.GetAllForInteractionSessions(ctx, common.GetTenantFromContext(ctx), ids, loadContent)
	if err != nil {
		return nil, err
	}
	interactionEventEntities := neo4jentity.InteractionEventEntities{}
	for _, v := range interactionEvents {
		interactionEventEntity := neo4jmapper.MapDbPropsToInteractionEventEntity(v.Props)
		interactionEventEntity.DataloaderKey = v.LinkedNodeId
		interactionEventEntities = append(interactionEventEntities, *interactionEventEntity)
	}
	return &interactionEventEntities, nil
}

func (s *interactionEventService) GetInteractionEventsForMeetings(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetInteractionEventsForMeetings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	interactionEvents, err := s.services.Neo4jRepositories.InteractionEventReadRepository.GetAllForMeetings(ctx, common.GetTenantFromContext(ctx), ids, loadContent)
	if err != nil {
		return nil, err
	}
	interactionEventEntities := neo4jentity.InteractionEventEntities{}
	for _, v := range interactionEvents {
		interactionEventEntity := neo4jmapper.MapDbPropsToInteractionEventEntity(v.Props)
		interactionEventEntity.DataloaderKey = v.LinkedNodeId
		interactionEventEntities = append(interactionEventEntities, *interactionEventEntity)
	}
	return &interactionEventEntities, nil
}

func (s *interactionEventService) GetInteractionEventsForIssues(ctx context.Context, issueIds []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetInteractionEventsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	interactionEvents, err := s.services.Neo4jRepositories.InteractionEventReadRepository.GetAllForIssues(ctx, common.GetTenantFromContext(ctx), issueIds, loadContent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	interactionEventEntities := neo4jentity.InteractionEventEntities{}
	for _, v := range interactionEvents {
		interactionEventEntity := neo4jmapper.MapDbPropsToInteractionEventEntity(v.Props)
		interactionEventEntity.DataloaderKey = v.LinkedNodeId
		interactionEventEntities = append(interactionEventEntities, *interactionEventEntity)
	}
	span.LogFields(log.Int("result count", len(interactionEventEntities)))
	return &interactionEventEntities, nil
}

func (s *interactionEventService) GetSentByParticipantsForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.InteractionEventParticipants, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetSentByParticipantsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	records, err := s.services.Neo4jRepositories.InteractionEventReadRepository.GetSentByFor(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToInteractionEventParticipants(records)

	span.LogFields(log.Int("result count", len(interactionEventParticipants)))

	return &interactionEventParticipants, nil
}

func (s *interactionEventService) GetSentToParticipantsForInteractionEvents(ctx context.Context, ids []string) (*neo4jentity.InteractionEventParticipants, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetSentToParticipantsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	records, err := s.services.Neo4jRepositories.InteractionEventReadRepository.GetSentToFor(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToInteractionEventParticipants(records)

	span.LogFields(log.Int("result count", len(interactionEventParticipants)))

	return &interactionEventParticipants, nil
}

func (s *interactionEventService) GetReplyToInteractionsEventForInteractionEvents(ctx context.Context, ids []string, loadContent bool) (*neo4jentity.InteractionEventEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetReplyToInteractionsEventForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	records, err := s.services.Neo4jRepositories.InteractionEventReadRepository.GetReplyToFor(ctx, common.GetTenantFromContext(ctx), ids, loadContent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	interactionEvents := neo4jentity.InteractionEventEntities{}
	for _, v := range records {
		event := neo4jmapper.MapDbPropsToInteractionEventEntity(v.Props)
		event.DataloaderKey = v.LinkedNodeId
		interactionEvents = append(interactionEvents, *event)
	}

	return &interactionEvents, nil
}

func (s *interactionEventService) Create(ctx context.Context, data *InteractionEventCreateData) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	interactionEventId, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		interactionEventId, err := s.CreateInTx(ctx, tx, data)
		return interactionEventId, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	//TODO EDI
	//for _, v := range newInteractionEvent.SentBy {
	//	if v.ContactId != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, *v.ContactId)
	//	}
	//	if v.Email != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByEmail(ctx, *v.Email)
	//	}
	//	if v.PhoneNumber != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByPhoneNumber(ctx, *v.PhoneNumber)
	//	}
	//}
	//for _, v := range newInteractionEvent.SentTo {
	//	if v.ContactId != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, *v.ContactId)
	//	}
	//	if v.Email != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByEmail(ctx, *v.Email)
	//	}
	//	if v.PhoneNumber != nil {
	//		s.services.OrganizationService.UpdateLastTouchpointByPhoneNumber(ctx, *v.PhoneNumber)
	//	}
	//}

	return interactionEventId.(*string), nil
}

func (s *interactionEventService) CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, newInteractionEvent *InteractionEventCreateData) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.createInteractionEventInDBTxWork")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	interactionEventId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelInteractionEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	err = s.services.Neo4jRepositories.InteractionEventWriteRepository.CreateInTx(ctx, tx, tenant, interactionEventId, *newInteractionEvent.InteractionEventEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if newInteractionEvent.SessionIdentifier != nil {
		sessionExists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsByIdInTx(ctx, tx, tenant, *newInteractionEvent.SessionIdentifier, commonModel.NodeLabelInteractionSession)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if !sessionExists {
			tracing.TraceErr(span, errors.New("session not found"))
			return nil, errors.New("session not found")
		}

		err = s.services.Neo4jRepositories.CommonWriteRepository.LinkEntityWithEntityInTx(ctx, tx, tenant, interactionEventId, commonModel.NodeLabelInteractionEvent, commonModel.PART_OF, nil, *newInteractionEvent.SessionIdentifier, commonModel.NodeLabelInteractionSession)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}
	if newInteractionEvent.ExternalSystem != nil && newInteractionEvent.ExternalSystem.ExternalSystemId != "" && newInteractionEvent.ExternalSystem.Relationship.ExternalId != "" {
		err := s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, tenant, interactionEventId, commonModel.NodeLabelInteractionEvent, neo4jmodel.ExternalSystem{
			ExternalId:       newInteractionEvent.ExternalSystem.Relationship.ExternalId,
			ExternalSystemId: newInteractionEvent.ExternalSystem.ExternalSystemId.String(),
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}
	if newInteractionEvent.MeetingIdentifier != nil {
		meetingExists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsByIdInTx(ctx, tx, tenant, *newInteractionEvent.MeetingIdentifier, commonModel.NodeLabelMeeting)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if !meetingExists {
			tracing.TraceErr(span, errors.New("meeting not found"))
			return nil, errors.New("meeting not found")
		}

		err = s.services.Neo4jRepositories.CommonWriteRepository.LinkEntityWithEntityInTx(ctx, tx, tenant, interactionEventId, commonModel.NodeLabelInteractionEvent, commonModel.PART_OF, nil, *newInteractionEvent.MeetingIdentifier, commonModel.NodeLabelMeeting)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}
	if newInteractionEvent.RepliesTo != nil {
		parentExists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsByIdInTx(ctx, tx, tenant, *newInteractionEvent.RepliesTo, commonModel.NodeLabelInteractionEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if !parentExists {
			tracing.TraceErr(span, errors.New("parent not found"))
			return nil, errors.New("parent not found")
		}

		err = s.services.Neo4jRepositories.CommonWriteRepository.LinkEntityWithEntityInTx(ctx, tx, tenant, interactionEventId, commonModel.NodeLabelInteractionEvent, commonModel.REPLIES_TO, nil, *newInteractionEvent.RepliesTo, commonModel.NodeLabelInteractionEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	for _, sentBy := range newInteractionEvent.SentBy {
		err := s.linkInteractionEventParticipantInTx(ctx, tx, tenant, interactionEventId, sentBy, commonModel.SENT_BY, nil)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	for _, sentTo := range newInteractionEvent.SentTo {
		relationshipType := "TO"
		err := s.linkInteractionEventParticipantInTx(ctx, tx, tenant, interactionEventId, sentTo, commonModel.SENT_TO, &relationshipType)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	for _, sentCc := range newInteractionEvent.SentCc {
		relationshipType := "CC"
		err := s.linkInteractionEventParticipantInTx(ctx, tx, tenant, interactionEventId, sentCc, commonModel.SENT_TO, &relationshipType)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	for _, sentBcc := range newInteractionEvent.SentBcc {
		relationshipType := "BCC"
		err := s.linkInteractionEventParticipantInTx(ctx, tx, tenant, interactionEventId, sentBcc, commonModel.SENT_TO, &relationshipType)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	return &interactionEventId, nil
}

func (s *interactionEventService) linkInteractionEventParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId string, linkWIthData InteractionEventParticipantData, relationship commonModel.EntityRelation, relationshipType *string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.linkInteractionEventParticipantInTx")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var linkWithId string
	var linkWithLabel commonModel.EntityType

	if linkWIthData.UserId != nil {
		linkWithId = *linkWIthData.UserId
		linkWithLabel = commonModel.USER
	} else if linkWIthData.ContactId != nil {
		linkWithId = *linkWIthData.ContactId
		linkWithLabel = commonModel.CONTACT
	} else if linkWIthData.Email != nil {
		linkWithLabel = commonModel.EMAIL

		emailId, err := s.services.Neo4jRepositories.EmailReadRepository.GetEmailIdIfExists(ctx, tenant, *linkWIthData.Email)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if emailId != "" {
			linkWithId = emailId
		} else {
			//TODO create and use inTx method
			createdEmailId, err := s.services.EmailService.Merge(ctx, neo4jentity.EmailEntity{
				Email: *linkWIthData.Email,
			}, nil)

			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}

			if createdEmailId == nil {
				tracing.TraceErr(span, errors.New("failed to create email"))
				return errors.New("failed to create email")
			}

			linkWithId = *createdEmailId
		}

	} else if linkWIthData.PhoneNumber != nil {
		linkWithLabel = commonModel.PHONE_NUMBER

		phoneNumberId, err := s.services.Neo4jRepositories.PhoneNumberReadRepository.GetPhoneNumberIdIfExists(ctx, tenant, *linkWIthData.Email)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if phoneNumberId != "" {
			linkWithId = phoneNumberId
		} else {
			//TODO create and use inTx method
		}
	} else {
		tracing.TraceErr(span, errors.New("no link with data provided"))
		return errors.New("no link with data provided")
	}

	var relationshipProperties *map[string]interface{}

	if relationshipType != nil {
		relationshipProperties = &map[string]interface{}{
			"type": *relationshipType,
		}
	}

	err := s.services.Neo4jRepositories.CommonWriteRepository.LinkEntityWithEntityInTx(ctx, tx, tenant, interactionEventId, commonModel.INTERACTION_EVENT, relationship, relationshipProperties, linkWithId, linkWithLabel)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *interactionEventService) convertDbNodesToInteractionEventParticipants(records []*utils.DbNodeWithRelationAndId) neo4jentity.InteractionEventParticipants {
	interactionEventParticipants := neo4jentity.InteractionEventParticipants{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, commonModel.NodeLabelEmail) {
			participant := neo4jmapper.MapDbNodeToEmailEntity(v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelPhoneNumber) {
			participant := neo4jmapper.MapDbNodeToPhoneNumberEntity(v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelUser) {
			participant := neo4jmapper.MapDbNodeToUserEntity(v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelContact) {
			participant := neo4jmapper.MapDbNodeToContactEntity(v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelOrganization) {
			participant := neo4jmapper.MapDbNodeToOrganizationEntity(v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, commonModel.NodeLabelJobRole) {
			participant := neo4jmapper.MapDbNodeToJobRoleEntity(v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		}
	}
	return interactionEventParticipants
}

func (s *interactionEventService) mapDbRelationshipToParticipantDetails(relationship dbtype.Relationship) neo4jentity.InteractionEventParticipantDetails {
	props := utils.GetPropsFromRelationship(relationship)
	details := neo4jentity.InteractionEventParticipantDetails{
		Type: utils.GetStringPropOrEmpty(props, "type"),
	}
	return details
}
