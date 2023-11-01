package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/exp/slices"
)

type InteractionEventService interface {
	InteractionEventLinkAttachment(ctx context.Context, noteID string, attachmentID string) (*entity.InteractionEventEntity, error)
	GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error)
	GetInteractionEventsForMeetings(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error)
	GetInteractionEventsForIssues(ctx context.Context, issueIds []string) (*entity.InteractionEventEntities, error)
	GetSentByParticipantsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventParticipants, error)
	GetSentToParticipantsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventParticipants, error)
	GetInteractionEventById(ctx context.Context, id string) (*entity.InteractionEventEntity, error)
	GetInteractionEventByEventIdentifier(ctx context.Context, eventIdentifier string) (*entity.InteractionEventEntity, error)
	GetReplyToInteractionsEventForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error)

	Create(ctx context.Context, newInteractionEvent *InteractionEventCreateData) (*entity.InteractionEventEntity, error)

	mapDbNodeToInteractionEventEntity(node dbtype.Node) *entity.InteractionEventEntity
}

type ParticipantAddressData struct {
	Email       *string
	PhoneNumber *string
	ContactId   *string
	UserId      *string
	Type        *string
}

type InteractionEventCreateData struct {
	InteractionEventEntity *entity.InteractionEventEntity
	SessionIdentifier      *string
	MeetingIdentifier      *string
	RepliesTo              *string
	SentBy                 []ParticipantAddressData
	SentTo                 []ParticipantAddressData
	Source                 entity.DataSource
	SourceOfTruth          entity.DataSource
}

type interactionEventService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewInteractionEventService(log logger.Logger, repositories *repository.Repositories, services *Services) InteractionEventService {
	return &interactionEventService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *interactionEventService) InteractionEventLinkAttachment(ctx context.Context, noteID string, attachmentID string) (*entity.InteractionEventEntity, error) {
	node, err := s.services.AttachmentService.LinkNodeWithAttachment(ctx, repository.LINKED_WITH_INTERACTION_EVENT, nil, attachmentID, noteID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToInteractionEventEntity(*node), nil
}

func (s *interactionEventService) Create(ctx context.Context, newInteractionEvent *InteractionEventCreateData) (*entity.InteractionEventEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	interactionEventDbNode, err := session.ExecuteWrite(ctx, s.createInteractionEventInDBTxWork(ctx, newInteractionEvent))
	if err != nil {
		return nil, err
	}

	for _, v := range newInteractionEvent.SentBy {
		if v.ContactId != nil {
			s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, *v.ContactId)
		}
		if v.Email != nil {
			s.services.OrganizationService.UpdateLastTouchpointSyncByEmail(ctx, *v.Email)
		}
		if v.PhoneNumber != nil {
			s.services.OrganizationService.UpdateLastTouchpointSyncByPhoneNumber(ctx, *v.PhoneNumber)
		}
	}
	for _, v := range newInteractionEvent.SentTo {
		if v.ContactId != nil {
			s.services.OrganizationService.UpdateLastTouchpointSyncByContactId(ctx, *v.ContactId)
		}
		if v.Email != nil {
			s.services.OrganizationService.UpdateLastTouchpointSyncByEmail(ctx, *v.Email)
		}
		if v.PhoneNumber != nil {
			s.services.OrganizationService.UpdateLastTouchpointSyncByPhoneNumber(ctx, *v.PhoneNumber)
		}
	}

	return s.mapDbNodeToInteractionEventEntity(*interactionEventDbNode.(*dbtype.Node)), nil
}

