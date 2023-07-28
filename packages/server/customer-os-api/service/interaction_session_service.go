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
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/exp/slices"
	"time"
)

type InteractionSessionService interface {
	GetInteractionSessionsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionSessionEntities, error)
	InteractionSessionLinkAttachment(ctx context.Context, noteID string, attachmentID string) (*entity.InteractionSessionEntity, error)
	GetInteractionSessionById(ctx context.Context, id string) (*entity.InteractionSessionEntity, error)
	Create(ctx context.Context, newInteractionSession *InteractionSessionCreateData) (*entity.InteractionSessionEntity, error)
	GetInteractionSessionBySessionIdentifier(ctx context.Context, sessionIdentifier string) (*entity.InteractionSessionEntity, error)
	GetAttendedByParticipantsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionSessionParticipants, error)

	mapDbNodeToInteractionSessionEntity(node dbtype.Node) *entity.InteractionSessionEntity
}

type InteractionSessionCreateData struct {
	InteractionSessionEntity *entity.InteractionSessionEntity
	AttendedBy               []ParticipantAddressData
}

type interactionSessionService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewInteractionSessionService(log logger.Logger, repositories *repository.Repositories, services *Services) InteractionSessionService {
	return &interactionSessionService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *interactionSessionService) InteractionSessionLinkAttachment(ctx context.Context, noteID string, attachmentID string) (*entity.InteractionSessionEntity, error) {
	node, err := s.services.AttachmentService.LinkNodeWithAttachment(ctx, repository.LINKED_WITH_INTERACTION_SESSION, nil, attachmentID, noteID)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToInteractionSessionEntity(*node), nil
}

func (s *interactionSessionService) GetAttendedByParticipantsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionSessionParticipants, error) {
	records, err := s.repositories.InteractionSessionRepository.GetAttendedByParticipantsForInteractionSessions(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToInteractionSessionParticipants(records)

	return &interactionEventParticipants, nil
}

func (s *interactionSessionService) Create(ctx context.Context, newInteractionSession *InteractionSessionCreateData) (*entity.InteractionSessionEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, s.createInteractionSessionInDBTxWork(ctx, newInteractionSession))
	if err != nil {
		return nil, err
	}

	for _, v := range newInteractionSession.AttendedBy {
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

	return s.mapDbNodeToInteractionSessionEntity(*queryResult.(*dbtype.Node)), nil
}

func (s *interactionSessionService) createInteractionSessionInDBTxWork(ctx context.Context, newInteractionSession *InteractionSessionCreateData) func(tx neo4j.ManagedTransaction) (any, error) {
	return func(tx neo4j.ManagedTransaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		interactionEventDbNode, err := s.repositories.InteractionSessionRepository.Create(ctx, tx, common.GetTenantFromContext(ctx), newInteractionSession.InteractionSessionEntity)
		if err != nil {
			return nil, err
		}
		var interactionSessionId = utils.GetPropsFromNode(*interactionEventDbNode)["id"].(string)

		for _, attendedBy := range newInteractionSession.AttendedBy {
			if attendedBy.ContactId != nil {
				err := s.repositories.InteractionSessionRepository.LinkWithAttendedByParticipantInTx(ctx, tx, tenant, entity.CONTACT, interactionSessionId, *attendedBy.ContactId, attendedBy.Type)
				if err != nil {
					return nil, err
				}
			} else if attendedBy.UserId != nil {
				err := s.repositories.InteractionSessionRepository.LinkWithAttendedByParticipantInTx(ctx, tx, tenant, entity.USER, interactionSessionId, *attendedBy.UserId, attendedBy.Type)
				if err != nil {
					return nil, err
				}
			} else if attendedBy.Email != nil {
				exists, err := s.repositories.EmailRepository.Exists(ctx, tenant, *attendedBy.Email)
				if err != nil {
					return nil, err
				}

				curTime := utils.Now()
				if !exists {
					_, err = s.services.ContactService.Create(ctx, &ContactCreateData{
						ContactEntity: &entity.ContactEntity{CreatedAt: &curTime, FirstName: "", LastName: ""},
						EmailEntity:   mapper.MapEmailInputToEntity(&model.EmailInput{Email: *attendedBy.Email}),
						Source:        entity.DataSourceOpenline,
						SourceOfTruth: entity.DataSourceOpenline,
					})
				}
				err = s.repositories.InteractionSessionRepository.LinkWithAttendedByEmailInTx(ctx, tx, tenant, interactionSessionId, *attendedBy.Email, attendedBy.Type)
				if err != nil {
					return nil, err
				}

			} else if attendedBy.PhoneNumber != nil {
				exists, err := s.repositories.PhoneNumberRepository.Exists(ctx, tenant, *attendedBy.PhoneNumber)
				if err != nil {
					return nil, err
				}

				curTime := utils.Now()
				if !exists {
					_, err = s.services.ContactService.Create(ctx, &ContactCreateData{
						ContactEntity:     &entity.ContactEntity{CreatedAt: &curTime, FirstName: "", LastName: ""},
						PhoneNumberEntity: mapper.MapPhoneNumberInputToEntity(&model.PhoneNumberInput{PhoneNumber: *attendedBy.PhoneNumber}),
						Source:            entity.DataSourceOpenline,
						SourceOfTruth:     entity.DataSourceOpenline,
					})
				}
				err = s.repositories.InteractionSessionRepository.LinkWithAttendedByPhoneNumberInTx(ctx, tx, tenant, interactionSessionId, *attendedBy.PhoneNumber, attendedBy.Type)
				if err != nil {
					return nil, err
				}

			}

		}
		return interactionEventDbNode, nil
	}
}

func (s *interactionSessionService) GetInteractionSessionById(ctx context.Context, id string) (*entity.InteractionSessionEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (e:InteractionSession_%s {id:$id}) RETURN e`,
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

	return s.mapDbNodeToInteractionSessionEntity(queryResult.(dbtype.Node)), nil
}

func (s *interactionSessionService) GetInteractionSessionBySessionIdentifier(ctx context.Context, sessionIdentifier string) (*entity.InteractionSessionEntity, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (e:InteractionSession_%s {identifier:$identifier}) RETURN e`,
			common.GetTenantFromContext(ctx)),
			map[string]interface{}{
				"identifier": sessionIdentifier,
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

	return s.mapDbNodeToInteractionSessionEntity(queryResult.(dbtype.Node)), nil
}

func (s *interactionSessionService) GetInteractionSessionsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionSessionEntities, error) {
	interactionSessions, err := s.repositories.InteractionSessionRepository.GetAllForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}
	interactionSessionEntities := entity.InteractionSessionEntities{}
	for _, v := range interactionSessions {
		interactionSessionEntity := s.mapDbNodeToInteractionSessionEntity(*v.Node)
		interactionSessionEntity.DataloaderKey = v.LinkedNodeId
		interactionSessionEntities = append(interactionSessionEntities, *interactionSessionEntity)
	}
	return &interactionSessionEntities, nil
}

