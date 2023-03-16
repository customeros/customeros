package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/exp/slices"
	"golang.org/x/net/context"
)

type InteractionEventService interface {
	GetInteractionEventsForInteractionSessions(ctx context.Context, ids []string) (*entity.InteractionEventEntities, error)
	GetSentByParticipantsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventParticipants, error)
	GetSentToParticipantsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventParticipants, error)

	mapDbNodeToInteractionEventEntity(node dbtype.Node) *entity.InteractionEventEntity
}

type interactionEventService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewInteractionEventService(repositories *repository.Repositories, services *Services) InteractionEventService {
	return &interactionEventService{
		repositories: repositories,
		services:     services,
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

func (s *interactionEventService) GetSentByParticipantsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventParticipants, error) {
	records, err := s.repositories.InteractionEventRepository.GetSentByParticipantsForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToInteractionEventParticipants(records)

	return &interactionEventParticipants, nil
}

func (s *interactionEventService) GetSentToParticipantsForInteractionEvents(ctx context.Context, ids []string) (*entity.InteractionEventParticipants, error) {
	records, err := s.repositories.InteractionEventRepository.GetSentToParticipantsForInteractionEvents(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	interactionEventParticipants := s.convertDbNodesToInteractionEventParticipants(records)

	return &interactionEventParticipants, nil
}

func (s *interactionEventService) mapDbNodeToInteractionEventEntity(node dbtype.Node) *entity.InteractionEventEntity {
	props := utils.GetPropsFromNode(node)
	interactionEventEntity := entity.InteractionEventEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		EventIdentifier: utils.GetStringPropOrEmpty(props, "identifier"),
		Channel:         utils.GetStringPropOrEmpty(props, "channel"),
		Content:         utils.GetStringPropOrEmpty(props, "content"),
		ContentType:     utils.GetStringPropOrEmpty(props, "contentType"),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		Source:          entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionEventEntity
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