func (s *interactionEventService) createInteractionEventInDBTxWork(ctx context.Context, newInteractionEvent *InteractionEventCreateData) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		interactionEventDbNode, err := s.repositories.InteractionEventRepository.Create(ctx, tx, tenant, *newInteractionEvent.InteractionEventEntity, newInteractionEvent.Source, newInteractionEvent.SourceOfTruth)
		if err != nil {
			return nil, err
		}
		var interactionEventId = utils.GetPropsFromNode(*interactionEventDbNode)["id"].(string)

		if newInteractionEvent.SessionIdentifier != nil {
			err := s.repositories.InteractionEventRepository.LinkWithPartOfXXInTx(ctx, tx, tenant, interactionEventId, *newInteractionEvent.SessionIdentifier, repository.PART_OF_INTERACTION_SESSION)
			if err != nil {
				return nil, err
			}
		}
		if newInteractionEvent.InteractionEventEntity.ExternalId != nil && newInteractionEvent.InteractionEventEntity.ExternalSystemId != nil {
			err := s.repositories.InteractionEventRepository.LinkWithExternalSystemInTx(ctx, tx, tenant, interactionEventId, *newInteractionEvent.InteractionEventEntity.ExternalId, *newInteractionEvent.InteractionEventEntity.ExternalSystemId)
			if err != nil {
				return nil, err
			}
		}
		if newInteractionEvent.MeetingIdentifier != nil {
			err := s.repositories.InteractionEventRepository.LinkWithPartOfXXInTx(ctx, tx, tenant, interactionEventId, *newInteractionEvent.MeetingIdentifier, repository.PART_OF_MEETING)
			if err != nil {
				return nil, err
			}
		}
		if newInteractionEvent.RepliesTo != nil {
			err := s.repositories.InteractionEventRepository.LinkWithRepliesToInTx(ctx, tx, tenant, interactionEventId, *newInteractionEvent.RepliesTo)
			if err != nil {
				return nil, err
			}
		}

		for _, sentTo := range newInteractionEvent.SentTo {
			if sentTo.ContactId != nil {
				err := s.repositories.InteractionEventRepository.LinkWithSentXXParticipantInTx(ctx, tx, tenant, entity.CONTACT, interactionEventId, *sentTo.ContactId, sentTo.Type, repository.SENT_TO)
				if err != nil {
					return nil, err
				}
			} else if sentTo.UserId != nil {
				err := s.repositories.InteractionEventRepository.LinkWithSentXXParticipantInTx(ctx, tx, tenant, entity.USER, interactionEventId, *sentTo.UserId, sentTo.Type, repository.SENT_TO)
				if err != nil {
					return nil, err
				}
			} else if sentTo.Email != nil {
				exists, err := s.repositories.EmailRepository.Exists(ctx, tenant, *sentTo.Email)
				if err != nil {
					return nil, err
				}

				curTime := utils.Now()
				if !exists {
					_, err = s.services.ContactService.Create(ctx, &ContactCreateData{
						ContactEntity: &entity.ContactEntity{CreatedAt: &curTime, FirstName: "", LastName: ""},
						EmailEntity:   mapper.MapEmailInputToEntity(&model.EmailInput{Email: *sentTo.Email}),
						Source:        entity.DataSourceOpenline,
					})
				}
				err = s.repositories.InteractionEventRepository.LinkWithSentXXEmailInTx(ctx, tx, tenant, interactionEventId, *sentTo.Email, sentTo.Type, repository.SENT_TO)
				if err != nil {
					return nil, err
				}

			} else if sentTo.PhoneNumber != nil {
				exists, err := s.repositories.PhoneNumberRepository.Exists(ctx, tenant, *sentTo.PhoneNumber)
				if err != nil {
					return nil, err
				}

				curTime := utils.Now()
				if !exists {
					_, err = s.services.ContactService.Create(ctx, &ContactCreateData{
						ContactEntity:     &entity.ContactEntity{CreatedAt: &curTime, FirstName: "", LastName: ""},
						PhoneNumberEntity: mapper.MapPhoneNumberInputToEntity(&model.PhoneNumberInput{PhoneNumber: *sentTo.PhoneNumber}),
						Source:            entity.DataSourceOpenline,
					})
				}
				err = s.repositories.InteractionEventRepository.LinkWithSentXXPhoneNumberInTx(ctx, tx, tenant, interactionEventId, *sentTo.PhoneNumber, sentTo.Type, repository.SENT_TO)
				if err != nil {
					return nil, err
				}

			}

		}

		for _, sentBy := range newInteractionEvent.SentBy {
			if sentBy.ContactId != nil {
				err := s.repositories.InteractionEventRepository.LinkWithSentXXParticipantInTx(ctx, tx, tenant, entity.CONTACT, interactionEventId, *sentBy.ContactId, sentBy.Type, repository.SENT_BY)
				if err != nil {
					return nil, err
				}
			} else if sentBy.UserId != nil {
				err := s.repositories.InteractionEventRepository.LinkWithSentXXParticipantInTx(ctx, tx, tenant, entity.USER, interactionEventId, *sentBy.UserId, sentBy.Type, repository.SENT_BY)
				if err != nil {
					return nil, err
				}
			} else if sentBy.Email != nil {
				exists, err := s.repositories.EmailRepository.Exists(ctx, tenant, *sentBy.Email)
				if err != nil {
					return nil, err
				}

				curTime := utils.Now()
				if !exists {
					_, err = s.services.ContactService.Create(ctx, &ContactCreateData{
						ContactEntity: &entity.ContactEntity{CreatedAt: &curTime, FirstName: "", LastName: ""},
						EmailEntity:   mapper.MapEmailInputToEntity(&model.EmailInput{Email: *sentBy.Email}),
						Source:        entity.DataSourceOpenline,
					})
				}
				err = s.repositories.InteractionEventRepository.LinkWithSentXXEmailInTx(ctx, tx, tenant, interactionEventId, *sentBy.Email, sentBy.Type, repository.SENT_BY)
				if err != nil {
					return nil, err
				}

			} else if sentBy.PhoneNumber != nil {
				exists, err := s.repositories.PhoneNumberRepository.Exists(ctx, tenant, *sentBy.PhoneNumber)
				if err != nil {
					return nil, err
				}

				curTime := utils.Now()
				if !exists {
					_, err = s.services.ContactService.Create(ctx, &ContactCreateData{
						ContactEntity:     &entity.ContactEntity{CreatedAt: &curTime, FirstName: "", LastName: ""},
						PhoneNumberEntity: mapper.MapPhoneNumberInputToEntity(&model.PhoneNumberInput{PhoneNumber: *sentBy.PhoneNumber}),
						Source:            entity.DataSourceOpenline,
					})
				}
				err = s.repositories.InteractionEventRepository.LinkWithSentXXPhoneNumberInTx(ctx, tx, tenant, interactionEventId, *sentBy.PhoneNumber, sentBy.Type, repository.SENT_BY)
				if err != nil {
					return nil, err
				}

			}

		}

		return interactionEventDbNode, nil
	}
}

