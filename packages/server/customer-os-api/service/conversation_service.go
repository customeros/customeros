package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"github.com/sirupsen/logrus"
	"reflect"
)

type ConversationService interface {
	CreateNewConversation(ctx context.Context, userIds, contactIds []string, input *entity.ConversationEntity) (*entity.ConversationEntity, error)
	UpdateConversation(ctx context.Context, userIds, contactIds []string, input *entity.ConversationEntity, skipMessageCountIncrement bool) (*entity.ConversationEntity, error)
	CloseConversation(ctx context.Context, conversationId string, sourceOfTruth entity.DataSource) (*entity.ConversationEntity, error)
	GetConversationsForUser(ctx context.Context, userId string, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error)
	GetConversationsForContact(ctx context.Context, contactId string, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error)
}

type conversationService struct {
	repository *repository.Repositories
}

func NewConversationService(repository *repository.Repositories) ConversationService {
	return &conversationService{
		repository: repository,
	}
}

func (s *conversationService) getNeo4jDriver() *neo4j.Driver {
	return s.repository.Drivers.Neo4jDriver
}

func (s *conversationService) CreateNewConversation(ctx context.Context, userIds, contactIds []string, input *entity.ConversationEntity) (*entity.ConversationEntity, error) {
	if len(userIds) == 0 && len(contactIds) == 0 {
		msg := "Missing participants for new conversation"
		logrus.Error(msg)
		return nil, errors.New(msg)
	}
	if len(input.Id) == 0 {
		newUuid, _ := uuid.NewRandom()
		input.Id = newUuid.String()
	}

	session := utils.NewNeo4jReadSession(*s.getNeo4jDriver())
	defer session.Close()

	dbNodePtr, err := s.repository.ConversationRepository.Create(session, common.GetContext(ctx).Tenant, userIds, contactIds, *input)
	if err != nil {
		return nil, err
	}
	conversationEntity := s.mapDbNodeToConversationEntity(*dbNodePtr)
	return conversationEntity, nil
}

func (s *conversationService) UpdateConversation(ctx context.Context, userIds, contactIds []string, input *entity.ConversationEntity, skipMessageCountIncrement bool) (*entity.ConversationEntity, error) {
	session := utils.NewNeo4jReadSession(*s.getNeo4jDriver())
	defer session.Close()

	dbNodePtr, err := s.repository.ConversationRepository.Update(session, common.GetContext(ctx).Tenant, userIds, contactIds, skipMessageCountIncrement, *input)
	if err != nil {
		return nil, err
	}
	conversationEntity := s.mapDbNodeToConversationEntity(*dbNodePtr)
	return conversationEntity, nil
}

func (s *conversationService) CloseConversation(ctx context.Context, conversationId string, sourceOfTruth entity.DataSource) (*entity.ConversationEntity, error) {
	session := utils.NewNeo4jReadSession(*s.getNeo4jDriver())
	defer session.Close()

	dbNodePtr, err := s.repository.ConversationRepository.Close(session, common.GetContext(ctx).Tenant, conversationId, mapper.MapConversationStatusFromModel(model.ConversationStatusClosed), sourceOfTruth)
	if err != nil {
		return nil, err
	}
	conversationEntity := s.mapDbNodeToConversationEntity(*dbNodePtr)
	return conversationEntity, nil
}

func (s *conversationService) GetConversationsForUser(ctx context.Context, userId string, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(*s.getNeo4jDriver())
	defer session.Close()

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.ConversationEntity{}))
	if err != nil {
		return nil, err
	}

	conversationDbNodesWithTotalCount, err := s.repository.ConversationRepository.GetPaginatedConversationsForUser(
		session,
		common.GetContext(ctx).Tenant,
		userId,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(conversationDbNodesWithTotalCount.Count)

	conversationEntities := entity.ConversationEntities{}

	for _, v := range conversationDbNodesWithTotalCount.Nodes {
		conversationEntity := *s.mapDbNodeToConversationEntity(*v.Node)
		conversationEntities = append(conversationEntities, conversationEntity)
	}
	paginatedResult.SetRows(&conversationEntities)
	return &paginatedResult, nil
}

func (s *conversationService) GetConversationsForContact(ctx context.Context, contactId string, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(*s.getNeo4jDriver())
	defer session.Close()

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}
	cypherSort, err := buildSort(sortBy, reflect.TypeOf(entity.ConversationEntity{}))
	if err != nil {
		return nil, err
	}

	conversationDbNodesWithTotalCount, err := s.repository.ConversationRepository.GetPaginatedConversationsForContact(
		session,
		common.GetContext(ctx).Tenant,
		contactId,
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		cypherSort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(conversationDbNodesWithTotalCount.Count)

	conversationEntities := entity.ConversationEntities{}

	for _, v := range conversationDbNodesWithTotalCount.Nodes {
		conversationEntity := *s.mapDbNodeToConversationEntity(*v.Node)
		conversationEntities = append(conversationEntities, conversationEntity)
	}
	paginatedResult.SetRows(&conversationEntities)
	return &paginatedResult, nil
}

func (s *conversationService) mapDbNodeToConversationEntity(dbNode dbtype.Node) *entity.ConversationEntity {
	props := utils.GetPropsFromNode(dbNode)
	conversationEntity := entity.ConversationEntity{
		Id:                 utils.GetStringPropOrEmpty(props, "id"),
		StartedAt:          utils.GetTimePropOrEpochStart(props, "startedAt"),
		UpdatedAt:          utils.GetTimePropOrEpochStart(props, "updatedAt"),
		EndedAt:            utils.GetTimePropOrNil(props, "endedAt"),
		Channel:            utils.GetStringPropOrEmpty(props, "channel"),
		Subject:            utils.GetStringPropOrEmpty(props, "subject"),
		Status:             utils.GetStringPropOrEmpty(props, "status"),
		MessageCount:       utils.GetInt64PropOrZero(props, "messageCount"),
		Source:             entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		ThreadId:           utils.GetStringPropOrEmpty(props, "threadId"),
		InitiatorFirstName: utils.GetStringPropOrEmpty(props, "initiatorFirstName"),
		InitiatorLastName:  utils.GetStringPropOrEmpty(props, "initiatorLastName"),
		InitiatorType:      utils.GetStringPropOrEmpty(props, "initiatorType"),
		InitiatorUsername:  utils.GetStringPropOrEmpty(props, "initiatorUsername"),
	}
	return &conversationEntity
}
