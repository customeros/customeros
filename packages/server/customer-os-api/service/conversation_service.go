package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ConversationService interface {
	CreateNewConversation(ctx context.Context, userId string, contactId string, conversationId *string) (*entity.ConversationEntity, error)
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
		strUuid := newUuid.String()
		conversationId = &strUuid
	}
	record, err := s.repository.ConversationRepository.Create(common.GetContext(ctx).Tenant, userId, contactId, *conversationId)
	if err != nil {
		return nil, err
	}
	conversationEntity := s.mapDbNodeToConversationEntity(record.(dbtype.Node))
	return conversationEntity, nil
}

func (s *conversationService) mapDbNodeToConversationEntity(dbNode dbtype.Node) *entity.ConversationEntity {
	props := utils.GetPropsFromNode(dbNode)
	conversationEntity := entity.ConversationEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		StartedAt: utils.GetTimePropOrNow(props, "started"),
	}
	return &conversationEntity
}
