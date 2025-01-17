package interaction_event

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/exp/slices"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type interactionEventService struct {
	neo4j *neoRepo.Repositories
	email interfaces.EmailService
}

func NewInteractionEventService(neo4j *neoRepo.Repositories, email interfaces.EmailService) interfaces.InteractionEventService {
	return &interactionEventService{
		neo4j: neo4j,
		email: email,
	}
}

func (s *interactionEventService) SetEmailService(email interfaces.EmailService) {
	s.email = email
}

func (s *interactionEventService) IsInitialized() bool {
	if s.neo4j == nil || s.email == nil {
		return false
	}
	return true
}

func (s *interactionEventService) GetById(ctx context.Context, id string) (*neo4jentity.InteractionEventEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetById")
	defer span.Finish()

	byId, err := s.neo4j.CommonReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), id, commonModel.NodeLabelInteractionEvent)
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

	interactionEvents, err := s.neo4j.InteractionEventReadRepository.GetAllForInteractionSessions(ctx, common.GetTenantFromContext(ctx), ids, loadContent)
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

	interactionEvents, err := s.neo4j.InteractionEventReadRepository.GetAllForMeetings(ctx, common.GetTenantFromContext(ctx), ids, loadContent)
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

	interactionEvents, err := s.neo4j.InteractionEventReadRepository.GetAllForIssues(ctx, common.GetTenantFromContext(ctx), issueIds, loadContent)
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

	records, err := s.neo4j.InteractionEventReadRepository.GetSentByFor(ctx, common.GetTenantFromContext(ctx), ids)
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

	records, err := s.neo4j.InteractionEventReadRepository.GetSentToFor(ctx, common.GetTenantFromContext(ctx), ids)
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

	records, err := s.neo4j.InteractionEventReadRepository.GetReplyToFor(ctx, common.GetTenantFromContext(ctx), ids, loadContent)
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

func (s *interactionEventService) Create(ctx context.Context, data *interfaces.InteractionEventCreateData) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *s.neo4j.Neo4jDriver)
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
	//for _, v := range newInteractionEvent.SentTo {
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

	return interactionEventId.(*string), nil
}

func (s *interactionEventService) CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, newInteractionEvent *interfaces.InteractionEventCreateData) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.createInteractionEventInDBTxWork")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	interactionEventId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelInteractionEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	err = s.neo4j.InteractionEventWriteRepository.CreateInTx(ctx, tx, tenant, interactionEventId, *newInteractionEvent.InteractionEventEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if newInteractionEvent.SessionIdentifier != nil {
		sessionExists, err := s.neo4j.CommonReadRepository.ExistsByIdInTx(ctx, &tx, tenant, *newInteractionEvent.SessionIdentifier, commonModel.NodeLabelInteractionSession)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if !sessionExists {
			tracing.TraceErr(span, errors.New("session not found"))
			return nil, errors.New("session not found")
		}

		err = s.neo4j.CommonWriteRepository.Link(ctx, &tx, tenant, repository.LinkDetails{
			FromEntityId:           interactionEventId,
			FromEntityType:         commonModel.INTERACTION_EVENT,
			Relationship:           commonModel.PART_OF,
			RelationshipProperties: nil,
			ToEntityId:             *newInteractionEvent.SessionIdentifier,
			ToEntityType:           commonModel.INTERACTION_SESSION,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}
	if newInteractionEvent.ExternalSystem != nil && newInteractionEvent.ExternalSystem.ExternalSystemId != "" && newInteractionEvent.ExternalSystem.Relationship.ExternalId != "" {
		err := s.neo4j.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, &tx, tenant, interactionEventId, commonModel.NodeLabelInteractionEvent, neo4jmodel.ExternalSystem{
			ExternalId:       newInteractionEvent.ExternalSystem.Relationship.ExternalId,
			ExternalSystemId: newInteractionEvent.ExternalSystem.ExternalSystemId.String(),
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}
	if newInteractionEvent.MeetingIdentifier != nil {
		meetingExists, err := s.neo4j.CommonReadRepository.ExistsByIdInTx(ctx, &tx, tenant, *newInteractionEvent.MeetingIdentifier, commonModel.NodeLabelMeeting)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if !meetingExists {
			tracing.TraceErr(span, errors.New("meeting not found"))
			return nil, errors.New("meeting not found")
		}

		err = s.neo4j.CommonWriteRepository.Link(ctx, &tx, tenant, repository.LinkDetails{
			FromEntityId:           interactionEventId,
			FromEntityType:         commonModel.INTERACTION_EVENT,
			Relationship:           commonModel.PART_OF,
			RelationshipProperties: nil,
			ToEntityId:             *newInteractionEvent.MeetingIdentifier,
			ToEntityType:           commonModel.MEETING,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}
	if newInteractionEvent.RepliesTo != nil {
		parentExists, err := s.neo4j.CommonReadRepository.ExistsByIdInTx(ctx, &tx, tenant, *newInteractionEvent.RepliesTo, commonModel.NodeLabelInteractionEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if !parentExists {
			tracing.TraceErr(span, errors.New("parent not found"))
			return nil, errors.New("parent not found")
		}

		err = s.neo4j.CommonWriteRepository.Link(ctx, &tx, tenant, repository.LinkDetails{
			FromEntityId:   interactionEventId,
			FromEntityType: commonModel.INTERACTION_EVENT,
			Relationship:   commonModel.REPLIES_TO,
			ToEntityId:     *newInteractionEvent.RepliesTo,
			ToEntityType:   commonModel.INTERACTION_EVENT,
		})
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

func (s *interactionEventService) linkInteractionEventParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId string, linkWIthData interfaces.InteractionEventParticipantData, relationship commonModel.EntityRelation, relationshipType *string) error {
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

		emailId, err := s.neo4j.EmailReadRepository.GetEmailIdIfExists(ctx, &tx, tenant, *linkWIthData.Email)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if emailId != "" {
			linkWithId = emailId
		} else {
			// TODO create and use inTx method
			createdEmailId, err := s.email.Merge(ctx, nil, "",
				interfaces.EmailFields{
					Email:     *linkWIthData.Email,
					AppSource: "",
				},
				nil)
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

		phoneNumberId, err := s.neo4j.PhoneNumberReadRepository.GetPhoneNumberIdIfExists(ctx, tenant, *linkWIthData.Email)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if phoneNumberId != "" {
			linkWithId = phoneNumberId
		} else {
			// TODO create and use inTx method
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

	err := s.neo4j.CommonWriteRepository.Link(ctx, &tx, tenant, repository.LinkDetails{
		FromEntityId:           interactionEventId,
		FromEntityType:         commonModel.INTERACTION_EVENT,
		Relationship:           relationship,
		RelationshipProperties: relationshipProperties,
		ToEntityId:             linkWithId,
		ToEntityType:           linkWithLabel,
	})
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