// createdAt takes priority over startedAt
func (s *interactionSessionService) migrateStartedAt(props map[string]any) time.Time {
	if props["createdAt"] != nil {
		return utils.GetTimePropOrNow(props, "createdAt")
	}
	if props["startedAt"] != nil {
		return utils.GetTimePropOrNow(props, "startedAt")
	}
	return time.Now()
}

func (s *interactionSessionService) mapDbNodeToInteractionSessionEntity(node dbtype.Node) *entity.InteractionSessionEntity {
	props := utils.GetPropsFromNode(node)
	interactionSessionEntity := entity.InteractionSessionEntity{
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:         s.migrateStartedAt(props),
		UpdatedAt:         utils.GetTimePropOrNow(props, "updatedAt"),
		SessionIdentifier: utils.GetStringPropOrNil(props, "identifier"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		Status:            utils.GetStringPropOrEmpty(props, "status"),
		Type:              utils.GetStringPropOrNil(props, "type"),
		Channel:           utils.GetStringPropOrNil(props, "channel"),
		ChannelData:       utils.GetStringPropOrNil(props, "channelData"),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		Source:            entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionSessionEntity
}

func (s *interactionSessionService) mapDbRelationshipToParticipantDetails(relationship dbtype.Relationship) entity.InteractionSessionParticipantDetails {
	props := utils.GetPropsFromRelationship(relationship)
	details := entity.InteractionSessionParticipantDetails{
		Type: utils.GetStringPropOrEmpty(props, "type"),
	}
	return details
}

func (s *interactionSessionService) convertDbNodesToInteractionSessionParticipants(records []*utils.DbNodeWithRelationAndId) entity.InteractionSessionParticipants {
	interactionEventParticipants := entity.InteractionSessionParticipants{}
	for _, v := range records {
		if slices.Contains(v.Node.Labels, entity.NodeLabel_Email) {
			participant := s.services.EmailService.mapDbNodeToEmailEntity(*v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_PhoneNumber) {
			participant := s.services.PhoneNumberService.mapDbNodeToPhoneNumberEntity(*v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_User) {
			participant := s.services.UserService.mapDbNodeToUserEntity(*v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		} else if slices.Contains(v.Node.Labels, entity.NodeLabel_Contact) {
			participant := s.services.ContactService.mapDbNodeToContactEntity(*v.Node)
			participant.InteractionSessionParticipantDetails = s.mapDbRelationshipToParticipantDetails(*v.Relationship)
			participant.DataloaderKey = v.LinkedNodeId
			interactionEventParticipants = append(interactionEventParticipants, participant)
		}
	}
	return interactionEventParticipants
}

func MapInteractionSessionParticipantInputToAddressData(input []*model.InteractionSessionParticipantInput) []ParticipantAddressData {
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

func (s *interactionSessionService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}
