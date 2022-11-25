package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"reflect"
	"time"
)

type ConversationService interface {
	CreateNewConversation(ctx context.Context, userId string, contactId string, conversationId *string) (*entity.ConversationEntity, error)
	GetConversationsForUser(ctx context.Context, userId string, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error)
	GetConversationsForContact(ctx context.Context, contactId string, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error)
}

type conversationService struct {
	repository *repository.RepositoryContainer
}

func NewConversationService(repository *repository.RepositoryContainer) ConversationService {
	return &conversationService{
		repository: repository,
	}
}

func (s *conversationService) CreateNewConversation(ctx context.Context, userId string, contactId string, conversationId *string) (*entity.ConversationEntity, error) {
	if conversationId == nil {
		newUuid, _ := uuid.NewRandom()
		conversationId = utils.StringPtr(newUuid.String())
	}
	record, err := s.repository.ConversationRepository.Create(common.GetContext(ctx).Tenant, userId, contactId, *conversationId)
	if err != nil {
		return nil, err
	}
	conversationEntity := s.mapDbNodeToConversationEntity(utils.NodePtr(record.(dbtype.Node)))
	return conversationEntity, nil
}

func (s *conversationService) GetConversationsForUser(ctx context.Context, userId string, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(*s.repository.Drivers.Neo4jDriver)
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
		conversationEntity := *s.mapDbNodeToConversationEntity(v.Node)
		conversationEntity.UserId = v.UserId
		conversationEntity.ContactId = v.ContactId
		conversationEntities = append(conversationEntities, conversationEntity)
	}
	paginatedResult.SetRows(&conversationEntities)
	return &paginatedResult, nil
}

func (s *conversationService) GetConversationsForContact(ctx context.Context, contactId string, page, limit int, sortBy []*model.SortBy) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(*s.repository.Drivers.Neo4jDriver)
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
		conversationEntity := *s.mapDbNodeToConversationEntity(v.Node)
		conversationEntity.UserId = v.UserId
		conversationEntity.ContactId = v.ContactId
		conversationEntities = append(conversationEntities, conversationEntity)
	}
	paginatedResult.SetRows(&conversationEntities)
	return &paginatedResult, nil
}

func (s *conversationService) mapDbNodeToConversationEntity(dbNode *dbtype.Node) *entity.ConversationEntity {
	props := utils.GetPropsFromNode(*dbNode)
	conversationEntity := entity.ConversationEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		StartedAt: props["startedAt"].(time.Time),
	}
	return &conversationEntity
}