func (s *interactionEventService) GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error) {
	interactionEvents, err := s.repositories.InteractionEventRepository.GetAllForInteractionSessions(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	interactionEventEntities := entity.InteractionEventEntities{}
	for _, v := range interactionEvents {
		interactionEventEntity := s.mapDbNodeToInteractionEventEntity(*v.Node)
		interactionEventEntity.DataloaderKey = v.LinkedNodeId
		interactionEventEntities = append(interactionEventEntities, *interactionEventEntity)
	}
	return &interactionEventEntities, nil
}

func (s *interactionEventService) GetInteractionEventsForMeetings(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error) {
	interactionEvents, err := s.repositories.InteractionEventRepository.GetAllForMeetings(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	interactionEventEntities := entity.InteractionEventEntities{}
	for _, v := range interactionEvents {
		interactionEventEntity := s.mapDbNodeToInteractionEventEntity(*v.Node)
		interactionEventEntity.DataloaderKey = v.LinkedNodeId
		interactionEventEntities = append(interactionEventEntities, *interactionEventEntity)
	}
	return &interactionEventEntities, nil
}

func (s *interactionEventService) GetInteractionEventsForIssues(ctx context.Context, issueIds []string) (*entity.InteractionEventEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetInteractionEventsForIssues")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("issueIds", issueIds))

	interactionEvents, err := s.repositories.InteractionEventRepository.GetAllForIssues(ctx, common.GetTenantFromContext(ctx), issueIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	interactionEventEntities := entity.InteractionEventEntities{}
	for _, v := range interactionEvents {
		interactionEventEntity := s.mapDbNodeToInteractionEventEntity(*v.Node)
		interactionEventEntity.DataloaderKey = v.LinkedNodeId
		interactionEventEntities = append(interactionEventEntities, *interactionEventEntity)
	}
	span.LogFields(log.Int("result count", len(interactionEventEntities)))
	return &interactionEventEntities, nil
}

func (s *interactionEventService) GetSentByParticipantsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventParticipants, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetSentByParticipantsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	records, err := s.repositories.InteractionEventRepository.GetSentByParticipantsForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToInteractionEventParticipants(records)

	span.LogFields(log.Int("result count", len(interactionEventParticipants)))

	return &interactionEventParticipants, nil
}

func (s *interactionEventService) GetSentToParticipantsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventParticipants, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventService.GetSentToParticipantsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	records, err := s.repositories.InteractionEventRepository.GetSentToParticipantsForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToInteractionEventParticipants(records)

	span.LogFields(log.Int("result count", len(interactionEventParticipants)))

	return &interactionEventParticipants, nil
}

func (s *interactionEventService) GetReplyToInteractionsEventForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error) {
	records, err := s.repositories.InteractionEventRepository.GetReplyToInteractionEventsForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEvents := s.convertDbNodesToInteractionEvent(records)

	return &interactionEvents, nil
}

func (s *interactionEventService) GetInteractionEventById(ctx context.Context, id string) (*entity.InteractionEventEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (e:InteractionEvent_%s {id:$id}) RETURN e`,
			common.GetTenantFromContext(ctx)),
			map[string]interface{}{
				"id": id,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToInteractionEventEntity(queryResult.(dbtype.Node)), nil
}

func (s *interactionEventService) GetInteractionEventByEventIdentifier(ctx context.Context, eventIdentifier string) (*entity.InteractionEventEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (e:InteractionEvent_%s {identifier:$identifier}) RETURN e`,
			common.GetTenantFromContext(ctx)),
			map[string]interface{}{
				"identifier": eventIdentifier,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToInteractionEventEntity(queryResult.(dbtype.Node)), nil
}

func (s *interactionEventService) mapDbNodeToInteractionEventEntity(node dbtype.Node) *entity.InteractionEventEntity {
	props := utils.GetPropsFromNode(node)
	createdAt := utils.GetTimePropOrEpochStart(props, "createdAt")
	interactionEventEntity := entity.InteractionEventEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:       &createdAt,
		EventIdentifier: utils.GetStringPropOrEmpty(props, "identifier"),
		Channel:         utils.GetStringPropOrNil(props, "channel"),
		ChannelData:     utils.GetStringPropOrNil(props, "channelData"),
		EventType:       utils.GetStringPropOrNil(props, "eventType"),
		Content:         utils.GetStringPropOrEmpty(props, "content"),
		ContentType:     utils.GetStringPropOrEmpty(props, "contentType"),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		Source:          entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionEventEntity
}

func (s *interactionEventService) convertDbNodesToInteractionEvent(records []*utils.DbNodeAndId) entity.InteractionEventEntities {
	interactionEvents := entity.InteractionEventEntities{}
	for _, v := range records {
		event := s.mapDbNodeToInteractionEventEntity(*v.Node)
		event.DataloaderKey = v.LinkedNodeId
		interactionEvents = append(interactionEvents, *event)

	}
	return interactionEvents
}

func (s *interactionEventService) convertDbNodesToInteractionEventParticipants(records []*utils.DbNodeWithRelationAndId) entity.InteractionEventParticipants {
	interactionEventParticipants := entity.InteractionEventParticipants{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, entity.NodeLabel_Email) {
			participant := s.services.EmailService.mapDbNodeToEmailEntity(*v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_PhoneNumber) {
			participant := s.services.PhoneNumberService.mapDbNodeToPhoneNumberEntity(*v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_User) {
			participant := s.services.UserService.mapDbNodeToUserEntity(*v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Contact) {
			participant := s.services.ContactService.mapDbNodeToContactEntity(*v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Organization) {
			participant := s.services.OrganizationService.mapDbNodeToOrganizationEntity(*v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_JobRole) {
			participant := s.services.JobRoleService.mapDbNodeToJobRoleEntity(*v.Node)
			participant.InteractionEventParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		}
	}
	return interactionEventParticipants
}

func (s *interactionEventService) mapDbRelationshipToParticipantDetails(relationship dbtype.Relationship) entity.InteractionEventParticipantDetails {
	props := utils.GetPropsFromRelationship(relationship)
	details := entity.InteractionEventParticipantDetails{
		Type: utils.GetStringPropOrEmpty(props, "type"),
	}
	return details
}

func MapInteractionEventParticipantInputToAddressData(input []*model.InteractionEventParticipantInput) []ParticipantAddressData {
	var inputData []ParticipantAddressData
	for _, participant := range input {
		inputData = append(inputData, ParticipantAddressData{
			Email:       participant.Email,
			PhoneNumber: participant.PhoneNumber,
			UserId:      participant.UserID,
			ContactId:   participant.ContactID,
			Type:        participant.Type,
		})
	}
	return inputData
}

func (s *interactionEventService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}
